package main

import (
	"bytes"
	"errors"
	"fmt"
)

type term struct {
	typ      byte
	name     string
	marked   bool
	terminal tFunc
}

type Rule struct {
	lvalue term // non-terminal
	rvalue []term
}

func generateRules(it Iterator) (*Rule, error) {
	pd := it.Data()

	// TODO do const names
	base := pd.GetLabel("exprBase")
	cycle := pd.GetLabel("exprCycle")
	expr1 := pd.GetLabel("spRvalues")
	expr2 := pd.GetLabel("orRvalues")

	typ := byte('Z')
	if (!base.isEmpty() || expr1.isOnly()) && expr2.isEmpty() {
		typ = 'N'
	} else if !base.isEmpty() && expr2.isOnly() {
		typ = 'L'
	} else if cycle.isOnly() {
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

	if isCapital(lid[0]) {
		rule.lvalue.marked = true
	}

	rvalue := make([]term, 0)
LP:
	for _, item := range pd.GetAll("rvalue") {
		if len(item) < 1 {
			return nil, errors.New("rvalue empty")
		}
		// check is rid any
		endSymbol := byte(0)
		includeFlag := false
		if isTermAny(item, &endSymbol, &includeFlag) {
			rvalue = append(rvalue, term{typ: 'T', name: string(item), terminal: termAny(endSymbol, includeFlag)})
			continue LP
		}

		for _, rid := range pd.GetAll("rid") {
			if bytes.Equal(item, rid) {
				rvalue = append(rvalue, term{name: string(item)})
				continue LP
			}
		}
		// str: "a", "absdf"
		if len(item) < 3 {
			return nil, errors.New("str must have at least one character")
		}
		for _, str := range pd.GetAll("string") {
			if bytes.Equal(item, str) {
				value := removeQuotes(str)
				// todo generate name
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
	lvalueFuncs := []*function{}
	var initial *function
	count := 0
	lvalues := make(map[string]*function, 0)
	for _, rule := range rules {
		f := &function{}
		f.typ = rule.lvalue.typ
		f.name = rule.lvalue.name
		f.marked = rule.lvalue.marked

		lvalues[f.name] = f

		for _, rvalue := range rule.rvalue {
			f.funcs = append(f.funcs, &function{
				typ:      rvalue.typ,
				name:     rvalue.name,
				terminal: rvalue.terminal,
				marked:   rvalue.marked,
			})
		}

		lvalueFuncs = append(lvalueFuncs, f)
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
	for _, f := range lvalueFuncs {
		for _, subf := range f.funcs {
			if subf.isTerminal() {
				continue
			}

			f, ok := lvalues[subf.name]
			if !ok {
				return nil, fmt.Errorf("not resolved entity %s", subf.name)
			}
			subf.typ = f.typ

			if f.marked {
				subf.marked = true
			}

			subf.funcs = append(subf.funcs, f.funcs...)
		}
	}

	if err := checkBNF(initial); err != nil {
		return nil, err
	}
	return initial, nil
}

func isCapital(b byte) bool {
	return b >= 'A' && b <= 'Z'
}

func isTermAny(rid []byte, end *byte, includeEnd *bool) bool {
	// any(c)
	if len(rid) == 6 && bytes.HasPrefix(rid, anyPrefix) {
		if rid[3] == '(' && rid[5] == ')' {
			*end = rid[4]
			*includeEnd = false
			return true
		}
		if rid[3] == '[' && rid[5] == ']' {
			*end = rid[4]
			*includeEnd = true
			return true
		}
	}
	return false
}
