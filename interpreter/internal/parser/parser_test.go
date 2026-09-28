package parser_test

import (
	"testing"

	"mojolang/internal/ast"
	"mojolang/internal/lexer"
	"mojolang/internal/parser"
)

func parse(input string) (*ast.Program, error) {
	l := lexer.New(input)
	tokens, err := l.Tokenize()
	if err != nil {
		return nil, err
	}
	p := parser.New(tokens)
	return p.ParseProgram()
}

func TestParserValidProgramsTableDriven(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expectedCount int
	}{
		{
			name:          "Empty Input",
			input:         "",
			expectedCount: 0,
		},
		{
			name:          "Whitespace & Newlines Only",
			input:         "\n\n\n",
			expectedCount: 0,
		},
		{
			name:          "Omong Statement",
			input:         `omong("Sugeng rawuh")`,
			expectedCount: 1,
		},
		{
			name:          "Assign Statement",
			input:         `x = 10 + 20 * 30`,
			expectedCount: 1,
		},
		{
			name:          "Parenthesized Expression",
			input:         `x = (10 + 20) * 30`,
			expectedCount: 1,
		},
		{
			name: "Yen Statement With Liyane",
			input: `
yen x >= 10 {
	omong("Besar")
} liyane {
	omong("Kecil")
}
`,
			expectedCount: 1,
		},
		{
			name: "Baleni Loop",
			input: `
baleni i saka 1 tekan 5 {
	omong(i)
}
`,
			expectedCount: 1,
		},
		{
			name: "Gawe Function With Return",
			input: `
gawe tambah(a, b) {
	balik a + b
}
`,
			expectedCount: 1,
		},
		{
			name:          "Array Literal Empty and Multi",
			input:         "arr1 = []\narr2 = [10, 20, 30]\nmatriks = [[1, 2], [3, 4]]",
			expectedCount: 3,
		},
		{
			name:          "Index Access and Mutation",
			input:         "omong(arr[0])\narr[1] = 99\nomong(matriks[0][1])",
			expectedCount: 3,
		},
		// --- Operator Logika v0.4 ---
		{
			name:          "Logic OR utawa",
			input:         "x = bener utawa salah",
			expectedCount: 1,
		},
		{
			name:          "Logic AND lan",
			input:         "x = bener lan salah",
			expectedCount: 1,
		},
		{
			name:          "Logic NOT ora",
			input:         "x = ora bener",
			expectedCount: 1,
		},
		{
			name:          "Logic NOT double ora",
			input:         "x = ora ora bener",
			expectedCount: 1,
		},
		{
			name:          "Logic Precedence ora before lan",
			input:         "x = ora bener lan salah",
			expectedCount: 1,
		},
		{
			name:          "Logic Precedence lan before utawa",
			input:         "x = bener utawa salah lan salah",
			expectedCount: 1,
		},
		{
			name:          "Logic Parentheses Override Precedence",
			input:         "x = (bener utawa salah) lan salah",
			expectedCount: 1,
		},
		{
			name:          "Logic In Yen Condition",
			input:         "yen bener lan bener { omong(\"OK\") }",
			expectedCount: 1,
		},
		{
			name:          "Logic With Omong",
			input:         "omong(bener lan salah)",
			expectedCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prog, err := parse(tt.input)
			if err != nil {
				t.Fatalf("Parse failed unexpectedly: %v", err)
			}
			if len(prog.Statements) != tt.expectedCount {
				t.Errorf("Statement count mismatch: got %d, want %d", len(prog.Statements), tt.expectedCount)
			}
		})
	}
}

func TestParserSyntaxErrorCasesTableDriven(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectedMsg string
		expLine     int
	}{
		{
			name:        "Omong Without Left Paren",
			input:       `omong "Halo"`,
			expectedMsg: "Sintaks 'omong' harus diikuti '('",
			expLine:     1,
		},
		{
			name:        "Omong Without Right Paren",
			input:       `omong("Halo"`,
			expectedMsg: "Kurung tutup ')' tidak ditemukan setelah ekspresi 'omong'",
			expLine:     1,
		},
		{
			name:        "Baleni Without Iterator",
			input:       `baleni saka 1 tekan 5 {}`,
			expectedMsg: "Sintaks 'baleni' harus diikuti variabel iterator",
			expLine:     1,
		},
		{
			name:        "Baleni Without Saka",
			input:       `baleni i 1 tekan 5 {}`,
			expectedMsg: "Sintaks 'baleni' membutuhkan kata kunci 'saka'",
			expLine:     1,
		},
		{
			name:        "Baleni Without Tekan",
			input:       `baleni i saka 1 5 {}`,
			expectedMsg: "Sintaks 'baleni' membutuhkan kata kunci 'tekan'",
			expLine:     1,
		},
		{
			name:        "Gawe Without Name",
			input:       `gawe (a, b) {}`,
			expectedMsg: "Deklarasi fungsi 'gawe' membutuhkan nama fungsi",
			expLine:     1,
		},
		{
			name:        "Gawe Without Right Paren",
			input:       `gawe tes(a, b {}`,
			expectedMsg: "Kurung tutup ')' tidak ditemukan pada deklarasi fungsi",
			expLine:     1,
		},
		{
			name:        "Return Without Expression",
			input:       "gawe tes() {\n\tbalik\n}",
			expectedMsg: "'balik' harus diikuti ekspresi",
			expLine:     2,
		},
		{
			name:        "Block Without Opening Brace",
			input:       "yen bener omong(\"ok\")",
			expectedMsg: "Blok kode harus diawali '{'",
			expLine:     1,
		},
		{
			name:        "Block Without Closing Brace",
			input:       "yen bener { omong(\"ok\")",
			expectedMsg: "Blok kode tidak ditutup dengan '}'",
			expLine:     1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parse(tt.input)
			if err == nil {
				t.Fatalf("Expected syntax error containing %q, got nil", tt.expectedMsg)
			}
			if posErr, ok := err.(interface{ GetPosition() (int, int) }); ok {
				line, _ := posErr.GetPosition()
				if line != tt.expLine {
					t.Errorf("Line mismatch: got Line %d, want Line %d (error: %v)", line, tt.expLine, err)
				}
			}
		})
	}
}
