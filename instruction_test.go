package main

import (
	"fmt"
	"testing"
)

func TestExecuteSimple(t *testing.T) {
	f := function{
		typ: "N",
		funcs: []function{
			{typ: "T", terminal: func() bool { fmt.Println("A"); return false }},
			{typ: "T", terminal: func() bool { fmt.Println("B"); return false }},
		},
	}
	execute(f)
}

func TestExecuteComplex(t *testing.T) {
	f := function{
		typ: "N",
		funcs: []function{
			{typ: "N", funcs: []function{
				{typ: "N", funcs: []function{
					{typ: "T", terminal: func() bool { fmt.Println("C"); return false }},
				}},
			},
			},
			{typ: "T", terminal: func() bool { fmt.Println("D"); return false }},
		},
	}
	execute(f)
}

func TestExecuteLogic1(t *testing.T) {
	f := function{typ: "N"}
	f.funcs = append(f.funcs, function{
		typ: "L",
		funcs: []function{
			{typ: "T", terminal: func() bool { fmt.Println("A"); return false }},
			{typ: "T", terminal: func() bool { fmt.Println("B"); return false }},
		},
	},
	)
	fmt.Println(execute(f))
}

func TestExecuteLogic2(t *testing.T) {
	f := function{typ: "N"}
	f.funcs = append(f.funcs, function{
		typ: "L",
		funcs: []function{
			{typ: "T", terminal: func() bool { fmt.Println("A"); return true }},
			{typ: "T", terminal: func() bool { fmt.Println("B"); return false }},
		},
	},
	)
	fmt.Println(execute(f))
}
