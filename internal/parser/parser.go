package parser

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"sic_assembler/internal/instruction"
)

func Parse(reader io.Reader) ([]instruction.Statement, error) {
	scanner := bufio.NewScanner(reader)
	instructions := []instruction.Statement{}

	for scanner.Scan() {
		line := scanner.Text()
		ins := parseLine(line)
		if ins != nil {
			instructions = append(instructions, *ins)
		}
	}

	if err := scanner.Err(); err != nil {
		return instructions, err
	}
	return instructions, nil
}

func parseLine(line string) *instruction.Statement {
	orig := line
	line = stripComment(line)
	pattern := regexp.MustCompile(`(?P<label>\w+)?\s+(?P<mnemonic>\w+)\s+(?P<args>.+)?`)
	matches := pattern.FindStringSubmatch(line)
	if matches == nil {
		fmt.Printf("No match '%v' (line): '%v' \n", orig, (line))
		return nil
	}

	label := matches[1]
	mnemonic := matches[2]
	args := matches[3]
	return instruction.NewStatement(label, mnemonic, args)
}

func stripComment(line string) string {
	commentPattern := regexp.MustCompile(`\s*\..*`)
	return commentPattern.ReplaceAllString(line, "")
}
