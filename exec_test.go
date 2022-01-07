package main

import (
	"testing"
)

func assert(t *testing.T, b bool, s ...interface{}) {
	err, ok := s[0].(error)
	if ok {
		t.Errorf("assert failed %v\n", err.Error())
	}
	if !b {
		t.Errorf("assert failed %v\n", s)
	}
}

func TestExecuteSimple(t *testing.T) {
	n := 0
	f := function{
		typ: 'N',
		funcs: []*function{
			{typ: 'T', terminal: func() bool { n++; return false }},
			{typ: 'T', terminal: func() bool { n++; return false }},
		},
	}
	ret := execute(&f)
	assert(t, n == 2)
	assert(t, !ret)
}

func TestExecuteComplex1(t *testing.T) {
	n := 0
	f := function{
		typ: 'N',
		funcs: []*function{
			{typ: 'N', funcs: []*function{
				{typ: 'N', funcs: []*function{
					{typ: 'T', terminal: func() bool { n++; return false }},
				}},
			},
			},
			{typ: 'T', terminal: func() bool { n++; return false }},
		},
	}
	ret := execute(&f)
	assert(t, n == 2)
	assert(t, !ret)
}

func TestExecuteComplex2(t *testing.T) {
	n := 0
	f := function{typ: 'N'}
	f2 := function{typ: 'N'}
	f2.funcs = append(f2.funcs, &function{
		typ: 'N',
		funcs: []*function{
			{typ: 'N'},
			{typ: 'N'},
		},
	},
	)
	f2.funcs = append(f2.funcs, &function{typ: 'T', terminal: func() bool { n++; return false }})
	f.funcs = append(f.funcs, &f2)
	ret := execute(&f)
	assert(t, n == 1)
	assert(t, !ret)
}

func TestExecuteLogic1(t *testing.T) {
	n := 0
	f := function{typ: 'N'}
	f.funcs = append(f.funcs, &function{
		typ: 'L',
		funcs: []*function{
			{typ: 'T', terminal: func() bool { n++; return false }},
			{typ: 'T', terminal: func() bool { n++; return false }},
		},
	},
	)
	ret := execute(&f)
	assert(t, n == 1)
	assert(t, !ret)
}

func TestExecuteLogic2(t *testing.T) {
	n := 0
	f := function{typ: 'N'}
	f.funcs = append(f.funcs, &function{
		typ: 'L',
		funcs: []*function{
			{typ: 'T', terminal: func() bool { n++; return true }},
			{typ: 'T', terminal: func() bool { n++; return false }},
		},
	},
	)
	ret := execute(&f)
	assert(t, n == 2)
	assert(t, !ret)
}

func TestExecuteLogic3(t *testing.T) {
	n := 0
	f := function{typ: 'N'}
	f.funcs = append(f.funcs, &function{
		typ: 'L',
		funcs: []*function{
			{typ: 'T', terminal: func() bool { n++; return true }},
			{typ: 'T', terminal: func() bool { n++; return true }},
			{typ: 'T', terminal: func() bool { n++; return false }},
		},
	},
	)
	ret := execute(&f)
	assert(t, n == 3)
	assert(t, !ret)
}

func TestExecuteLogic4(t *testing.T) {
	n := 0
	f := function{typ: 'N'}
	f2 := function{typ: 'L'}
	f2.funcs = append(f2.funcs, &function{
		typ: 'L',
		funcs: []*function{
			{typ: 'T', terminal: func() bool { n++; return true }},
			{typ: 'T', terminal: func() bool { n++; return true }},
		},
	},
	)
	f2.funcs = append(f2.funcs, &function{typ: 'T', terminal: func() bool { n++; return false }})
	f.funcs = append(f.funcs, &f2)
	ret := execute(&f)
	assert(t, n == 3)
	assert(t, !ret)
}

func TestExecuteLogicBad1(t *testing.T) {
	n := 0
	f := function{typ: 'N'}
	f2 := function{typ: 'L'}
	f2.funcs = append(f2.funcs, &function{
		typ: 'L',
		funcs: []*function{
			{typ: 'T', terminal: func() bool { n++; return true }},
			{typ: 'T', terminal: func() bool { n++; return true }},
		},
	},
	)
	f2.funcs = append(f2.funcs, &function{typ: 'T', terminal: func() bool { n++; return true }})
	f.funcs = append(f.funcs, &f2)
	ret := execute(&f)
	assert(t, n == 3)
	assert(t, ret)
}

func TestExecuteLogicBad2(t *testing.T) {
	n := 0
	f := function{typ: 'N'}
	f.funcs = append(f.funcs, &function{
		typ: 'L',
		funcs: []*function{
			{typ: 'T', terminal: func() bool { n++; return true }},
			{typ: 'T', terminal: func() bool { n++; return true }},
			{typ: 'T', terminal: func() bool { n++; return true }},
		},
	},
	)
	ret := execute(&f)
	assert(t, n == 3)
	assert(t, ret)
}

func TestExecuteRecursive(t *testing.T) {
	// F = T | F
	n := 0
	f := function{typ: 'L'}
	f.funcs = append(f.funcs, &function{typ: 'T', terminal: func() bool { n++; return n < 3 }})
	f.funcs = append(f.funcs, &f)
	execute(&f)
	assert(t, n == 3)
}
