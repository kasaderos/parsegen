package parsegen

import (
	"errors"
)

// IteratorStruct contains common fields to iterate by bytes slice
type IteratorStruct struct {
	curr byte
	ind  int
	eof  bool
	err  error
}

// CommonIterator implements Iterator interface
type CommonIterator struct {
	*ParsedData
	IteratorStruct
}

// Iterator is intended to iterate by bytes slice
type Iterator interface {
	// GC moves the pointer(index) to the next character
	// if current character is last it raises flag EOF.
	// The current symbol remains unchanged.
	GC()
	// CC returns the current char
	CC() byte
	// GP returns the current pointer(index)
	GP() int
	// EOF returns true returns true
	// if the slice of bytes is exhausted (GP() == lenght of data)
	EOF() bool
	// BT returns a pointer to the given position.
	// This also changes the current character to the character under the pointer.
	BT(int)
	Labeler
	Data() Data
}

// NewIterator is a constructor of Iterator
func NewIterator(data []byte) (Iterator, error) {
	if len(data) == 0 {
		return nil, errors.New("data is empty")
	}
	it := IteratorStruct{
		curr: data[0],
	}
	pd := &ParsedData{data, make(map[string]Label)}
	return &CommonIterator{pd, it}, nil
}

func IsSpace(b byte) bool {
	return b == ' ' || b == '\n' || b == '\t' || b == '\r'
}

func (it *CommonIterator) GC() {
	if it.ind+1 >= len(it.data) {
		it.eof = true
		it.ind = it.ind + 1
		return
	}
	it.ind++
	it.curr = it.data[it.ind]
}

func (it *CommonIterator) EOF() bool {
	return it.eof
}

func (it *CommonIterator) CC() byte {
	return it.curr
}

func (it *CommonIterator) GP() int {
	return it.ind
}

func (it *CommonIterator) BT(ind int) {
	it.ind = ind
	it.curr = it.data[ind]
}

func (it *CommonIterator) Data() Data {
	return it.ParsedData
}
