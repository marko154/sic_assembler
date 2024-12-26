package symtable

import (
	"fmt"
	"maps"
	"sic_assembler/internal/expr"
)

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

func (s *SymTable) Get(label string) int {
	return s.resolve(expr.Label(label), map[expr.Label]bool{})
}

func (symtab *SymTable) resolve(expression expr.Expr, seen map[expr.Label]bool) int {
	switch v := expression.(type) {
	case expr.Label:
		if _, ok := seen[v]; ok {
			panic(fmt.Sprintf("cyclical reference detected for label '%v': ", v))
		}
		nextSeen := maps.Clone(seen)
		nextSeen[v] = true
		if value, ok := symtab.table[string(v)]; ok {
			return symtab.resolve(value, nextSeen)
		}
		panic("undefined symbol")
	case expr.Number:
		return int(v)
	case expr.BinOp:
		switch v.Op {
		case "+":
			return symtab.resolve(v.Left, seen) + symtab.resolve(v.Right, seen)
		case "-":
			return symtab.resolve(v.Left, seen) - symtab.resolve(v.Right, seen)
		case "*":
			return symtab.resolve(v.Left, seen) * symtab.resolve(v.Right, seen)
		case "/":
			return symtab.resolve(v.Left, seen) / symtab.resolve(v.Right, seen)
		}
	}
	panic("invalid address type")
}

func (s *SymTable) Has(label string) bool {
	_, ok := s.table[label]
	return ok
}
