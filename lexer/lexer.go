package lexer

import "monkey/token"

// 只能读取ASCLL码，
// TODO：utf-8,emoji
// TODO：input的类型，以及记录行号等问题


type Lexer struct {
	input string 	// 输入的需要进行词法分析的
	position int	// 当前字符的位置
	readPosition int	// 当前位置的下一个位置
	ch byte		// 当前字符
}

func New(input string) *Lexer {
	l := &Lexer{
		input: input,
	}
	// 在创建对象的时候就读取一起char
	l.readChar()
	return l
}


// 一个个的去读取input中的字符，把字符写入ch字段中
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}


// 词法分析的主要函数
func (l *Lexer) NextToken() token.Token {
	var t *token.Token
	switch l.ch {
	case '=':
		t = token.NewToken(token.ASSIGN, l.ch)
	case '+':
		t = token.NewToken(token.PLUS, l.ch)
	case '(':
		t = token.NewToken(token.LPAREN, l.ch)
	case ')':
		t = token.NewToken(token.RPAREN, l.ch)
	case '{':
		t = token.NewToken(token.LBRACE, l.ch)
	case '}':
		t = token.NewToken(token.RBRACE, l.ch)
	case ',':
		t = token.NewToken(token.COMMA, l.ch)
	case ';':
		t = token.NewToken(token.SEMICOLON, l.ch)
	case 0:
		//t.Literal = " "
		//t.Type = token.EOF
		t = &token.Token{
			Literal: "",
			Type: token.EOF,
		}
	}
	l.readChar()
	return *t
}