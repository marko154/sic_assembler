package statement

import (
	"fmt"
	"sic_assembler/internal/symtable"
)

type InstructionF4 struct {
	Label     string
	Mnemonic  string
	Opcode    byte
	Operand   AddressOperand
	IsIndexed bool
}

func NewInstructionF4(label, mnemonic string, opcode byte, operand AddressOperand) *InstructionF4 {
	return &InstructionF4{
		Label:    label,
		Mnemonic: mnemonic,
		Opcode:   opcode,
		Operand:  operand,
	}
}

func (i *InstructionF4) EmitCode(symtab symtable.SymTable, base, locctr int) []byte {
	// ni xbpe
	byte1 := i.Opcode | i.Operand.Mode
	address := i.resolveAddress(symtab)

	byte2 := byte(0)
	byte3 := byte(0)
	byte4 := byte(0)

	// direct (absolute)
	byte2 |= 0x0F & byte(address>>16)
	byte3 = byte(address >> 8)
	byte4 = byte(address)

	// TODO: create M record for this (only if address was label)
	if i.IsIndexed {
		byte2 |= 0x80
	}
	byte2 |= 0x10
	return []byte{byte1, byte2, byte3, byte4}
}

// TODO: this is duplicated in instructionF3.go
func (i *InstructionF4) resolveAddress(symtab symtable.SymTable) int {
	switch v := i.Operand.Address.(type) {
	case Label:
		if address, ok := symtab.Get(string(v)); ok {
			return address
		}
		panic(fmt.Sprintf("undefined symbol: %s", v))
	case Number:
		return int(v)
	}
	panic("invalid address type")
}

func (i *InstructionF4) GetLabel() string {
	return i.Label
}

func (i *InstructionF4) GetLocctr(prevLocctr int) int {
	return prevLocctr + 4
}
