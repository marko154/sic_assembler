package symtable

import "sic_assembler/internal/expr"

const (
	UNKNOWN = -1
)

type SymTable struct {
	table  map[string]expr.Expr
	isExpr map[string]bool
}

func NewSymTable() *SymTable {
	return &SymTable{
		table:  make(map[string]expr.Expr),
		isExpr: make(map[string]bool),
	}
}

func (s *SymTable) Set(label string, address expr.Expr, isExpr bool) {
	s.table[label] = address
	s.isExpr[label] = isExpr
}

func (s *SymTable) IsExpr(label string) bool {
	return s.isExpr[label]
}

func (s *SymTable) Get(label string) (int, bool) {
	value, ok := s.table[label]
	return s.resolve(value), ok
}

func (symtab *SymTable) resolve(expression expr.Expr) int {
	switch v := expression.(type) {
	case expr.Label:
		if value, ok := symtab.Get(string(v)); ok {
			return value
		}
		panic("undefined symbol")
	case expr.Number:
		return int(v)
	case expr.BinOp:
		switch v.Op {
		case "+":
			return symtab.resolve(v.Left) + symtab.resolve(v.Right)
		case "-":
			return symtab.resolve(v.Left) - symtab.resolve(v.Right)
		case "*":
			return symtab.resolve(v.Left) * symtab.resolve(v.Right)
		case "/":
			return symtab.resolve(v.Left) / symtab.resolve(v.Right)
		}
	}
	panic("invalid address type")
}

func (s *SymTable) Has(label string) bool {
	_, ok := s.table[label]
	return ok
}
