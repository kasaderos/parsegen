package main

import "fmt"

func back(stack *Stack, it Iterator, ret *code) {
	for !stack.Empty() {
		f := stack.Top().f
		i := stack.Top().i
		start := stack.Top().start
		buf := stack.Top().buf
		stack.Pop()
		fmt.Println("back: top", f.name)

		switch f.typ {
		case 'L':
			if *ret == missed && f.hasNext(i) {
				*ret = zero
				stack.Push(Frame{f, i + 1, start, it.GP()})
				it.BT(start)
				return
			}
		case 'C':
			if *ret == zero {
				stack.Push(Frame{f, 0, start, it.GP()})
				return
			}
			it.BT(buf)
			*ret = zero
		case 'N':
			if *ret == zero && f.hasNext(i) {
				stack.Push(Frame{f, i + 1, start, it.GP()})
				return
			}
		}

		if f.marked && *ret != missed {
			it.SetStart(f.name, start)
			it.SetEnd(f.name, it.GP())
		}
	}
}

func execute(f *function, it Iterator) code {
	stack := &Stack{}
	stack.Push(Frame{f, 0, it.GP(), 0})
	// return value register
	ret := zero

	for !stack.Empty() {
		f := stack.Top().f
		i := stack.Top().i
		fmt.Println("front: top", f.name)
		switch f.typ {
		case 'N', 'L', 'C':
			// if non terminal is not empty push
			if f.existFunc(i) {
				stack.Push(Frame{f.funcs[i], 0, it.GP(), it.GP()})
			} else {
				back(stack, it, &ret)
			}
		case 'T':
			ret = f.call(it)
			fmt.Println("front: ret", ret.String())
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
