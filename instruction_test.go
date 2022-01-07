package main

import (
	"fmt"
	"testing"
)

func assert(t *testing.T, b bool) {
	if !b {
		t.Errorf("assert failed\n")
	}
}

func TestExecuteSimple(t *testing.T) {
	n := 0
	f := function{
		typ: "N",
		funcs: []function{
			{typ: "T", terminal: func() bool { n++; return false }},
			{typ: "T", terminal: func() bool { n++; return false }},
		},
	}
	execute(f)
	assert(t, n == 2)
}

func TestExecuteComplex1(t *testing.T) {
	n := 0
	f := function{
		typ: "N",
		funcs: []function{
			{typ: "N", funcs: []function{
				{typ: "N", funcs: []function{
					{typ: "T", terminal: func() bool { n++; return false }},
				}},
			},
			},
			{typ: "T", terminal: func() bool { n++; return false }},
		},
	}
	execute(f)
	assert(t, n == 2)
}

func TestExecuteComplex2(t *testing.T) {
	n := 0
	f := function{typ: "N"}
	f2 := function{typ: "N"}
	f2.funcs = append(f2.funcs, function{
		typ: "N",
		funcs: []function{
			{typ: "N"},
			{typ: "N"},
		},
	},
	)
	f2.funcs = append(f2.funcs, function{typ: "T", terminal: func() bool { n++; return false }})
	f.funcs = append(f.funcs, f2)
	fmt.Println(execute(f))
	assert(t, n == 1)
}

func TestExecuteLogic1(t *testing.T) {
	n := 0
	f := function{typ: "N"}
	f.funcs = append(f.funcs, function{
		typ: "L",
		funcs: []function{
			{typ: "T", terminal: func() bool { n++; return false }},
			{typ: "T", terminal: func() bool { n++; return false }},
		},
	},
	)
	fmt.Println(execute(f))
	assert(t, n == 1)
}

func TestExecuteLogic2(t *testing.T) {
	n := 0
	f := function{typ: "N"}
	f.funcs = append(f.funcs, function{
		typ: "L",
		funcs: []function{
			{typ: "T", terminal: func() bool { n++; return true }},
			{typ: "T", terminal: func() bool { n++; return false }},
		},
	},
	)
	fmt.Println(execute(f))
	assert(t, n == 2)
}

func TestExecuteLogic3(t *testing.T) {
	n := 0
	f := function{typ: "N"}
	f.funcs = append(f.funcs, function{
		typ: "L",
		funcs: []function{
			{typ: "T", terminal: func() bool { n++; return true }},
			{typ: "T", terminal: func() bool { n++; return true }},
			{typ: "T", terminal: func() bool { n++; return false }},
		},
	},
	)
	fmt.Println(execute(f))
	assert(t, n == 3)
}

func TestExecuteLogic4(t *testing.T) {
	n := 0
	f := function{typ: "N"}
	f2 := function{typ: "L"}
	f2.funcs = append(f2.funcs, function{
		typ: "L",
		funcs: []function{
			{typ: "T", terminal: func() bool { n++; return true }},
			{typ: "T", terminal: func() bool { n++; return true }},
		},
	},
	)
	f2.funcs = append(f2.funcs, function{typ: "T", terminal: func() bool { n++; return false }})
	f.funcs = append(f.funcs, f2)
	fmt.Println(execute(f))
	assert(t, n == 3)
}

func TestExecuteRecursive(t *testing.T) {
	// F = T | F
	n := 0
	f := function{typ: "L"}
	f.funcs = append(f.funcs, function{typ: "T", terminal: func() bool { n++; return n < 3 }})
	f.funcs = append(f.funcs, f)
	execute(f)
	assert(t, n == 3)
}
