package lexer_test

import (
	"testing"

	"mojolang/internal/lexer"
)

func mustTokenize(t *testing.T, source string) []lexer.Token {
	t.Helper()
	l := lexer.New(source)
	tokens, err := l.Tokenize()
	if err != nil {
		t.Fatalf("Tokenize gagal: %v", err)
	}
	return tokens
}

func TestLexerTableDriven(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []lexer.Token
	}{
		{
			name:  "Empty Source",
			input: "",
			expected: []lexer.Token{
				{Type: lexer.TOKEN_EOF, Literal: "", Line: 1, Column: 1},
			},
		},
		{
			name:  "Whitespace and Comments Only",
			input: "  \t // komentar \n  ",
			expected: []lexer.Token{
				{Type: lexer.TOKEN_NEWLINE, Literal: "\n", Line: 1, Column: 17},
				{Type: lexer.TOKEN_EOF, Literal: "", Line: 2, Column: 3},
			},
		},
		{
			name:  "Keywords All",
			input: "omong yen liyane baleni saka tekan gawe balik bener salah",
			expected: []lexer.Token{
				{Type: lexer.TOKEN_OMONG, Literal: "omong", Line: 1, Column: 1},
				{Type: lexer.TOKEN_YEN, Literal: "yen", Line: 1, Column: 7},
				{Type: lexer.TOKEN_LIYANE, Literal: "liyane", Line: 1, Column: 11},
				{Type: lexer.TOKEN_BALENI, Literal: "baleni", Line: 1, Column: 18},
				{Type: lexer.TOKEN_SAKA, Literal: "saka", Line: 1, Column: 25},
				{Type: lexer.TOKEN_TEKAN, Literal: "tekan", Line: 1, Column: 30},
				{Type: lexer.TOKEN_GAWE, Literal: "gawe", Line: 1, Column: 36},
				{Type: lexer.TOKEN_BALIK, Literal: "balik", Line: 1, Column: 41},
				{Type: lexer.TOKEN_BENER, Literal: "bener", Line: 1, Column: 47},
				{Type: lexer.TOKEN_SALAH, Literal: "salah", Line: 1, Column: 53},
				{Type: lexer.TOKEN_EOF, Literal: "", Line: 1, Column: 58},
			},
		},
		{
			name:  "Logic Keywords lan utawa ora",
			input: "lan utawa ora",
			expected: []lexer.Token{
				{Type: lexer.TOKEN_LAN, Literal: "lan", Line: 1, Column: 1},
				{Type: lexer.TOKEN_UTAWA, Literal: "utawa", Line: 1, Column: 5},
				{Type: lexer.TOKEN_ORA, Literal: "ora", Line: 1, Column: 11},
				{Type: lexer.TOKEN_EOF, Literal: "", Line: 1, Column: 14},
			},
		},
		{
			name:  "Logic Keywords With Position",
			input: "bener lan salah",
			expected: []lexer.Token{
				{Type: lexer.TOKEN_BENER, Literal: "bener", Line: 1, Column: 1},
				{Type: lexer.TOKEN_LAN, Literal: "lan", Line: 1, Column: 7},
				{Type: lexer.TOKEN_SALAH, Literal: "salah", Line: 1, Column: 11},
				{Type: lexer.TOKEN_EOF, Literal: "", Line: 1, Column: 16},
			},
		},
		{
			name:  "Identifier Similar to Keywords Not Confused",
			input: "lancar orang utawari oratan",
			expected: []lexer.Token{
				{Type: lexer.TOKEN_IDENTIFIER, Literal: "lancar", Line: 1, Column: 1},
				{Type: lexer.TOKEN_IDENTIFIER, Literal: "orang", Line: 1, Column: 8},
				{Type: lexer.TOKEN_IDENTIFIER, Literal: "utawari", Line: 1, Column: 14},
				{Type: lexer.TOKEN_IDENTIFIER, Literal: "oratan", Line: 1, Column: 22},
				{Type: lexer.TOKEN_EOF, Literal: "", Line: 1, Column: 28},
			},
		},
		{
			name:  "Numbers and Floats",
			input: "123 45.67 0",
			expected: []lexer.Token{
				{Type: lexer.TOKEN_NUMBER, Literal: "123", Line: 1, Column: 1},
				{Type: lexer.TOKEN_NUMBER, Literal: "45.67", Line: 1, Column: 5},
				{Type: lexer.TOKEN_NUMBER, Literal: "0", Line: 1, Column: 11},
				{Type: lexer.TOKEN_EOF, Literal: "", Line: 1, Column: 12},
			},
		},
		{
			name:  "String Escapes",
			input: `"hello\nworld\t\"MOJA\"\\"`,
			expected: []lexer.Token{
				{Type: lexer.TOKEN_STRING, Literal: "hello\nworld\t\"MOJA\"\\", Line: 1, Column: 1},
				{Type: lexer.TOKEN_EOF, Literal: "", Line: 1, Column: 27},
			},
		},
		{
			name:  "Operators",
			input: "+ - * / = == != > < >= <=",
			expected: []lexer.Token{
				{Type: lexer.TOKEN_PLUS, Literal: "+", Line: 1, Column: 1},
				{Type: lexer.TOKEN_MINUS, Literal: "-", Line: 1, Column: 3},
				{Type: lexer.TOKEN_STAR, Literal: "*", Line: 1, Column: 5},
				{Type: lexer.TOKEN_SLASH, Literal: "/", Line: 1, Column: 7},
				{Type: lexer.TOKEN_ASSIGN, Literal: "=", Line: 1, Column: 9},
				{Type: lexer.TOKEN_EQ, Literal: "==", Line: 1, Column: 11},
				{Type: lexer.TOKEN_NEQ, Literal: "!=", Line: 1, Column: 14},
				{Type: lexer.TOKEN_GT, Literal: ">", Line: 1, Column: 17},
				{Type: lexer.TOKEN_LT, Literal: "<", Line: 1, Column: 19},
				{Type: lexer.TOKEN_GTE, Literal: ">=", Line: 1, Column: 21},
				{Type: lexer.TOKEN_LTE, Literal: "<=", Line: 1, Column: 24},
				{Type: lexer.TOKEN_EOF, Literal: "", Line: 1, Column: 26},
			},
		},
		{
			name:  "Delimiters and Brackets",
			input: "( ) { } [ ] ,",
			expected: []lexer.Token{
				{Type: lexer.TOKEN_LPAREN, Literal: "(", Line: 1, Column: 1},
				{Type: lexer.TOKEN_RPAREN, Literal: ")", Line: 1, Column: 3},
				{Type: lexer.TOKEN_LBRACE, Literal: "{", Line: 1, Column: 5},
				{Type: lexer.TOKEN_RBRACE, Literal: "}", Line: 1, Column: 7},
				{Type: lexer.TOKEN_LBRACKET, Literal: "[", Line: 1, Column: 9},
				{Type: lexer.TOKEN_RBRACKET, Literal: "]", Line: 1, Column: 11},
				{Type: lexer.TOKEN_COMMA, Literal: ",", Line: 1, Column: 13},
				{Type: lexer.TOKEN_EOF, Literal: "", Line: 1, Column: 14},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens := mustTokenize(t, tt.input)
			if len(tokens) != len(tt.expected) {
				t.Fatalf("Token count mismatch: got %d, want %d\nGot: %v\nWant: %v", len(tokens), len(tt.expected), tokens, tt.expected)
			}
			for i, tok := range tokens {
				exp := tt.expected[i]
				if tok.Type != exp.Type || tok.Literal != exp.Literal || tok.Line != exp.Line || tok.Column != exp.Column {
					t.Errorf("Token[%d] mismatch:\nGot:  %v\nWant: %v", i, tok, exp)
				}
			}
		})
	}
}

func TestLexerErrorCases(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectedMsg string
		expLine     int
		expCol      int
	}{
		{
			name:        "Unclosed String",
			input:       `"string tanpa penutup`,
			expectedMsg: "String tidak ditutup",
			expLine:     1,
			expCol:      1,
		},
		{
			name:        "Unknown Character @",
			input:       "x = @ 10",
			expectedMsg: "Karakter tidak dikenal: '@'",
			expLine:     1,
			expCol:      5,
		},
		{
			name:        "Exclamation Without Equal",
			input:       "yen x ! 10",
			expectedMsg: "Token tidak dikenali: '!'",
			expLine:     1,
			expCol:      7,
		},
		{
			name:        "Invalid Escape Sequence",
			input:       `"test \k"`,
			expectedMsg: "Escape sequence tidak dikenal: '\\k'",
			expLine:     1,
			expCol:      7,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			_, err := l.Tokenize()
			if err == nil {
				t.Fatalf("Expected error containing %q, got nil", tt.expectedMsg)
			}
			if posErr, ok := err.(interface{ GetPosition() (int, int) }); ok {
				line, col := posErr.GetPosition()
				if line != tt.expLine || col != tt.expCol {
					t.Errorf("Position mismatch: got Line %d, Col %d; want Line %d, Col %d", line, col, tt.expLine, tt.expCol)
				}
			}
		})
	}
}
