package main

var S = function{}

type tFunc func() bool

// lvalue = rvalue ~
// function = instructions
type function struct {
	typ      string
	name     string
	terminal tFunc
	funcs    []function
}

func (f *function) call() bool {
	return f.terminal()
}

func (f *function) isTerminal() bool {
	return f.terminal != nil
}

func execute(f function) bool {
	stack := &Stack{}
	stack.Push(Frame{f, 0})
	// register of return value
	ret := false

	for !stack.Empty() {
		f := stack.Top().f
		i := stack.Top().i
		switch f.typ {
		case "N":
			// if non terminal is not empty push
			if len(f.funcs) > 0 {
				stack.Push(Frame{f.funcs[0], 0})
			}
		case "T":
			ret = f.call()
			for !stack.Empty() {
				f = stack.Top().f
				i = stack.Top().i
				if ret && (f.typ == "C" || f.typ == "L") {
					break
				}
				if f.typ == "N" && i+1 < len(f.funcs) {
					stack.Pop()
					stack.Push(Frame{f.funcs[i+1], i + 1})
					break
				}
				stack.Pop()
			}
		case "L":
			if ret && i+1 < len(f.funcs) {
				ret = false
				stack.Pop()
				stack.Push(Frame{f, i + 1})
				stack.Push(Frame{f.funcs[i+1], i + 1})
			} else {
				stack.Push(Frame{f.funcs[0], 0})
			}
			// case "C":
		}
	}
	return ret
}

/*
	F()
*/

// LOGIC():
//    G(A())
//    G(B())

// A():
//   print('A')

// B():
//   print('B')

// C():
//   print('C')
