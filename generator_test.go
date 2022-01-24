package parsegen

import (
	"testing"
)

func TestGenerateFunction1(t *testing.T) {
	rules := []*rule{
		{term{typ: 'N', name: "S"}, []term{
			{typ: 'N', name: "A"},
			{typ: 'N', name: "B"},
			{typ: 'N', name: "C"},
		}},
		{term{typ: 'N', name: "A"}, []term{
			{typ: 'T', name: "Terminal", terminal: func(it Iterator) code { /*fmt.Println("A");*/ return zero }},
		}},
		{term{typ: 'N', name: "B"}, []term{
			{typ: 'T', name: "Terminal", terminal: func(it Iterator) code { /*fmt.Println("B");*/ return zero }},
		}},
		{term{typ: 'N', name: "C"}, []term{
			{typ: 'N', name: "D"},
			{typ: 'N', name: "E"},
		}},
		{term{typ: 'N', name: "D"}, []term{
			{typ: 'T', name: "Terminal", terminal: func(it Iterator) code { /*fmt.Println("D");*/ return zero }},
		}},
		{term{typ: 'N', name: "E"}, []term{
			{typ: 'T', name: "Terminal", terminal: func(it Iterator) code { /*fmt.Println("E");*/ return zero }},
		}},
	}
	f, err := generateFunction(rules)
	assert(t, err == nil, err)
	assert(t, len(f.funcs) == 3, "len != 3")
	assert(t, len(f.funcs[0].funcs) == 1, "len sub != 1")
	printTree(f)
}

func TestGenerateFunction2(t *testing.T) {
	rules := []*rule{
		{term{typ: 'N', name: "S"}, []term{
			{typ: 'L', name: "A"},
			{typ: 'N', name: "B"},
		}},
		{term{typ: 'L', name: "A"}, []term{
			{typ: 'T', name: "Terminal", terminal: func(it Iterator) code { /*fmt.Println("A");*/ return zero }},
			{typ: 'N', name: "C"},
		}},
		{term{typ: 'N', name: "B"}, []term{
			{typ: 'T', name: "Terminal", terminal: func(it Iterator) code { /*fmt.Println("B");*/ return zero }},
		}},
		{term{typ: 'N', name: "C"}, []term{
			{typ: 'T', name: "Terminal", terminal: func(it Iterator) code { /*fmt.Println("B");*/ return zero }},
		}},
	}
	f, err := generateFunction(rules)
	assert(t, err == nil, err)
	printTree(f)
}

func TestBacktrackLogic(t *testing.T) {
	rules := []*rule{
		{term{typ: 'N', name: "S"}, []term{
			{typ: 'L', name: "A", marked: true},
		}},
		{term{typ: 'L', name: "A", marked: true}, []term{
			{typ: 'T', name: "T1", terminal: termStr("AA")},
			{typ: 'T', name: "T2", terminal: termStr("AB")},
		}},
	}
	f, err := generateFunction(rules)
	assert(t, err == nil, err)
	it, err := NewIterator([]byte("AB"))
	assert(t, err == nil, err)
	ret := execute(f, it)
	assert(t, ret == zero || ret == eof, "ret == err")
	printTree(f)
}

func TestBacktrackCycle(t *testing.T) {
	rules := []*rule{
		{term{typ: 'N', name: "S"}, []term{
			{typ: 'C', name: "A", marked: true},
		}},
		{term{typ: 'C', name: "A", marked: true}, []term{
			{typ: 'T', name: "T1", terminal: termStr("AB")},
		}},
	}
	f, err := generateFunction(rules)
	assert(t, err == nil, err)
	it, err := NewIterator([]byte("ABABBC"))
	assert(t, err == nil, err)
	ret := execute(f, it)
	assert(t, ret == zero, "ret == err")
}

func TestExecuteCycleData(t *testing.T) {
	rules := []*rule{
		{term{typ: 'N', name: "S"}, []term{
			term{typ: 'C', name: "SP", marked: true},
			{typ: 'C', name: "A", marked: true},
			term{typ: 'C', name: "SP", marked: true},
			{typ: 'C', name: "A", marked: true},
			term{typ: 'C', name: "SP", marked: true},
			{typ: 'C', name: "A", marked: true},
			term{typ: 'C', name: "SP", marked: true},
		}},
		{term{typ: 'C', name: "SP", marked: true}, []term{
			{typ: 'T', name: "sp", terminal: termSpace()},
		}},
		{term{typ: 'C', name: "A", marked: true}, []term{
			{typ: 'T', name: "T1", terminal: termStr("AB")},
		}},
	}
	f, err := generateFunction(rules)
	assert(t, err == nil, err)
	it, err := NewIterator([]byte("  AB  AB  AB  "))
	assert(t, err == nil, err)
	ret := execute(f, it)
	// printTree(f)
	assert(t, ret == zero, "ret == err")
}

func TestBaseRule(t *testing.T) {
	parser, err := Generate([]byte(
		"S = \"Hello World\" \" \" \"!!!\";",
	))
	assert(t, err == nil, err)
	_, err = parser.Parse([]byte("Hello World !!!"))
	assert(t, err == nil, err)
}

func TestCaseRule(t *testing.T) {
	parser, err := Generate([]byte(
		"S = \"!!\" | \"Hello World\" | \"!\" ;",
	))
	assert(t, err == nil, err)
	_, err = parser.Parse([]byte("Hello World"))
	assert(t, err == nil, err)
}

func TestCycleRule(t *testing.T) {
	parser, err := Generate([]byte(
		"S = { \"Hello World;\" } ;",
	))
	assert(t, err == nil, err)
	_, err = parser.Parse([]byte("Hello World;Hello World;"))
	assert(t, err == nil, err)
}

func TestTwoRules(t *testing.T) {
	parser, err := Generate([]byte(
		"S = A \"!!!\" ;" +
			"A = { \"Hello World\" } ;",
	))
	assert(t, err == nil, err)
	_, err = parser.Parse([]byte("Hello World!!!"))
	assert(t, err == nil, err)
	_, err = parser.Parse([]byte("!!!"))
	assert(t, err == nil, err)
}

func TestHttpGetRequest(t *testing.T) {
	parser, err := Generate([]byte(
		"S = Method SP Url SP StatusOk;" +
			"Method = any(0x20);" +
			"SP = 0x20 ;" +
			"Url = any(0x20);" +
			"StatusOk = integer;" +
			"integer = digit digits;" +
			"digit = 0x30-39;" +
			"digits = { digit };",
	))
	assert(t, err == nil, err)
	_, err = parser.Parse([]byte("GET https://google.com 200"))
	assert(t, err == nil, err)
}

func TestBasicAny(t *testing.T) {
	end := byte(0)
	included := false
	assert(t, isTermAny([]byte("any(0x0d)"), &end, &included))
	assert(t, end == '\r' && !included)

	assert(t, isTermAny([]byte("any[0x0d]"), &end, &included))
	assert(t, end == '\r' && included)
}
