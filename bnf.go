package main

/*
	из теории трансляции
	Условия на входную БНФ
	1. Язык, порождаемый грамматикой желательно должна быть однозначной (КЗ)
		т.е. должен полностью описан и проверен
		т.к. проблема определения, является ли заданная КС грамматика однозначной,
		является алгоритмически неразрешимой.
	2. Желательно приведенной, т.е.
		2.1. отсутствие бесплодных сущностей
		2.2. отсутсвие недостижимых (неиспользуемых) правил
	TODO реализовать алгоритм удаления 1 и 2 (фича генератора парсера)
*/
//  Урезанный БНФ version 1.0.1:
//  TODO in 2.0.0:
//  - многозначные варианты
//  - откаты
//  - добавить вспомогательные правила [] {}
//  - проверка и упрощения правил (приведение правил)
//
// 	0. любое правило состоит из левой и правой части разделенные знаком '=' (апострофы не считаются)
// 	1. правила бывают двух типов:
//     R = A B C      группированные
//     R = A | B | C  вариантные
// 	2. сущности определяются однозначно, т.е. должны быть уникальными
// 		пример: A = B C
// 				B = "a"
//              C = "c1"
//              C = "c2"
//           C определяется дважды
//      (* empty - зарезервирована)
// 	3. через знак или (|) можно определить некоторые подходящие варианты
//  ограничения:
//		1. каждый вариант разделяется |
//      2. вариант должен содержать только одну сущность
//      3*. желательно чтобы каждый вариант был "очевидно" различим т.е.
//          R = A | B ;
//          A = "AA"
//          B = "AB"
//          входное слово "AB"
//          здесь A и B неочевидно различимы, т.к. парсинг проводится слева направо
//          сначала будет проверка A затем B
//          для предотвращения такого будет вызываться "откат" на начальный символ при входе в вариант
// 	4. циклы отсутствуют, но есть рекурсия:
//               правильная рекурсия: А = B | A
//  ограничения:
//      1. 3.1
//      2. 3.2
//      3. рекурсивная сущность определяется в конце
//  5. 1. В наборе правил каждый идентификатор должен быть полностью определен,
//        представляет некую строку (cм. п.2)
//     2. В наборе первое правило начинается с 'S'
//  6. если нужно определить необязательные символы, т.е. R = [ A ] B
//       можно заменить на эквивалентные правила
//       R = R1 | B ;
//       R1 = A B ;
//  7. если нужно определить цикл R = { A } (ноль и более раз A)
//     то
//       R = A | R ;
//     рекурсия работает пока A возвращает false
//  8. выделяемые сущности в <>
//      данные выделенных байт есть начало и конец lvalue
//                  lvalue = rvalue1 rvalue2
//                  выделяется все что входит [rvalue1, rvalue2]
//	9. Пример BNF
//      (* пробельные символы игнорируются заисключением пробелов в строках)
//      определение правила согласно нашему BNF:
//      S = rule | S
//      rule = lvalue "=" expr ";" ;
//      expr = rvalue expr1
//      expr1 = exprT1 | exprT2 ;
//      exprT1 = rvalue | exprT1;
//      exprT2 = rvalue1 | exprT2 ;
//      rvalue1 = "|" rvalue
//      lvalue = highlighted | id ;
// 		highlighted = "<" id ">" ;
//      rvalue = id | string
//
//      где id - идентификатор из латинских букв и чисел (ид начинается с лат букв)
//          string - любой набор символов в кавычках "" (example "This is a string")
/*

 */
//      rule = lvalue "=" expr ";" ;
//      expr = expr1 | expr2 | expr3 ;
//      expr1 = rvalue { SP { SP } rvalue }
//      expr2 = "{" rvalue "}"
//      expr3 = rvalue { "|" rvalue }
//      lvalue = id ;
//      rvalue = id | string ;
/*
	label:
		name string
		i, j  int   // i included, j not

	lex:
		name string
		data []byte

	get_parsed_data(
		[]entity, - list of entities (lvalues)
		bool,     - if true include strings
	) map[entity][]label

	get_slices(
		map[entity][]label, - result from get_parsed_data
		[]byte               - data
		bool                 - if true slice initial data, else allocate memory for each
	) map[entity][]lex
*/

const lid = "lid"
const rid = "rid"
const rstring = "string"

func bnfparser(it Iterator) (*function, error) {
	rules := []*Rule{
		{term{typ: 'N', name: "S", marked: true}, []term{
			{typ: 'N', name: "rule", marked: true},
		}},
		// rule = lvalue "=" expr ";" ;
		{term{typ: 'N', name: "rule", marked: true}, []term{
			{typ: 'C', name: "SP"},
			{typ: 'T', name: lid, marked: true, terminal: termID()},
			{typ: 'C', name: "SP"},
			{typ: 'T', name: "assign", marked: true, terminal: termStr("=")},
			{typ: 'C', name: "SP"},
			{typ: 'L', name: "expr", marked: true},
			{typ: 'C', name: "SP"},
			{typ: 'T', name: "end", marked: true, terminal: termStr(";")},
		}},
		{term{typ: 'L', name: "expr", marked: true}, []term{
			{typ: 'N', name: "expr1", marked: true},
			{typ: 'N', name: "expr2", marked: true},
			{typ: 'N', name: "expr3", marked: true},
		}},
		// rvalue SP rvalue { SP rvalue }
		{term{typ: 'N', name: "expr1", marked: true}, []term{
			{typ: 'L', name: "rvalue", marked: true},
			{typ: 'N', name: "SP1"},
			{typ: 'L', name: "rvalue", marked: true},
			{typ: 'C', name: "rvalues", marked: true},
		}},
		{term{typ: 'C', name: "rvalues"}, []term{
			{typ: 'N', name: "SP1"},
			{typ: 'L', name: "rvalue"},
		}},
		// "{" rvalue "}"
		{term{typ: 'N', name: "expr2", marked: true}, []term{
			{typ: 'T', name: "openBrace", marked: true, terminal: termStr("{")},
			{typ: 'C', name: "SP"},
			{typ: 'L', name: "rvalue", marked: true},
			{typ: 'C', name: "SP"},
			{typ: 'T', name: "closeBrace", marked: true, terminal: termStr("}")},
		}},
		// rvalue "|" rvalue { "|" rvalue }
		{term{typ: 'N', name: "expr3", marked: true}, []term{
			{typ: 'L', name: "rvalue"},
			{typ: 'N', name: "rvalue1", marked: true},
			{typ: 'C', name: "rvalue1S", marked: true},
		}},
		{term{typ: 'N', name: "rvalue1", marked: true}, []term{
			{typ: 'C', name: "SP"},
			{typ: 'T', name: "signOr", marked: true, terminal: termStr("|")},
			{typ: 'C', name: "SP"},
			{typ: 'L', name: "rvalue", marked: true},
		}},
		{term{typ: 'C', name: "rvalue1S", marked: true}, []term{
			{typ: 'N', name: "rvalue1"},
		}},
		{term{typ: 'L', name: "rvalue", marked: true}, []term{
			{typ: 'T', name: rid, marked: true, terminal: termID()},
			{typ: 'T', name: rstring, marked: true, terminal: termAnyQuoted()},
		}},
		{term{typ: 'C', name: "SP"}, []term{
			{typ: 'T', name: "sp", terminal: termSpace()},
		}},
		{term{typ: 'N', name: "SP1"}, []term{
			{typ: 'T', name: "sp", terminal: termSpace()},
			{typ: 'C', name: "SP"},
		}},
	}
	f, err := generateFunction(rules)
	if err != nil {
		return nil, err
	}
	return f, nil
}
