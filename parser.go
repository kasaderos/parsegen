package main

import (
	"bytes"
	"errors"
	"io"
)

var ErrEmptyLine = errors.New("empty line")

type Type int32

// if typ == string -> value of string
type entity struct {
	typ Type

	// cycle parameters
	minRepeats int // -2 ~ 0..inf is kind of rule: ( ENTITY )*
	maxRepeats int // -1 ~ inf
	value      []byte
}

// Rule is a set of word like these:
// S ::= A B C
// lvalue S
// rvalue A B C
// type Rule struct {
// 	Marked bool
// 	lvalue entity
// 	rvalue []entity
// }

// \n\n\n
func readLine(data *[]byte) ([]byte, error) {
	ind := bytes.Index(*data, []byte("\n"))
	if ind < 0 {
		return nil, io.EOF
	}
	line := (*data)[:ind]
	*data = (*data)[ind+1:]
	return line, nil
}

// type state int32

// const (
// 	Init state = iota
// 	Final
// )

type Parser struct {
	line []byte
	i    int
	eof  bool

	err error
}

func NewParser(line []byte) *Parser {
	return &Parser{line, 0, false, nil}
}

func (p *Parser) cc() byte {
	return p.line[p.i]
}

func (p *Parser) gp() int {
	return p.i
}

func (p *Parser) slice(start, end int) []byte {
	return p.line[start:end]
}

func (p *Parser) gc() bool {
	if p.i >= len(p.line) {
		p.eof = true
		return true
	}
	p.i++
	return false
}

// ____FAFDFD__
func (p *Parser) Space(must bool) bool {
	exist := false
	for !p.eof {
		char := p.cc()
		if char == ' ' {
			exist = true
		} else if must && !exist {
			p.err = errors.New("not space at line %d TODO")
			return true
		} else {
			return p.eof
		}
		p.gc()
	}
	return p.eof
}

func isAlpha(b byte) bool {
	return b >= 'a' && b <= 'z' || b >= 'A' && b <= 'Z'
}

func (p *Parser) Entity(lvalue bool) bool {
	start := p.gp()
	for !p.eof {
		char := p.cc()
		if !isAlpha(char) {
			break
		}
		p.gc()
	}
	end := p.gp()

	if start == end {
		if lvalue {
			p.err = errors.New("lvalue can't be empty")
			return true
		}
		return p.eof
	}

	// value := p.slice(start, end)
	// if lvalue {
	// p.Rule.lvalue.value = value
	// } else {
	// p.Rule.rvalue = append(p.Rule.rvalue, entity{value: value})
	// }
	return p.eof
}

func (p *Parser) Lvalue() bool {
	return p.Entity(true)
}

func (p *Parser) Delimeter() bool {
	first := p.cc()
	p.gc()
	second := p.cc()
	p.gc()

	if first != ':' || second != '=' {
		p.err = errors.New("invalid delimeter")
		return true
	}

	return p.eof
}

func (p *Parser) Rvalues() bool {
	return p.Space(true) || p.Rvalue() || p.Rvalues()
}

// Any parses words that starts with 'starts' and ends with 'ends'
// If any has two or more ends, it stops in first 'ends' byte
func (p *Parser) Any(starts, ends byte) bool {
	if p.cc() != starts {
		p.err = errors.New("invalid start of 'any' entity")
		return true
	}
	p.gc()

	start := p.gp()
	end := start
	ended := false
	for !p.eof {
		if p.cc() == ends {
			ended = true
			p.gc()
			break
		}
		p.gc()
		end = p.gp()
	}

	if !ended {
		p.err = errors.New("invalid end of 'any' entity")
		return true
	}

	if start == end {
		// feature 'any' might be empty inside
		return p.eof
	}

	// newEntity := entity{
	// 	typ:   String,
	// 	value: p.slice(start, end),
	// }
	// p.Rule.rvalue = append(p.Rule.rvalue, newEntity)
	return p.eof
}

func (p *Parser) String() bool {
	// without quotes inside (symbol '"')
	return p.Any('"', '"')
}

func (p *Parser) Rvalue() bool {
	if !p.Entity(false) || !p.String() {
		return p.eof
	}
	return true
}

// // input data is a line without '\n'
// // rule := space lvalue space ":=" rvalues
// // lvalue := entity
// // rvalues := (space rvalue)(1:n)*
// // rvalue := entity | string
// // entity := word
// // string := <"> any <"> without quotes
// func lexic(line []byte) (*Rule, error) {
// 	return nil, nil
// 	// p := NewParser(line)

// 	// p.Space(false)
// 	// p.Lvalue()
// 	// p.Space(true)
// 	// p.Delimeter()
// 	// p.Rvalues()
// 	// return p.Rule, p.err
// }

// func getRule(line []byte) (*Rule, error) {
// 	if len(line) > 0 {
// 		return lexic(line)
// 	}
// 	return nil, ErrEmptyLine
// }

// // getRules parses from given bnf.
// func getRules(bnf []byte) ([]*Rule, error) {
// 	rules := make([]*Rule, 0)
// 	for {
// 		line, err := readLine(&bnf)
// 		if err == io.EOF {
// 			break
// 		}

// 		rule, err := getRule(line)
// 		if err != nil && err != ErrEmptyLine {
// 			return nil, err
// 		}
// 		rules = append(rules, rule)
// 	}

// 	// if err := syntax(rules); err != nil {
// 	// return nil, err
// 	// }
// 	return rules, nil
// }

func genCC(ptr **byte) func() byte {
	return func() byte {
		return **ptr
	}
}

func genGC(line []byte) (func(), *byte) {
	i := 0
	ptr := &line[0]
	return func() {
		if i >= len(line) {
			return
		}
		i++
		ptr = &line[i]
	}, ptr
}
