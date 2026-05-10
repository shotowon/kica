package parser

import (
	"fmt"
	"slices"

	"github.com/shotowon/kica/internal/ast"
	"github.com/shotowon/kica/internal/gears"
	"github.com/shotowon/kica/internal/tokens"
)

type parser struct {
	files    []tokens.File
	currFile int
	loc      gears.Location
	pos      int
}

func New(files []tokens.File) *parser {
	p := new(parser)
	p.files = files
	p.loc = gears.ZeroLoc()
	p.pos = 0

	return p
}

func (p *parser) currentFile() *tokens.File {
	return &p.files[p.currFile]
}

func (p *parser) setCurrentFile(id int) {
	p.currFile = id
}

func (p *parser) Parse() ([]ast.File, error) {
	files := make([]ast.File, len(p.files))
	for i := range p.files {
		p.setCurrentFile(i)

		stmts, err := p.parseFile()
		if err != nil {
			return nil, err
		}

		files[i] = *ast.NewFile(stmts)
	}

	return files, nil
}

func (p *parser) parseFile() ([]ast.Stmt, error) {
	stmts := make([]ast.Stmt, 0)

	for p.curr().Kind != tokens.EOF {
		stmt, err := p.parseStmt()
		if err != nil {
			return nil, err
		}

		stmts = append(stmts, stmt)
	}

	return stmts, nil
}

func (p *parser) parseStmt() (ast.Stmt, error) {
	if p.curr() == nil {
		return nil, fmt.Errorf("unexpected EOF at %s", p.loc)
	}

	switch p.curr().Kind {
	case tokens.Print:
		return p.parsePrintStmt()
	case tokens.Function:
		return p.parseFuncDefStmt()
	case tokens.Var:
		return p.parseVarDefStmt()
	case tokens.Struct:
		return p.parseStructDefStmt()
	case tokens.Return:
		return p.parseReturnStmt()
	}

	return p.parseExprStmt()
}

func (p *parser) parseReturnStmt() (ast.Stmt, error) {
	if err := p.expect(tokens.Return); err != nil {
		return nil, err
	}

	loc := p.curr().Location
	p.advance()

	if err := p.expect(tokens.Semicolon); err == nil {
		loc.End = p.curr().Location.End
		p.advance()
		return &ast.ReturnStmt{
			Value: nil,
			Loc:   loc,
		}, nil
	}

	returnValue, err := p.parseExpr()
	if err != nil {
		return nil, err
	}

	if err := p.expect(tokens.Semicolon); err != nil {
		return nil, err
	}

	loc.End = p.curr().Location.End
	p.advance()

	return &ast.ReturnStmt{
		Value: returnValue,
		Loc:   loc,
	}, nil
}

func (p *parser) parseStructDefStmt() (ast.Stmt, error) {
	if err := p.expect(tokens.Struct); err != nil {
		return nil, err
	}

	loc := p.curr().Location
	p.advance()

	if err := p.expect(tokens.ID); err != nil {
		return nil, err
	}

	structName := p.curr()
	p.advance()
	if err := p.expect(tokens.LBrace); err != nil {
		return nil, err
	}

	p.advance()

	fields := make(map[string]ast.Type)

	for p.expect(tokens.RBrace) != nil {
		if err := p.expect(tokens.ID); err != nil {
			return nil, err
		}
		fieldName := p.curr()
		p.advance()
		if err := p.expect(tokens.Colon); err != nil {
			return nil, err
		}
		p.advance()

		fieldType, err := p.parseTypeLabel()
		if err != nil {
			return nil, err
		}
		if err := p.expect(tokens.Semicolon); err != nil {
			return nil, err
		}

		p.advance()

		fields[fieldName.Literal] = fieldType
	}

	if err := p.expect(tokens.RBrace); err != nil {
		return nil, err
	}

	loc.End = p.curr().Location.End
	p.advance()

	return &ast.StructDefStmt{
		Name:   structName,
		Fields: fields,
		Loc:    loc,
	}, nil
}

func (p *parser) parseFuncDefStmt() (ast.Stmt, error) {
	if err := p.expect(tokens.Function); err != nil {
		return nil, err
	}

	loc := p.curr().Location
	p.advance()
	if err := p.expect(tokens.ID); err != nil {
		return nil, err
	}
	funcName := p.curr()

	p.advance()
	if err := p.expect(tokens.LParen); err != nil {
		return nil, err
	}

	p.advance()
	params := make([]ast.NTMapping, 0)

	for p.expect(tokens.RParen) != nil {
		if err := p.expect(tokens.ID); err != nil {
			return nil, err
		}
		paramName := p.curr()
		p.advance()

		if err := p.expect(tokens.Colon); err != nil {
			return nil, err
		}
		p.advance()

		paramType, err := p.parseTypeLabel()
		if err != nil {
			return nil, err
		}

		params = append(params, ast.NTMapping{
			Name: *paramName,
			Ty:   paramType,
		})

		if p.expect(tokens.Comma) == nil {
			p.advance()
		}
	}

	if err := p.expect(tokens.RParen); err != nil {
		return nil, err
	}
	p.advance()

	var (
		returnType ast.Type
	)

	if err := p.expect(tokens.Minus); err == nil {
		if err2 := p.expect_peekn(tokens.GT, 1); err2 != nil {
			return nil, err
		}

		p.advance()
		p.advance()
		returnType, err = p.parseTypeLabel()
		if err != nil {
			return nil, err
		}
	}

	stmts, err := p.parseBlock()
	if err != nil {
		return nil, err
	}

	return &ast.FuncDefStmt{
		Loc:        loc,
		Name:       funcName,
		Body:       stmts,
		Params:     params,
		ReturnType: returnType,
	}, nil
}

func (p *parser) parseVarDefStmt() (ast.Stmt, error) {
	if err := p.expect(tokens.Var); err != nil {
		return nil, err
	}
	loc := p.curr().Location
	p.advance()
	if err := p.expect(tokens.ID); err != nil {
		return nil, err
	}
	name := p.curr()
	p.advance()

	if err := p.expect(tokens.Eq); err != nil {
		return nil, err
	}
	p.advance()

	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}

	if err := p.expect(tokens.Semicolon); err != nil {
		return nil, err
	}
	loc.End = p.curr().Location.End
	p.advance()

	return &ast.VarDefStmt{
		Name:  name,
		Value: expr,
		Loc:   loc,
	}, nil
}

func (p *parser) parseBlock() ([]ast.Stmt, error) {
	if err := p.expect(tokens.LBrace); err != nil {
		return nil, err
	}
	p.advance()

	stmts := make([]ast.Stmt, 0)
	for p.expect(tokens.RBrace) != nil {
		stmt, err := p.parseStmt()
		if err != nil {
			return nil, err
		}
		stmts = append(stmts, stmt)
	}

	if err := p.expect(tokens.RBrace); err != nil {
		return nil, err
	}
	p.advance()

	return stmts, nil
}

func (p *parser) parsePrintStmt() (ast.Stmt, error) {
	if err := p.expect(tokens.Print); err != nil {
		return nil, err
	}
	loc := p.curr().Location
	p.advance()

	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}

	if err := p.expect(tokens.Semicolon); err != nil {
		return nil, err
	}

	loc.End = p.curr().Location.End
	p.advance()

	return &ast.PrintStmt{
		Value: expr,
		Loc:   loc,
	}, nil
}

func (p *parser) parseExprStmt() (ast.Stmt, error) {
	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}

	loc := expr.Location()
	if err := p.expect(tokens.Semicolon); err != nil {
		return nil, err
	}

	loc.End = p.curr().Location.End
	p.advance()

	return &ast.ExprStmt{
		Value: expr,
		Loc:   loc,
	}, nil
}

func (p *parser) parseExpr() (ast.Expr, error) {
	return p.parseBinary(0)
}

func (p *parser) parseBinary(minBP int) (ast.Expr, error) {
	lhs, err := p.parseUnary()
	if err != nil {
		return nil, err
	}

	for err = p.expect_one_of(
		tokens.Plus,
		tokens.Minus,
		tokens.Mul,
		tokens.Div,
	); err == nil; err = p.expect_one_of(
		tokens.Plus,
		tokens.Minus,
		tokens.Mul,
		tokens.Div,
	) {
		op := p.curr()
		lBP, rBP := infixBP(op.Kind)
		if lBP == 0 {
			return nil, fmt.Errorf("expected operation at %s got %s", op.Location, op.Kind)
		}

		if lBP < minBP {
			break
		}

		p.advance()
		rhs, err := p.parseBinary(rBP)
		if err != nil {
			return nil, err
		}

		lhs = &ast.BinaryExpr{LHS: lhs, RHS: rhs, Op: op}
	}

	return lhs, nil
}

func (p *parser) parseUnary() (ast.Expr, error) {
	curr := p.curr()
	if err := p.expect_one_of(tokens.Plus, tokens.Minus); err != nil {
		return p.parseAtom()
	}

	op := curr
	p.advance()
	expr, err := p.parseAtom()
	if err != nil {
		return nil, err
	}

	return &ast.UnaryExpr{
		Value: expr,
		Op:    op,
	}, nil
}

func (p *parser) parseAtom() (ast.Expr, error) {
	curr := p.curr()
	if p.expect(tokens.ID) == nil {
		if p.expect_peekn(tokens.LParen, 1) == nil {
			return p.parseFuncCall()
		}
	}

	if err := p.expect_one_of(tokens.Int, tokens.Float, tokens.String, tokens.ID); err != nil {
		return nil, err
	}

	p.advance()

	var e ast.AtomExpr
	e.Value = curr
	return &e, nil
}

func (p *parser) parseFuncCall() (ast.Expr, error) {
	if err := p.expect(tokens.ID); err != nil {
		return nil, err
	}

	name := p.curr()

	loc := name.Location
	p.advance()

	if err := p.expect(tokens.LParen); err != nil {
		return nil, err
	}
	p.advance()

	args := make([]ast.Expr, 0)

	for p.expect(tokens.RParen) != nil {
		arg, err := p.parseExpr()
		if err != nil {
			return nil, err
		}

		args = append(args, arg)
		if p.expect(tokens.Comma) == nil {
			p.advance()
		}
	}

	if err := p.expect(tokens.RParen); err != nil {
		return nil, err
	}
	loc.End = p.curr().Location.End
	p.advance()

	return &ast.FunctionCallExpr{
		Name: name,
		Args: args,
		Loc:  loc,
	}, nil
}

func (p *parser) curr() *tokens.Token {
	if p.pos >= len(p.currentFile().Tokens()) {
		return nil
	}

	return p.currentFile().Tokens()[p.pos]
}

func (p *parser) peekn(n int) *tokens.Token {
	if p.pos+n >= len(p.currentFile().Tokens()) {
		return nil
	}

	return p.currentFile().Tokens()[p.pos+n]
}

func (p *parser) advance() {
	p.pos++

	if p.curr() != nil {
		p.loc = p.curr().Location
	}
}

func (p *parser) expect(kind tokens.TokenKind) error {
	curr := p.curr()
	if curr == nil {
		return fmt.Errorf("unexpected EOF after %s", p.loc)
	}

	if curr.Kind != kind {
		return fmt.Errorf("expected '%s' at %s, got '%s'", kind, curr.Location, curr.Kind)
	}

	return nil
}

func (p *parser) parseTypeLabel() (ast.Type, error) {
	curr := p.curr()
	switch curr.Kind {
	case tokens.TLInt:
		p.advance()
		return &ast.IntT{
			Loc: curr.Location,
		}, nil
	case tokens.TLFloat:
		p.advance()
		return &ast.FloatT{
			Loc: curr.Location,
		}, nil
	case tokens.TLString:
		p.advance()
		return &ast.StringT{
			Loc: curr.Location,
		}, nil
	case tokens.ID:
		p.advance()
		return &ast.NamedT{
			Tok: curr,
		}, nil
	}

	return nil, nil
}

func (p *parser) expect_peekn(kind tokens.TokenKind, n int) error {
	peek := p.peekn(n)
	if peek == nil {
		return fmt.Errorf("unexpected EOF after %s", p.loc)
	}

	if peek.Kind != kind {
		return fmt.Errorf("expected '%s' at %s, got '%s'", kind, peek.Location, peek.Kind)
	}

	return nil
}

func (p *parser) expect_one_of(kinds ...tokens.TokenKind) error {
	curr := p.curr()
	if curr == nil {
		return fmt.Errorf("unexpected EOF after %s", p.loc)
	}

	if slices.Contains(kinds, curr.Kind) {
		return nil
	}

	return fmt.Errorf("expected one of [%s] at %s, got %s", kinds, curr.Location, curr.Kind)
}

func infixBP(kind tokens.TokenKind) (int, int) {
	switch kind {
	case tokens.Plus, tokens.Minus:
		return 1, 2
	case tokens.Mul, tokens.Div:
		return 3, 4
	}

	return 0, 0
}
