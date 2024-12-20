package assembler

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"io"
)

type RecordWriter struct {
	writer *bufio.Writer
}

func NewRecordWriter(writer io.Writer) *RecordWriter {
	return &RecordWriter{
		writer: bufio.NewWriter(writer),
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
	// TODO: implement
}

func (w *RecordWriter) WriteTRecord(address int, bytes []byte) {
	w.writeLine(
		fmt.Sprintf("T%06X%02X%s", address, len(bytes), hex.EncodeToString(bytes)),
	)
}
