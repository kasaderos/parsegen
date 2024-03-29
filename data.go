package parsegen

import "fmt"

// Label contains entity bounds of parsing input
type Label struct {
	i, j []int // i included, j not
}

// IsOnly returns true if label has exactly one bound
func (lbls *Label) IsOnly() bool {
	return len(lbls.i) == 1 && len(lbls.j) == 1
}

func (lbls *Label) IsEmpty() bool {
	return len(lbls.i) < 1 && len(lbls.j) < 1
}

// GetBound returns the boundary of the i-th occurrence of the entity in the data array
// If label don't have i-th occurrence it returns zeros.
func (lbls *Label) GetBound(i int) (start int, end int) {
	if i < len(lbls.i) && i < len(lbls.j) {
		return lbls.i[i], lbls.j[i]
	}
	return 0, 0
}

type lex struct {
	name  string
	value []byte
}

// ParsedData implements Data interface
type ParsedData struct {
	data   []byte
	labels map[string]Label // lvalue : labels
}

// Data contains indices of exported entities.
// For example:
// S = Hello " " World;
// Hello = "hello" ;
// World = "world!!!" ;
// Hello, World entities are exported (titled)
// Then, with the input data "hello world!!!", we have indices,
// Hello : [0, 5), World : [5, 13)
//
// If entites cycled then we have indices [i0, j0], [i1, j1] ...
// Example :
// S = CycledAB ;
// CycledAB = { AB } ;
// AB = "AB" ;
// input : "ABAB"
// AB contains indices [0, 2], [2, 4]
type Data interface {
	// Get returns the first parsed subbytes of exported entity.
	// IMPORTANT: This function SLICES input slice bytes.
	// The result is READ ONLY data.
	// If you want to change you have to cop
	Get(string) []byte
	// GetAll returns all entries of exported entity.
	GetAll(string) [][]byte
	// Clean cleans label buffers
	Clean()
	// Print displays exported parsed entities
	Print()
	// GetLabel returns a label from given key.
	GetLabel(string) Label
}

// Labeler allows add bounds of entities.
type Labeler interface {
	// AppendStart appends a start index
	AppendStart(string, int)
	// AppendEnd appends an end index
	AppendEnd(string, int)
}

func (pd *ParsedData) Get(entity string) []byte {
	lbl, ok := pd.labels[entity]
	if !ok || len(lbl.i) < 1 || len(lbl.j) < 1 {
		return nil
	}
	return pd.data[lbl.i[0]:lbl.j[0]]
}

func (pd *ParsedData) GetLabel(key string) Label {
	lbl, _ := pd.labels[key]
	return lbl
}

func (pd *ParsedData) GetAll(entity string) [][]byte {
	lbl, ok := pd.labels[entity]
	if !ok || len(lbl.i) < 1 || len(lbl.j) < 1 {
		return nil
	}
	all := make([][]byte, 0)
	for k := 0; k < len(lbl.i) && k < len(lbl.j); k++ {
		all = append(all, pd.data[lbl.i[k]:lbl.j[k]])
	}
	return all
}

func (pd *ParsedData) AppendStart(name string, ind int) {
	labels := pd.labels[name]
	labels.i = append(labels.i, ind)
	pd.labels[name] = labels
}

func (pd *ParsedData) AppendEnd(name string, ind int) {
	labels := pd.labels[name]
	labels.j = append(labels.j, ind)
	pd.labels[name] = labels
}

func (pd *ParsedData) Print() {
	for key, value := range pd.labels {
		if value.IsEmpty() || key == "S" {
			continue
		}
		fmt.Printf("%s : \n", key)
		for i := 0; i < len(value.i); i++ {
			fmt.Printf("\t\t%s\n", pd.data[value.i[i]:value.j[i]])
		}
	}
}

func (pd *ParsedData) Clean() {
	for kv := range pd.labels {
		delete(pd.labels, kv)
	}
}
