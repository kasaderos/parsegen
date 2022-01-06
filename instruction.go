package main

var S = function{}

type tFunc func() bool

// instruction:
//    simple (terminal)
//    func   (non-terminal)
//    for
//    if else
//    call another func
type instruction struct {
	typ string
	// if f is simple call f()
	// if f is non-terminal move to f
	f function
}

// lvalue = rvalue ~
// function = instructions
type function struct {
	instructions []instruction
	name         string
	terminal     tFunc
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
	ret := false

	for !stack.Empty() {
		top := stack.Top()
		next := top.f.instructions[top.i]
		switch next.typ {
		case "N":
			// fmt.Println("pushed N")
			if len(next.f.instructions) > 0 {
				stack.Push(Frame{next.f, 0})
			}
		case "T":
			// fmt.Println("called T")
			ret = next.f.call()
			for !stack.Empty() {
				top = stack.Top()
				typ := top.f.instructions[top.i].typ
				if typ == "L" || typ == "C" {
					break
				}
				if !ret && top.i+1 < len(top.f.instructions) {
					stack.Pop()
					stack.Push(Frame{f, top.i + 1})
					break
				}
				stack.Pop()
			}
		case "L":
			if ret && top.i+1 < len(next.f.instructions) {
				ret = false
				stack.Pop()
				stack.Push(Frame{next.f, top.i + 1})
			}
		case "C":
		}
	}
	return ret
}

/*
	F()
*/

// func():
//    A() | B()

// A():
//   print('A')

// B():
//   print('B')

// C():
//   print('C')
