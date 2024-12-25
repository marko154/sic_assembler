package statement

import (
	"fmt"
	"sic_assembler/internal/symtable"
)

type InstructionF3 struct {
	Label     string
	Mnemonic  string
	Opcode    byte
	Operand   AddressOperand
	IsIndexed bool
	Source    string
}

func NewInstructionF3(label, mnemonic string, opcode byte, operand AddressOperand, source string) *InstructionF3 {
	return &InstructionF3{
		Label:    label,
		Mnemonic: mnemonic,
		Opcode:   opcode,
		Operand:  operand,
		Source:   source,
	}
}

func (i *InstructionF3) EmitCode(symtab symtable.SymTable, base, pc int, relocTable map[int]int) []byte {
	byte1 := i.Opcode | i.Operand.Mode
	byte2, byte3 := i.resolveAddress(symtab, base, pc, relocTable)
	if i.IsIndexed {
		byte2 |= 0x80
	}
	return []byte{byte1, byte2, byte3}
}

func (i *InstructionF3) resolveAddress(symtab symtable.SymTable, base, pc int, relocTable map[int]int) (byte, byte) {
	switch v := i.Operand.Address.(type) {
	case Label:
		if address, ok := symtab.Get(string(v)); ok {
			return i.resolveLabel(address, base, pc, relocTable)
		}
		panic(fmt.Sprintf("undefined symbol: %s", v))
	case Number:
		return byte((v >> 8) & 0x0F), byte(v)
	}
	if i.Mnemonic == "RSUB" {
		return 0, 0
	}
	fmt.Printf("%+v\n", i)
	panic("invalid address type")
}

func (i *InstructionF3) resolveLabel(address, base, locctr int, relocTable map[int]int) (byte, byte) {
	byte2 := byte(0)
	byte3 := byte(0)
	if -2048 <= address && address <= 2047 {
		// pc relative
		offset := address - locctr
		byte2 |= 0x0F & byte(offset>>8)
		byte3 = byte(offset)
		byte2 |= 0x20
	} else if base != -1 && 0 <= address-base && address-base < 4096 {
		// base relative
		offset := address - base
		byte2 |= 0x0F & byte(offset>>8)
		byte3 = byte(offset)
		byte2 |= 0x40
	} else if 0 <= address && address < 4096 {
		// direct (absolute)
		byte2 |= 0x0F & byte(address>>8)
		byte3 = byte(address)
		relocTable[locctr-2] = 3
	} else {
		// TODO: optional - try SIC FORMAT which has 15 bits
		panic(fmt.Sprintf("address out of range '%v'", address))
	}

	if i.IsIndexed {
		byte2 |= 0x80
	}
	return byte2, byte3

}

func (i *InstructionF3) GetLabel() string {
	return i.Label
}

func (i *InstructionF3) GetLocctr(prevLocctr int) int {
	return prevLocctr + 3
}

func (i *InstructionF3) GetSource() string {
	return i.Source
}
