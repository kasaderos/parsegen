package main

// import "errors"

// var ErrNotMatched = errors.New("not matched")

// type tNode int32

// const (
// 	Root tNode = iota
// 	Term
// 	NonTerm
// 	Optional
// 	Cycle
// 	Group
// )

// type tKey int32

// const (
// 	S tKey = iota
// 	SP
// 	LVAL
// 	STR
// 	RVALS
// 	RVAL
// 	ENT
// 	ANY
// )

// type Node struct {
// 	tState
// 	child []Node
// }

// type Tree struct {
// 	root Node
// }

// type tState struct {
// 	*Parser
// 	key tKey
// 	typ tNode
// 	fn  fsmState
// }

// type tRule struct {
// 	lVal    tKey
// 	rGroups [][]tState
// }

// func Space(p *Parser) {
// 	if p.cc() != ' ' {
// 		p.err = ErrNotMatched
// 		return
// 	}
// 	p.gc()
// }

// func String(p *Parser) {
// 	if p.cc() != '"' {
// 		p.err = ErrNotMatched
// 		return
// 	}
// 	p.gc()
// 	start := p.gp()
// 	for !p.eof && p.cc() != '"' {
// 		p.gc()
// 	}
// 	p.result = p.slice(start, p.gp())
// 	p.gc()
// }

// func Any(p *Parser, end byte) {
// 	start := p.gp()
// 	for !p.eof && p.cc() != end {
// 		p.gc()
// 	}
// 	p.result = p.slice(start, p.gp())
// }

// func Entity(p *Parser) {
// 	start := p.gp()
// 	for !p.eof && isAlpha(p.cc()) {
// 		p.gc()
// 	}
// 	p.result = p.slice(start, p.gp())
// }

// func Delimeter(p *Parser) {
// 	if p.cc() != '=' {
// 		p.err = ErrNotMatched
// 		return
// 	}
// 	p.gc()
// }

// func EndRule(p *Parser) {
// 	if p.cc() != ';' {
// 		p.err = ErrNotMatched
// 		return
// 	}
// 	p.gc()
// }

// /*
// 	1. key1 := key12 .. key1n
// 	   key2 := key22 .. key2n
// */

// // func (t *Tree) Init(rules []tRule) {
// // 	for _, rule := range rules {
// // 		t.insert(rule)
// // 	}
// // }

// // func (t *Tree) insert(rule tRule) {
// // 	if rule.lVal == S {
// // 		for _, group := range rule.rGroups {
// // 			for _, tState := range group {
// // 				// Node{tState, nil}
// // 			}
// // 		}
// // 	}
// // }

// // rule := space lvalue space ":=" rvalues
// // lvalue := entity
// // rvalues := (space rvalue)(1:n)*
// // rvalue := entity | string
// // entity := word
// // string := <"> any <"> without quotes
// /*
// 	rule
// 	/   \          \    \     \
//   space  lvalue  space  "=" rvalues
//                             /                  \
// 						 group1(entity)  group2(string);
// 	T Tree
// 	T.insert("T", space)
// 	T.insert("T", lvalue)
// 	T.insert("T", space)
// 	T.insert("T", str(":="))

// 	S = { rule }
// 	rule = { spaces } <entity> { spaces } ":=" { spaces } expr { spaces } ";"
// 	expr = rvalues { "|" rvalues }
// 	rvalues = rvalue { rvalue }
// 	<rvalue> := entity | string | "[" expr "]" |

// 	cccccccc ... abcabcabc

// 	S := ABC
// 	ABC := (str)*(1:3)
// 	str := "abc"

// 	S = { rule }
// 	rule = lvalue "=" rvalues
// 	lvalue = entity
// 	rvalues = rvalue {"|" rvalue}
// 	rvalue = entity | string | "[" rvalues "]"
// 										 			 _____    _______
// 			                            			|      ^  |      ^
// 	s -> rule -> lvalue                 "=" -> rvalues -> rvalue -> "|"
// 					|          		     |		            |
// 					-> { SP2 } -> entity-> 					-> entity
// 															|
// 															-> string()
// 															|
// 															-> "["
// 		________
// 	  	v	    |
// 	s -> q -> "(" -> ")" -> e
//     |
// 	e
// 	задача проверки определения всех сущностей
// 	задача на проверку однозначности неразрешима
// 	но пример неправильного правила
// 	A := (SP)* " := "
// 	"   := "

// */
