package parser

import (
	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
	"fmt"
)


type Parser struct {
	l *lexer.Lexer
	currentToken token.Token
	peekToken token.Token
	errors []string
}


func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l: l,
		errors: []string{},
	}
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
		



	default:
		return nil
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
	for !p.currentTokenIs(token.SEMICOLON) { 
		// TODO: 需要在这里添加对表达式的处理
		p.nextToken() 
	}
	return &stmt
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