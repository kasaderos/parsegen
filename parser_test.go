package main

import (
	"io"
	"os"
	"testing"
)

func TestParser(t *testing.T) {
	f, err := bnfFunction(nil)
	assert(t, err == nil, err)
	printTree(f)
}

func TestSDPParser(t *testing.T) {
	f, err := os.Open("sdp.bnf")
	assert(t, err == nil)
	bnf, err := io.ReadAll(f)
	assert(t, err == nil)

	parser, err := Generate(bnf)
	assert(t, err == nil, err)
	_ = parser

}
