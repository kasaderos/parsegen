package main

import (
	"fmt"
	"testing"
)

func TestGenerateFunction1(t *testing.T) {
	rules := []*Rule{
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
	rules := []*Rule{
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
	rules := []*Rule{
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
	lbls := it.Data().labels
	fmt.Println(lbls)
	assert(t, len(lbls["A"].i) > 0 && lbls["A"].i[0] == 0 && lbls["A"].j[0] == 2)
	assert(t, ret == zero || ret == eof, "ret == err")
	printTree(f)
}

func TestBacktrackCycle(t *testing.T) {
	rules := []*Rule{
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
	lbls := it.Data().labels
	assert(t, lbls["A"].i[0] == 0 && lbls["A"].j[0] == 4)
	assert(t, ret == zero, "ret == err")
}

func TestExecuteCycleData(t *testing.T) {
	rules := []*Rule{
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
	lbls := it.Data().labels
	fmt.Println(lbls)
	assert(t, len(lbls["A"].i) == 3 && len(lbls["SP"].j) == 4)
	assert(t, ret == zero, "ret == err")
}

func TestBaseRule(t *testing.T) {
	parser, err := Generate([]byte(
		"S = \"Hello World\" \" \" \"!!!\";",
	))
	assert(t, err == nil, err)
	pd, err := parser.Parse([]byte("Hello World !!!"))
	assert(t, err == nil, err)
	fmt.Println(pd.labels)
}

func TestCaseRule(t *testing.T) {
	parser, err := Generate([]byte(
		"S = \"!!\" | \"Hello World\" | \"!\" ;",
	))
	assert(t, err == nil, err)
	pd, err := parser.Parse([]byte("Hello World"))
	assert(t, err == nil, err)
	fmt.Println(pd.labels)
}

func TestCycleRule(t *testing.T) {
	parser, err := Generate([]byte(
		"S = { \"Hello World;\" } ;",
	))
	assert(t, err == nil, err)
	pd, err := parser.Parse([]byte("Hello World;Hello World;"))
	assert(t, err == nil, err)
	fmt.Println(pd.labels)
}

func TestTwoRules(t *testing.T) {
	parser, err := Generate([]byte(
		"S = A \"!!!\" ;" +
			"A = { \"Hello World\" } ;",
	))
	assert(t, err == nil, err)
	pd, err := parser.Parse([]byte("Hello World!!!"))
	assert(t, err == nil, err)
	fmt.Println(pd.labels)
	pd, err = parser.Parse([]byte("!!!"))
	assert(t, err == nil, err)
	fmt.Println(pd.labels)
}

func TestHttpGetRequest(t *testing.T) {
	parser, err := Generate([]byte(
		"S = Method SP Url SP StatusOk;" +
			"Method = any(0x20);" +
			"SP = 0x20 ;" +
			"Url = any(0x20);" +
			"StatusOk = integer;",
	))
	assert(t, err == nil, err)
	pd, err := parser.Parse([]byte("GET https://google.com 200"))
	assert(t, err == nil, err)
	fmt.Println(pd.labels)
}

func TestBasicAny(t *testing.T) {
	end := byte(0)
	included := false
	assert(t, isTermAny([]byte("any(0x0d)"), &end, &included))
	assert(t, end == '\r' && !included)

	assert(t, isTermAny([]byte("any[0x0d]"), &end, &included))
	assert(t, end == '\r' && included)
}

/*
	Parser generator based on short BNF rules (SBNF) - experimental tool to parsing data by BNF rules.

	Example:
	....

	Guide:

	1. Rule types:
		rule1 = A1 A2 A3 ... An ;         - A1 than A2 than A3 .. than An, rule1 in [N], Ai in ([N], [T])
		rule2 = A1 | A2 | A3 ... | An ;   - A1 or A2 .. or An,             rule2 in [L], Ai in ([N], [T], [C])
		rule3 = { A } ;                   - zero or more times A,          rule3 in [C], A in ([N], [T], [C])

	2. Entities on the right side can be:
		A = "This is a string" ;     - string
		A = any(0x3b) ;                 - any character ended with ';'(0x3b) not including ';' (this character can be any '/', '@', ..)
		A = any[0x3b] ;                 - also like any(;) BUT ';' is included to A
		A = integer ;                - integer, 100500
		A = B ;                      - another entity, more about see p.5

	3. Parsing starts from the initial 'S' (start rule)
	    S = A B;
		A = "A";
		B = "B";

	4. Each rule MUST end with ';' symbol

	5. Each entity[N] MUST defined finally with terminals[T]
		[T]: string, any(c), any[c], integer.
		[N]: identifier like C variables, first character is letter, than digits or letters

	6. Exported entities must start with capital letters:
		S = prefix Passwd ;
		prefix = "password: " ; // not exported
		Passwd = "defcon99"   ; // exported

	7. If entity is cycled like A = { B } you can get exported data with: GetMany(key) [][]byte
	   If entity has only one instance you can use: Get(key) []byte

	Remarks 1.
	Exprimental utility, depending on the rules, can generate a "bad" parser that parses ambiguously or
	goes into an infinite loop. As rules for determining the stopping or uniqueness of given rules,
	this is an algorithmically unsolvable problem. Therefore, the user checks the rules himself.

	Remarks 2.
	Checks have been added to avoid common mistakes like recursion (A = B; B = A;)

	More complex examples ...
	reference bnf of 1-5: ...

	TODO

*/
// URI           = scheme ":" hier-part [ "?" query ] [ "#" fragment ]

//    hier-part     = "//" authority path-abempty
//                  / path-absolute
//                  / path-rootless
//                  / path-empty

//    URI-reference = URI / relative-ref

//    absolute-URI  = scheme ":" hier-part [ "?" query ]

//    relative-ref  = relative-part [ "?" query ] [ "#" fragment ]

//    relative-part = "//" authority path-abempty
//                  / path-absolute
//                  / path-noscheme
//                  / path-empty

//    scheme        = ALPHA *( ALPHA / DIGIT / "+" / "-" / "." )

//    authority     = [ userinfo "@" ] host [ ":" port ]
//    userinfo      = *( unreserved / pct-encoded / sub-delims / ":" )
//    host          = IP-literal / IPv4address / reg-name
//    port          = *DIGIT

//    IP-literal    = "[" ( IPv6address / IPvFuture  ) "]"

//    IPvFuture     = "v" 1*HEXDIG "." 1*( unreserved / sub-delims / ":" )

//    IPv6address   =                            6( h16 ":" ) ls32
//                  /                       "::" 5( h16 ":" ) ls32
//                  / [               h16 ] "::" 4( h16 ":" ) ls32
//                  / [ *1( h16 ":" ) h16 ] "::" 3( h16 ":" ) ls32
//                  / [ *2( h16 ":" ) h16 ] "::" 2( h16 ":" ) ls32
//                  / [ *3( h16 ":" ) h16 ] "::"    h16 ":"   ls32
//                  / [ *4( h16 ":" ) h16 ] "::"              ls32
//                  / [ *5( h16 ":" ) h16 ] "::"              h16
//                  / [ *6( h16 ":" ) h16 ] "::"

//    h16           = 1*4HEXDIG
//    ls32          = ( h16 ":" h16 ) / IPv4address
//    IPv4address   = dec-octet "." dec-octet "." dec-octet "." dec-octet

// Berners-Lee, et al.         Standards Track                    [Page 49]

// RFC 3986                   URI Generic Syntax               January 2005

//    dec-octet     = DIGIT                 ; 0-9
//                  / %x31-39 DIGIT         ; 10-99
//                  / "1" 2DIGIT            ; 100-199
//                  / "2" %x30-34 DIGIT     ; 200-249
//                  / "25" %x30-35          ; 250-255

//    reg-name      = *( unreserved / pct-encoded / sub-delims )

//    path          = path-abempty    ; begins with "/" or is empty
//                  / path-absolute   ; begins with "/" but not "//"
//                  / path-noscheme   ; begins with a non-colon segment
//                  / path-rootless   ; begins with a segment
//                  / path-empty      ; zero characters

//    path-abempty  = *( "/" segment )
//    path-absolute = "/" [ segment-nz *( "/" segment ) ]
//    path-noscheme = segment-nz-nc *( "/" segment )
//    path-rootless = segment-nz *( "/" segment )
//    path-empty    = 0<pchar>

//    segment       = *pchar
//    segment-nz    = 1*pchar
//    segment-nz-nc = 1*( unreserved / pct-encoded / sub-delims / "@" )
//                  ; non-zero-length segment without any colon ":"

//    pchar         = unreserved / pct-encoded / sub-delims / ":" / "@"

//    query         = *( pchar / "/" / "?" )

//    fragment      = *( pchar / "/" / "?" )

//    pct-encoded   = "%" HEXDIG HEXDIG

//    unreserved    = ALPHA / DIGIT / "-" / "." / "_" / "~"
//    reserved      = gen-delims / sub-delims
//    gen-delims    = ":" / "/" / "?" / "#" / "[" / "]" / "@"
//    sub-delims    = "!" / "$" / "&" / "'" / "(" / ")"
//                  / "*" / "+" / "," / ";" / "="
// */
