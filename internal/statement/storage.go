package statement

import "sic_assembler/internal/symtable"

type Storage struct {
	Label    string
	Mnemonic string
	Data     []byte
}

func NewStorage(label, mnemonic string, data []byte) *Storage {
	return &Storage{
		Label:    label,
		Mnemonic: mnemonic,
		Data:     data,
	}
}

func (i *Storage) EmitCode(symtable.SymTable, int, int) []byte {
	return i.Data
}

func (i *Storage) GetLabel() string {
	return i.Label
}

func (i *Storage) GetLocctr(prevLocctr int) int {
	return prevLocctr + len(i.Data)
}
