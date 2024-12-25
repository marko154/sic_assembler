package statement

import "sic_assembler/internal/symtable"

type InstuctionF1 struct {
	Label    string
	Mnemonic string
	Opcode   byte
	Source   string
}

func NewInstructionF1(label, mnemonic string, opcode byte, source string) *InstuctionF1 {
	return &InstuctionF1{
		Label:    label,
		Mnemonic: mnemonic,
		Opcode:   opcode,
		Source:   source,
	}
}

func (i *InstuctionF1) EmitCode(symtable.SymTable, int, int, map[int]int) []byte {
	return []byte{i.Opcode}
}

func (i *InstuctionF1) GetLabel() string {
	return i.Label
}

func (i *InstuctionF1) GetLocctr(prevLocctr int) int {
	return prevLocctr + 1
}

func (i *InstuctionF1) GetSource() string {
	return i.Source
}
