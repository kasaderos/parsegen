package main

func back(stack *Stack, it Iterator, ret *bool) {
	for !stack.Empty() {
		f := stack.Top().f
		i := stack.Top().i
		stack.Pop()
		switch f.typ {
		case 'L':
			if f.marked {
				f.appendEnds(it.GP())
			}
			if *ret && f.hasNext(i) {
				*ret = false
				stack.Push(Frame{f, i + 1})
				stack.Push(Frame{f.funcs[i+1], 0})
				return
			}
		case 'N':
			if f.marked {
				f.appendEnds(it.GP())
			}
			if !*ret && f.hasNext(i) {
				stack.Push(Frame{f, i + 1})
				stack.Push(Frame{f.funcs[i+1], 0})
				return
			}
		}
	}
}

func execute(f *function, it Iterator) bool {
	stack := &Stack{}
	stack.Push(Frame{f, 0})
	// return value register
	ret := false

	for !stack.Empty() {
		f := stack.Top().f
		i := stack.Top().i
		switch f.typ {
		case 'N', 'L':
			if f.marked {
				f.appendStarts(it.GP())
			}
			// if non terminal is not empty push
			if f.existFunc(i) {
				stack.Push(Frame{f.funcs[i], i})
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
