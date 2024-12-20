package statement

import "sic_assembler/internal/symtable"

type InstructionF2 struct {
	Label    string
	Mnemonic string
	Opcode   byte
	Operand1 byte
	Operand2 byte
}

func NewInstructionF2(label, mnemonic string, opcode, operand1, operand2 byte) *InstructionF2 {
	return &InstructionF2{
		Label:    label,
		Mnemonic: mnemonic,
		Opcode:   opcode,
		Operand1: operand1,
		Operand2: operand2,
	}
}

func (i *InstructionF2) EmitCode(symtable.SymTable, int, int) []byte {
	return []byte{i.Opcode, (i.Operand1 << 4) | i.Operand2}
}

func (i *InstructionF2) GetLabel() string {
	return i.Label
}

func (i *InstructionF2) GetLocctr(prevLocctr int) int {
	return prevLocctr + 2
}