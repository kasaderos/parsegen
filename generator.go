package main

import (
	"errors"
	"fmt"
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

func generateParser(lexes map[string][]lex) (*Parser, error) {
	// TODO BNFData
	rules, err := generateRules(&BNFData{})
	if err != nil {
		return nil, err
	}

	f, err := generateFunction(rules)
	if err != nil {
		return nil, err
	}

	return &Parser{f}, nil
}

func generateRules(bnfData *BNFData) ([]Rule, error) {
	return nil, nil
}

// syntax analyzer
func generateFunction(rules []Rule) (*function, error) {
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
				marked:   rvalue.marked,
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
		funcNames := make(map[string]struct{}, 0)
		funcNames[""] = struct{}{}
		for i, subf := range f.funcs {
			if _, ok := funcNames[subf.name]; ok {
				return nil, errors.New("all rvalue must be unique")
			}

			if subf.isTerminal() {
				if subf.existFunc(0) {
					return nil, errors.New("terminal had sub funcs")
				}
				continue
			}

			funcNames[subf.name] = struct{}{}
			found := false
			// find subf from all funcs
			for _, fn := range funcs {
				// if it's found append
				if subf.name == fn.name && subf.typ == fn.typ {
					subf.funcs = append(subf.funcs, fn.funcs...)
					found = true
					break
				}
				// TODO below code don't work because top does this
				// but without checks

				// if it's the same func it must be in the end
				if subf.name == f.name {
					// must be typ L
					if f.typ != 'L' {
						return nil, errors.New("invalid recursive call")
					}
					// defined in the end
					if i != len(f.funcs)-1 {
						return nil, errors.New("self func must defined in the end")
					}
				}
			}
			if !found {
				return nil, errors.New(fmt.Sprintf("not resolved entity %s", subf.name))
			}
			f.funcs[i] = subf
		}
	}
	return initial, nil
}
