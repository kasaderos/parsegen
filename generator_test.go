package main

import (
	"fmt"
	"testing"
)

func TestGenerateFunction1(t *testing.T) {
	rules := []Rule{
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
	rules := []Rule{
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

func TestWithExec1(t *testing.T) {
	n := 0
	rules := []Rule{
		{term{typ: 'N', name: "S"}, []term{
			{typ: 'N', name: "A"},
			{typ: 'N', name: "B"},
			{typ: 'N', name: "C"},
		}},
		{term{typ: 'N', name: "A"}, []term{
			{typ: 'T', name: "Terminal", terminal: func(it Iterator) code { n++; return zero }},
		}},
		{term{typ: 'N', name: "B"}, []term{
			{typ: 'T', name: "Terminal", terminal: func(it Iterator) code { n++; return zero }},
		}},
		{term{typ: 'N', name: "C"}, []term{
			{typ: 'N', name: "D"},
			{typ: 'N', name: "E"},
		}},
		{term{typ: 'N', name: "D"}, []term{
			{typ: 'T', name: "Terminal", terminal: func(it Iterator) code { n++; return zero }},
		}},
		{term{typ: 'N', name: "E"}, []term{
			{typ: 'T', name: "Terminal", terminal: func(it Iterator) code { n++; return zero }},
		}},
	}
	f, err := generateFunction(rules)
	assert(t, err == nil, err)
	assert(t, n != 4, fmt.Sprintf("n != 4; n = %d", n))
	it, _ := NewIterator([]byte("Example"), true)
	ret := execute(f, it)
	assert(t, ret == zero, "ret == err or oef")
}

func TestExecuteCycle1(t *testing.T) {
	n := 0
	// S = A
	// A = {Terminal}, A - cycle
	rules := []Rule{
		{term{typ: 'N', name: "S"}, []term{
			{typ: 'C', name: "A"},
		}},
		{term{typ: 'C', name: "A"}, []term{
			{typ: 'T', name: "Terminal", terminal: func(it Iterator) code {
				if n >= 3 {
					return missed
				}
				n++
				fmt.Println("A")
				return zero
			}},
		}},
	}
	f, err := generateFunction(rules)
	assert(t, err == nil, err)
	it, _ := NewIterator([]byte("Example"), true)
	ret := execute(f, it)
	assert(t, ret == zero, "ret true")
	assert(t, n == 3, fmt.Sprintf("n != 3; n = %d", n))
	printTree(f)
}

func TestBacktrackLogic(t *testing.T) {
	rules := []Rule{
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
	it, err := NewIterator([]byte("AB"), true)
	assert(t, err == nil, err)
	ret := execute(f, it)
	lbls := it.Data().labels
	assert(t, lbls["A"].i == 0 && lbls["A"].j == 2)
	assert(t, ret == zero || ret == eof, "ret == err")
	printTree(f)
}

func TestBacktrackCycle(t *testing.T) {
	rules := []Rule{
		{term{typ: 'N', name: "S"}, []term{
			{typ: 'C', name: "A", marked: true},
		}},
		{term{typ: 'C', name: "A", marked: true}, []term{
			{typ: 'T', name: "T1", terminal: termStr("AB")},
		}},
	}
	f, err := generateFunction(rules)
	assert(t, err == nil, err)
	it, err := NewIterator([]byte("ABABBC"), true)
	assert(t, err == nil, err)
	ret := execute(f, it)
	lbls := it.Data().labels
	assert(t, lbls["A"].i == 0 && lbls["A"].j == 4)
	assert(t, ret == zero, "ret == err")
}

// func TestGenerateRules(t *testing.T) {
// 	// parser, err := Generate([]byte(
// 	// 	"S = A | S;" +
// 	// 		"A = \"-\" ;",
// 	// ))
// 	// Uri

// 	parser, err := Generate([]byte(
// 		"S = URI ;" +
// 			"URI = scheme \":\" path ;" +
// 			"scheme = Any(:) ;" +
// 			"path = Any() ;",
// 	))
// 	assert(t, err == nil, err)
// 	_ = parser
// 	/*
// 		[A] B [C]
// 		V1 | V2 | V3
// 		V1 = A B
// 		V2 = B C
// 		V3 = A B C
// 	*/
// }

// func TestWithExec2(t *testing.T) {
// 	n := 0
// 	rules := []Rule{
// 		{term{typ: 'N', name: "S"}, []term{
// 			{typ: 'L', name: "A"},
// 			{typ: 'N', name: "B"},
// 		}},
// 		{term{typ: 'L', name: "A"}, []term{
// 			{typ: 'T', name: "Terminal", terminal: func(it Iterator) bool { n++; fmt.Println("A"); return false }},
// 			{typ: 'N', name: "C"},
// 		}},
// 		{term{typ: 'N', name: "C"}, []term{
// 			{typ: 'T', name: "Terminal", terminal: func(it Iterator) bool { fmt.Println("C"); return false }},
// 		}},
// 		{term{typ: 'N', name: "B"}, []term{
// 			{typ: 'T', name: "Terminal", terminal: func(it Iterator) bool { n++; fmt.Println("B"); return false }},
// 		}},
// 	}
// 	f, err := generateFunction(rules)
// 	assert(t, err == nil, err)
// 	it, _ := NewIterator([]byte("Example"), true)
// 	ret := execute(f, it)
// 	assert(t, !ret, "ret true")
// 	assert(t, n == 2, fmt.Sprintf("n != 2; n = %d", n))
// }

// func TestWithExec3(t *testing.T) {
// 	n := 0
// 	rules := []Rule{
// 		{term{typ: 'N', name: "S"}, []term{
// 			{typ: 'L', name: "A"},
// 			{typ: 'N', name: "B"},
// 		}},
// 		{term{typ: 'L', name: "A"}, []term{
// 			{typ: 'T', name: "Terminal", terminal: func(it Iterator) bool { n++; fmt.Println("A"); return true }},
// 			{typ: 'N', name: "C"},
// 		}},
// 		{term{typ: 'N', name: "C"}, []term{
// 			{typ: 'T', name: "Terminal", terminal: func(it Iterator) bool { n++; fmt.Println("C"); return false }},
// 		}},
// 		{term{typ: 'N', name: "B"}, []term{
// 			{typ: 'T', name: "Terminal", terminal: func(it Iterator) bool { n++; fmt.Println("B"); return false }},
// 		}},
// 	}
// 	f, err := generateFunction(rules)
// 	assert(t, err == nil, err)
// 	it, _ := NewIterator([]byte("Example"), true)
// 	ret := execute(f, it)
// 	assert(t, !ret, "ret true")
// 	assert(t, n == 3, fmt.Sprintf("n != 3; n = %d", n))
// }

// /*
// // Generate([]byte(
// 	// 	"S = URI" +
// 	// 		"URI = scheme \":\" hierPart Query Fragment ;" +
// 	// 		"hierPart = hierPart1 | pathAbsolute | pathRootLess | pathEmpty ;" +
// 	// 		"hierPart1 = \"//\" authority path-abempty ;" +
// 	// 		"scheme = ALPHA scheme1 ;" +
// 	// 		"scheme1 = ALPHA | DIGIT | \"+\" | \"-\" | \".\" | empty | scheme1 ;" +
// 	// 		"authority = UserInfoHostPort | UserInfoHost | HostPort  ;" +
// 	// 		"UserInfoHostPort = userinfo \"@\" host \":\" port ;"+
// 	// 		"UserInfoHost = userinfo \"@\" host ;" +
// 	// 		"HostPort = host \":\" port ;" +
// 	// ))
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

// // "S = SIPRequestLine ;" +
// // 	"SIPRequestLine = Method SP RequestUri SP SIPVersion CRLF ;" +
// // 	"RequestUri = SIPUri | SIPSUri | absoluteUri ;" +
// // 	"absoluteURI = absoluteURIHierPart | opaquePart ;" +
// // 	"absoluteURIHierPart = scheme \":\" hierPart ;" +
// // 	"absoluteURIOpaquePart = scheme \":\" opaquePart ;" +
// // 	"hierPart = hierPart1  ;" +
// // 	"hierPart1 = NetPath | AbsPath ;" +
// // 	"NetPath = "
