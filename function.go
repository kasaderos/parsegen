package main

import "fmt"

var S = function{}

type tFunc func(Iterator) bool

// lvalue = rvalue ~
// function = instructions
type function struct {
	typ      byte
	name     string
	terminal tFunc
	funcs    []*function
}

func (f *function) call(it Iterator) bool {
	return f.terminal(it)
}

func (f *function) isTerminal() bool {
	return f.typ == 'T' //&& f.terminal != nil
}

func (f *function) hasNext(current int) bool {
	return current+1 < len(f.funcs)
}

func (f *function) existFunc(current int) bool {
	return current < len(f.funcs)
}

func recPrint(f *function, i int) {
	if f.isTerminal() {
		fmt.Printf("%s [T]\n", f.name)
	} else {
		fmt.Printf("%s-- [%c]\n", f.name, f.typ)
	}
	for _, subf := range f.funcs {
		fmt.Printf("\t")
		for j := 0; j < i-1; j++ {
			fmt.Printf("\t")
		}
		fmt.Printf("-")
		if subf.name != f.name {
			recPrint(subf, i+1)
		} else {
			fmt.Printf("%s\n", f.name)
		}
	}
}
func printTree(f *function) {
	recPrint(f, 1)
}
