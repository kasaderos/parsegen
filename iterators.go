package parsegen

import (
	"errors"
	"log"
)

type IteratorStruct struct {
	curr byte
	ind  int
	eof  bool
	err  error
}

func (it *IteratorStruct) SetError(s string) {
	it.err = errors.New(s)
	log.Println(it.err)
}

type CommonIterator struct {
	*ParsedData
	IteratorStruct
}

type Iterator interface {
	GC()
	CC() byte
	GP() int
	EOF() bool
	BT(int)
	SetError(string)
	Labeler
	Data() Data
}

func NewIterator(data []byte) (Iterator, error) {
	if len(data) == 0 {
		return nil, errors.New("data is empty")
	}
	it := IteratorStruct{
		curr: data[0],
	}
	pd := &ParsedData{data, make(map[string]label)}
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
