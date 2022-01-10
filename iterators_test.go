package main

import (
	"fmt"
	"testing"
)

// add tests for each iterator
func TestSimpleIterator(t *testing.T) {
	data := []byte("abc | defg | hij | lmn   ")
	inds := []int{0, 1, 2, 4, 6, 7, 8, 9, 11, 13, 14, 15, 17, 19, 20, 21}
	it, err := NewIterator(data, false)
	assert(t, err == nil, err)
	j := 0
	for _, c := range []byte("abc|defg|hij|lmn") {
		assert(t, c == it.CC(), it.CC(), fmt.Sprintf("%c %c", rune(c), rune(it.CC())))
		assert(t, j < len(inds) && inds[j] == it.GP())
		j++
		it.GC()
	}
	assert(t, it.EOF(), "not eof")
	assert(t, inds[len(inds)-1] == it.GP()-1, it.GP())
}

func TestCommonIterator(t *testing.T) {
	data := []byte("abc | defg | hij | lmn   ")
	it, err := NewIterator(data, true)
	assert(t, err == nil, err)
	for i, c := range data {
		assert(t, c == it.CC(), it.CC())
		assert(t, i == it.GP(), it.GP())
		it.GC()
	}
	assert(t, it.EOF(), "not eof")
}
