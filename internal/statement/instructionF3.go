package statement

import (
	"fmt"
	"sic_assembler/internal/symtable"
)

type InstructionF3 struct {
	Label     string
	Mnemonic  string
	Opcode    byte
	Operand   AddressOperand
	IsIndexed bool
}

func NewInstructionF3(label, mnemonic string, opcode byte, operand AddressOperand) *InstructionF3 {
	return &InstructionF3{
		Label:    label,
		Mnemonic: mnemonic,
		Opcode:   opcode,
		Operand:  operand,
	}
}

func (i *InstructionF3) EmitCode(symtab symtable.SymTable, base, locctr int) []byte {
	// ni x bp e
	byte1 := i.Opcode | i.Operand.Mode
	address := i.resolveAddress(symtab)

	byte2 := byte(0)
	byte3 := byte(0)

	if -2048 <= address && address <= 2047 {
		// pc relative
		offset := address - locctr
		byte2 |= 0x0F & byte(offset>>8)
		byte3 = byte(offset)
		byte2 |= 0x20
	} else if 0 <= address-base && address-base < 4096 {
		// base relative
		offset := address - base
		byte2 |= 0x0F & byte(offset>>8)
		byte3 = byte(offset)
		byte2 |= 0x40
	} else if 0 <= address && address < 4096 {
		// direct (absolute) TODO: create M record for this (only if address was label)
		byte2 |= 0x0F & byte(address>>8)
		byte3 = byte(address)
	} else {
		// TODO: optional - SIC FORMAT
		panic(fmt.Sprintf("address out of range %v", address))
	}

	if i.IsIndexed {
		byte2 |= 0x80
	}
	return []byte{byte1, byte2, byte3}
}

func (i *InstructionF3) resolveAddress(symtab symtable.SymTable) int {
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

func (i *InstructionF3) GetLabel() string {
	return i.Label
}

func (i *InstructionF3) GetLocctr(prevLocctr int) int {
	return prevLocctr + 3
}
