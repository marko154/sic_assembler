package assembler

import (
	"bufio"
	"fmt"
	"io"
	"maps"
	"sic_assembler/internal/expr"
	"sic_assembler/internal/parser"
	"sic_assembler/internal/statement"
	"sic_assembler/internal/symtable"
	"slices"
	"unicode"
)

type Assembler struct {
	program   []statement.IStatement
	symtable  *symtable.SymTable
	writer    *bufio.Writer
	lstWriter *bufio.Writer
}

func NewAssembler(writer, lsWriter io.Writer) *Assembler {
	return &Assembler{
		program:   []statement.IStatement{},
		symtable:  symtable.NewSymTable(),
		writer:    bufio.NewWriter(writer),
		lstWriter: bufio.NewWriter(lsWriter),
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
			if equ, ok := stmt.(*statement.EQU); ok {
				expr := parser.ParseEQUExpr(equ.Expr, locctr)
				a.symtable.Set(label, expr, true)
			} else {
				a.symtable.Set(label, expr.Number(locctr), false)
			}

		}
		locctr = stmt.GetLocctr(locctr)
	}
}

// second pass - resolve symbols with unknown addresses
func (a *Assembler) resolve() {
	base := -1
	locctr := 0
	builder := NewTRecordBuilder()
	var hRecord *HRecord
	relocationTable := map[int]int{}

	for _, stmt := range a.program {
		if dir, ok := stmt.(*statement.Directive); ok {
			switch dir.Mnemonic {
			case "START":
				hRecord = NewHRecord(dir.Label, dir.ResolveOperand(a.symtable, locctr))
			case "BASE":
				base = dir.ResolveOperand(a.symtable, locctr)
			case "NOBASE":
				base = -1
			case "LTORG":
				panic("LTORG directive not supported")
			case "END":
				erecord := &ERecord{dir.ResolveOperand(a.symtable, locctr)}
				a.WriteRecords(hRecord, builder.GetRecords(), relocationTable, erecord)
			}
		}
		prevLocctr := locctr
		locctr = stmt.GetLocctr(locctr)
		bytes := stmt.EmitCode(*a.symtable, base, locctr, relocationTable)
		a.writeDebugInfo(prevLocctr, bytes, stmt)
		builder.WriteCode(prevLocctr, bytes)

		fmt.Printf("writing bytes for instruction: %X\n", bytes)
	}

	a.writer.Flush()
	a.lstWriter.Flush()
}

func (a *Assembler) WriteRecords(
	hrecord *HRecord,
	trecords []*TRecord,
	relocationTable map[int]int,
	erecord *ERecord,
) {
	records := []Record{hrecord}
	for _, record := range trecords {
		hrecord.endAddress = record.address + len(record.text)
		records = append(records, record)
	}
	addresses := slices.Collect(maps.Keys(relocationTable))
	slices.Sort(addresses)
	for _, address := range addresses {
		records = append(records, &MRecord{address: address, nibbles: relocationTable[address]})
	}

	records = append(records, erecord)

	for _, record := range records {
		line := record.Serialize()
		_, err := a.writer.WriteString(line + "\n")
		if err != nil {
			panic(err)
		}
	}
}

func (a *Assembler) writeDebugInfo(locctr int, bytes []byte, stmt statement.IStatement) {
	// address  emitted_code    "source code line" -> filename.lst:line_number
	output := fmt.Sprintf("%X", bytes)
	source := stmt.GetSource()
	// i don't know why the lines are formatted wrong, and i dont really care
	if unicode.IsSpace(rune(source[0])) {
		source = "\t" + source
	}
	_, err := a.lstWriter.WriteString(
		fmt.Sprintf("%06X\t%-10s\t%s\n", locctr, output, source),
	)
	if err != nil {
		panic(fmt.Sprintf("Error writing to lst file: %s", err))
	}
}
