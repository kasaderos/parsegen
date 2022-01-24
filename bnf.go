package parsegen

/*
Below are the rules for parsing single bnf rule

S = rule ;
rule = lvalue "=" expr ;
expr = exprBase exprCycle ;
exprBase = SP rvalue lastPart ;
lastPart = endPart | spRvalues | orRvalues ;
spRvalues = spRvalue spRvalueCycle endPart ;
orRvalues = orRvalue orRvalueCycle endPart ;
spRvalue = { SP1 rvalue } ;
spRvalueCycle = { spRvalue } ;
orRvalue = or SP rvalue ;
orRvalueCycle = { or rvalue } ;
or =  SP  "|" ;
expr2 = "{" rvalue "}" endPart ;
rvalue = basicAnyTerm | hexTerm | id | string ;
basicAnyTerm = "any" hexPart ;
hexPart = notIncluded | included ;
notIncluded = "(" hexTerm ")" ;
included = "[" hexTerm "]" ;
hexTerm = "0x" hexdig hexdig interv ;
hexdig = digit | 0x0d | 0x09 | 0x0a | 0x20 ;
digit = 0x30-0x39 ;
interv = intervHexPart | empty ;
intervHexPart = "-" hexdig hexdig ;
SP = { sp } ;
sp = 0x0d | 0x0a | 0x09 | 0x20 ;
SP1 = sp SP ;
endPart = SP ";" ;
lvalue = SP id ;
assign = SP "=" ;
*/
const lidTerm = "lid"
const ridTerm = "rid"
const stringTerm = "string"
const basicAnyTerm = "basicTerm"
const hexTerm = "hexTerm"

// bnfFunction returns a function that parses single BNF rule.
func bnfFunction(it Iterator) (*function, error) {
	rules := []*rule{
		// S = rule ;
		{term{typ: 'N', name: "S", marked: true}, []term{
			{typ: 'N', name: "rule", marked: true},
		}},
		// rule = lvalue "=" expr ;
		{term{typ: 'N', name: "rule", marked: true}, []term{
			{typ: 'N', name: "lvalue"},
			{typ: 'N', name: "assign"},
			{typ: 'L', name: "expr"},
		}},
		// expr = exprBase exprCycle ;
		{term{typ: 'L', name: "expr", marked: true}, []term{
			{typ: 'N', name: "exprBase"},
			{typ: 'N', name: "exprCycle"},
		}},
		// exprBase = SP rvalue lastPart ;
		{term{typ: 'N', name: "exprBase", marked: true}, []term{
			{typ: 'C', name: "SP"},
			{typ: 'L', name: "rvalue"},
			{typ: 'L', name: "lastPart"},
		}},
		// lastPart = endPart | spRvalues | orRvalues ;
		{term{typ: 'L', name: "lastPart"}, []term{
			{typ: 'N', name: "endPart"},
			{typ: 'N', name: "spRvalues"},
			{typ: 'N', name: "orRvalues"},
		}},
		// spRvalues = spRvalue spRvalueCycle endPart ;
		{term{typ: 'N', name: "spRvalues", marked: true}, []term{
			{typ: 'N', name: "spRvalue"},
			{typ: 'C', name: "spRvalueCycle"},
			{typ: 'N', name: "endPart"},
		}},
		// orRvalues = orRvalue orRvalueCycle endPart ;
		{term{typ: 'N', name: "orRvalues", marked: true}, []term{
			{typ: 'N', name: "orRvalue"},
			{typ: 'C', name: "orRvalueCycle"},
			{typ: 'N', name: "endPart"},
		}},
		// spRvalue = { SP1 rvalue } ;
		{term{typ: 'N', name: "spRvalue"}, []term{
			{typ: 'N', name: "SP1"},
			{typ: 'L', name: "rvalue"},
		}},
		// spRvalueCycle = { spRvalue } ;
		{term{typ: 'C', name: "spRvalueCycle"}, []term{
			{typ: 'N', name: "spRvalue"},
		}},
		// orRvalue = or SP rvalue ;
		{term{typ: 'N', name: "orRvalue"}, []term{
			{typ: 'N', name: "or"},
			{typ: 'C', name: "SP"},
			{typ: 'L', name: "rvalue"},
		}},
		// orRvalueCycle = { or rvalue } ;
		{term{typ: 'C', name: "orRvalueCycle"}, []term{
			{typ: 'N', name: "orRvalue"},
		}},
		// or =  SP  "|" ;
		{term{typ: 'N', name: "or"}, []term{
			{typ: 'C', name: "SP"},
			{typ: 'T', name: "signOr", terminal: termStr("|")},
		}},
		// expr2 = "{" rvalue "}" endPart ;
		{term{typ: 'N', name: "exprCycle", marked: true}, []term{
			{typ: 'C', name: "SP"},
			{typ: 'T', name: "openBrace", terminal: termStr("{")},
			{typ: 'C', name: "SP"},
			{typ: 'L', name: "rvalue"},
			{typ: 'C', name: "SP"},
			{typ: 'T', name: "closeBrace", terminal: termStr("}")},
			{typ: 'N', name: "endPart"},
		}},
		// rvalue = basicAnyTerm | hexTerm | id | string ;
		// basicAnyTerm = "any" hexPart ;
		// hexPart = notIncluded | included ;
		// notIncluded = "(" hexTerm ")" ;
		// included = "[" hexTerm "]" ;
		// hexTerm = "0x" hexdig hexdig interv ;
		// hexdig = digit | 0x0d | 0x09 | 0x0a | 0x20 ;
		// digit = 0x30-0x39 ;
		// interv = intervHexPart | empty ;
		// intervHexPart = "-" hexdig hexdig ;
		{term{typ: 'L', name: "rvalue", marked: true}, []term{
			{typ: 'T', name: basicAnyTerm, marked: true, terminal: termBasicAny()},
			{typ: 'T', name: hexTerm, marked: true, terminal: basicHex()}, // 0xff or 0x00-0xff
			{typ: 'T', name: ridTerm, marked: true, terminal: termID()},
			{typ: 'T', name: stringTerm, marked: true, terminal: termAnyQuoted()},
		}},
		// SP = { sp } ;
		// sp = 0x0d | "\n" | "\t" | " "
		{term{typ: 'C', name: "SP"}, []term{
			{typ: 'T', name: "sp", terminal: termSpace()},
		}},
		// SP1 = sp SP
		{term{typ: 'N', name: "SP1"}, []term{
			{typ: 'T', name: "sp", terminal: termSpace()},
			{typ: 'C', name: "SP"},
		}},
		// endPart = SP ";"
		{term{typ: 'N', name: "endPart"}, []term{
			{typ: 'C', name: "SP"},
			{typ: 'T', name: "end", marked: true, terminal: termStr(";")},
			{typ: 'C', name: "SP"},
		}},
		// lvalue = SP id ;
		{term{typ: 'N', name: "lvalue"}, []term{
			{typ: 'C', name: "SP"},
			{typ: 'T', name: lidTerm, marked: true, terminal: termID()},
		}},
		// assign = SP "=" ;
		{term{typ: 'N', name: "assign"}, []term{
			{typ: 'C', name: "SP"},
			{typ: 'T', name: "assignMark", marked: true, terminal: termStr("=")},
		}},
	}
	f, err := generateFunction(rules)
	if err != nil {
		return nil, err
	}
	return f, nil
}
