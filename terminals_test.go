package parsegen

import (
	"testing"
)

func TestStr(t *testing.T) {
	s := "abcd bcd"
	name := "String"
	rules := []*Rule{
		{term{typ: 'N', name: "S"}, []term{
			{typ: 'N', name: name, marked: true},
		}},
		{term{typ: 'N', name: name, marked: true}, []term{
			{typ: 'T', name: "str", terminal: termStr(s)},
		}},
	}
	f, err := generateFunction(rules)
	assert(t, err == nil, err)
	it, err := NewIterator([]byte(s))
	assert(t, err == nil, err)
	ret := execute(f, it)
	assert(t, ret == eof)
}

func TestAnyQuoted(t *testing.T) {
	// TODO add bad tests and check error
	s := "\"abcd bcd\""
	name := "QuotedString"
	rules := []*Rule{
		{term{typ: 'N', name: "S"}, []term{
			{typ: 'N', name: name, marked: true},
		}},
		{term{typ: 'N', name: name, marked: true}, []term{
			{typ: 'T', name: "quoted", terminal: termAnyQuoted()},
		}},
	}
	f, err := generateFunction(rules)
	assert(t, err == nil, err)
	it, err := NewIterator([]byte(s))
	assert(t, err == nil, err)
	ret := execute(f, it)
	assert(t, ret == eof)
}

func TestID(t *testing.T) {
	// TODO add bad tests and check error
	s := "abcd"
	name := "ID"
	rules := []*Rule{
		{term{typ: 'N', name: "S"}, []term{
			{typ: 'N', name: name, marked: true},
		}},
		{term{typ: 'N', name: name, marked: true}, []term{
			{typ: 'T', name: "id", terminal: termID()},
		}},
	}
	f, err := generateFunction(rules)
	assert(t, err == nil, err)
	it, err := NewIterator([]byte(s))
	assert(t, err == nil, err)
	ret := execute(f, it)
	assert(t, ret == eof)
}

func TestCombined1(t *testing.T) {
	// TODO add bad tests and check error
	s := "Hello \"World\""
	id := "ID"
	str := "String"
	quoted := "Quoted"

	rules := []*Rule{
		{term{typ: 'N', name: "S"}, []term{
			{typ: 'N', name: id, marked: true},
			{typ: 'N', name: str, marked: true},
			{typ: 'N', name: quoted, marked: true},
		}},
		{term{typ: 'N', name: id, marked: true}, []term{
			{typ: 'T', terminal: termID()},
		}},
		{term{typ: 'N', name: str, marked: true}, []term{
			{typ: 'T', terminal: termStr(" ")},
		}},
		{term{typ: 'N', name: quoted, marked: true}, []term{
			{typ: 'T', terminal: termAnyQuoted()},
		}},
	}
	f, err := generateFunction(rules)
	assert(t, err == nil, err)
	it, err := NewIterator([]byte(s))
	assert(t, err == nil, err)
	ret := execute(f, it)
	assert(t, ret == eof)
}

func TestCombined2(t *testing.T) {
	// TODO add bad tests and check error
	s := "GET dafadf 200 OK"
	method := "Method"
	url := "Url"
	code := "Code"
	space := "Space"
	message := "Message"

	// S = Method " " Url " " Code " " Message
	// Method = AnySpace
	// Url = AnySpace
	// Code = Number
	// Message = AnySpace

	// S:
	//   Method()
	//   CheckStr(" ")
	//   Url()
	//   CheckStr(" ")
	//   Code()
	//   CheckStr(" ")
	//   Message()

	// Method:
	//    AnySpace()

	rules := []*Rule{
		{term{typ: 'N', name: "S"}, []term{
			{typ: 'N', name: method, marked: true},
			{typ: 'T', name: space, terminal: termStr(" ")},
			{typ: 'N', name: url, marked: true},
			{typ: 'T', name: space, terminal: termStr(" ")},
			{typ: 'N', name: code, marked: true},
			{typ: 'T', name: space, terminal: termStr(" ")},
			{typ: 'N', name: message, marked: true},
		}},
		{term{typ: 'N', name: method, marked: true}, []term{
			{typ: 'T', name: "tMethod", terminal: termAny(' ', false)},
		}},
		{term{typ: 'N', name: url, marked: true}, []term{
			{typ: 'T', name: "tUrl", terminal: termAny(' ', false)},
		}},
		{term{typ: 'N', name: code, marked: true}, []term{
			{typ: 'T', name: "tCode", terminal: termInteger()},
		}},
		{term{typ: 'N', name: message, marked: true}, []term{
			{typ: 'T', name: "tMessage", terminal: termAny(0, false)},
		}},
	}
	f, err := generateFunction(rules)
	assert(t, err == nil, err)
	it, err := NewIterator([]byte(s))
	assert(t, err == nil, err)
	ret := execute(f, it)
	assert(t, ret == eof)
}

func TestTermHex(t *testing.T) {
	// integer
	parser, err := Generate([]byte(
		"S = A ;" +
			"A = 0x30-39 A1 ;" +
			"A1 = { 0x30-39 } ;",
	))
	assert(t, err == nil, err)
	_, err = parser.Parse([]byte("121313"))
	assert(t, err == nil, err)
}

func TestTermEmpty(t *testing.T) {
	// ["+"] 7
	parser, err := Generate([]byte(
		"S = A 0x37 ;" +
			"A = 0x2b | 0x2d | empty ;",
	))
	assert(t, err == nil, err)
	_, err = parser.Parse([]byte("+7"))
	assert(t, err == nil, err)
	_, err = parser.Parse([]byte("-7"))
	assert(t, err == nil, err)
	_, err = parser.Parse([]byte("7"))
	assert(t, err == nil, err)
}

func TestCycleTermEmpty(t *testing.T) {
	parser, err := Generate([]byte(
		"S = \" \" A \" \" ;" +
			"A = { empty } ;",
	))
	assert(t, err == nil, err)
	_, err = parser.Parse([]byte("  "))
	assert(t, err == nil, err)
}
