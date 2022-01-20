package main

import (
	"bytes"
	"errors"
	"fmt"
	"unicode"
)

type term struct {
	typ      byte
	name     string
	marked   bool
	terminal tFunc
}

type BNFData struct{}

type Rule struct {
	lvalue term // non-terminal
	rvalue []term
}

func generateRules(it Iterator) (*Rule, error) {
	pd := it.Data()

	expr1 := pd.GetLabel("expr1")
	expr2 := pd.GetLabel("expr2")
	expr3 := pd.GetLabel("expr3")

	typ := byte('Z')
	if expr1.isOnly() {
		typ = 'N'
	} else if expr2.isOnly() {
		typ = 'L'
	} else if expr3.isOnly() {
		typ = 'C'
	}

	if typ == 'Z' {
		return nil, errors.New("invalid expr")
	}
	lid := string(pd.Get(lid))
	if len(lid) < 1 {
		return nil, errors.New("lvalue empty")
	}
	rule := &Rule{
		lvalue: term{typ: typ, name: lid},
	}

	if unicode.IsTitle(rune(lid[0])) {
		rule.lvalue.marked = true
	}

	rvalue := make([]term, 0)
LP:
	for _, item := range pd.GetAll("rvalue") {
		if len(item) < 1 {
			return nil, errors.New("rvalue empty")
		}
		for _, rid := range pd.GetAll("rid") {
			if bytes.Compare(item, rid) == 0 {
				rvalue = append(rvalue, term{name: string(item)})
				continue LP
			}
		}
		// str: "a", "absdf"
		if len(item) < 3 {
			return nil, errors.New("str must have at least one character")
		}
		for _, str := range pd.GetAll("string") {
			if bytes.Compare(item, str) == 0 {
				value := removeQuotes(str)
				rvalue = append(rvalue, term{typ: 'T', name: value, terminal: termStr(value)})
				continue LP
			}
		}
		return nil, errors.New("invalid rvalue")
	}

	rule.rvalue = rvalue
	return rule, nil
}

func removeQuotes(s []byte) string {
	return string(s[1 : len(s)-1])
}

// syntax analyzer
func generateFunction(rules []*Rule) (*function, error) {
	funcs := []*function{}
	var initial *function
	count := 0
	for _, rule := range rules {
		f := &function{}
		f.typ = rule.lvalue.typ
		f.name = rule.lvalue.name
		f.marked = rule.lvalue.marked
		for _, rvalue := range rule.rvalue {
			f.funcs = append(f.funcs, &function{
				typ:      rvalue.typ,
				name:     rvalue.name,
				terminal: rvalue.terminal,
				marked:   rvalue.marked, // if lvalue marked => all rvalues marked
			})
		}
		funcs = append(funcs, f)

		// must exactly one
		if f.name == "S" {
			initial = f
			count++
		}
	}
	if count != 1 {
		return nil, errors.New("initial rule 'S' not found or there're more than one")
	}

	// brute force O(n^3)
	// rules number ~ 100
	for _, f := range funcs {
		if f.isTerminal() {
			if f.existFunc(0) {
				return nil, errors.New("terminal had sub funcs")
			}
			continue
		}
		// all sub funcs must be unique
		// can't be unnamed
		// funcNames := make(map[string]struct{}, 0)
		// funcNames[""] = struct{}{}
		for i, subf := range f.funcs {
			// if _, ok := funcNames[subf.name]; ok {
			// return nil, errors.New("all rvalue must be unique")
			// }

			if subf.isTerminal() {
				if subf.existFunc(0) {
					return nil, errors.New("terminal had sub funcs")
				}
				continue
			}

			if subf.isCycle() && subf.hasNext(0) {
				return nil, errors.New("cycle has more than one rvalue")
			}

			// funcNames[subf.name] = struct{}{}
			found := false
			// find subf from all funcs
			for _, fn := range funcs {
				// if it's found append
				if subf.name == fn.name {
					subf.funcs = append(subf.funcs, fn.funcs...)
					found = true
					break
				}
				// TODO check leaf are terminals
				// use printTree func
			}
			if !found {
				return nil, errors.New(fmt.Sprintf("not resolved entity %s", subf.name))
			}
			f.funcs[i] = subf
		}
	}
	return initial, nil
}
