package parsegen

import "fmt"

var Debug = false

// back roll back a stack with ret code that terminal returned
func back(stack *Stack, it Iterator, ret *code) {
	for !stack.Empty() {
		f := stack.Top().f
		i := stack.Top().i
		start := stack.Top().start
		buf := stack.Top().buf
		stack.Pop()

		if Debug {
			fmt.Printf("back: top %s '%c'\n", f.name, it.CC())
		}

		switch f.typ {
		case 'L':
			// if it missed call next function
			if *ret == missed && f.hasNext(i) {
				*ret = zero
				stack.Push(Frame{f, i + 1, start, buf})
				it.BT(start)
				return
			}
		case 'C':
			// push while ret is zero (normal)
			// store current index in the stack buf
			if *ret == zero {
				stack.Push(Frame{f, 0, start, it.GP()})
				return
			}
			// backtrack to last zero ended entity
			if *ret == missed {
				it.BT(buf)
			}
			*ret = zero
		case 'N':
			// if it's ok,  push next function
			if (*ret == exit || *ret == zero) && f.hasNext(i) {
				stack.Push(Frame{f, i + 1, start, it.GP()})
				return
			}
		}

		// add bounds to label
		if f.marked && *ret != missed {
			it.AppendStart(f.name, start)
			it.AppendEnd(f.name, it.GP())
		}
	}
}

// execute executes function with "recursive descent method" on a stack
// Example:
//   foo = bar1 bar2 ;
//   foo[N]:    // foo is a function with typ 'N'
//      bar1[N]  // another function that foo contains
//		bar2[N]
//   executes bar1, then bar2.
//   Logic of rollback implemented in back function (see back)
func execute(f *function, it Iterator) code {
	stack := &Stack{}
	stack.Push(Frame{f, 0, it.GP(), 0})
	// return value register
	ret := zero

	for !stack.Empty() {
		f := stack.Top().f
		i := stack.Top().i

		if Debug {
			fmt.Printf("front: top, %s '%c'\n", f.name, it.CC())
		}
		switch f.typ {
		case 'N', 'L', 'C':
			// if non terminal is not empty push
			if f.existFunc(i) {
				stack.Push(Frame{f.funcs[i], 0, it.GP(), it.GP()})
			} else {
				back(stack, it, &ret)
			}
		case 'T':
			// call terminal function
			ret = f.call(it)

			if Debug {
				fmt.Printf("front: ret, %s '%c'\n", ret.String(), it.CC())
			}
			back(stack, it, &ret)
		default:
			// unknown type, just back
			back(stack, it, &ret)
		}
	}
	return ret
}
