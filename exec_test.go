package main

import (
	"testing"
)

func assert(t *testing.T, b bool, s ...interface{}) {
	if s != nil {
		err, ok := s[0].(error)
		if ok {
			t.Fatalf("assert failed %v\n", err.Error())
		}
	}
	if !b {
		t.Fatalf("assert failed %v\n", s)
	}
}

func TestExecuteSimple(t *testing.T) {
	n := 0
	f := function{
		typ: 'N',
		funcs: []*function{
			{typ: 'T', terminal: func(it Iterator) code { n++; return zero }},
			{typ: 'T', terminal: func(it Iterator) code { n++; return zero }},
		},
	}
	it, _ := NewIterator([]byte("Example"))
	ret := execute(&f, it)
	assert(t, n == 2)
	assert(t, ret == zero)
}

func TestExecuteComplex1(t *testing.T) {
	n := 0
	f := function{
		typ: 'N',
		funcs: []*function{
			{typ: 'N', funcs: []*function{
				{typ: 'T', terminal: func(it Iterator) code { n++; return zero }},
			}},
			{typ: 'N', funcs: []*function{
				{typ: 'T', terminal: func(it Iterator) code { n++; return zero }},
			}},
			{typ: 'N', funcs: []*function{
				{typ: 'T', terminal: func(it Iterator) code { n++; return zero }},
			}},
			{typ: 'T', terminal: func(it Iterator) code { n++; return zero }},
		},
	}
	it, _ := NewIterator([]byte("Example"))
	ret := execute(&f, it)
	assert(t, n == 4, "n != 4")
	assert(t, ret == zero)
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
	f2.funcs = append(f2.funcs, &function{typ: 'T', terminal: func(it Iterator) code { n++; return zero }})
	f.funcs = append(f.funcs, &f2)
	it, _ := NewIterator([]byte("Example"))
	ret := execute(&f, it)
	assert(t, n == 1)
	assert(t, ret == zero)
}

func TestExecuteLogic1(t *testing.T) {
	n := 0
	f := function{typ: 'N'}
	f.funcs = append(f.funcs, &function{
		typ: 'L',
		funcs: []*function{
			{typ: 'T', terminal: func(it Iterator) code { n++; return zero }},
			{typ: 'T', terminal: func(it Iterator) code { n++; return zero }},
		},
	},
	)
	it, _ := NewIterator([]byte("Example"))
	ret := execute(&f, it)
	assert(t, n == 1)
	assert(t, ret == zero)
}

func TestExecuteLogic2(t *testing.T) {
	n := 0
	f := function{typ: 'N'}
	f.funcs = append(f.funcs, &function{
		typ: 'L',
		funcs: []*function{
			{typ: 'T', terminal: func(it Iterator) code { n++; return missed }},
			{typ: 'T', terminal: func(it Iterator) code { n++; return zero }},
		},
	},
	)
	it, _ := NewIterator([]byte("Example"))
	ret := execute(&f, it)
	assert(t, n == 2)
	assert(t, ret == zero)
}

func TestExecuteLogic3(t *testing.T) {
	n := 0
	f := function{typ: 'N'}
	f.funcs = append(f.funcs, &function{
		typ: 'L',
		funcs: []*function{
			{typ: 'T', terminal: func(it Iterator) code { n++; return missed }},
			{typ: 'T', terminal: func(it Iterator) code { n++; return missed }},
			{typ: 'T', terminal: func(it Iterator) code { n++; return zero }},
		},
	},
	)
	it, _ := NewIterator([]byte("Example"))
	ret := execute(&f, it)
	assert(t, n == 3)
	assert(t, ret == zero)
}

func TestExecuteLogic4(t *testing.T) {
	n := 0
	f := function{typ: 'N'}
	f2 := function{typ: 'L'}
	f2.funcs = append(f2.funcs, &function{
		typ: 'L',
		funcs: []*function{
			{typ: 'T', terminal: func(it Iterator) code { n++; return missed }},
			{typ: 'T', terminal: func(it Iterator) code { n++; return missed }},
		},
	},
	)
	f2.funcs = append(f2.funcs, &function{typ: 'T', terminal: func(it Iterator) code { n++; return zero }})
	f.funcs = append(f.funcs, &f2)
	it, _ := NewIterator([]byte("Example"))
	ret := execute(&f, it)
	assert(t, n == 3)
	assert(t, ret == zero)
}

func TestExecuteLogicBad1(t *testing.T) {
	n := 0
	f := function{typ: 'N', name: "S"}
	f2 := function{typ: 'L', name: "L"}
	f2.funcs = append(f2.funcs, &function{
		typ: 'L', name: "A",
		funcs: []*function{
			{typ: 'T', name: "B", terminal: func(it Iterator) code { n++; return missed }},
			{typ: 'T', name: "C", terminal: func(it Iterator) code { n++; return missed }},
		},
	},
	)
	f2.funcs = append(f2.funcs, &function{typ: 'T', name: "D", terminal: func(it Iterator) code { n++; return missed }})
	f.funcs = append(f.funcs, &f2)
	it, _ := NewIterator([]byte("Example"))
	ret := execute(&f, it)
	assert(t, n == 3)
	assert(t, ret == missed)
}

func TestExecuteLogicBad2(t *testing.T) {
	n := 0
	f := function{typ: 'N'}
	f.funcs = append(f.funcs, &function{
		typ: 'L',
		funcs: []*function{
			{typ: 'T', terminal: func(it Iterator) code { n++; return missed }},
			{typ: 'T', terminal: func(it Iterator) code { n++; return missed }},
			{typ: 'T', terminal: func(it Iterator) code { n++; return missed }},
		},
	},
	)
	it, _ := NewIterator([]byte("Example"))
	ret := execute(&f, it)
	assert(t, n == 3)
	assert(t, ret == missed)
}

func TestExecuteCycle(t *testing.T) {
	// F = T | F
	n := 0
	f := function{typ: 'C'}
	f.funcs = append(f.funcs, &function{typ: 'T', terminal: func(it Iterator) code {
		if n >= 3 {
			return missed
		}
		n++
		return zero
	}})
	f.funcs = append(f.funcs, &f)

	it, _ := NewIterator([]byte("Example"))
	execute(&f, it)
	assert(t, n == 3, n)
}
