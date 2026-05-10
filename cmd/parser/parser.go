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

	l, err := lexer.New(os.Args[1:])
	tokenFiles, err := l.Lex()
	if err != nil {
		panic(err)
	}

	p := parser.New(tokenFiles)
	astFiles, err := p.Parse()

	if err != nil {
		panic(err)
	}

	for _, file := range astFiles {
		for _, stmt := range file.Statements() {
			fmt.Println(stmt)
		}
	}
}
