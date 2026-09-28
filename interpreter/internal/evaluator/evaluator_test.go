package evaluator_test

import (
	"bytes"
	"strings"
	"testing"

	"mojolang/internal/evaluator"
	"mojolang/internal/lexer"
	"mojolang/internal/parser"
)

func run(input string) (string, error) {
	l := lexer.New(input)
	tokens, err := l.Tokenize()
	if err != nil {
		return "", err
	}
	p := parser.New(tokens)
	prog, err := p.ParseProgram()
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	ev := evaluator.New(&buf)
	env := evaluator.NewEnvironment()
	res := ev.Eval(prog, env)

	if errObj, ok := res.(*evaluator.RuntimeError); ok {
		return buf.String(), errObj
	}

	return strings.TrimSuffix(buf.String(), "\n"), nil
}

func TestEvaluatorSuccessTableDriven(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Empty Program",
			input:    "",
			expected: "",
		},
		{
			name:     "Omong String Coercion",
			input:    "umur = 20\nomong(\"Umur: \" + umur)",
			expected: "Umur: 20",
		},
		{
			name:     "Arithmetic Expression",
			input:    "x = (10 + 5) * 2 - 4 / 2\nomong(x)",
			expected: "28",
		},
		{
			name:     "Boolean Equality Comparison",
			input:    "omong(10 > 5)\nomong(10 == 20)\nomong(\"MOJA\" == \"MOJA\")",
			expected: "bener\nsalah\nbener",
		},
		{
			name:     "Yen Liyane Condition",
			input:    "umur = 16\nyen umur >= 17 {\n\tomong(\"Cukup\")\n} liyane {\n\tomong(\"Durung\")\n}",
			expected: "Durung",
		},
		{
			name:     "Baleni Loop Accumulator",
			input:    "total = 0\nbaleni i saka 1 tekan 5 {\n\ttotal = total + i\n}\nomong(total)",
			expected: "15",
		},
		{
			name:     "Baleni Zero Iteration",
			input:    "omong(\"Start\")\nbaleni i saka 5 tekan 1 {\n\tomong(\"Fail\")\n}\nomong(\"End\")",
			expected: "Start\nEnd",
		},
		{
			name:     "Function Scope & Mutation",
			input:    "x = 10\ngawe ubah() {\n\tx = 99\n}\nubah()\nomong(x)",
			expected: "99",
		},
		{
			name:     "Recursive Factorial",
			input:    "gawe faktorial(n) {\n\tyen n <= 1 {\n\t\tbalik 1\n\t} liyane {\n\t\tbalik n * faktorial(n - 1)\n\t}\n}\nomong(faktorial(5))",
			expected: "120",
		},
		{
			name:     "Array Literal and Printing",
			input:    `angka = [10, 20, 30]` + "\n" + `campuran = [10, "halo", bener]` + "\n" + `omong(angka)` + "\n" + `omong(campuran)`,
			expected: "[10, 20, 30]\n[10, \"halo\", bener]",
		},
		{
			name:     "Array Index Access and Expression Index",
			input:    `angka = [10, 20, 30]` + "\n" + `omong(angka[0])` + "\n" + `omong(angka[1 + 1])`,
			expected: "10\n30",
		},
		{
			name:     "Nested Array Access",
			input:    `matriks = [[1, 2], [3, 4]]` + "\n" + `omong(matriks[0][1])` + "\n" + `omong(matriks[1][0])`,
			expected: "2\n3",
		},
		{
			name:     "Array Element Mutation",
			input:    `angka = [10, 20, 30]` + "\n" + `angka[1] = 99` + "\n" + `omong(angka[1])` + "\n" + `omong(angka)`,
			expected: "99\n[10, 99, 30]",
		},
		{
			name:     "Panjang Builtin Function",
			input:    `angka = [10, 20, 30]` + "\n" + `kosong = []` + "\n" + `omong(panjang(angka))` + "\n" + `omong(panjang(kosong))` + "\n" + `omong(panjang("MOJA"))`,
			expected: "3\n0\n4",
		},
		// --- Operator Logika v0.4 ---
		{
			name:     "Logic AND bener lan bener",
			input:    "omong(bener lan bener)",
			expected: "bener",
		},
		{
			name:     "Logic AND bener lan salah",
			input:    "omong(bener lan salah)",
			expected: "salah",
		},
		{
			name:     "Logic AND salah lan bener",
			input:    "omong(salah lan bener)",
			expected: "salah",
		},
		{
			name:     "Logic AND salah lan salah",
			input:    "omong(salah lan salah)",
			expected: "salah",
		},
		{
			name:     "Logic OR bener utawa bener",
			input:    "omong(bener utawa bener)",
			expected: "bener",
		},
		{
			name:     "Logic OR bener utawa salah",
			input:    "omong(bener utawa salah)",
			expected: "bener",
		},
		{
			name:     "Logic OR salah utawa bener",
			input:    "omong(salah utawa bener)",
			expected: "bener",
		},
		{
			name:     "Logic OR salah utawa salah",
			input:    "omong(salah utawa salah)",
			expected: "salah",
		},
		{
			name:     "Logic NOT ora bener",
			input:    "omong(ora bener)",
			expected: "salah",
		},
		{
			name:     "Logic NOT ora salah",
			input:    "omong(ora salah)",
			expected: "bener",
		},
		{
			name:     "Logic NOT double ora",
			input:    "omong(ora ora bener)",
			expected: "bener",
		},
		{
			name:     "Logic Variable Assignment",
			input:    "a = bener\nb = salah\nomong(a lan b)\nomong(a utawa b)\nomong(ora a)",
			expected: "salah\nbener\nsalah",
		},
		{
			name:     "Logic Precedence lan tighter than utawa",
			input:    "omong(bener utawa salah lan salah)",
			expected: "bener",
		},
		{
			name:     "Logic Precedence ora before lan",
			input:    "omong(ora bener lan salah)",
			expected: "salah",
		},
		{
			name:     "Logic Parentheses Override",
			input:    "omong((bener utawa salah) lan salah)",
			expected: "salah",
		},
		{
			name:     "Logic With Comparison In Yen",
			input:    `umur = 20` + "\n" + `duwe_ktp = bener` + "\n" + `yen umur >= 17 lan duwe_ktp {` + "\n\t" + `omong("Syarat terpenuhi")` + "\n}",
			expected: "Syarat terpenuhi",
		},
		{
			name:     "Logic Short Circuit AND skip right on false",
			input:    "x = 0\nomong(salah lan (x / x > 0))",
			expected: "salah",
		},
		{
			name:     "Logic Short Circuit OR skip right on true",
			input:    "x = 0\nomong(bener utawa (x / x > 0))",
			expected: "bener",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := run(tt.input)
			if err != nil {
				t.Fatalf("Runtime error unexpected: %v", err)
			}
			if out != tt.expected {
				t.Errorf("Output mismatch:\nGot:  %q\nWant: %q", out, tt.expected)
			}
		})
	}
}

func TestEvaluatorRuntimeErrorCasesTableDriven(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectedMsg string
	}{
		{
			name:        "Division By Zero",
			input:       "x = 10 / 0",
			expectedMsg: "Pembagian dengan nol tidak diperbolehkan",
		},
		{
			name:        "Undefined Variable",
			input:       "omong(variabel_gaib)",
			expectedMsg: "Variabel 'variabel_gaib' belum dideklarasikan",
		},
		{
			name:        "Invalid Arithmetic Type",
			input:       `x = "teks" - 5`,
			expectedMsg: "Operator '-' hanya mendukung tipe Number",
		},
		{
			name:        "Iterator Variable Mutation",
			input:       "baleni i saka 1 tekan 5 {\n\ti = 10\n}",
			expectedMsg: "Variabel iterator 'i' tidak dapat diubah di dalam loop",
		},
		{
			name:        "Return Outside Function",
			input:       "balik 42",
			expectedMsg: "'balik' hanya dapat digunakan di dalam fungsi",
		},
		{
			name:        "Wrong Argument Count",
			input:       "gawe sapa(nama) {\n\tomong(nama)\n}\nsapa()",
			expectedMsg: "Fungsi 'sapa' butuh 1 argumen, diberikan 0",
		},
		{
			name:        "Calling Non Function",
			input:       "x = 10\nx()",
			expectedMsg: "'x' bukan merupakan fungsi",
		},
		{
			name:        "Non Boolean Condition In Yen",
			input:       `yen "bukan_boolean" { omong("ok") }`,
			expectedMsg: "Kondisi 'yen' harus bernilai Boolean",
		},
		{
			name:        "Non Number Bound In Baleni",
			input:       `baleni i saka "satu" tekan 5 {}`,
			expectedMsg: "Batas awal perulangan 'baleni' harus berupa Number",
		},
		{
			name:        "Array Index Out Of Bounds",
			input:       "angka = [10, 20]\nomong(angka[5])",
			expectedMsg: "Indeks array di luar rentang (posisi 5, total elemen 2)",
		},
		{
			name:        "Array Index Negative",
			input:       "angka = [10, 20]\nomong(angka[-1])",
			expectedMsg: "Indeks array tidak boleh negatif",
		},
		{
			name:        "Array Index Float",
			input:       "angka = [10, 20]\nomong(angka[1.5])",
			expectedMsg: "Indeks array harus berupa bilangan bulat (integer)",
		},
		{
			name:        "Array Index Non Number",
			input:       `angka = [10, 20]` + "\n" + `omong(angka["siji"])`,
			expectedMsg: "Indeks array harus berupa tipe Number",
		},
		{
			name:        "Index Access On Non Array",
			input:       "val = 10\nomong(val[0])",
			expectedMsg: "Akses indeks '[...]' hanya berlaku untuk tipe Array",
		},
		{
			name:        "Panjang Function Invalid Type",
			input:       "omong(panjang(100))",
			expectedMsg: "Fungsi 'panjang' hanya mendukung tipe Array atau String",
		},
		// --- Error Operator Logika v0.4 ---
		{
			name:        "Logic AND Left Non Boolean",
			input:       "omong(10 lan bener)",
			expectedMsg: "Operand kiri operator 'lan' harus bertipe Boolean",
		},
		{
			name:        "Logic AND Right Non Boolean",
			input:       "omong(bener lan 10)",
			expectedMsg: "Operand kanan operator 'lan' harus bertipe Boolean",
		},
		{
			name:        "Logic OR Left Non Boolean",
			input:       `omong("teks" utawa bener)`,
			expectedMsg: "Operand kiri operator 'utawa' harus bertipe Boolean",
		},
		{
			name:        "Logic NOT Non Boolean",
			input:       "omong(ora 42)",
			expectedMsg: "Operator 'ora' hanya mendukung tipe Boolean",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := run(tt.input)
			if err == nil {
				t.Fatalf("Expected runtime error containing %q, got nil", tt.expectedMsg)
			}
			if !strings.Contains(err.Error(), tt.expectedMsg) {
				t.Errorf("Error message mismatch:\nGot:  %v\nWant substring: %q", err, tt.expectedMsg)
			}
		})
	}
}
