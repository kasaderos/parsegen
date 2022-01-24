package parsegen

import (
	"fmt"
	"log"
)

type Parser struct {
	f *function
}

// PrintTree prints a graph of parser
func (p *Parser) PrintTree() {
	printTree(p.f)
}

// Parse parses data array by BNF.
// If the data structure does not match the bnf of the parser, it returns an error.
func (p *Parser) Parse(data []byte) (Data, error) {
	dataIt, err := NewIterator(data)
	if err != nil {
		return nil, err
	}

	if execute(p.f, dataIt) == missed {
		return nil, fmt.Errorf("[exec] data error, stopped at %d, '%c'", dataIt.GP(), dataIt.CC())
	}
	return dataIt.Data(), nil
}

// ParseAll parses data until it reaches the end
func (p *Parser) ParseAll(data []byte) (Data, error) {
	dataIt, err := NewIterator(data)
	if err != nil {
		return nil, err
	}

	for !dataIt.EOF() {
		if execute(p.f, dataIt) == missed {
			return nil, fmt.Errorf("[exec] data error, stopped at %d, '%c'", dataIt.GP(), dataIt.CC())
		}
		dataIt.Data().Print()
	}
	return dataIt.Data(), nil
}

// Generate creates a new parser by given bnf
func Generate(bnf []byte) (*Parser, error) {
	it, err := NewIterator(bnf)
	if err != nil {
		return nil, err
	}

	f, err := bnfFunction(it)
	if err != nil {
		return nil, err
	}

	rules := make([]*rule, 0)

	for !it.EOF() {
		ret := execute(f, it)
		if ret == missed {
			return nil, fmt.Errorf("[gen] not bnf rule, stopped at %d '%c'", it.GP(), it.CC())
		}

		rule, err := generateRule(it)
		if err != nil {
			return nil, err
		}
		it.Data().Clean()
		rules = append(rules, rule)
		log.Println("[gen] appended rule", rule.lvalue.name)
	}
	f, err = generateFunction(rules)
	if err != nil {
		return nil, err
	}

	return &Parser{f}, nil
}
