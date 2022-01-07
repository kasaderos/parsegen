package main

import (
	"fmt"
	"testing"
)

func TestGenerateFunction1(t *testing.T) {
	rules := []Rule{
		{term{typ: 'N', name: "S"}, []term{
			{typ: 'N', name: "A"},
			{typ: 'N', name: "B"},
			{typ: 'N', name: "C"},
		}},
		{term{typ: 'N', name: "A"}, []term{
			{'T', "Terminal", func() bool { /*fmt.Println("A");*/ return false }},
		}},
		{term{typ: 'N', name: "B"}, []term{
			{'T', "Terminal", func() bool { /*fmt.Println("B");*/ return false }},
		}},
		{term{typ: 'N', name: "C"}, []term{
			{typ: 'N', name: "D"},
		}},
		{term{typ: 'N', name: "D"}, []term{
			{'T', "Terminal", func() bool { /*fmt.Println("D");*/ return false }},
		}},
	}
	f, err := generateFunction(rules)
	assert(t, err == nil, "error")
	assert(t, len(f.funcs) == 3, "len != 3")
	assert(t, len(f.funcs[0].funcs) == 1, "len sub != 1")
	printTree(f)
	fmt.Println(f)
}
