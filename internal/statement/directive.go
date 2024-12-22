package statement

import "sic_assembler/internal/symtable"

type Directive struct {
	Label    string
	Mnemonic string
	Operand  Address
}

func NewDirective(label, mnemonic string, operand Address) *Directive {
	return &Directive{
		Label:    label,
		Mnemonic: mnemonic,
		Operand:  operand,
	}
}

func (i *Directive) EmitCode(symtable.SymTable, int, int, map[int]int) []byte {
	return []byte{}
}

func (i *Directive) GetLabel() string {
	return i.Label
}

func (i *Directive) ResolveOperand(symtab *symtable.SymTable) int {
	switch v := i.Operand.(type) {
	case Label:
		if value, ok := symtab.Get(string(v)); ok {
			return value
		}
		panic("undefined symbol")
	case Number:
		return int(v)
	}
	panic("invalid address type")
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
