package main

import "fmt"

func back(stack *Stack, it Iterator, ret *bool) {
	if it.GP() >= 23 {
		a := 0
		_ = a
	}
	for !stack.Empty() {
		f := stack.Top().f
		i := stack.Top().i
		start := stack.Top().start
		stack.Pop()
		fmt.Println("top", f.name)
		switch f.typ {
		case 'L':
			if *ret && f.hasNext(i) && f.funcs[i+1].name != f.name {
				*ret = false
				stack.Push(Frame{f, i + 1, start, 0})
				// it.BT(start)
				stack.Push(Frame{f.funcs[i+1], 0, it.GP(), 0})
				return
			}

			if !*ret && f.isRecursive() {
				stack.Push(Frame{f, i + 1, start, 0})
				stack.Push(Frame{f, 0, it.GP(), 0})
				return
			}
		case 'N':
			if start < it.GP() && f.hasNext(i) {
				stack.Push(Frame{f, i + 1, start, 0})
				stack.Push(Frame{f.funcs[i+1], 0, it.GP(), 0})
				return
			}
		}

		if f.marked && start < it.GP() {
			it.AddStart(f.name, start)
			it.AddEnd(f.name, it.GP())
		}
	}
}

func execute(f *function, it Iterator) bool {
	stack := &Stack{}
	stack.Push(Frame{f, 0, it.GP(), 0})
	// return value register
	ret := false

	for !stack.Empty() {
		f := stack.Top().f
		i := stack.Top().i
		fmt.Println("top", f.name)
		switch f.typ {
		case 'N', 'L':
			// if non terminal is not empty push
			if f.existFunc(i) {
				stack.Push(Frame{f.funcs[i], i, it.GP(), 0})
			} else {
				back(stack, it, &ret)
			}
		case 'T':
			ret = f.call(it)
			back(stack, it, &ret)
		default:
			back(stack, it, &ret)
		}
	}
	return ret
}

/*
	S = A B
	A = "T1"
	B = "T2"

	S = A | B
	A = "AA"
	B = "AB"
*/
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
