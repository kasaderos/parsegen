package main

import "fmt"

type label struct {
	i, j []int // i included, j not
}

func (lbls *label) isOnly() bool {
	return len(lbls.i) == 1 && len(lbls.j) == 1
}

func (lbls *label) isEmpty() bool {
	return len(lbls.i) < 1 && len(lbls.j) < 1
}

type lex struct {
	name  string
	value []byte
}

type ParsedData struct {
	data   []byte
	labels map[string]label // lvalue : labels
}

type Data interface {
	Data() *ParsedData
	Reset()
}

type Labeler interface {
	SetStart(string, int)
	SetEnd(string, int)
	Get(string) []byte
	GetAll(string) [][]byte
}

// read only
func (pd *ParsedData) Get(key string) []byte {
	lbl, ok := pd.labels[key]
	if !ok || len(lbl.i) < 1 || len(lbl.j) < 1 {
		return nil
	}
	return pd.data[lbl.i[0]:lbl.j[0]]
}

func (pd *ParsedData) GetLabel(key string) label {
	lbl, _ := pd.labels[key]
	return lbl
}

func (pd *ParsedData) GetAll(key string) [][]byte {
	lbl, ok := pd.labels[key]
	if !ok || len(lbl.i) < 1 || len(lbl.j) < 1 {
		return nil
	}
	all := make([][]byte, 0)
	for k := 0; k < len(lbl.i) && k < len(lbl.j); k++ {
		all = append(all, pd.data[lbl.i[k]:lbl.j[k]])
	}
	return all
}

func (pd *ParsedData) SetStart(name string, ind int) {
	labels := pd.labels[name]
	labels.i = append(labels.i, ind)
	pd.labels[name] = labels
}

func (pd *ParsedData) SetEnd(name string, ind int) {
	labels := pd.labels[name]
	labels.j = append(labels.j, ind)
	pd.labels[name] = labels
}

func (pd *ParsedData) Data() *ParsedData {
	return pd
}

func (pd *ParsedData) Print() {
	for key, value := range pd.labels {
		if value.isEmpty() {
			continue
		}
		fmt.Printf("%s : \n", key)
		for i := 0; i < len(value.i); i++ {
			fmt.Printf("\t\t%s\n", pd.data[value.i[0]:value.j[0]])
		}
	}
}

func (pd *ParsedData) Reset() {
	for kv := range pd.labels {
		delete(pd.labels, kv)
	}
}
