package parsegen

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
	f := function{
		typ: 'N',
		funcs: []*function{
			{typ: 'T', name: "T1", terminal: termStr("Exam")},
			{typ: 'T', name: "T2", terminal: termStr("ple")},
		},
	}
	it, _ := NewIterator([]byte("Example"))
	ret := execute(&f, it)
	assert(t, ret == eof)
}

func TestExecuteComplex1(t *testing.T) {
	f := function{
		typ: 'N',
		funcs: []*function{
			{typ: 'N', funcs: []*function{
				{typ: 'T', terminal: termStr("1")},
			}},
			{typ: 'N', funcs: []*function{
				{typ: 'T', terminal: termStr("2")},
			}},
			{typ: 'N', funcs: []*function{
				{typ: 'T', terminal: termStr("3")},
			}},
			{typ: 'T', terminal: termStr("4")},
		},
	}
	it, _ := NewIterator([]byte("1234"))
	ret := execute(&f, it)
	assert(t, ret == eof)
}

func TestExecuteComplex2(t *testing.T) {
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
	f2.funcs = append(f2.funcs, &function{typ: 'T', terminal: termStr("1")})
	f.funcs = append(f.funcs, &f2)
	it, _ := NewIterator([]byte("1"))
	ret := execute(&f, it)
	assert(t, ret == eof)
}

func TestExecuteLogic1(t *testing.T) {
	f := function{typ: 'N'}
	f.funcs = append(f.funcs, &function{
		typ: 'L',
		funcs: []*function{
			{typ: 'T', terminal: termStr("1")},
			{typ: 'T', terminal: termStr("2")},
		},
	},
	)
	it, err := NewIterator([]byte("1"))
	assert(t, err == nil)
	ret := execute(&f, it)
	assert(t, ret == eof)
	it, err = NewIterator([]byte("2"))
	assert(t, err == nil)
	ret = execute(&f, it)
	assert(t, ret == eof)
}

func TestExecuteLogic3(t *testing.T) {
	f := function{typ: 'N'}
	f.funcs = append(f.funcs, &function{
		typ: 'L',
		funcs: []*function{
			{typ: 'T', terminal: termStr("1")},
			{typ: 'T', terminal: termStr("2")},
			{typ: 'T', terminal: termStr("3")},
		},
	},
	)
	it, _ := NewIterator([]byte("3"))
	ret := execute(&f, it)
	assert(t, ret == eof)
}

func TestExecuteLogic4(t *testing.T) {
	/*
		N:
			L:
				L:
					1 | 2
				3
	*/
	f := function{typ: 'N'}
	f2 := function{typ: 'L'}
	f2.funcs = append(f2.funcs, &function{
		typ: 'L',
		funcs: []*function{
			{typ: 'T', terminal: termStr("1")},
			{typ: 'T', terminal: termStr("2")},
		},
	},
	)
	f2.funcs = append(f2.funcs, &function{typ: 'T', terminal: termStr("3")})
	f.funcs = append(f.funcs, &f2)
	it, err := NewIterator([]byte("3"))
	assert(t, err == nil)
	ret := execute(&f, it)
	assert(t, ret == eof)
	it, err = NewIterator([]byte("4"))
	assert(t, err == nil)
	ret = execute(&f, it)
	assert(t, ret == missed)
}

func TestExecuteLogicBad(t *testing.T) {
	f := function{typ: 'N'}
	f.funcs = append(f.funcs, &function{
		typ: 'L',
		funcs: []*function{
			{typ: 'T', terminal: termStr("1")},
			{typ: 'T', terminal: termStr("2")},
			{typ: 'T', terminal: termStr("3")},
		},
	},
	)
	it, err := NewIterator([]byte("4"))
	assert(t, err == nil)
	ret := execute(&f, it)
	assert(t, ret == missed)
}

func TestExecuteCycle(t *testing.T) {
	f := function{typ: 'C'}
	f.funcs = append(f.funcs, &function{typ: 'T', terminal: termStr("1")})

	it, err := NewIterator([]byte("1111"))
	assert(t, err == nil)
	ret := execute(&f, it)
	assert(t, ret == zero, ret.String())
}
