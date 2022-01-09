package main

import (
	"errors"
)

type Parser struct {
	f *function
}

// TODO generate parser from ParsedData that contains bnf lexes
func NewParser(pd *ParsedData) (*Parser, error) {
	return generateParser(pd)
}

func (p *Parser) Parse(data []byte) (*ParsedData, error) {
	dataIt, err := NewIterator(data, true)
	if err != nil {
		return nil, err
	}

	if execute(p.f, dataIt) {
		return nil, errors.New("exec data error")
	}
	return dataIt.ParsedData(), nil
}

func Generate(bnf []byte) (*Parser, error) {
	it, err := NewIterator(bnf, false)
	if err != nil {
		return nil, err
	}

	f, _ := bnfparser(it)

	if execute(f, it) {
		// TODO add errors for exec
		return nil, errors.New("exec bnf error")
	}

	return NewParser(it.ParsedData())
}
