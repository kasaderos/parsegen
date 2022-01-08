package main

import "testing"

func TestParser(t *testing.T) {
	f, err := bnfparser(nil)
	assert(t, err == nil, err)
	printTree(f)
}
