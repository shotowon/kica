package lexer

import (
	"fmt"
	"os"
	"strings"
	"unicode"

	"github.com/shotowon/kica/internal/gears"
	"github.com/shotowon/kica/internal/tokens"
)

type lexer struct {
	sources [][]rune
	loc     gears.Location
	currSrc int
}

func New(sources []string) (*lexer, error) {
	l := new(lexer)

	if len(sources) < 1 {
		return nil, fmt.Errorf("source files not provided")
	}

	for _, src := range sources {
		contents, err := os.ReadFile(src)
		if err != nil {
			return nil, err
		}

		l.sources = append(l.sources, []rune(string(contents)))
	}

	l.loc = gears.ZeroLoc()
	return l, nil
}

func (l *lexer) Lex() ([]tokens.File, error) {
	files := make([]tokens.File, len(l.sources))

	for i := range l.sources {
		l.setCurrentSource(i)
		toks, err := l.lexFile()
		if err != nil {
			return nil, err
		}

		files[i] = *tokens.NewFile(toks)
	}

	return files, nil
}

func (l *lexer) lexFile() ([]*tokens.Token, error) {
	toks := make([]*tokens.Token, 0)
	for l.curr() != '0' {
		if unicode.IsSpace(l.curr()) {
			l.advance()
			continue
		}

		var (
			tok *tokens.Token
			err error
		)

		switch {
		case l.curr() == '"':
			tok, err = l.handleString()
		case unicode.IsLetter(l.curr()):
			tok, err = l.handleID()
		case unicode.IsDigit(l.curr()):
			tok, err = l.handleNumber()
		default:
			tok, err = l.handleSingleChar()
		}

		if err != nil {
			return nil, err
		}

		toks = append(toks, tok)
	}

	toks = append(toks, &tokens.Token{
		Kind: tokens.EOF,
	})

	return toks, nil
}

func (l *lexer) handleNumber() (*tokens.Token, error) {
	if l.curr() == '0' {
		return nil, fmt.Errorf("unexpected EOF at %s", l.loc)
	}

	if !unicode.IsDigit(l.curr()) {
		return nil, fmt.Errorf("expected digit at %s", l.loc)
	}

	tok := tokens.Token{}
	tok.Location = l.loc
	tok.Kind = tokens.Int

	literal := strings.Builder{}

	for unicode.IsDigit(l.curr()) {
		literal.WriteRune(l.curr())
		l.advance()

		if l.curr() == '.' && unicode.IsDigit(l.peek()) {
			tok.Kind = tokens.Float
			literal.WriteRune(l.curr())
			l.advance()
		}
	}

	tok.Location.End = l.loc.Offset
	tok.Literal = literal.String()

	return &tok, nil
}

func (l *lexer) handleSingleChar() (*tokens.Token, error) {
	if l.curr() == '0' {
		return nil, fmt.Errorf("unexpected EOF at %s", l.loc)
	}

	tok := tokens.Token{}
	tok.Location = l.loc
	tok.Kind = tokens.EOF

	switch l.curr() {
	case '(':
		tok.Kind = tokens.LParen
	case ')':
		tok.Kind = tokens.RParen
	case '{':
		tok.Kind = tokens.LBrace
	case '<':
		tok.Kind = tokens.LT
	case '>':
		tok.Kind = tokens.GT
	case '}':
		tok.Kind = tokens.RBrace
	case '+':
		tok.Kind = tokens.Plus
	case '-':
		tok.Kind = tokens.Minus
	case '*':
		tok.Kind = tokens.Mul
	case '/':
		tok.Kind = tokens.Div
	case '=':
		tok.Kind = tokens.Eq
	case ':':
		tok.Kind = tokens.Colon
	case ';':
		tok.Kind = tokens.Semicolon
	case ',':
		tok.Kind = tokens.Comma
	default:
		return nil, fmt.Errorf("invalid token at %s", l.loc)
	}

	tok.Literal = string(l.curr())
	l.advance()
	tok.Location.End = l.loc.Offset
	return &tok, nil
}

func (l *lexer) handleID() (*tokens.Token, error) {
	if l.curr() == '0' {
		return nil, fmt.Errorf("unexpected EOF at %s", l.loc)
	}

	if !unicode.IsLetter(l.curr()) {
		return nil, fmt.Errorf("expected letter for ID token at %s", l.loc)
	}

	tok := tokens.Token{}
	tok.Location = l.loc

	literal := strings.Builder{}
	literal.WriteRune(l.curr())
	l.advance()

	for unicode.IsLetter(l.curr()) || unicode.IsDigit(l.curr()) || l.curr() == '_' {
		if l.curr() == '0' {
			return nil, fmt.Errorf("unexpected EOF at %s", l.loc)
		}

		literal.WriteRune(l.curr())
		l.advance()
	}

	tok.Literal = literal.String()
	tok.Location.End = l.loc.Offset
	tok.Kind = tokens.IDtoKeyword(tok.Literal)

	return &tok, nil
}

func (l *lexer) handleString() (*tokens.Token, error) {
	if l.curr() == '0' {
		return nil, fmt.Errorf("unexpected EOF at %s", l.loc)
	}

	if l.curr() != '"' {
		return nil, fmt.Errorf("expected letter for ID token at %s", l.loc)
	}

	tok := tokens.Token{}
	tok.Location = l.loc
	tok.Kind = tokens.String

	literal := strings.Builder{}
	literal.WriteRune(l.curr())
	l.advance()

	if l.curr() == '0' {
		return nil, fmt.Errorf("unexpected EOF at %s", l.loc)
	}

	backSlashes := 0

	for {
		if l.curr() == '0' {
			return nil, fmt.Errorf("unexpected EOF at %s", l.loc)
		}

		if l.curr() == '"' {
			if backSlashes%2 == 0 {
				break
			}
		}

		if l.curr() == '\\' {
			backSlashes++
		} else {
			backSlashes = 0
		}

		literal.WriteRune(l.curr())
		l.advance()
	}

	if l.curr() == '0' {
		return nil, fmt.Errorf("unexpected EOF at %s", l.loc)
	}

	literal.WriteRune(l.curr())
	l.advance()

	tok.Literal = literal.String()
	tok.Location.End = l.loc.Offset
	return &tok, nil
}

func (l *lexer) advance() {
	if l.loc.Offset >= len(l.currentSource()) {
		return
	}

	if l.curr() == '\n' {
		l.loc.Line++
		l.loc.Col = 1
	} else {
		l.loc.Col++
	}

	l.loc.Offset++
}

func (l *lexer) curr() rune {
	if l.loc.Offset >= len(l.currentSource()) {
		return '0'
	}

	return l.currentSource()[l.loc.Offset]
}

func (l *lexer) peek() rune {
	if l.loc.Offset+1 >= len(l.currentSource()) {
		return '0'
	}

	return l.currentSource()[l.loc.Offset+1]
}

func (l *lexer) setCurrentSource(id int) {
	l.currSrc = id
}

func (l *lexer) currentSource() []rune {
	return l.sources[l.currSrc]
}
