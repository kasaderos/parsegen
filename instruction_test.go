package main

import (
	"fmt"
	"testing"
)

func TestExecuteSimple(t *testing.T) {
	f := function{
		instructions: []instruction{
			{"T", function{terminal: func() bool { fmt.Println("A"); return false }}},
			{"T", function{terminal: func() bool { fmt.Println("B"); return false }}},
		},
	}
	execute(f)

	f = function{
		instructions: []instruction{
			{"N", function{
				instructions: []instruction{
					{"N", function{
						instructions: []instruction{
							{"T", function{terminal: func() bool { fmt.Println("C"); return false }}},
						}},
					},
				}},
			},
			{"T", function{terminal: func() bool { fmt.Println("D"); return false }}},
		},
	}
	execute(f)
}

func TestExecuteLogic(t *testing.T) {
	f := function{}
	f.instructions = append(f.instructions, instruction{
		typ: "L",
		f: function{
			instructions: []instruction{
				{"T", function{terminal: func() bool { fmt.Println("A"); return false }}},
				{"T", function{terminal: func() bool { fmt.Println("B"); return false }}},
			},
		},
	})
	fmt.Println(execute(f))
}
