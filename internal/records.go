package assembler

import (
	"encoding/hex"
	"fmt"
	"strings"
)

const (
	T_RECORD_MAX_SIZE = 0x1E // not sure why
)

type Record interface {
	Serialize() string
}

type HRecord struct {
	name    string
	address int
	length  int
}

func NewHRecord(name string, address int) *HRecord {
	return &HRecord{
		name:    name,
		address: address,
		length:  0,
	}
}

func (r *HRecord) Serialize() string {
	name := r.name[:min(len(r.name), 6)]
	if len(name) == 0 {
		name = "null"
	}
	return fmt.Sprintf("H%-6s%06X%06X", name, r.address, r.length)
}

type ERecord struct {
	address int
}

func (r *ERecord) Serialize() string {
	return fmt.Sprintf("E%06X", r.address)
}

type TRecord struct {
	address int
	text    []byte
}

func (r *TRecord) Serialize() string {
	encoded := strings.ToUpper(hex.EncodeToString(r.text))
	return fmt.Sprintf("T%06X%02X%s", r.address, len(r.text), encoded)
}

type MRecord struct {
	address int
	nibbles int
}

func (r *MRecord) Serialize() string {
	return fmt.Sprintf("M%06X%02X", r.address, r.nibbles)
}

type TRecordBuilder struct {
	currTRecord *TRecord
	records     []*TRecord
}

func NewTRecordBuilder() *TRecordBuilder {
	return &TRecordBuilder{
		currTRecord: &TRecord{0, []byte{}},
	}
}

func (b *TRecordBuilder) GetRecords() []*TRecord {
	b.appendRecord()
	return b.records
}

func (w *TRecordBuilder) WriteCode(address int, bytes []byte) {
	// flush current T record if there is a gap
	fmt.Println("writing code at address", address, bytes)
	lastAddress := w.currTRecord.address + len(w.currTRecord.text)
	if address > lastAddress {
		w.appendRecord()
		w.currTRecord = &TRecord{address, []byte{}}
	}

	// append to last record, flush when full
	for len(bytes) > 0 {
		remainingSpace := T_RECORD_MAX_SIZE - len(w.currTRecord.text)
		bytesToWrite := min(remainingSpace, len(bytes))
		fmt.Println("buffering", bytesToWrite, "bytes")
		w.currTRecord.text = append(w.currTRecord.text, bytes[:bytesToWrite]...)
		bytes = bytes[bytesToWrite:]

		if len(w.currTRecord.text) == T_RECORD_MAX_SIZE {
			address := w.currTRecord.address + len(w.currTRecord.text)
			w.appendRecord()
			w.currTRecord = &TRecord{address, []byte{}}
		}
	}
}

func (w *TRecordBuilder) appendRecord() {
	bytes := w.currTRecord.text
	if len(bytes) == 0 {
		return
	}
	w.records = append(w.records, w.currTRecord)
}
