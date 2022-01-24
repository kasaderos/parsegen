## **parsegen** v1.0
<span style="color:blue">#parser #generator \#golang \#bnf \#parsing </span>
### parsegen - experimental tool for extracting data based on BNF rules.

### Installation:

```
    go get github.com/kasaderos/parsegen
```

###	Example:
```go
package main

import (
	"log"

	"github.com/kasaderos/parsegen"
)

func main() {
	parser, err := parsegen.Generate([]byte(
		"S = Method sp Url sp StatusOk;" +
			"Method = any(0x20);" +
			"sp = 0x20 ;" +
			"Url = any(0x20);" +
			"StatusOk = integer;" +
			"integer = digit digits;" +
			"digit = 0x30-39 ;" +
			"digits = { digit };",
	))
	if err != nil {
		log.Fatal(err)
	}

	pd, err := parser.Parse([]byte("GET https://google.com 200"))
	if err != nil {
		log.Fatal(err)
	}
	pd.Print()
}
```
### Output
```
Method :
                GET
Url :
                https://google.com
StatusOk :
                200
```

### Guide:

	1. Rule types:
        non-terminals:
		[N] group:
            rule1 = A1 A2 A3 ... An ;         - strictly ordered group,  Ai in ([N], [T], [C], [L]), i = 1..n

		[L] logic, :
            rule2 = A1 | A2 | A3 ... | An ;   - first left matching,     Ai in ([N], [T], [C], [L]), i = 1..n

		[C] cycle :
            rule3 = { A } ;                   - zero or more times A,     A in ([N], [T], [C], [L])

        terminals (not rule):
        [T] : 
            1. string: 
                "some string"
            2. any ended with hex without including : 
                any(0x20)  matches "aa" from "aa bb"
            3. any ended with hex with including 0x20: 
                any[0x3b]  matches "aa;" from "aa;bb"
            4. identifier consisting a-z, A-Z, '-', '_', 0-9 symbols:
                my-id, my-id1  
            5. hex :
                0x20  - is equivalent to  A = " "
            6. hex interval:
                0x30-32  - is equivalent to rule A = "0" | "1" | "2" 

            (* any(hex), "empty", hexes reserved)        
            (* special case for [L]
                S = ["-"] (minus or nothing) is equivalent to S = "-" | empty )           

	2. Parsing starts from the initial 'S' start rule:
	    S = A B ;
		A = "A" ;
		B = "B" ;

	3. Each rule MUST ends with ';' symbol

	4. Each entity [N], [C], [L] MUST defined finally with terminals [T]

	5. Exported entities MUST start with capital letters:
		S = prefix Passwd ; 
		prefix = "password: " ; // not exported
		Passwd = "defcon99"   ; // exported

    (* S not exported by default)


### Remark 1.
Exprimental utility, depending on the rules, can generate a "bad" parser that parses ambiguously or
goes into an infinite loop. As rules for determining the stopping or uniqueness (by input) of given rules,
this is an algorithmically unsolvable problem. Therefore, the user checks the rules himself.

### Remark 2.
Checks have been added to avoid common mistakes like recursion (A = B; B = A;), BUT no checks for 
loops of non-terminals like (A = B; B = C; C = A;). Therefore, the graph of bnf rules should preferably be acyclic.

See more complex example in **parser_test.go**. Try additional tools  [bnfchecker](https://github.com/kasaderos/bnfchecker)


### Next versions v1.1-2.0
If this program helps to solve your problems, please donate. 
If the project gets a good donation, in the next version the following will be done:

- optimizations: reuse buffers, change some logic.
- add specific rules ```[ "may be" ], <m><n>entity, comments``` like in ABNF 
- add common terminals like ```digit, alpha, ...``` 


Donation:

   Yoomoney: 

        web: https://yoomoney.ru/

        ID: 4100117537690878 
