package main

type labels struct {
	i, j []int // i included, j not
}

type lex struct {
	name  string
	value []byte
}

type ParsedData struct {
	data   []byte
	labels map[string]labels // lvalue : labels
}

type Data interface {
	Data() *ParsedData
}

type Labeler interface {
	AddStart(string, int)
	AddEnd(string, int)
	GetLexes([]string, bool, bool) map[string][]lex
}

func (pd *ParsedData) GetLexes(entities []string, strIncluded bool, alloc bool) map[string][]lex {
	return nil
}

func (pd *ParsedData) AddStart(name string, ind int) {
	labels := pd.labels[name]
	labels.i = append(labels.i, ind)
	pd.labels[name] = labels
}

func (pd *ParsedData) AddEnd(name string, ind int) {
	labels := pd.labels[name]
	labels.j = append(labels.j, ind)
	pd.labels[name] = labels
}

func (pd *ParsedData) Data() *ParsedData {
	return pd
}
