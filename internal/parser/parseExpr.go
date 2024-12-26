package parser

import (
	"fmt"
	"sic_assembler/internal/expr"
	"strconv"
	"strings"
	"unicode"
)

// support +, -, *, /, parentheses

func ParseEQUExpr(input string, pc int) expr.Expr {
	input = strings.ReplaceAll(input, "(", " ( ")
	input = strings.ReplaceAll(input, ")", " ) ")
	input = strings.ReplaceAll(input, "+", " + ")
	input = strings.ReplaceAll(input, "-", " - ")
	input = strings.ReplaceAll(input, "*", " * ")

	tokens := strings.Fields(input)
	fmt.Printf("parseEQUExpr: %+q\n", tokens)

	_, expr := parseExpr(tokens, pc)
	return expr
}

/*
expr := factor | expr + factor | expr - factor
factor := term | term * factor | term / factor
term = * | label | int | ( expr )
*/

func parseExpr(tokens []string, pc int) ([]string, expr.Expr) {
	tokens, left := parseFactor(tokens, pc)

	for len(tokens) > 0 && (tokens[0] == "+" || tokens[0] == "-") {
		op := tokens[0]
		tokensLeft, right := parseFactor(tokens[1:], pc)
		tokens = tokensLeft
		left = expr.BinOp{Left: left, Op: op, Right: right}
	}
	return tokens, left
}

func parseFactor(tokens []string, pc int) ([]string, expr.Expr) {
	tokens, left := parseTerm(tokens, pc)

	for len(tokens) > 0 && (tokens[0] == "*" || tokens[0] == "/") {
		op := tokens[0]
		tokensLeft, right := parseTerm(tokens[1:], pc)
		tokens = tokensLeft
		left = expr.BinOp{Left: left, Op: op, Right: right}
	}

	return tokens, left
}

func parseTerm(tokens []string, pc int) ([]string, expr.Expr) {
	token := tokens[0]
	if token == "*" {
		return tokens[1:], expr.Number(pc)
	}
	if unicode.IsLetter(rune(token[0])) {
		return tokens[1:], expr.Label(token)
	}
	if unicode.IsDigit(rune(token[0])) {
		v, err := strconv.Atoi(token)
		if err != nil {
			panic(fmt.Sprintf("error parsing number in EQU expression: %s", err))
		}
		return tokens[1:], expr.Number(v)
	}
	if token == "(" {
		tokensLeft, expr := parseExpr(tokens[1:], pc)
		if tokensLeft[0] != ")" {
			panic("expected closing parenthesis")
		}
		return tokensLeft[1:], expr
	}
	panic("unexpected token in EQU expression")
}
