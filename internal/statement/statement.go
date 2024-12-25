package statement

import "sic_assembler/internal/symtable"

type IStatement interface {
	EmitCode(symtable.SymTable, int, int, map[int]int) []byte
	GetLabel() string
	GetLocctr(int) int
	GetSource() string
}

const (
	INDIRECT  byte = 0b10
	IMMEDIATE byte = 0b01
	NORMAL    byte = 0b11
	SIC       byte = 0b00
)

type AddressOperand struct {
	Mode    byte
	Address Address
}

type Address interface {
	isAddress()
}

type StorageOperand interface {
	isStorageOperand()
}

type Label string

func (l Label) isAddress() {}

type Number int

func (l Number) isAddress() {}

func (l Number) isStorageOperand() {}

type Data []byte

func (l Data) isStorageOperand() {}
