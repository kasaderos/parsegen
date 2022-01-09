package main

import (
	"errors"
	"unicode"
)

type SimpleIterator struct {
	data []byte
	curr byte
	ind  int
	eof  bool

	parsed *ParsedData
}

type CommonIterator struct {
	data []byte
	curr byte
	ind  int
	eof  bool

	parsed *ParsedData
}

type Iterator interface {
	GC()
	CC() byte
	GP() int
	EOF() bool
	ParsedData() *ParsedData
}

func NewIterator(data []byte, includeSpaces bool) (Iterator, error) {
	if len(data) == 0 {
		return nil, errors.New("data is empty")
	}
	if includeSpaces {
		return &CommonIterator{data: data, curr: data[0], parsed: &ParsedData{data, make(map[string][]label)}}, nil
	}
	return &SimpleIterator{data: data, curr: data[0], parsed: &ParsedData{data, make(map[string][]label)}}, nil
}

// with ignoring spaces
func (it *SimpleIterator) GC() {
	it.ind++
	for !it.eof {
		if it.ind >= len(it.data) {
			it.eof = true
			it.ind = len(it.data)
			return
		}
		if !unicode.IsSpace(rune(it.CC())) {
			break
		}
		it.ind++
	}
}

func (it *SimpleIterator) EOF() bool {
	return it.eof
}

func (it *SimpleIterator) ParsedData() *ParsedData {
	return it.parsed

}
func (it *SimpleIterator) CC() byte {
	return it.curr
}

func (it *SimpleIterator) GP() int {
	return it.ind
}

func (it *CommonIterator) GC() {
	if it.ind >= len(it.data) {
		it.eof = true
		return
	}
	it.ind++
}

func (it *CommonIterator) EOF() bool {
	return it.eof
}

func (it *CommonIterator) ParsedData() *ParsedData {
	return it.parsed

}
func (it *CommonIterator) CC() byte {
	return it.curr
}

func (it *CommonIterator) GP() int {
	return it.ind
}
