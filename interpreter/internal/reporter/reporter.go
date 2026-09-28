// Package reporter provides formatted error reporting with code snippets and caret pointers for MOJA.
package reporter

import (
	"fmt"
	"regexp"
	"strings"
)

// PositionalError is an interface implemented by LexError, ParseError, and RuntimeError.
type PositionalError interface {
	error
	GetPosition() (line int, column int)
	GetErrorType() string
}

var prefixRegex = regexp.MustCompile(`^\[MOJA Error\] Baris \d+(, Kolom \d+)?: `)

// FormatError formats a positional error into a human-readable diagnostic message
// complete with file name, line number, column number, code snippet, and caret (^) indicator.
func FormatError(filename string, source string, line int, col int, errType string, message string) string {
	if filename == "" {
		filename = "main.mojo"
	}
	if errType == "" {
		errType = "Error"
	}

	// Clean up raw string prefix if present
	cleanMsg := prefixRegex.ReplaceAllString(message, "")

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Error: %s\n", errType))
	sb.WriteString(fmt.Sprintf("File: %s\n", filename))
	sb.WriteString(fmt.Sprintf("Line %d, Column %d\n\n", line, col))

	lines := strings.Split(strings.ReplaceAll(source, "\r\n", "\n"), "\n")

	if line >= 1 && line <= len(lines) {
		lineStr := lines[line-1]

		// Expand tabs to 4 spaces for display and accurate pointer calculation
		expandedLine := ""
		colOffset := 0
		colIndex := 0

		for i, r := range []rune(lineStr) {
			if i+1 == col {
				colIndex = colOffset
			}
			if r == '\t' {
				expandedLine += "    "
				colOffset += 4
			} else {
				expandedLine += string(r)
				colOffset++
			}
		}
		if col > len([]rune(lineStr)) {
			colIndex = colOffset
		}

		lineNumStr := fmt.Sprintf("%d", line)
		padding := strings.Repeat(" ", len(lineNumStr))

		sb.WriteString(fmt.Sprintf("  %s | %s\n", lineNumStr, expandedLine))
		caretSpaces := strings.Repeat(" ", colIndex)
		sb.WriteString(fmt.Sprintf("  %s | %s^\n\n", padding, caretSpaces))
	}

	sb.WriteString(cleanMsg)
	return sb.String()
}
