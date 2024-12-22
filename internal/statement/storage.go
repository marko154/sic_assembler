package statement

import "sic_assembler/internal/symtable"

type Storage struct {
	Label    string
	Mnemonic string
	Operand  StorageOperand
}

func NewStorage(label, mnemonic string, operand StorageOperand) *Storage {
	return &Storage{
		Label:    label,
		Mnemonic: mnemonic,
		Operand:  operand,
	}
}

func (i *Storage) EmitCode(symtable.SymTable, int, int, map[int]int) []byte {
	if i.Mnemonic == "RESW" || i.Mnemonic == "RESB" {
		return []byte{}
	}
	switch v := i.Operand.(type) {
	case Data:
		return []byte(v)
	}
	panic("invalid storage operand")
}

func (i *Storage) GetLabel() string {
	return i.Label
}

func (i *Storage) GetLocctr(prevLocctr int) int {
	if i.Mnemonic == "RESW" {
		v := i.Operand.(Number)
		return prevLocctr + 3*int(v)
	}
	if i.Mnemonic == "RESB" {
		v := i.Operand.(Number)
		return prevLocctr + int(v)
	}
	bytes := i.Operand.(Data)
	return prevLocctr + len(bytes)
}
