package parser

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"io"
	"regexp"
	"sic_assembler/internal/statement"
	"strconv"
	"strings"
	"unicode"
)

func Parse(reader io.Reader) ([]statement.IStatement, error) {
	scanner := bufio.NewScanner(reader)
	statements := []statement.IStatement{}

	for scanner.Scan() {
		line := scanner.Text()
		statement := parseLine(line)
		if statement != nil {
			statements = append(statements, statement)
			fmt.Printf("statement: %+v\n", statement)
		}
	}

	if err := scanner.Err(); err != nil {
		return statements, err
	}
	return statements, nil
}

func parseLine(line string) statement.IStatement {
	sourceLine := line
	line = stripComment(line)
	pattern := regexp.MustCompile(`(?P<label>\w+)?\s+(?P<mnemonic>\+?\w+)\s*(?P<args>.+)?`)
	matches := pattern.FindStringSubmatch(line)
	if matches == nil {
		fmt.Printf("No match '%v' (line): '%v' \n", sourceLine, (line))
		return nil
	}

	label := matches[1]
	mnemonic := matches[2]
	args := []string{}

	for _, arg := range strings.Split(matches[3], ",") {
		args = append(args, strings.TrimSpace(arg))
	}
	return parseIStatement(label, mnemonic, args, sourceLine)
}

func stripComment(line string) string {
	commentPattern := regexp.MustCompile(`\s*\..*`)
	return commentPattern.ReplaceAllString(line, "")
}

func parseIStatement(label, mnemonic string, args []string, sourceLine string) statement.IStatement {
	isExtended := false
	if mnemonic[0] == '+' {
		mnemonic = mnemonic[1:]
		isExtended = true
	}

	if opInfo, ok := statement.OPTAB[mnemonic]; ok {
		format := opInfo.Format
		if isExtended {
			format = 4
		}
		switch format {
		case 1:
			return statement.NewInstructionF1(label, mnemonic, opInfo.Opcode, sourceLine)
		case 2:
			return parseInstructionF2(label, mnemonic, args, opInfo.Opcode, sourceLine)
		case 3:
			return parseInstructionF3(label, mnemonic, args, opInfo.Opcode, sourceLine)
		case 4:
			return parseInstructionF4(label, mnemonic, args, opInfo.Opcode, sourceLine)
		}
	}

	if _, ok := statement.DIRECTIVES[mnemonic]; ok {
		return parseDirective(label, mnemonic, args, sourceLine)
	}

	if _, ok := statement.STORAGE[mnemonic]; ok {
		return parseStorage(label, mnemonic, args, sourceLine)
	}

	return nil
}

func parseInstructionF2(label, mnemonic string, args []string, opcode byte, source string) statement.IStatement {
	operand1 := parseRegister(args[0])
	operand2 := byte(0)
	if len(args) == 2 {
		operand2 = parseRegister(args[1])
	}
	return statement.NewInstructionF2(label, mnemonic, opcode, operand1, operand2, source)
}

func parseRegister(reg string) byte {
	if regIdx, ok := statement.REGTAB[reg]; ok {
		return regIdx
	}
	literal, err := strconv.ParseInt(reg, 10, 8)
	if err != nil {
		panic(err)
	}
	return byte(literal)
}

func parseInstructionF3(label, mnemonic string, args []string, opcode byte, source string) statement.IStatement {
	var operand statement.AddressOperand
	if len(args) >= 1 {
		operand = parseAddressOperand(args[0])
	}
	stmt := statement.NewInstructionF3(label, mnemonic, opcode, operand, source)
	if len(args) == 2 {
		stmt.IsIndexed = true
	}
	return stmt
}

func parseInstructionF4(label, mnemonic string, args []string, opcode byte, source string) statement.IStatement {
	var operand statement.AddressOperand
	if len(args) >= 1 {
		operand = parseAddressOperand(args[0])
	}
	stmt := statement.NewInstructionF4(label, mnemonic, opcode, operand, source)
	if len(args) == 2 {
		stmt.IsIndexed = true
	}
	return stmt
}

func parseAddressOperand(operand string) statement.AddressOperand {
	mode := statement.NORMAL
	if len(operand) == 0 {
		return statement.AddressOperand{
			Mode:    mode,
			Address: statement.Number(0),
		}
	}
	if operand[0] == '#' {
		operand = operand[1:]
		mode = statement.IMMEDIATE
	} else if operand[0] == '@' {
		operand = operand[1:]
		mode = statement.INDIRECT
	}

	return statement.AddressOperand{
		Mode:    mode,
		Address: parseAddress(operand),
	}
}

func parseAddress(address string) statement.Address {

	if unicode.IsLetter(rune(address[0])) {
		return statement.Label(address)
	}
	return parseNumber(address)
}

func parseNumber(address string) statement.Number {
	literal, err := strconv.ParseInt(address, 10, 32)
	if err != nil {
		panic(err)
	}
	return statement.Number(literal)
}

func parseDirective(label, mnemonic string, args []string, source string) statement.IStatement {
	return statement.NewDirective(label, mnemonic, parseAddress(args[0]), source)
}

func parseStorage(label, mnemonic string, args []string, source string) statement.IStatement {
	return statement.NewStorage(label, mnemonic, parseStorageOperand(mnemonic, args[0]), source)
}

func parseStorageOperand(mnemonic, arg string) statement.StorageOperand {
	if mnemonic == "BYTE" || mnemonic == "WORD" {
		return parseData(mnemonic, arg)
	}
	return parseNumber(arg)
}

func parseData(mnemonic, arg string) statement.Data {
	value, err := strconv.ParseInt(arg, 10, 24)
	if err == nil {
		if mnemonic == "BYTE" {
			return []byte{byte(value)}
		}
		// WORD
		return []byte{byte(value >> 16), byte(value >> 8), byte(value)}
	}
	/*
		byte[] parseData()
		– C'<chars>' … ASCII encoding
		– X'<hex>' … hex encoding
		– num … 24 bit number (WORD representation)
	*/
	switch arg[0] {
	case 'C':
		return []byte(arg[1:])
	case 'X':
		data, err := hex.DecodeString(arg[1:])
		if err != nil {
			panic(err)
		}
		return []byte(data)
	}
	panic("Invalid literal")
}
