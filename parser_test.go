package main

import "testing"

func TestParser(t *testing.T) {
	f := bnfparser(nil)
	printTree(f)
}
