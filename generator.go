package parsegen

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

func generateRule(it Iterator) (*Rule, error) {
	pd := it.Data()

	// TODO do const names
	base := pd.GetLabel("exprBase")
	cycle := pd.GetLabel("exprCycle")
	expr1 := pd.GetLabel("spRvalues")
	expr2 := pd.GetLabel("orRvalues")

	typ := byte('Z')
	if (!base.IsEmpty() || expr1.IsOnly()) && expr2.IsEmpty() {
		typ = 'N'
	} else if !base.IsEmpty() && expr2.IsOnly() {
		typ = 'L'
	} else if cycle.IsOnly() {
		typ = 'C'
	}

	if typ == 'Z' {
		return nil, fmt.Errorf("[gen] invalid expr, '%s'", pd.Get("expr"))
	}
	lid := string(pd.Get(lidTerm))
	if len(lid) < 1 {
		return nil, errors.New("[gen] got empty lvalue")
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
		if Debug {
			fmt.Println("[debug] rvalue", string(item))
		}
		if len(item) < 1 {
			return nil, errors.New("[gen] got empty rvalue")
		}
		// check is rid any
		endSymbol := byte(0)
		includeFlag := false
		if isTermAny(item, &endSymbol, &includeFlag) {
			rvalue = append(rvalue, term{typ: 'T', name: string(item), terminal: termAny(endSymbol, includeFlag)})
			continue LP
		}

		if isTermEmpty(item) {
			rvalue = append(rvalue, term{typ: 'T', name: string(item), terminal: termEmpty()})
			continue LP
		}

		start, end := byte(0), byte(0)
		if isHex(item, &start, &end) {
			if start == end {
				rvalue = append(rvalue, term{typ: 'T', name: string(item), terminal: termHex(start)})
			} else {
				rvalue = append(rvalue, term{typ: 'T', name: string(item), terminal: termHexes(start, end)})
			}
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
			return nil, fmt.Errorf("[gen] string must have at least one character, '%s'", item)
		}
		for _, str := range pd.GetAll("string") {
			if bytes.Equal(item, str) {
				value := removeQuotes(str)
				// todo generate name
				rvalue = append(rvalue, term{typ: 'T', name: value, terminal: termStr(value)})
				continue LP
			}
		}
		return nil, fmt.Errorf("[gen] invalid rvalue %s", item)
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
		return nil, errors.New("[gen] initial rule 'S' not found or there're more than one")
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
				return nil, fmt.Errorf("[gen] not resolved entity %s", subf.name)
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

func isTermAny(rvalue []byte, end *byte, includeEnd *bool) bool {
	// any(c)
	if len(rvalue) == 9 && bytes.HasPrefix(rvalue, anyPrefix) {
		if rvalue[3] == '(' && rvalue[8] == ')' {
			*includeEnd = false
		} else if rvalue[3] == '[' && rvalue[8] == ']' {
			*includeEnd = true
		} else {
			return false
		}
		b1, b2 := byte(0), byte(0)
		if isHex(rvalue[4:8], &b1, &b2) {
			*end = b1
			return b1 == b2
		}
	}
	return false
}

func isTermInteger(rvalue []byte) bool {
	return bytes.Equal(rvalue, integerTerm)
}

func isTermEmpty(rvalue []byte) bool {
	return bytes.Equal(rvalue, emptyTerm)
}
