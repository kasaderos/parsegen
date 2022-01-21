package main

import (
	"errors"
	"fmt"
)

type Parser struct {
	f *function
}

func (p *Parser) Parse(data []byte) (*ParsedData, error) {
	dataIt, err := NewIterator(data)
	if err != nil {
		return nil, err
	}

	if execute(p.f, dataIt) == missed {
		return nil, errors.New("exec data error")
	}
	return dataIt.Data(), nil
}

func Generate(bnf []byte) (*Parser, error) {
	it, err := NewIterator(bnf)
	if err != nil {
		return nil, err
	}

	f, err := bnfFunction(it)
	if err != nil {
		return nil, err
	}

	rules := make([]*Rule, 0)
	ret := execute(f, it)
	fmt.Println(it.Data().labels)
	if ret == missed {
		return nil, errors.New("not bnf rule")
	}

	rule, err := generateRules(it)
	if err != nil {
		return nil, err
	}
	rules = append(rules, rule)

	f, err = generateFunction(rules)
	printTree(f)
	if err != nil {
		return nil, err
	}

	return &Parser{f}, nil
}
