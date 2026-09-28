// Package parser implements the parser for MOJA.
// It converts token sequences produced by the lexer into an AST.
package parser

import (
	"fmt"
	"strconv"

	"mojolang/internal/ast"
	"mojolang/internal/lexer"
)

// Precedence levels for operator parsing
const (
	_ int = iota
	LOWEST
	LOGICAL_OR  // utawa
	LOGICAL_AND // lan
	EQUALS      // == != > < >= <=
	SUM         // + -
	PRODUCT     // * /
	PREFIX      // -X  ora X
	CALL        // f(x)
	INDEX       // arr[i]
)

var precedences = map[lexer.TokenType]int{
	lexer.TOKEN_UTAWA:    LOGICAL_OR,
	lexer.TOKEN_LAN:      LOGICAL_AND,
	lexer.TOKEN_EQ:       EQUALS,
	lexer.TOKEN_NEQ:      EQUALS,
	lexer.TOKEN_GT:       EQUALS,
	lexer.TOKEN_LT:       EQUALS,
	lexer.TOKEN_GTE:      EQUALS,
	lexer.TOKEN_LTE:      EQUALS,
	lexer.TOKEN_PLUS:     SUM,
	lexer.TOKEN_MINUS:    SUM,
	lexer.TOKEN_STAR:     PRODUCT,
	lexer.TOKEN_SLASH:    PRODUCT,
	lexer.TOKEN_LPAREN:   CALL,
	lexer.TOKEN_LBRACKET: INDEX,
}

// ParseError represents a syntax or parse error.
type ParseError struct {
	Line    int
	Column  int
	Message string
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("[MOJA Error] Baris %d, Kolom %d: %s", e.Line, e.Column, e.Message)
}

func (e *ParseError) GetPosition() (int, int) {
	return e.Line, e.Column
}

func (e *ParseError) GetErrorType() string {
	return "Syntax Error"
}

// Parser holds the parser state.
type Parser struct {
	tokens  []lexer.Token
	current int
}

// New creates a new Parser with the provided tokens.
func New(tokens []lexer.Token) *Parser {
	return &Parser{
		tokens:  tokens,
		current: 0,
	}
}

// ParseProgram parses the entire token stream into an ast.Program.
func (p *Parser) ParseProgram() (*ast.Program, error) {
	program := &ast.Program{Statements: []ast.Statement{}}

	for !p.isAtEnd() {
		// Ignore extra newlines between statements
		if p.match(lexer.TOKEN_NEWLINE) {
			continue
		}
		if p.peek().Type == lexer.TOKEN_EOF {
			break
		}

		stmt, err := p.parseStatement()
		if err != nil {
			return nil, err
		}
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
	}

	return program, nil
}

// parseStatement parses a single statement.
func (p *Parser) parseStatement() (ast.Statement, error) {
	tok := p.peek()

	switch tok.Type {
	case lexer.TOKEN_OMONG:
		return p.parseOmongStatement()
	case lexer.TOKEN_YEN:
		return p.parseYenStatement()
	case lexer.TOKEN_BALENI:
		return p.parseBaleniStatement()
	case lexer.TOKEN_GAWE:
		return p.parseGaweStatement()
	case lexer.TOKEN_BALIK:
		return p.parseReturnStatement()
	case lexer.TOKEN_IDENTIFIER:
		// Check if it is an assignment: IDENTIFIER = expr
		if p.peekNext().Type == lexer.TOKEN_ASSIGN {
			return p.parseAssignStatement()
		}
		return p.parseExpressionStatement()
	default:
		return p.parseExpressionStatement()
	}
}

// parseOmongStatement: omong(expr)
func (p *Parser) parseOmongStatement() (ast.Statement, error) {
	tok := p.advance() // consume 'omong'

	if !p.match(lexer.TOKEN_LPAREN) {
		nextTok := p.peek()
		return nil, &ParseError{Line: nextTok.Line, Column: nextTok.Column, Message: "Sintaks 'omong' harus diikuti '('"}
	}

	expr, err := p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}

	if !p.match(lexer.TOKEN_RPAREN) {
		nextTok := p.peek()
		return nil, &ParseError{Line: nextTok.Line, Column: nextTok.Column, Message: "Kurung tutup ')' tidak ditemukan setelah ekspresi 'omong'"}
	}

	// Consume optional trailing newline
	p.match(lexer.TOKEN_NEWLINE)

	return &ast.OmongStatement{
		Line:   tok.Line,
		Column: tok.Column,
		Value:  expr,
	}, nil
}

// parseAssignStatement: IDENTIFIER = expression
func (p *Parser) parseAssignStatement() (ast.Statement, error) {
	idTok := p.advance() // consume IDENTIFIER
	p.advance()          // consume '='

	val, err := p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}

	p.match(lexer.TOKEN_NEWLINE)

	return &ast.AssignStatement{
		Line:   idTok.Line,
		Column: idTok.Column,
		Name: &ast.Identifier{
			Line:   idTok.Line,
			Column: idTok.Column,
			Value:  idTok.Literal,
		},
		Value: val,
	}, nil
}

// parseYenStatement: yen expression block (liyane block)?
func (p *Parser) parseYenStatement() (ast.Statement, error) {
	tok := p.advance() // consume 'yen'

	cond, err := p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}

	consequence, err := p.parseBlockStatement()
	if err != nil {
		return nil, err
	}

	var alternative *ast.BlockStatement
	p.skipNewlines()

	if p.check(lexer.TOKEN_LIYANE) {
		p.advance() // consume 'liyane'
		alternative, err = p.parseBlockStatement()
		if err != nil {
			return nil, err
		}
	}

	p.match(lexer.TOKEN_NEWLINE)

	return &ast.YenStatement{
		Line:        tok.Line,
		Column:      tok.Column,
		Condition:   cond,
		Consequence: consequence,
		Alternative: alternative,
	}, nil
}

// parseBaleniStatement: baleni IDENTIFIER saka expr_start tekan expr_end block
func (p *Parser) parseBaleniStatement() (ast.Statement, error) {
	tok := p.advance() // consume 'baleni'

	if !p.check(lexer.TOKEN_IDENTIFIER) {
		nextTok := p.peek()
		return nil, &ParseError{Line: nextTok.Line, Column: nextTok.Column, Message: "Sintaks 'baleni' harus diikuti variabel iterator"}
	}
	idTok := p.advance()

	if !p.match(lexer.TOKEN_SAKA) {
		nextTok := p.peek()
		return nil, &ParseError{Line: nextTok.Line, Column: nextTok.Column, Message: "Sintaks 'baleni' membutuhkan kata kunci 'saka'"}
	}

	startExpr, err := p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}

	if !p.match(lexer.TOKEN_TEKAN) {
		nextTok := p.peek()
		return nil, &ParseError{Line: nextTok.Line, Column: nextTok.Column, Message: "Sintaks 'baleni' membutuhkan kata kunci 'tekan'"}
	}

	endExpr, err := p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}

	body, err := p.parseBlockStatement()
	if err != nil {
		return nil, err
	}

	p.match(lexer.TOKEN_NEWLINE)

	return &ast.BaleniStatement{
		Line:   tok.Line,
		Column: tok.Column,
		Iterator: &ast.Identifier{
			Line:   idTok.Line,
			Column: idTok.Column,
			Value:  idTok.Literal,
		},
		Start: startExpr,
		End:   endExpr,
		Body:  body,
	}, nil
}

// parseGaweStatement: gawe IDENTIFIER(params) block
func (p *Parser) parseGaweStatement() (ast.Statement, error) {
	tok := p.advance() // consume 'gawe'

	if !p.check(lexer.TOKEN_IDENTIFIER) {
		nextTok := p.peek()
		return nil, &ParseError{Line: nextTok.Line, Column: nextTok.Column, Message: "Deklarasi fungsi 'gawe' membutuhkan nama fungsi"}
	}
	nameTok := p.advance()

	if !p.match(lexer.TOKEN_LPAREN) {
		nextTok := p.peek()
		return nil, &ParseError{Line: nextTok.Line, Column: nextTok.Column, Message: "Deklarasi fungsi membutuhkan '(' setelah nama fungsi"}
	}

	var params []*ast.Identifier
	if !p.check(lexer.TOKEN_RPAREN) {
		for {
			if !p.check(lexer.TOKEN_IDENTIFIER) {
				nextTok := p.peek()
				return nil, &ParseError{Line: nextTok.Line, Column: nextTok.Column, Message: "Parameter fungsi harus berupa identifier"}
			}
			paramTok := p.advance()
			params = append(params, &ast.Identifier{
				Line:   paramTok.Line,
				Column: paramTok.Column,
				Value:  paramTok.Literal,
			})

			if p.check(lexer.TOKEN_COMMA) {
				p.advance()
			} else {
				break
			}
		}
	}

	if !p.match(lexer.TOKEN_RPAREN) {
		nextTok := p.peek()
		return nil, &ParseError{Line: nextTok.Line, Column: nextTok.Column, Message: "Kurung tutup ')' tidak ditemukan pada deklarasi fungsi"}
	}

	body, err := p.parseBlockStatement()
	if err != nil {
		return nil, err
	}

	p.match(lexer.TOKEN_NEWLINE)

	return &ast.GaweStatement{
		Line:   tok.Line,
		Column: tok.Column,
		Name: &ast.Identifier{
			Line:   nameTok.Line,
			Column: nameTok.Column,
			Value:  nameTok.Literal,
		},
		Params: params,
		Body:   body,
	}, nil
}

// parseReturnStatement: balik expression NEWLINE
func (p *Parser) parseReturnStatement() (ast.Statement, error) {
	tok := p.advance() // consume 'balik'

	if p.check(lexer.TOKEN_NEWLINE) || p.check(lexer.TOKEN_EOF) || p.check(lexer.TOKEN_RBRACE) {
		return nil, &ParseError{Line: tok.Line, Column: tok.Column, Message: "'balik' harus diikuti ekspresi"}
	}

	val, err := p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}

	p.match(lexer.TOKEN_NEWLINE)

	return &ast.ReturnStatement{
		Line:   tok.Line,
		Column: tok.Column,
		Value:  val,
	}, nil
}

// parseExpressionStatement parses a standalone expression statement or assignment statement.
func (p *Parser) parseExpressionStatement() (ast.Statement, error) {
	tok := p.peek()
	expr, err := p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}

	if p.check(lexer.TOKEN_ASSIGN) {
		assignTok := p.advance() // consume '='
		val, err := p.parseExpression(LOWEST)
		if err != nil {
			return nil, err
		}
		p.match(lexer.TOKEN_NEWLINE)

		if id, ok := expr.(*ast.Identifier); ok {
			return &ast.AssignStatement{
				Line:   id.Line,
				Column: id.Column,
				Name:   id,
				Value:  val,
			}, nil
		}

		if idxExpr, ok := expr.(*ast.IndexExpression); ok {
			return &ast.IndexAssignStatement{
				Line:   idxExpr.Line,
				Column: idxExpr.Column,
				Left:   idxExpr,
				Value:  val,
			}, nil
		}

		return nil, &ParseError{
			Line:    assignTok.Line,
			Column:  assignTok.Column,
			Message: "Target penugasan '=' tidak valid",
		}
	}

	p.match(lexer.TOKEN_NEWLINE)

	return &ast.ExpressionStatement{
		Line:       tok.Line,
		Column:     tok.Column,
		Expression: expr,
	}, nil
}

// parseBlockStatement: { NEWLINE statement* }
func (p *Parser) parseBlockStatement() (*ast.BlockStatement, error) {
	p.skipNewlines()

	if !p.check(lexer.TOKEN_LBRACE) {
		nextTok := p.peek()
		return nil, &ParseError{Line: nextTok.Line, Column: nextTok.Column, Message: "Blok kode harus diawali '{'"}
	}
	tok := p.advance() // consume '{'

	block := &ast.BlockStatement{
		Line:       tok.Line,
		Column:     tok.Column,
		Statements: []ast.Statement{},
	}

	for !p.check(lexer.TOKEN_RBRACE) && !p.isAtEnd() {
		if p.match(lexer.TOKEN_NEWLINE) {
			continue
		}
		stmt, err := p.parseStatement()
		if err != nil {
			return nil, err
		}
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
	}

	if !p.match(lexer.TOKEN_RBRACE) {
		nextTok := p.peek()
		return nil, &ParseError{Line: nextTok.Line, Column: nextTok.Column, Message: "Blok kode tidak ditutup dengan '}'"}
	}

	return block, nil
}

// parseExpression parses expressions using Pratt parsing methodology.
func (p *Parser) parseExpression(precedence int) (ast.Expression, error) {
	tok := p.peek()

	var left ast.Expression
	var err error

	switch tok.Type {
	case lexer.TOKEN_IDENTIFIER:
		p.advance()
		if p.check(lexer.TOKEN_LPAREN) {
			left, err = p.parseCallExpression(&ast.Identifier{Line: tok.Line, Column: tok.Column, Value: tok.Literal})
		} else {
			left = &ast.Identifier{Line: tok.Line, Column: tok.Column, Value: tok.Literal}
		}
	case lexer.TOKEN_NUMBER:
		p.advance()
		val, parseErr := strconv.ParseFloat(tok.Literal, 64)
		if parseErr != nil {
			return nil, &ParseError{Line: tok.Line, Column: tok.Column, Message: fmt.Sprintf("Format angka tidak valid: %s", tok.Literal)}
		}
		left = &ast.NumberLiteral{Line: tok.Line, Column: tok.Column, Value: val, Raw: tok.Literal}
	case lexer.TOKEN_STRING:
		p.advance()
		left = &ast.StringLiteral{Line: tok.Line, Column: tok.Column, Value: tok.Literal}
	case lexer.TOKEN_BENER:
		p.advance()
		left = &ast.BooleanLiteral{Line: tok.Line, Column: tok.Column, Value: true}
	case lexer.TOKEN_SALAH:
		p.advance()
		left = &ast.BooleanLiteral{Line: tok.Line, Column: tok.Column, Value: false}
	case lexer.TOKEN_MINUS:
		p.advance()
		right, subErr := p.parseExpression(PREFIX)
		if subErr != nil {
			return nil, subErr
		}
		left = &ast.PrefixExpression{Line: tok.Line, Column: tok.Column, Operator: "-", Right: right}
	case lexer.TOKEN_ORA:
		p.advance()
		right, subErr := p.parseExpression(PREFIX)
		if subErr != nil {
			return nil, subErr
		}
		left = &ast.PrefixExpression{Line: tok.Line, Column: tok.Column, Operator: "ora", Right: right}
	case lexer.TOKEN_LPAREN:
		p.advance()
		left, err = p.parseExpression(LOWEST)
		if err != nil {
			return nil, err
		}
		if !p.match(lexer.TOKEN_RPAREN) {
			nextTok := p.peek()
			return nil, &ParseError{Line: nextTok.Line, Column: nextTok.Column, Message: "Kurung tutup ')' tidak ditemukan"}
		}
	case lexer.TOKEN_LBRACKET:
		left, err = p.parseArrayLiteral()
	default:
		return nil, &ParseError{Line: tok.Line, Column: tok.Column, Message: fmt.Sprintf("Ekspresi tidak valid: '%s'", tok.Literal)}
	}

	if err != nil {
		return nil, err
	}

	for !p.isAtEnd() && precedence < p.peekPrecedence() {
		opTok := p.peek()
		if opTok.Type == lexer.TOKEN_NEWLINE || opTok.Type == lexer.TOKEN_EOF {
			break
		}

		if opTok.Type == lexer.TOKEN_LBRACKET {
			left, err = p.parseIndexExpression(left)
			if err != nil {
				return nil, err
			}
			continue
		}

		p.advance()
		right, infixErr := p.parseExpression(precedences[opTok.Type])
		if infixErr != nil {
			return nil, infixErr
		}

		left = &ast.InfixExpression{
			Line:     opTok.Line,
			Column:   opTok.Column,
			Left:     left,
			Operator: opTok.Literal,
			Right:    right,
		}
	}

	return left, nil
}

func (p *Parser) parseArrayLiteral() (ast.Expression, error) {
	tok := p.advance() // consume '['

	var elements []ast.Expression
	p.skipNewlines()

	if !p.check(lexer.TOKEN_RBRACKET) {
		for {
			p.skipNewlines()
			el, err := p.parseExpression(LOWEST)
			if err != nil {
				return nil, err
			}
			elements = append(elements, el)

			p.skipNewlines()
			if p.check(lexer.TOKEN_COMMA) {
				p.advance()
			} else {
				break
			}
		}
	}

	p.skipNewlines()
	if !p.match(lexer.TOKEN_RBRACKET) {
		nextTok := p.peek()
		return nil, &ParseError{Line: nextTok.Line, Column: nextTok.Column, Message: "Kurung siku tutup ']' tidak ditemukan pada array"}
	}

	return &ast.ArrayLiteral{
		Line:     tok.Line,
		Column:   tok.Column,
		Elements: elements,
	}, nil
}

func (p *Parser) parseIndexExpression(left ast.Expression) (ast.Expression, error) {
	tok := p.advance() // consume '['
	p.skipNewlines()

	idx, err := p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}

	p.skipNewlines()
	if !p.match(lexer.TOKEN_RBRACKET) {
		nextTok := p.peek()
		return nil, &ParseError{Line: nextTok.Line, Column: nextTok.Column, Message: "Kurung siku tutup ']' tidak ditemukan pada akses indeks"}
	}

	return &ast.IndexExpression{
		Line:   tok.Line,
		Column: tok.Column,
		Left:   left,
		Index:  idx,
	}, nil
}

// parseCallExpression parses function call expression f(a, b)
func (p *Parser) parseCallExpression(fn *ast.Identifier) (ast.Expression, error) {
	p.advance() // consume '('

	var args []ast.Expression
	if !p.check(lexer.TOKEN_RPAREN) {
		for {
			arg, err := p.parseExpression(LOWEST)
			if err != nil {
				return nil, err
			}
			args = append(args, arg)

			if p.check(lexer.TOKEN_COMMA) {
				p.advance()
			} else {
				break
			}
		}
	}

	if !p.match(lexer.TOKEN_RPAREN) {
		nextTok := p.peek()
		return nil, &ParseError{Line: nextTok.Line, Column: nextTok.Column, Message: "Kurung tutup ')' tidak ditemukan pada pemanggilan fungsi"}
	}

	return &ast.CallExpression{
		Line:      fn.Line,
		Column:    fn.Column,
		Function:  fn,
		Arguments: args,
	}, nil
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peek().Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) check(t lexer.TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Type == t
}

func (p *Parser) match(t lexer.TokenType) bool {
	if p.check(t) {
		p.advance()
		return true
	}
	return false
}

func (p *Parser) advance() lexer.Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.tokens[p.current-1]
}

func (p *Parser) peek() lexer.Token {
	if p.current >= len(p.tokens) {
		lastCol := 1
		lastLine := 1
		if len(p.tokens) > 0 {
			lastTok := p.tokens[len(p.tokens)-1]
			lastLine = lastTok.Line
			lastCol = lastTok.Column + len(lastTok.Literal)
		}
		return lexer.Token{Type: lexer.TOKEN_EOF, Line: lastLine, Column: lastCol}
	}
	return p.tokens[p.current]
}

func (p *Parser) peekNext() lexer.Token {
	if p.current+1 >= len(p.tokens) {
		return lexer.Token{Type: lexer.TOKEN_EOF, Line: p.peek().Line, Column: p.peek().Column}
	}
	return p.tokens[p.current+1]
}

func (p *Parser) isAtEnd() bool {
	if p.current >= len(p.tokens) {
		return true
	}
	return p.tokens[p.current].Type == lexer.TOKEN_EOF
}

func (p *Parser) skipNewlines() {
	for p.check(lexer.TOKEN_NEWLINE) {
		p.advance()
	}
}
