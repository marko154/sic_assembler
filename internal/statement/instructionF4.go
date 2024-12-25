package statement

import (
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

func (i *InstructionF4) EmitCode(symtab symtable.SymTable, base, pc int, relocTable map[int]int) []byte {
	// TODO: handle literal and label like instructionF3
	byte1 := i.Opcode | i.Operand.Mode
	byte2, byte3, byte4 := i.resolveAddress(symtab, pc, relocTable)
	if i.IsIndexed {
		byte2 |= 0x80
	}
	return []byte{byte1, byte2, byte3, byte4}
}

func (i *InstructionF4) resolveAddress(symtab symtable.SymTable, pc int, relocTable map[int]int) (byte, byte, byte) {
	switch v := i.Operand.Address.(type) {
	case Label:
		if address, ok := symtab.Get(string(v)); ok {
			return i.resolveLabel(address, pc, relocTable)
		}
		panic("undefined symbol")
	case Number:
		return byte((v >> 16) & 0x0F), byte((v >> 8) & 0xFF), byte(v)
	}
	panic("invalid address type")
}

// TODO: this is duplicated in instructionF3.go
func (i *InstructionF4) resolveLabel(address, pc int, relocTable map[int]int) (byte, byte, byte) {
	byte2 := byte(0)
	byte3 := byte(0)
	byte4 := byte(0)

	// format 4 only supports direct (absolute) addressing
	byte2 |= 0x0F & byte(address>>16)
	byte3 = byte(address >> 8)
	byte4 = byte(address)
	relocTable[pc-3] = 5

	byte2 |= 0x10 // x bit
	return byte2, byte3, byte4
}

func (i *InstructionF4) GetLabel() string {
	return i.Label
}

func (i *InstructionF4) GetLocctr(prevLocctr int) int {
	return prevLocctr + 4
}
