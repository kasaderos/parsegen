package main

import "fmt"

var S = function{}

type tFunc func() bool

// lvalue = rvalue ~
// function = instructions
type function struct {
	typ      byte
	name     string
	terminal tFunc
	funcs    []*function
}

func (f *function) call() bool {
	return f.terminal()
}

func (f *function) isTerminal() bool {
	return f.typ == 'T' && f.terminal != nil
}

func (f *function) hasNext(current int) bool {
	return current+1 < len(f.funcs)
}

func (f *function) existFunc(current int) bool {
	return current < len(f.funcs)
}

func printTree(f *function) {
	fmt.Println(f.name)
	for _, subf := range f.funcs {
		fmt.Printf("\t")
		printTree(subf)
	}
}
