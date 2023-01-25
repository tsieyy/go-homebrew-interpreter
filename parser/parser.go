package parser

import (
	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
	"fmt"
	"strconv"
)


type Parser struct {
	l *lexer.Lexer
	currentToken token.Token
	peekToken token.Token
	errors []string
	// 用来检查遇到词法单元的时候，使用哪个解析函数
	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns map[token.TokenType]infixParseFn

}

// 语言的优先级
const ( 
	_ int = iota 
	LOWEST 
	EQUALS // == 
	LESSGREATER // > or < 
	SUM // + 
	PRODUCT // * 
	PREFIX // -X or !X 
	CALL // myFunction(X) 
)
// 优先级对应表	
var precedences = map[token.TokenType]int { 
	token.EQ: EQUALS, 
	token.NOT_EQ: EQUALS, 
	token.LT: LESSGREATER, 
	token.GT: LESSGREATER, 
	token.PLUS: SUM, 
	token.MINUS: SUM, 
	token.SLASH: PRODUCT, 
	token.ASTERISK: PRODUCT, 
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
		l: l,
		errors: []string{},
	}
	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)

	p.infixParseFns = make(map[token.TokenType]infixParseFn) 
	p.registerInfix(token.PLUS, p.parseInfixExpression) 
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression) 
	p.registerInfix(token.ASTERISK, p.parseInfixExpression) 
	p.registerInfix(token.EQ, p.parseInfixExpression) 
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression) 
	p.registerInfix(token.LT, p.parseInfixExpression) 
	p.registerInfix(token.GT, p.parseInfixExpression)

	p.nextToken()
	p.nextToken()
	return p
}


// 错误处理
func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t token.TokenType) { 
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", 	t, p.peekToken.Type) 
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
	for p.currentToken.Type != token.EOF {
		if stmt := p.parseStatement(); stmt != nil {
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
func (p *Parser) parseLetStatement() ast.Statement {
	stmt := ast.LetStatement{
		Token: p.currentToken,
	}
	if !p.expectPeek(token.IDENT) {
		//p.peekError(p.currentToken.Type)
		return nil
	}
	stmt.Name = &ast.Identifier{
		Token: p.currentToken,
		Value: p.currentToken.Literal,
	}
	if !p.expectPeek(token.ASSIGN) {
		return nil
	}
	// TODO：表达式的处理
	for !p.currentTokenIs(token.SEMICOLON) { 
		// TODO: 需要在这里添加对表达式的处理
		p.nextToken() 
	}
	return &stmt
}


// 解析 returnStatement
func (p *Parser) parseReturnStatement() ast.Statement {
	stmt := &ast.ReturnStatement{
		Token: p.currentToken,
	}
	p.nextToken()

	// TODO:先跳过表达式的处理
	for !p.currentTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}


// 解析 expressionStatement
func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
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



//工具函数，用于比较currentToken Type是否与传入的tokenType相同
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
	infixParseFn func(ast.Expression) ast.Expression
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
	return &ast.Identifier{
		Token: p.currentToken,
		Value: p.currentToken.Literal,
	}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
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
	pe := &ast.PrefixExpression{
		Token: p.currentToken,
		Operator: p.currentToken.Literal,
	}
	p.nextToken()
	pe.Right = p.parseExpression(PREFIX)
	return pe
}


func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{ 
		Token: p.currentToken, 
		Operator: p.currentToken.Literal, 
		Left: left, 
	}
	precedence := p.currentPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)
	return expression
}