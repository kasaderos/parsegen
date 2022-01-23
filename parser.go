package main

import (
	"errors"
	"fmt"
	"log"
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

func (p *Parser) ParseAll(data []byte) (*ParsedData, error) {
	dataIt, err := NewIterator(data)
	if err != nil {
		return nil, err
	}

	for !dataIt.EOF() {
		if execute(p.f, dataIt) == missed {
			return nil, errors.New("exec data error")
		}
		dataIt.Data().Print()
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

	for !it.EOF() {
		ret := execute(f, it)
		fmt.Println(it.Data().labels)
		if ret == missed {
			return nil, errors.New("not bnf rule")
		}

		rule, err := generateRule(it)
		if err != nil {
			return nil, err
		}
		it.Data().Reset()
		rules = append(rules, rule)
		log.Println("appended rule", rule.lvalue.name)
	}
	f, err = generateFunction(rules)
	if err != nil {
		return nil, err
	}
	printTree(f)

	return &Parser{f}, nil
}
