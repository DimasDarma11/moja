package reporter_test

import (
	"strings"
	"testing"

	"mojolang/internal/reporter"
)

func TestFormatError(t *testing.T) {
	src := "x = 10\nomong(\"Halo\"\ny = 20"
	formatted := reporter.FormatError("test.mojo", src, 2, 13, "Syntax Error", "Expected ')' after expression.")

	if !strings.Contains(formatted, "Error: Syntax Error") {
		t.Errorf("Expected 'Error: Syntax Error', got:\n%s", formatted)
	}
	if !strings.Contains(formatted, "File: test.mojo") {
		t.Errorf("Expected 'File: test.mojo', got:\n%s", formatted)
	}
	if !strings.Contains(formatted, "Line 2, Column 13") {
		t.Errorf("Expected 'Line 2, Column 13', got:\n%s", formatted)
	}
	if !strings.Contains(formatted, "2 | omong(\"Halo\"") {
		t.Errorf("Expected line snippet, got:\n%s", formatted)
	}
	if !strings.Contains(formatted, "^") {
		t.Errorf("Expected caret pointer '^', got:\n%s", formatted)
	}
}
