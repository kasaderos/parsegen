package main

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

func (it *IteratorStruct) HasError() bool {
	return it.err != nil
}

type SimpleIterator struct {
	*ParsedData
	IteratorStruct
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
	HasError() bool
	Labeler
	Data
}

func NewIterator(data []byte, includeSpaces bool) (Iterator, error) {
	if len(data) == 0 {
		return nil, errors.New("data is empty")
	}
	it := IteratorStruct{
		curr: data[0],
	}
	pd := &ParsedData{data, make(map[string]labels)}
	if includeSpaces {
		return &CommonIterator{pd, it}, nil
	}
	return &SimpleIterator{pd, it}, nil
}

func IsSpace(b byte) bool {
	return b == ' ' || b == '\n' || b == '\t'
}

// with ignoring spaces
func (it *SimpleIterator) GC() {
	if it.eof {
		return
	}
	ind := it.ind
	for !it.eof {
		if it.ind+1 >= len(it.data) {
			it.eof = true
			// not included
			it.ind = ind + 1
			return
		}
		it.ind++
		if !IsSpace(it.data[it.ind]) {
			break
		}
	}
	it.curr = it.data[it.ind]
}

func (it *SimpleIterator) EOF() bool {
	return it.eof
}

func (it *SimpleIterator) CC() byte {
	return it.curr
}

func (it *SimpleIterator) GP() int {
	return it.ind
}

func (it *SimpleIterator) BT(ind int) {
	it.ind = ind
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
}
