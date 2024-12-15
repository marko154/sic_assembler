package assembler

import (
	"fmt"
	"io"
	"sic_assembler/internal/instruction"
	"sic_assembler/internal/parser"
)

type Assembler struct {
	program  []instruction.Statement
	symtable *SymTable
}

func NewAssembler() *Assembler {
	return &Assembler{
		program:  []instruction.Statement{},
		symtable: NewSymTable(),
	}
}

func (a *Assembler) Assemble(reader io.Reader) {
	instructions, err := parser.Parse(reader)
	if err != nil {
		panic(err)
	}
	a.program = instructions
	fmt.Println(a.program)
	a.markSymbols()
	a.resolve()
}

// first pass
func (a *Assembler) markSymbols() {
	locctr := 0

	for _, ins := range a.program {
		if ins.Label != "" {
			// TODO: handle duplicate symbols
			a.symtable.Set(ins.Label, locctr)
		}

		/*
			if isdirective {
				no idea ? handle it
			}

			if is
		*/
	}
}

// second pass - resolve symbols with unknown addresses
func (a *Assembler) resolve() {
}
