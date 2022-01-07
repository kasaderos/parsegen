package main

// /*
// 	1. Парсинг bnf
// 	2. Сбор сущностей
// 	3. Генерация кода
// */

// /*
// 	из теории трансляции
// 	Условия на входную БНФ
// 	1. Язык, порождаемый грамматикой желательно должна быть однозначной (КЗ)
// 		т.е. должен полностью описан и проверен
// 		т.к. проблема определения, является ли заданная КС грамматика однозначной,
// 		является алгоритмически неразрешимой.
// 	2. Желательно приведенной, т.е.
// 		2.1. отсутствие бесплодных сущностей
// 		2.2. отсутсвие недостижимых (неиспользуемых) правил
// 	TODO реализовать алгоритм удаления 1 и 2 (фича генератора парсера)
// */
// /*
// 	bnf:

// 	0. любое правило состоит из левой и правой части, левая ::= правая
// 	1. строки/символы в кавычках или апострофах "", ''
// 	2. выделяемая сущность titled,
// 	3. невыделяемые сущности определяются до конца (в конечном итоге будут состоять из набора символов)
// 		пример: A ::= B
// 				B ::= 'a'
// 				A - есть символ 'a'
// 	4. выделяемая сущность не определяется через := (не является левой)
// 	5. через знак или (|) можно определить некоторые подходящие варианты
// 	6. то что повторяется 0 и более раз обертывается в скобки со звездочкой ( <ident> )*
// 	7. рекурсии нет, т.к. есть повторения (пункт 6.)
// 	8. то что необязательно в прямые скобки [ <ident> ]
// 	9. слово result зарезервирована
// 	10. имена сущностей задается пользователем. При этом область видимости также контролируется пользователем

// B :=
// */
// /*
// 	bnf

// 	text := A " " B " " A
// 	A    := word

// 	<B>    := word

// 	text := (word)*

// 	text():
// 		s = []byte
// 		if !A() {

// 		}
// 		s.write()
// 		if !B() {

// 		}

// S = A B C
// A = 'A'
// A -> B -> C -> A

//
// 	ABABBB
// 	text := A
// 	A    := ABABAA | ABABBB
// 	aaa abb
// 	( word " " )* -> func (){
// 		for {
// 			if word() {
// 				return
// 			}
// 			if checkStr(" ") {
// 				return
// 			}
// 		}
// 	}
// */

// func isLvalue(word []byte) bool {
// 	return false
// }

// func isRvalue(word []byte) bool {
// 	return false
// }

// func isDelimeter(word []byte) bool {
// 	return false
// }

// // func lexic(bnfWords [][]byte) (*Rule, error) {
// // 	if len(bnfWords) < 3 {
// // 		return nil, errors.New("not enough elements in bnf rule")
// // 	}
// // 	first := bnfWords[0]
// // 	second := bnfWords[1]
// // 	if !isLvalue(first) {
// // 		return nil, errors.New("not lvalue")
// // 	}
// // 	if !isDelimeter(second) {
// // 		return nil, errors.New("no delimeter between lvalue and rvalue")
// // 	}
// // 	for i := 2; i < len(bnfWords); i++ {
// // 		if !isRvalue(bnfWords[i]) {
// // 			return nil, errors.New(fmt.Sprintf("word %s is not rvalue", bnfWords[i]))
// // 		}
// // 	}
// // 	return nil, nil
// // }

// func syntax(rules []*Rule) error {
// 	return nil
// }
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
