package statement

import "sic_assembler/internal/symtable"

type Directive struct {
	Label    string
	Mnemonic string
	Operand  Address
	Source   string
}

func NewDirective(label, mnemonic string, operand Address, source string) *Directive {
	return &Directive{
		Label:    label,
		Mnemonic: mnemonic,
		Operand:  operand,
		Source:   source,
	}
}

func (i *Directive) EmitCode(symtable.SymTable, int, int, map[int]int) []byte {
	return []byte{}
}

func (i *Directive) GetLabel() string {
	return i.Label
}

func (i *Directive) ResolveOperand(symtab *symtable.SymTable, pc int) int {
	switch v := i.Operand.(type) {
	case Label:
		if value, ok := symtab.Get(string(v)); ok {
			return value
		}
		panic("undefined symbol")
	case Number:
		return int(v)
	}
	panic("invalid operand type")
}

func (i *Directive) GetLocctr(prevLocctr int) int {
	if i.Mnemonic == "START" || i.Mnemonic == "ORG" {
		value, ok := i.Operand.(Number)
		if !ok {
			panic("invalid operand, expected number, got label")
		}
		return int(value)
	}
	return prevLocctr
}

func (i *Directive) GetSource() string {
	return i.Source
}
