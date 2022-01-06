package main

type Frame struct {
	f function
	i int
}

type Stack struct {
	p []Frame
}

func (stck *Stack) Push(s Frame) {
	stck.p = append(stck.p, s)
}

func (stck *Stack) Pop() {
	stck.p = stck.p[:len(stck.p)-1]
}

func (stck *Stack) Top() Frame {
	return stck.p[len(stck.p)-1]
}

func (stck *Stack) Empty() bool {
	return len(stck.p) == 0
}
