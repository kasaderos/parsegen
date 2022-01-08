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
			{typ: 'N', name: "E"},
		}},
		{term{typ: 'N', name: "D"}, []term{
			{'T', "Terminal", func() bool { /*fmt.Println("D");*/ return false }},
		}},
		{term{typ: 'N', name: "E"}, []term{
			{'T', "Terminal", func() bool { /*fmt.Println("E");*/ return false }},
		}},
	}
	f, err := generateFunction(rules)
	assert(t, err == nil, err)
	assert(t, len(f.funcs) == 3, "len != 3")
	assert(t, len(f.funcs[0].funcs) == 1, "len sub != 1")
	printTree(f)
}

func TestGenerateFunction2(t *testing.T) {
	rules := []Rule{
		{term{typ: 'N', name: "S"}, []term{
			{typ: 'L', name: "A"},
			{typ: 'N', name: "B"},
		}},
		{term{typ: 'L', name: "A"}, []term{
			{'T', "Terminal", func() bool { /*fmt.Println("A");*/ return false }},
			{typ: 'N', name: "C"},
		}},
		{term{typ: 'N', name: "B"}, []term{
			{'T', "Terminal", func() bool { /*fmt.Println("B");*/ return false }},
		}},
		{term{typ: 'N', name: "C"}, []term{
			{'T', "Terminal", func() bool { /*fmt.Println("B");*/ return false }},
		}},
	}
	f, err := generateFunction(rules)
	assert(t, err == nil, err)
	printTree(f)
}

func TestWithExec1(t *testing.T) {
	n := 0
	rules := []Rule{
		{term{typ: 'N', name: "S"}, []term{
			{typ: 'N', name: "A"},
			{typ: 'N', name: "B"},
			{typ: 'N', name: "C"},
		}},
		{term{typ: 'N', name: "A"}, []term{
			{'T', "Terminal", func() bool { n++; return false }},
		}},
		{term{typ: 'N', name: "B"}, []term{
			{'T', "Terminal", func() bool { n++; return false }},
		}},
		{term{typ: 'N', name: "C"}, []term{
			{typ: 'N', name: "D"},
			{typ: 'N', name: "E"},
		}},
		{term{typ: 'N', name: "D"}, []term{
			{'T', "Terminal", func() bool { n++; return false }},
		}},
		{term{typ: 'N', name: "E"}, []term{
			{'T', "Terminal", func() bool { n++; return false }},
		}},
	}
	f, err := generateFunction(rules)
	assert(t, err == nil, err)
	assert(t, n != 4, fmt.Sprintf("n != 4; n = %d", n))
	ret := execute(f)
	assert(t, !ret, "ret true")
}

func TestWithExec2(t *testing.T) {
	n := 0
	rules := []Rule{
		{term{typ: 'N', name: "S"}, []term{
			{typ: 'L', name: "A"},
			{typ: 'N', name: "B"},
		}},
		{term{typ: 'L', name: "A"}, []term{
			{'T', "Terminal", func() bool { n++; fmt.Println("A"); return false }},
			{typ: 'N', name: "C"},
		}},
		{term{typ: 'N', name: "C"}, []term{
			{'T', "Terminal", func() bool { fmt.Println("C"); return false }},
		}},
		{term{typ: 'N', name: "B"}, []term{
			{'T', "Terminal", func() bool { n++; fmt.Println("B"); return false }},
		}},
	}
	f, err := generateFunction(rules)
	assert(t, err == nil, err)
	ret := execute(f)
	assert(t, !ret, "ret true")
	assert(t, n == 2, fmt.Sprintf("n != 2; n = %d", n))
}

func TestWithExec3(t *testing.T) {
	n := 0
	rules := []Rule{
		{term{typ: 'N', name: "S"}, []term{
			{typ: 'L', name: "A"},
			{typ: 'N', name: "B"},
		}},
		{term{typ: 'L', name: "A"}, []term{
			{'T', "Terminal", func() bool { n++; fmt.Println("A"); return true }},
			{typ: 'N', name: "C"},
		}},
		{term{typ: 'N', name: "C"}, []term{
			{'T', "Terminal", func() bool { n++; fmt.Println("C"); return false }},
		}},
		{term{typ: 'N', name: "B"}, []term{
			{'T', "Terminal", func() bool { n++; fmt.Println("B"); return false }},
		}},
	}
	f, err := generateFunction(rules)
	assert(t, err == nil, err)
	ret := execute(f)
	assert(t, !ret, "ret true")
	assert(t, n == 3, fmt.Sprintf("n != 3; n = %d", n))
}

func TestExecuteRecursive1(t *testing.T) {
	n := 0
	rules := []Rule{
		{term{typ: 'N', name: "S"}, []term{
			{typ: 'L', name: "A"},
		}},
		{term{typ: 'L', name: "A"}, []term{
			{'T', "Terminal", func() bool {
				if n >= 3 {
					return false
				}
				n++
				fmt.Println("A")
				return true
			}},
			{typ: 'L', name: "A"},
		}},
	}
	f, err := generateFunction(rules)
	assert(t, err == nil, err)
	ret := execute(f)
	assert(t, !ret, "ret true")
	assert(t, n == 3, fmt.Sprintf("n != 3; n = %d", n))
}
