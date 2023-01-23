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





// 词法分析的主要函数
func (l *Lexer) NextToken() token.Token {
	var t *token.Token
	l.skipWhitespace()
	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			t = &token.Token{
				Literal: string(ch) + string(l.ch),
				Type: token.EQ,
			}
		} else {
			t = token.NewToken(token.ASSIGN, l.ch)
		}
	case '+':
		t = token.NewToken(token.PLUS, l.ch)
	case '-':
		t = token.NewToken(token.MINUS, l.ch)
	case '*':
		t = token.NewToken(token.ASTERISK, l.ch)
	case '/':
		t = token.NewToken(token.SLASH, l.ch)
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			t = &token.Token{
				Literal: string(ch) + string(l.ch),
				Type: token.NOT_EQ,
			}
		} else {
			t = token.NewToken(token.BANG, l.ch)
		}
	case '<':
		t = token.NewToken(token.LT, l.ch)
	case '>':
		t = token.NewToken(token.GT, l.ch)
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
	default:
		if isLetter(l.ch) {
			t = &token.Token{
				Literal: l.readIdentifier(),
				//Type: token.LookupIdent(t.Literal),
			}
			t.Type = token.LookupIdent(t.Literal)
			return *t
		} else if isNumber(l.ch) {
			t = &token.Token{
				Type: token.INT,
				Literal: l.readNumber(),
			}
			return *t
		} else {
			t = token.NewToken(token.ILLEGAL, l.ch)
		}
	}
	l.readChar()
	return *t
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


// 窥探下一个字符，来判断是否是 == ！= 这种；
func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} 
	return l.input[l.readPosition]
}



// 读取标识符/关键字，后续还需要对其进行判断，但这里只需要读一下
func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}


// 读取数字，同上方函数相似，把数字转换为token
func (l *Lexer) readNumber() string {
	position := l.position
	for isNumber(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}



// 当读取到的是空格时，跳过
func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' { 
		l.readChar() 
	}
}


// 一个工具方法，判断byte是否为字母
func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

// 工具方法，判断byte是否为数字
func isNumber(ch byte) bool {
	return '0' <= ch && ch <= '9'
}