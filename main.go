package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

var InputFile = flag.String("f", "file", "file where described BNF of parser")

func main() {
	flag.Parse()

	file, err := os.Open(*InputFile)
	if err != nil {
		panic(err)
	}

	bnf, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}

	if err = genParser(bnf); err != nil {
		fmt.Println(err)
	}
}
