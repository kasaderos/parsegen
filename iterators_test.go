package main

import (
	"testing"
)

func TestCommonIterator(t *testing.T) {
	data := []byte("abc | defg | hij | lmn   ")
	it, err := NewIterator(data)
	assert(t, err == nil, err)
	for i, c := range data {
		assert(t, c == it.CC(), it.CC())
		assert(t, i == it.GP(), it.GP())
		it.GC()
	}
	assert(t, it.EOF(), "not eof")
}
