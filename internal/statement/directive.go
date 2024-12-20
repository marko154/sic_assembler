package statement

import "sic_assembler/internal/symtable"

type Directive struct {
	Label    string
	Mnemonic string
	Operand  int
}

func NewDirective(label, mnemonic string, operand int) *Directive {
	return &Directive{
		Label:    label,
		Mnemonic: mnemonic,
		Operand:  operand,
	}
}

func (i *Directive) EmitCode(symtable.SymTable, int, int) []byte {
	return []byte{}
}

func (i *Directive) GetLabel() string {
	return i.Label
}

func (i *Directive) GetLocctr(prevLocctr int) int {
	if i.Mnemonic == "START" || i.Mnemonic == "ORG" {
		return i.Operand
	}
	return prevLocctr
}
