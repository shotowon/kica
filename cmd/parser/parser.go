package main

import (
	"fmt"
	"os"

	"github.com/shotowon/kica/internal/lexer"
	"github.com/shotowon/kica/internal/parser"
)

func main() {
	if len(os.Args) < 2 {
		panic(fmt.Sprintf("usage: %s <source-code>", os.Args[0]))
	}

	file, err := os.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	l := lexer.New(string(file))
	tokens, err := l.Lex()
	if err != nil {
		panic(err)
	}

	p := parser.New(tokens)
	stmts, err := p.Parse()
	if err != nil {
		panic(err)
	}

	for _, stmt := range stmts {
		fmt.Println(stmt)
	}
}
