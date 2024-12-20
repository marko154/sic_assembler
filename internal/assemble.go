package assembler

import (
	"fmt"
	"io"
	"sic_assembler/internal/parser"
	"sic_assembler/internal/statement"
	"sic_assembler/internal/symtable"
)

type Assembler struct {
	program  []statement.IStatement
	symtable *symtable.SymTable
	rw       *RecordWriter
}

func NewAssembler(writer io.Writer) *Assembler {
	return &Assembler{
		program:  []statement.IStatement{},
		symtable: symtable.NewSymTable(),
		rw:       NewRecordWriter(writer),
	}
}

// BASE addressing only if BASE offset from base directive is small enough

func (a *Assembler) Assemble(reader io.Reader) {
	statements, err := parser.Parse(reader)
	if err != nil {
		panic(err)
	}
	a.program = statements
	a.activate()
	a.resolve()
}

// first pass
func (a *Assembler) activate() {
	locctr := 0

	for _, stmt := range a.program {
		label := stmt.GetLabel()
		if label != "" {
			if a.symtable.Has(label) {
				panic(fmt.Sprintf("Duplicate symbol '%s'", label))
			}
			a.symtable.Set(label, locctr)
		}
		locctr = stmt.GetLocctr(locctr)
	}
}

// second pass - resolve symbols with unknown addresses
func (a *Assembler) resolve() {
	base := 0
	locctr := 0

	// check if any symbol remains unresolved
	for _, stmt := range a.program {
		if dir, ok := stmt.(*statement.Directive); ok {
			switch dir.Mnemonic {
			case "START":
				a.rw.WriteHRecord(dir.Operand)
			case "END":
				a.rw.WriteERecord(dir.Operand)
			case "BASE":
				base = dir.Operand
			case "NOBASE":
				base = 0
			case "LTORG":
				panic("LTORG directive not supported")
			case "EQU":
				// TODO: add to symtable
				panic("EQU directive not supported")
			}
		}

		locctr = stmt.GetLocctr(locctr)
		bytes := stmt.EmitCode(*a.symtable, base, locctr)
		a.rw.WriteCode(locctr, bytes)
	}
}
