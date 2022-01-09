package main

type label struct {
	name string
	i, j int // i included, j not
}

type lex struct {
	name  string
	value []byte
}

type ParsedData struct {
	data   []byte
	labels map[string][]label // lvalue : labels
}

func (pd *ParsedData) GetLexes(entities []string, strIncluded bool, alloc bool) map[string][]lex {
	return nil
}
