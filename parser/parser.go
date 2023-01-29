package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
	"strconv"
)

type Parser struct {
	l            *lexer.Lexer
	currentToken token.Token
	peekToken    token.Token
	errors       []string
	// 用来检查遇到词法单元的时候，使用哪个解析函数
	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

// 语言的优先级
const (
	_ int = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X
	CALL        // myFunction(X)
)

// 优先级对应表
var precedences = map[token.TokenType]int{
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
	token.LPAREN:   CALL,
}

// 优先级辅助函数
func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}
func (p *Parser) currentPrecedence() int {
	if p, ok := precedences[p.currentToken.Type]; ok {
		return p
	}
	return LOWEST
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}
	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	p.registerPrefix(token.LPAREN, p.parseGroupExpression)
	p.registerPrefix(token.IF, p.parseIfExpression)
	p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral)
	p.registerPrefix(token.STRING, p.parseStringLiteral)

	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	p.registerInfix(token.LPAREN, p.parseCallExpression)

	p.nextToken()
	p.nextToken()
	return p
}

// 错误处理
func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

// 和 lexer 中的nextToken相似
func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

// 主要方法
func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}
	for !p.currentTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	return program
}

// 解析语句
func (p *Parser) parseStatement() ast.Statement {
	switch p.currentToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()

	default:
		return p.parseExpressionStatement()
	}
}

// 解析letStatement
func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.currentToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	p.nextToken()

	stmt.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// 解析 returnStatement
func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.currentToken}

	p.nextToken()

	stmt.ReturnValue = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// 解析 expressionStatement
func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	// untrace(trace("parseExpressionStatement"))
	stmt := &ast.ExpressionStatement{
		Token: p.currentToken,
	}
	stmt.Expression = p.parseExpression(LOWEST)
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

// 解析表达式的主要函数
func (p *Parser) parseExpression(precedence int) ast.Expression {
	// untrace(trace("parseExpression"))
	prefix := p.prefixParseFns[p.currentToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.currentToken.Type)
		return nil
	}
	leftExp := prefix()
	for !p.currentTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}
		p.nextToken()
		leftExp = infix(leftExp)
	}
	return leftExp
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

// 工具函数，用于比较currentToken Type是否与传入的tokenType相同
func (p *Parser) currentTokenIs(t token.TokenType) bool {
	return p.currentToken.Type == t
}

// 同上
func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

// 查看下一个token是否是期待的，用在parser中
func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}

}

// 前缀解析函数 和 中缀解析函数
type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

// 辅助函数，用于注册对应的解析函数
func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}
func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

/*
	解析函数：
	函数在开始解析表达式时，当前 curToken 必须
	是所关联的词法单元类型，返回分析的表达式结果时，curToken 是当前表达式类型
	中的最后一个词法单元。切勿将词法单元前移得太远。
*/
// 解析Identifier
func (p *Parser) parseIdentifier() ast.Expression {
	// untrace(trace("parseIdentifier"))
	return &ast.Identifier{
		Token: p.currentToken,
		Value: p.currentToken.Literal,
	}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	// untrace(trace("parseIntegerLiteral"))
	il := &ast.IntegerLiteral{
		Token: p.currentToken,
	}
	value, err := strconv.ParseInt(p.currentToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.currentToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}
	il.Value = value
	return il
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	// untrace(trace("parsePrefixExpression"))
	pe := &ast.PrefixExpression{
		Token:    p.currentToken,
		Operator: p.currentToken.Literal,
	}
	p.nextToken()
	pe.Right = p.parseExpression(PREFIX)
	return pe
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	// untrace(trace("parseInfixExpression"))
	expression := &ast.InfixExpression{
		Token:    p.currentToken,
		Operator: p.currentToken.Literal,
		Left:     left,
	}
	precedence := p.currentPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)
	return expression
}

func (p *Parser) parseBoolean() ast.Expression {
	// untrace(trace("parseBoolean"))
	return &ast.Boolean{
		Token: p.currentToken,
		Value: p.currentTokenIs(token.TRUE),
	}
}

func (p *Parser) parseGroupExpression() ast.Expression {
	p.nextToken()
	exp := p.parseExpression(LOWEST)
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return exp
}

func (p *Parser) parseIfExpression() ast.Expression {
	ie := &ast.IfExpression{
		Token: p.currentToken,
	}
	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	p.nextToken()
	ie.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	if !p.expectPeek(token.LBRACE) {
		return nil
	}
	ie.Consequence = p.parseBlockStatement()
	if p.peekTokenIs(token.ELSE) {
		p.nextToken()
		if !p.expectPeek(token.LBRACE) {
			return nil
		}
		ie.Alternative = p.parseBlockStatement()
	}
	return ie
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{
		Token: p.currentToken,
	}
	block.Statements = []ast.Statement{}
	p.nextToken()
	for !p.currentTokenIs(token.RBRACE) && !p.currentTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}
	return block
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	fl := &ast.FunctionLiteral{
		Token: p.currentToken,
	}
	if !p.expectPeek(token.LPAREN) {
		return nil
	}
	fl.Parameters = p.parseFunctionParameters()
	if !p.expectPeek(token.LBRACE) {
		return nil
	}
	fl.Body = p.parseBlockStatement()
	return fl
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	identifiers := []*ast.Identifier{}
	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return identifiers
	}

	p.nextToken()

	ident := &ast.Identifier{
		Token: p.currentToken,
		Value: p.currentToken.Literal,
	}
	identifiers = append(identifiers, ident)

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		ident := &ast.Identifier{
			Token: p.currentToken,
			Value: p.currentToken.Literal,
		}
		identifiers = append(identifiers, ident)
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return identifiers
}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	ce := &ast.CallExpression{
		Token:    p.currentToken,
		Function: function,
	}
	ce.Arguments = p.parseCallArguments()
	return ce
}

func (p *Parser) parseCallArguments() []ast.Expression {
	args := []ast.Expression{}
	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return args
	}
	p.nextToken()
	args = append(args, p.parseExpression(LOWEST))
	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		args = append(args, p.parseExpression(LOWEST))
	}
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return args
}


func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{
		Token: p.currentToken,
		Value: p.currentToken.Literal,
	}
}

