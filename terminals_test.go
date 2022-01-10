package main

import (
	"fmt"
	"testing"
)

func TestStr(t *testing.T) {
	// TODO add bad tests and check error
	s := "abcd bcd"
	name := "String"
	rules := []Rule{
		{term{typ: 'N', name: "S"}, []term{
			{typ: 'N', name: name, marked: true},
		}},
		{term{typ: 'N', name: name, marked: true}, []term{
			{typ: 'T', name: "str", terminal: termStr(s)},
		}},
	}
	f, err := generateFunction(rules)
	assert(t, err == nil, err)
	it, err := NewIterator([]byte(s), true)
	assert(t, err == nil, err)
	execute(f, it)
	i := it.Data().labels[name].i[0]
	j := it.Data().labels[name].j[0]
	assert(t, i == 0, i)
	assert(t, j == 7, j)
	assert(t, !it.HasError())
}

func TestAnyQuoted(t *testing.T) {
	// TODO add bad tests and check error
	s := "\"abcd bcd\""
	name := "QuotedString"
	rules := []Rule{
		{term{typ: 'N', name: "S"}, []term{
			{typ: 'N', name: name, marked: true},
		}},
		{term{typ: 'N', name: name, marked: true}, []term{
			{typ: 'T', name: "quoted", terminal: termAnyQuoted()},
		}},
	}
	f, err := generateFunction(rules)
	assert(t, err == nil, err)
	it, err := NewIterator([]byte(s), true)
	assert(t, err == nil, err)
	execute(f, it)
	i := it.Data().labels[name].i[0]
	j := it.Data().labels[name].j[0]
	assert(t, i == 0, i)
	assert(t, j == 9, j)
	assert(t, !it.HasError())
}

func TestID(t *testing.T) {
	// TODO add bad tests and check error
	s := "abcd"
	name := "ID"
	rules := []Rule{
		{term{typ: 'N', name: "S"}, []term{
			{typ: 'N', name: name, marked: true},
		}},
		{term{typ: 'N', name: name, marked: true}, []term{
			{typ: 'T', name: "id", terminal: termID()},
		}},
	}
	f, err := generateFunction(rules)
	assert(t, err == nil, err)
	it, err := NewIterator([]byte(s), true)
	assert(t, err == nil, err)
	execute(f, it)
	i := it.Data().labels[name].i[0]
	j := it.Data().labels[name].j[0]
	assert(t, i == 0, i)
	assert(t, j == 3, j)
	assert(t, !it.HasError())
}

func TestCombined1(t *testing.T) {
	// TODO add bad tests and check error
	s := "Hello \"World\""
	id := "ID"
	str := "String"
	quoted := "Quoted"

	rules := []Rule{
		{term{typ: 'N', name: "S"}, []term{
			{typ: 'N', name: id, marked: true},
			{typ: 'N', name: str, marked: true},
			{typ: 'N', name: quoted, marked: true},
		}},
		{term{typ: 'N', name: id, marked: true}, []term{
			{typ: 'T', name: id, terminal: termID()},
		}},
		{term{typ: 'N', name: str, marked: true}, []term{
			{typ: 'T', name: str, terminal: termStr(" ")},
		}},
		{term{typ: 'N', name: quoted, marked: true}, []term{
			{typ: 'T', name: quoted, terminal: termAnyQuoted()},
		}},
	}
	f, err := generateFunction(rules)
	assert(t, err == nil, err)
	it, err := NewIterator([]byte(s), true)
	assert(t, err == nil, err)
	execute(f, it)
	fmt.Println(it.Data().labels)
	assert(t, !it.HasError())
}

func TestCombined2(t *testing.T) {
	// TODO add bad tests and check error
	s := "GET / 200 OK"
	method := "Method"
	url := "Url"
	code := "Code"
	space := "Space"
	message := "Message"

	rules := []Rule{
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
	it, err := NewIterator([]byte(s), true)
	assert(t, err == nil, err)
	execute(f, it)
	fmt.Println(it.Data().labels)
	assert(t, !it.HasError())
}
