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

func back(stack *Stack, ret bool) {
	for !stack.Empty() {
		f := stack.Top().f
		i := stack.Top().i
		switch f.typ {
		case "T":
			stack.Pop()
		case "L":
			if ret {
				if i+1 < len(f.funcs) {
					stack.Pop()
					stack.Push(Frame{f, i + 1})
					stack.Push(Frame{f.funcs[i+1], 0})
					return
				}
			}
			stack.Pop()
		case "N":
			stack.Pop()
			if !ret && i+1 < len(f.funcs) {
				stack.Push(Frame{f.funcs[i+1], i + 1})
				return
			}
		}
	}
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
		case "N", "L":
			// if non terminal is not empty push
			if i < len(f.funcs) {
				stack.Push(Frame{f.funcs[i], i})
			} else {
				back(stack, ret)
			}
		case "T":
			ret = f.call()
			back(stack, ret)
		default:
			back(stack, ret)
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
