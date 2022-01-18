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
//      expr = rvalue expr1
//      expr1 = exprT1 | exprT2 ;
//      exprT1 = rvalue | empty | exprT1;
//      exprT2 = rvalue1 | empty | exprT2 ;
//      rvalue1 = "|" rvalue
//      lvalue = highlighted | id ;
// 		highlighted = "<" id ">" ;
//      rvalue = id | string
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

func bnfparser(it Iterator) (*function, error) {
	rules := []Rule{
		{term{typ: 'L', name: "S", marked: true}, []term{
			{typ: 'N', name: "rule", marked: true},
			{typ: 'L', name: "S", marked: true},
		}},
		{term{typ: 'N', name: "rule", marked: true}, []term{
			{typ: 'L', name: "lvalue", marked: true},
			{typ: 'T', name: "assign", marked: true, terminal: termStr("=")},
			{typ: 'N', name: "expr", marked: true},
			{typ: 'T', name: "end", marked: true, terminal: termStr(";")},
		}},
		{term{typ: 'L', name: "lvalue", marked: true}, []term{
			{typ: 'T', name: "lid", marked: true, terminal: termID()},
			{typ: 'N', name: "highlighted", marked: true},
		}},
		{term{typ: 'N', name: "highlighted", marked: true}, []term{
			{typ: 'T', name: "openH", marked: true, terminal: termStr("<")},
			{typ: 'T', name: "lid", marked: true, terminal: termID()},
			{typ: 'T', name: "openC", marked: true, terminal: termStr(">")},
		}},

		{term{typ: 'N', name: "expr", marked: true}, []term{
			{typ: 'L', name: "rvalue", marked: true},
			{typ: 'L', name: "expr1", marked: true},
		}},
		{term{typ: 'L', name: "expr1", marked: true}, []term{
			{typ: 'L', name: "exprT1", marked: true},
			{typ: 'L', name: "exprT2", marked: true},
		}},
		{term{typ: 'L', name: "exprT1", marked: true}, []term{
			{typ: 'L', name: "rvalue", marked: true},
			{typ: 'T', name: "empty", terminal: termEmpty()},
			{typ: 'L', name: "exprT1", marked: true},
		}},
		{term{typ: 'L', name: "exprT2", marked: true}, []term{
			{typ: 'N', name: "rvalue1", marked: true},
			{typ: 'T', name: "empty", terminal: termEmpty()},
			{typ: 'L', name: "exprT2", marked: true},
		}},
		{term{typ: 'N', name: "rvalue1", marked: true}, []term{
			{typ: 'T', name: "signOr", marked: true, terminal: termStr("|")},
			{typ: 'L', name: "rvalue", marked: true},
		}},
		{term{typ: 'L', name: "rvalue", marked: true}, []term{
			{typ: 'T', name: "rid", marked: true, terminal: termID()},
			{typ: 'T', name: "string", marked: true, terminal: termAnyQuoted()},
		}},
	}
	f, err := generateFunction(rules)
	if err != nil {
		return nil, err
	}
	return f, nil
}
