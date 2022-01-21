package main

import "testing"

func TestParser(t *testing.T) {
	f, err := bnfFunction(nil)
	assert(t, err == nil, err)
	printTree(f)
}
