package main

func back(stack *Stack, ret *bool) {
	for !stack.Empty() {
		f := stack.Top().f
		i := stack.Top().i
		stack.Pop()
		switch f.typ {
		case 'L':
			if *ret {
				if f.hasNext(i) {
					*ret = false
					stack.Push(Frame{f, i + 1})
					stack.Push(Frame{f.funcs[i+1], 0})
					return
				}
			}
		case 'N':
			if !*ret && f.hasNext(i) {
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
		case 'N', 'L':
			// if non terminal is not empty push
			if f.existFunc(i) {
				stack.Push(Frame{f.funcs[i], i})
			} else {
				back(stack, &ret)
			}
		case 'T':
			ret = f.call()
			back(stack, &ret)
		default:
			back(stack, &ret)
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
