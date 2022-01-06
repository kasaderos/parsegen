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
