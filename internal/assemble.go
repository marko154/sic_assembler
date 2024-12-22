package assembler

import (
	"bufio"
	"fmt"
	"io"
	"sic_assembler/internal/parser"
	"sic_assembler/internal/statement"
	"sic_assembler/internal/symtable"
)

type Assembler struct {
	program  []statement.IStatement
	symtable *symtable.SymTable
	writer   *bufio.Writer
}

func NewAssembler(writer io.Writer) *Assembler {
	return &Assembler{
		program:  []statement.IStatement{},
		symtable: symtable.NewSymTable(),
		writer:   bufio.NewWriter(writer),
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
	builder := NewTRecordBuilder()
	var hRecord *HRecord
	// TODO: m records
	relocationTable := map[int]int{}

	for _, stmt := range a.program {
		if dir, ok := stmt.(*statement.Directive); ok {
			switch dir.Mnemonic {
			case "START":
				hRecord = NewHRecord(dir.Label, dir.ResolveOperand(a.symtable))
			case "BASE":
				base = dir.ResolveOperand(a.symtable)
			case "NOBASE":
				base = 0
			case "LTORG":
				panic("LTORG directive not supported")
			case "EQU":
				// TODO: add to symtable
				panic("EQU directive not supported")
			case "END":
				erecord := &ERecord{dir.ResolveOperand(a.symtable)}
				a.WriteRecords(hRecord, builder.GetRecords(), relocationTable, erecord)
			}
		}
		prevLocctr := locctr
		locctr = stmt.GetLocctr(locctr)
		bytes := stmt.EmitCode(*a.symtable, base, locctr, relocationTable)
		fmt.Printf("writing bytes for instruction: %X\n", bytes)
		builder.WriteCode(prevLocctr, bytes)
	}
}

func (a *Assembler) WriteRecords(hrecord *HRecord, trecords []*TRecord, relocationTable map[int]int, erecord *ERecord) {
	records := []Record{hrecord}
	for _, record := range trecords {
		hrecord.endAddress = record.address + len(record.text)
		records = append(records, record)
	}

	for address, nibbles := range relocationTable {
		records = append(records, &MRecord{address, nibbles})
	}

	records = append(records, erecord)

	for _, record := range records {
		line := record.Serialize()
		_, err := a.writer.WriteString(line + "\n")
		if err != nil {
			panic(err)
		}
	}
	a.writer.Flush()
}
