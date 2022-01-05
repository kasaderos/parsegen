package main

type state struct {
	name       string
	start, end int
	gp         func() int
	f          func() bool
}

func (s *state) F() bool {
	s.start = s.gp()
	res := s.f()
	s.end = s.gp()
	return res
}

func (s state) IsTerminal() bool {
	return s.f != nil
}

func (s state) IsFinal() bool {
	return s.name == FinalState.name
}

type FSM struct {
	graph     *Graph
	current   *Node
	backtrack Stack
}

func (fsm *FSM) Start() {
	for !fsm.current.IsFinal() {
		if len(fsm.current.nexts) == 0 {
			// error must ended with "E"
			return
		}

		curr := fsm.current.nexts[0]
		if len(fsm.current.nexts) > 1 {
			fsm.backtrack.Push(Frame{fsm.current, 1})
		}

		if curr.IsTerminal() && curr.f() {
			if !fsm.hasBacktrack() {
				// error
				return
			}

			// backtrack and push next
			top := fsm.backtrack.Top()
			fsm.current = top.node.nexts[top.next]
			fsm.backtrack.Pop()
			// guarantee top.next < len(top.node.nexts)
			if top.next+1 < len(top.node.nexts) {
				fsm.backtrack.Push(Frame{fsm.current, top.next + 1})
			}
		} else {
			fsm.current = curr
		}
	}
}

func (fsm *FSM) check() bool {
	return false
}

func (fsm *FSM) hasBacktrack() bool {
	return !fsm.backtrack.Empty()
}
