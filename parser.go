package main

import (
	"fmt"
)

type Parser struct {
	f *function
}

// TODO generate parser from ParsedData that contains bnf lexes
func NewParser(it Iterator) (*Parser, error) {
	return generateParser(it)
}

func (p *Parser) Parse(data []byte) (*ParsedData, error) {
	dataIt, err := NewIterator(data, true)
	if err != nil {
		return nil, err
	}

	// if execute(p.f, dataIt) ==  {
	// 	return nil, errors.New("exec data error")
	// }
	return dataIt.Data(), nil
}

func Generate(bnf []byte) (*Parser, error) {
	it, err := NewIterator(bnf, true)
	if err != nil {
		return nil, err
	}

	f, err := bnfparser(it)
	if err != nil {
		return nil, err
	}
	printTree(f)

	execute(f, it)
	fmt.Println(it.Data().labels)

	return nil, nil
	// TODO
	// return NewParser(it)
}
