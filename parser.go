package main

import (
	"errors"
)

type Parser struct {
	f *function
}

func NewParser(f *ParsedData) *Parser {
	return nil
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

func generateParser(bnf []byte) (*Parser, error) {
	it, err := NewIterator(bnf, false)
	if err != nil {
		return nil, err
	}

	f, _ := bnfparser(it)

	if execute(f, it) {
		// TODO add errors for exec
		return nil, errors.New("exec bnf error")
	}

	return NewParser(it.ParsedData()), nil
}

func bnfparser(it Iterator) (*function, error) {
	rules := []Rule{
		{term{typ: 'N', name: "S"}, []term{
			{typ: 'N', name: "rule"},
		}},
		{term{typ: 'N', name: "rule"}, []term{
			{typ: 'L', name: "lvalue"},
			{typ: 'T', name: "assign", terminal: termStr("=")},
			{typ: 'N', name: "expr", terminal: termStr("=")},
			{typ: 'T', name: "end", terminal: termStr(";")},
		}},
		{term{typ: 'L', name: "lvalue"}, []term{
			{typ: 'T', name: "lid", terminal: termID()},
			{typ: 'N', name: "highlighted"},
		}},
		{term{typ: 'N', name: "highlighted"}, []term{
			{typ: 'T', name: "openH", terminal: termStr(">")},
			{typ: 'T', name: "lid", terminal: termID()},
			{typ: 'T', name: "openC", terminal: termStr(">")},
		}},

		{term{typ: 'N', name: "expr"}, []term{
			{typ: 'L', name: "rvalue"},
			{typ: 'L', name: "expr1"},
		}},
		{term{typ: 'L', name: "expr1"}, []term{
			{typ: 'L', name: "exprT1"},
			{typ: 'L', name: "exprT2"},
		}},
		{term{typ: 'L', name: "exprT1"}, []term{
			{typ: 'L', name: "rvalue"},
			{typ: 'T', name: "empty", terminal: termEmpty()},
			{typ: 'L', name: "exprT1"},
		}},
		{term{typ: 'L', name: "exprT2"}, []term{
			{typ: 'N', name: "rvalue1"},
			{typ: 'T', name: "empty", terminal: termEmpty()},
			{typ: 'L', name: "exprT2"},
		}},
		{term{typ: 'N', name: "rvalue1"}, []term{
			{typ: 'T', name: "signOr", terminal: termStr("|")},
			{typ: 'L', name: "rvalue"},
		}},
		{term{typ: 'L', name: "rvalue"}, []term{
			{typ: 'T', name: "rid", terminal: termID()},
			{typ: 'T', name: "string", terminal: termAnyQuoted()},
		}},
	}
	f, err := generateFunction(rules)
	if err != nil {
		return nil, err
	}
	return f, nil
}
