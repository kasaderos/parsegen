package main

type fsmFlag int32
type fsmState func()

const (
	EOF fsmFlag = iota
	MustExist
	ExitedWithError
	FSMFlags
)

type FSM struct {
	stack   []fsmState
	current fsmState
	flags   [FSMFlags]bool
}

func (q *FSM) Push(f fsmState) {
	q.stack = append(q.stack, f)
}

func (q *FSM) Empty() bool {
	return len(q.stack) == 0
}

func (q *FSM) Top() fsmState {
	return q.stack[len(q.stack)-1]
}

func (q *FSM) Pop() {
	if q.Empty() {
		q.current = nil
	}
	q.stack = q.stack[:len(q.stack)-1]
}

func (fsm *FSM) ChState(f fsmState) {
	fsm.current = f
}

func (fsm *FSM) NoErrors() bool {
	return fsm.flags[ExitedWithError]
}

func (fsm *FSM) Init(states []fsmState) {
	// fill in reverse order
	for i := len(states) - 1; i >= 0; i-- {
		fsm.stack = append(fsm.stack, states[i])
	}
}

func (fsm *FSM) Start() {
	for !fsm.Empty() && !fsm.NoErrors() {
		fsm.ChState(fsm.Top())
		fsm.Pop()

		// exit
		if fsm.current == nil {
			break
		}
		fsm.current()
	}
}
