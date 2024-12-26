package statement

import "sic_assembler/internal/symtable"

type EQU struct {
	Label  string
	Expr   string
	Source string
}

func NewEQU(label, expr, source string) *EQU {
	return &EQU{
		Label:  label,
		Expr:   expr,
		Source: source,
	}
}

func (i *EQU) EmitCode(symtable.SymTable, int, int, map[int]int) []byte {
	return []byte{}
}

func (i *EQU) GetLabel() string {
	return i.Label
}

func (i *EQU) GetLocctr(prevLocctr int) int {
	return prevLocctr
}

func (i *EQU) GetSource() string {
	return i.Source
}
