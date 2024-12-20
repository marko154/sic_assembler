package assembler

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"io"
)

const (
	T_RECORD_MAX_SIZE = (1 << 16) - 1
)

type TRecord struct {
	address int
	text    []byte
}

type RecordWriter struct {
	writer      *bufio.Writer
	currTRecord *TRecord
}

func NewRecordWriter(writer io.Writer) *RecordWriter {
	return &RecordWriter{
		writer:      bufio.NewWriter(writer),
		currTRecord: &TRecord{0, []byte{}},
	}
}

func (w *RecordWriter) Flush() error {
	return w.writer.Flush()
}

func (w *RecordWriter) writeLine(line string) {
	_, err := w.writer.WriteString(line + "\n")
	if err != nil {
		panic(err)
	}
}

func (w *RecordWriter) WriteHRecord(address int) {
	w.writeLine(
		fmt.Sprintf("H%06X\n", address),
	)
}

func (w *RecordWriter) WriteERecord(address int) {
	w.writeLine(
		fmt.Sprintf("E%06X\n", address),
	)
}

func (w *RecordWriter) WriteMRecord(address, nibbles int) {
	w.writeLine(
		fmt.Sprintf("M%06X%02X\n", address, nibbles),
	)
}

func (w *RecordWriter) WriteCode(address int, bytes []byte) {
	// flush current T record if there is a gap
	lastAddress := w.currTRecord.address + len(w.currTRecord.text)
	if address > lastAddress {
		w.flushTRecord()
		w.currTRecord = &TRecord{address, []byte{}}
	}

	// append to last record, flush when full
	remainingBytes := len(bytes)
	for remainingBytes > 0 {
		remainingSpace := T_RECORD_MAX_SIZE - len(w.currTRecord.text)
		bytesToWrite := min(remainingSpace, remainingBytes)
		remainingBytes -= bytesToWrite

		if len(w.currTRecord.text) == T_RECORD_MAX_SIZE {
			w.flushTRecord()
			address += T_RECORD_MAX_SIZE
			w.currTRecord = &TRecord{address, []byte{}}
		}
	}
}

func (w *RecordWriter) flushTRecord() {
	address := w.currTRecord.address
	bytes := w.currTRecord.text
	if len(bytes) == 0 {
		return
	}
	w.writeLine(
		fmt.Sprintf("T%06X%02X%s", address, len(bytes), hex.EncodeToString(bytes)),
	)
}
