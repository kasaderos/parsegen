package parsegen

import (
	"errors"
	"fmt"
)

var S = function{}

type code int

const (
	missed code = iota
	eof
	zero
	empty
)

type tFunc func(Iterator) code

var codes = []string{"missed", "eof", "zero", "empty"}

func (c code) String() string {
	return codes[c]
}

// lvalue = rvalue ~
// function = instructions
type function struct {
	typ      byte
	name     string
	terminal tFunc
	funcs    []*function

	marked bool
}

func (f *function) call(it Iterator) code {
	i := it.GP()
	ret := f.terminal(it)
	if it.EOF() {
		return eof
	}
	if !f.isEmpty() && i == it.GP() {
		return missed
	}
	return ret
}

func (f *function) isTerminal() bool {
	return f.typ == 'T' //&& f.terminal != nil
}

func (f *function) isEmpty() bool {
	return f.name == "empty" //&& f.terminal != nil
}

func (f *function) hasNext(current int) bool {
	return current+1 < len(f.funcs)
}

func (f *function) existFunc(current int) bool {
	return current < len(f.funcs)
}

func (f *function) isCycle() bool {
	return f.typ == 'C'
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

func checkBNF(f *function) error {
	if f.isTerminal() {
		if f.hasNext(0) {
			return errors.New("[check] terminal has child")
		}
		return nil
	}

	if f.isCycle() && f.hasNext(0) {
		return errors.New("[check] cycle has more child than one")
	}

	for _, fn := range f.funcs {
		if fn == f || fn.name == f.name {
			return errors.New("[check] found recursion " + f.name)
		}
		if err := checkBNF(fn); err != nil {
			return err
		}
	}
	return nil
}
