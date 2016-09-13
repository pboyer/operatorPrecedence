package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// operator precedence parsing
type exp interface {
	exp()
}

type lit struct {
	val int
}

type binOp struct {
	op       string
	lhs, rhs exp
}

func (*lit) exp()   {}
func (*binOp) exp() {}

var prec = map[string]int{
	"==": 10,
	"!=": 10,
	"+":  20,
	"-":  20,
	"*":  30,
	"/":  30,
}

var pos int
var tokens []string

func main() {
	str := "1 * 3 + 4"
	tokens = strings.Split(str, " ")
	e, err := parseExp()
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
	fmt.Println(str, "=", print(e), "=", eval(e))
}

func eval(e exp) int {
	switch et := e.(type) {
	case *binOp:
		switch et.op {
		case "+":
			return eval(et.lhs) + eval(et.rhs)
		case "-":
			return eval(et.lhs) - eval(et.rhs)
		case "*":
			return eval(et.lhs) * eval(et.rhs)
		case "/":
			return eval(et.lhs) / eval(et.rhs)
		}
	case *lit:
		return et.val
	}
	panic("Unknown exp type")
}

func print(e exp) string {
	switch et := e.(type) {
	case *binOp:
		return fmt.Sprintf("(%s %s %s)", print(et.lhs), et.op, print(et.rhs))
	case *lit:
		return fmt.Sprintf("%v", et.val)
	}
	panic("Unknown exp type")
}

func parseExp() (exp, error) {
	p, err := parsePrimary()
	if err != nil {
		return nil, err
	}
	return parseExp1(p, 0) //minPrec of 0 will cause all of the input to be consumed
}

func parseExp1(lhs exp, minPrec int) (exp, error) {
	la, ok := peek()
	for ok && prec[la] >= minPrec {
		op := la

		ok = consume()
		if !ok {
			return nil, fmt.Errorf("unexpected end of input")
		}

		rhs, err := parsePrimary()
		if !ok {
			return nil, err
		}

		// consume all of the higher precedence stuff, accumulating in rhs
		la, ok = peek()
		for ok && prec[la] >= prec[op] {
			rhs, err = parseExp1(rhs, prec[la])
			if err != nil {
				return nil, err
			}
			la, ok = peek()
		}

		lhs = &binOp{op, lhs, rhs}
	}
	return lhs, nil
}

func parsePrimary() (exp, error) {
	i, err := strconv.Atoi(tokens[pos])
	if err != nil {
		return nil, err
	}
	consume()
	return &lit{i}, nil
}

func consume() bool {
	pos++
	return pos <= len(tokens)
}

func peek() (string, bool) {
	if pos >= len(tokens) {
		return "", false
	}
	return tokens[pos], true
}
