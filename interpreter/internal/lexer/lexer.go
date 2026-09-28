// Package lexer implements the tokenizer (lexer) for MOJA.
// It reads source code text and produces a sequence of tokens with line and column positions.
//
// Terminasi pernyataan: newline (\n) dihasilkan sebagai token NEWLINE.
// Komentar // ... diabaikan.
package lexer

import "fmt"

// TokenType merepresentasikan jenis sebuah token.
type TokenType string

const (
	// Literal
	TOKEN_NUMBER     TokenType = "NUMBER"
	TOKEN_STRING     TokenType = "STRING"
	TOKEN_IDENTIFIER TokenType = "IDENTIFIER"

	// Keyword
	TOKEN_OMONG  TokenType = "OMONG"
	TOKEN_YEN    TokenType = "YEN"
	TOKEN_LIYANE TokenType = "LIYANE"
	TOKEN_BALENI TokenType = "BALENI"
	TOKEN_SAKA   TokenType = "SAKA"
	TOKEN_TEKAN  TokenType = "TEKAN"
	TOKEN_GAWE   TokenType = "GAWE"
	TOKEN_BALIK  TokenType = "BALIK"
	TOKEN_BENER  TokenType = "BENER"
	TOKEN_SALAH  TokenType = "SALAH"

	// Operator logika
	TOKEN_LAN   TokenType = "LAN"   // lan  = AND
	TOKEN_UTAWA TokenType = "UTAWA" // utawa = OR
	TOKEN_ORA   TokenType = "ORA"   // ora   = NOT (prefix unary)

	// Operator aritmatika
	TOKEN_PLUS  TokenType = "PLUS"  // +
	TOKEN_MINUS TokenType = "MINUS" // -
	TOKEN_STAR  TokenType = "STAR"  // *
	TOKEN_SLASH TokenType = "SLASH" // /

	// Operator perbandingan
	TOKEN_EQ  TokenType = "EQ"  // ==
	TOKEN_NEQ TokenType = "NEQ" // !=
	TOKEN_GT  TokenType = "GT"  // >
	TOKEN_LT  TokenType = "LT"  // <
	TOKEN_GTE TokenType = "GTE" // >=
	TOKEN_LTE TokenType = "LTE" // <=

	// Assignment
	TOKEN_ASSIGN TokenType = "ASSIGN" // =

	// Delimiter
	TOKEN_LPAREN   TokenType = "LPAREN"   // (
	TOKEN_RPAREN   TokenType = "RPAREN"   // )
	TOKEN_LBRACE   TokenType = "LBRACE"   // {
	TOKEN_RBRACE   TokenType = "RBRACE"   // }
	TOKEN_LBRACKET TokenType = "LBRACKET" // [
	TOKEN_RBRACKET TokenType = "RBRACKET" // ]
	TOKEN_COMMA    TokenType = "COMMA"    // ,

	// Struktur
	TOKEN_NEWLINE TokenType = "NEWLINE"
	TOKEN_EOF     TokenType = "EOF"
)

// keywords memetakan string ke TokenType keyword.
var keywords = map[string]TokenType{
	"omong":  TOKEN_OMONG,
	"yen":    TOKEN_YEN,
	"liyane": TOKEN_LIYANE,
	"baleni": TOKEN_BALENI,
	"saka":   TOKEN_SAKA,
	"tekan":  TOKEN_TEKAN,
	"gawe":   TOKEN_GAWE,
	"balik":  TOKEN_BALIK,
	"bener":  TOKEN_BENER,
	"salah":  TOKEN_SALAH,
	"lan":    TOKEN_LAN,
	"utawa":  TOKEN_UTAWA,
	"ora":    TOKEN_ORA,
}

// Token merepresentasikan satu unit leksikal.
type Token struct {
	Type    TokenType
	Literal string // Nilai teks asli dari source
	Line    int    // Nomor baris (1-based)
	Column  int    // Nomor kolom (1-based)
}

func (t Token) String() string {
	return fmt.Sprintf("Token{%s, %q, baris:%d, kolom:%d}", t.Type, t.Literal, t.Line, t.Column)
}

// LexError merepresentasikan error yang terjadi selama lexing.
type LexError struct {
	Line    int
	Column  int
	Message string
}

func (e *LexError) Error() string {
	return fmt.Sprintf("[MOJA Error] Baris %d, Kolom %d: %s", e.Line, e.Column, e.Message)
}

func (e *LexError) GetPosition() (int, int) {
	return e.Line, e.Column
}

func (e *LexError) GetErrorType() string {
	return "Lexical Error"
}

// Lexer membaca source code dan menghasilkan token.
type Lexer struct {
	source  []rune // Source code sebagai slice rune (mendukung UTF-8)
	current int    // Posisi baca saat ini
	line    int    // Nomor baris saat ini (1-based)
	column  int    // Nomor kolom saat ini (1-based)
}

// New membuat Lexer baru dari source code.
func New(source string) *Lexer {
	return &Lexer{
		source:  []rune(source),
		current: 0,
		line:    1,
		column:  1,
	}
}

// Tokenize membaca seluruh source dan mengembalikan daftar token.
func (l *Lexer) Tokenize() ([]Token, error) {
	var tokens []Token

	for {
		tok, err := l.nextToken()
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, tok)
		if tok.Type == TOKEN_EOF {
			break
		}
	}

	return tokens, nil
}

// nextToken membaca dan mengembalikan token berikutnya.
func (l *Lexer) nextToken() (Token, error) {
	l.skipWhitespaceAndComments()

	startLine := l.line
	startCol := l.column

	if l.isAtEnd() {
		return Token{Type: TOKEN_EOF, Literal: "", Line: startLine, Column: startCol}, nil
	}

	ch := l.advance()

	// Newline sebagai terminator pernyataan
	if ch == '\n' {
		return Token{Type: TOKEN_NEWLINE, Literal: "\n", Line: startLine, Column: startCol}, nil
	}

	// \r\n — skip \r, newline akan diproses di iterasi berikutnya
	if ch == '\r' {
		return l.nextToken()
	}

	// Number
	if isDigit(ch) {
		return l.readNumber(startLine, startCol)
	}

	// Identifier atau keyword
	if isAlpha(ch) {
		return l.readIdentifierOrKeyword(startLine, startCol), nil
	}

	// String
	if ch == '"' {
		return l.readString(startLine, startCol)
	}

	// Operator dan delimiter
	switch ch {
	case '+':
		return Token{Type: TOKEN_PLUS, Literal: "+", Line: startLine, Column: startCol}, nil
	case '-':
		return Token{Type: TOKEN_MINUS, Literal: "-", Line: startLine, Column: startCol}, nil
	case '*':
		return Token{Type: TOKEN_STAR, Literal: "*", Line: startLine, Column: startCol}, nil
	case '/':
		return Token{Type: TOKEN_SLASH, Literal: "/", Line: startLine, Column: startCol}, nil
	case '(':
		return Token{Type: TOKEN_LPAREN, Literal: "(", Line: startLine, Column: startCol}, nil
	case ')':
		return Token{Type: TOKEN_RPAREN, Literal: ")", Line: startLine, Column: startCol}, nil
	case '{':
		return Token{Type: TOKEN_LBRACE, Literal: "{", Line: startLine, Column: startCol}, nil
	case '}':
		return Token{Type: TOKEN_RBRACE, Literal: "}", Line: startLine, Column: startCol}, nil
	case '[':
		return Token{Type: TOKEN_LBRACKET, Literal: "[", Line: startLine, Column: startCol}, nil
	case ']':
		return Token{Type: TOKEN_RBRACKET, Literal: "]", Line: startLine, Column: startCol}, nil
	case ',':
		return Token{Type: TOKEN_COMMA, Literal: ",", Line: startLine, Column: startCol}, nil
	case '=':
		if l.peek() == '=' {
			l.advance()
			return Token{Type: TOKEN_EQ, Literal: "==", Line: startLine, Column: startCol}, nil
		}
		return Token{Type: TOKEN_ASSIGN, Literal: "=", Line: startLine, Column: startCol}, nil
	case '!':
		if l.peek() == '=' {
			l.advance()
			return Token{Type: TOKEN_NEQ, Literal: "!=", Line: startLine, Column: startCol}, nil
		}
		return Token{}, &LexError{Line: startLine, Column: startCol, Message: "Token tidak dikenali: '!'. Maksud Anda '!='?"}
	case '>':
		if l.peek() == '=' {
			l.advance()
			return Token{Type: TOKEN_GTE, Literal: ">=", Line: startLine, Column: startCol}, nil
		}
		return Token{Type: TOKEN_GT, Literal: ">", Line: startLine, Column: startCol}, nil
	case '<':
		if l.peek() == '=' {
			l.advance()
			return Token{Type: TOKEN_LTE, Literal: "<=", Line: startLine, Column: startCol}, nil
		}
		return Token{Type: TOKEN_LT, Literal: "<", Line: startLine, Column: startCol}, nil
	}

	return Token{}, &LexError{
		Line:    startLine,
		Column:  startCol,
		Message: fmt.Sprintf("Karakter tidak dikenal: '%c'", ch),
	}
}

// skipWhitespaceAndComments melewati spasi, tab, dan komentar //.
func (l *Lexer) skipWhitespaceAndComments() {
	for !l.isAtEnd() {
		ch := l.peek()
		switch ch {
		case ' ', '\t':
			l.advance()
		case '/':
			if l.peekNext() == '/' {
				// Komentar — lewati sampai akhir baris
				for !l.isAtEnd() && l.peek() != '\n' {
					l.advance()
				}
			} else {
				return
			}
		default:
			return
		}
	}
}

// readNumber membaca literal angka (integer atau desimal).
func (l *Lexer) readNumber(startLine, startCol int) (Token, error) {
	startPos := l.current - 1

	for !l.isAtEnd() && isDigit(l.peek()) {
		l.advance()
	}

	// Bagian desimal
	if !l.isAtEnd() && l.peek() == '.' && isDigit(l.peekNext()) {
		l.advance() // konsumsi '.'
		for !l.isAtEnd() && isDigit(l.peek()) {
			l.advance()
		}
	}

	lit := string(l.source[startPos:l.current])
	return Token{Type: TOKEN_NUMBER, Literal: lit, Line: startLine, Column: startCol}, nil
}

// readIdentifierOrKeyword membaca identifier dan mengecek apakah itu keyword.
func (l *Lexer) readIdentifierOrKeyword(startLine, startCol int) Token {
	startPos := l.current - 1

	for !l.isAtEnd() && isAlphaNumeric(l.peek()) {
		l.advance()
	}

	lit := string(l.source[startPos:l.current])

	if tokType, ok := keywords[lit]; ok {
		return Token{Type: tokType, Literal: lit, Line: startLine, Column: startCol}
	}
	return Token{Type: TOKEN_IDENTIFIER, Literal: lit, Line: startLine, Column: startCol}
}

// readString membaca literal string dengan tanda petik ganda.
func (l *Lexer) readString(startLine, startCol int) (Token, error) {
	var result []rune

	for {
		if l.isAtEnd() || l.peek() == '\n' {
			return Token{}, &LexError{
				Line:    startLine,
				Column:  startCol,
				Message: "String tidak ditutup, tanda petik penutup '\"' tidak ditemukan",
			}
		}

		ch := l.advance()

		if ch == '"' {
			break
		}

		if ch == '\\' {
			if l.isAtEnd() {
				return Token{}, &LexError{
					Line:    l.line,
					Column:  l.column,
					Message: "Escape sequence tidak lengkap di akhir string",
				}
			}
			escLine := l.line
			escCol := l.column
			esc := l.advance()
			switch esc {
			case '"':
				result = append(result, '"')
			case '\\':
				result = append(result, '\\')
			case 'n':
				result = append(result, '\n')
			case 't':
				result = append(result, '\t')
			default:
				return Token{}, &LexError{
					Line:    escLine,
					Column:  escCol - 1,
					Message: fmt.Sprintf("Escape sequence tidak dikenal: '\\%c'", esc),
				}
			}
			continue
		}

		result = append(result, ch)
	}

	return Token{Type: TOKEN_STRING, Literal: string(result), Line: startLine, Column: startCol}, nil
}

// advance membaca karakter saat ini dan memajukan posisi & kolom.
func (l *Lexer) advance() rune {
	ch := l.source[l.current]
	l.current++
	if ch == '\n' {
		l.line++
		l.column = 1
	} else {
		l.column++
	}
	return ch
}

// peek mengembalikan karakter saat ini tanpa memajukan posisi.
func (l *Lexer) peek() rune {
	if l.isAtEnd() {
		return 0
	}
	return l.source[l.current]
}

// peekNext mengembalikan karakter berikutnya (lookahead 2).
func (l *Lexer) peekNext() rune {
	if l.current+1 >= len(l.source) {
		return 0
	}
	return l.source[l.current+1]
}

// isAtEnd mengecek apakah posisi sudah di akhir source.
func (l *Lexer) isAtEnd() bool {
	return l.current >= len(l.source)
}

// isDigit mengecek apakah rune adalah digit ASCII.
func isDigit(ch rune) bool {
	return ch >= '0' && ch <= '9'
}

// isAlpha mengecek apakah rune adalah huruf atau underscore.
func isAlpha(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || ch == '_'
}

// isAlphaNumeric mengecek apakah rune adalah huruf, digit, atau underscore.
func isAlphaNumeric(ch rune) bool {
	return isAlpha(ch) || isDigit(ch)
}
