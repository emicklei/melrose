package dsl

import (
	"errors"
	"fmt"
	"strings"
)

const fourSpaces = "    "

// CleanStatements breaks a DSL program into statements ready for evaluation.
// Statements may be separated by newlines or semicolons. Indented lines are
// continuations of the preceding statement.
func CleanStatements(source string) ([]string, error) {
	physicalLines := strings.Split(source, "\n")
	logicalLines := make([]string, 0, len(physicalLines))
	previousLineNumber := -1
	lastExpressionIndex := -1

	for lineNumber, line := range physicalLines {
		code, comment := splitComment(line)
		trimmed := strings.TrimSpace(code)
		continuation := strings.HasPrefix(line, "\t") || strings.HasPrefix(line, fourSpaces)
		if trimmed == ")" && lastExpressionIndex >= 0 && parenthesisDepth(logicalLines[lastExpressionIndex]) > 0 {
			continuation = true
		}
		if trimmed == "" && comment != "" {
			logicalLines = append(logicalLines, comment)
			previousLineNumber = lineNumber
			continue
		}

		if continuation {
			if lastExpressionIndex < 0 {
				if len(logicalLines) > 0 && strings.TrimSpace(logicalLines[len(logicalLines)-1]) == "" {
					part := strings.TrimSpace(strings.ReplaceAll(code, "\t", " "))
					if part != "" {
						logicalLines = append(logicalLines, part)
						lastExpressionIndex = len(logicalLines) - 1
					}
					previousLineNumber = lineNumber
					continue
				}
				if len(logicalLines) > 0 && strings.HasPrefix(strings.TrimSpace(logicalLines[len(logicalLines)-1]), "//") {
					part := strings.TrimSpace(strings.ReplaceAll(code, "\t", " "))
					if part == "" {
						return nil, fmt.Errorf("syntax error, line with TAB [%d] must be part of expression", lineNumber+1)
					}
					logicalLines = append(logicalLines, part)
					lastExpressionIndex = len(logicalLines) - 1
					previousLineNumber = lineNumber
					continue
				}
				if len(logicalLines) == 0 && lineNumber > 0 {
					part := strings.TrimSpace(strings.ReplaceAll(code, "\t", " "))
					if part != "" {
						logicalLines = append(logicalLines, part)
						lastExpressionIndex = len(logicalLines) - 1
					}
					previousLineNumber = lineNumber
					continue
				}
				return nil, errors.New("syntax error, first line cannot start with TAB")
			}
			if lineNumber != previousLineNumber+1 {
				return nil, fmt.Errorf("syntax error, line with TAB [%d] must be part of expression", lineNumber+1)
			}
			_, previousComment := splitComment(logicalLines[lastExpressionIndex])
			if (comment != "" || previousComment != "") && parenthesisDepth(logicalLines[lastExpressionIndex]) > 0 {
				return nil, fmt.Errorf("syntax error, comment on line [%d] is inside an expression", lineNumber+1)
			}
			part := strings.TrimSpace(strings.ReplaceAll(code, "\t", " "))
			if part != "" {
				logicalLines[lastExpressionIndex] += " " + part
			}
			if comment != "" {
				logicalLines[lastExpressionIndex] += " " + comment
			}
			previousLineNumber = lineNumber
			continue
		}

		statement := strings.TrimSpace(strings.ReplaceAll(code, "\t", " "))
		if comment != "" {
			if statement != "" {
				statement += " "
			}
			statement += comment
		}
		logicalLines = append(logicalLines, statement)
		if trimmed != "" {
			lastExpressionIndex = len(logicalLines) - 1
		}
		previousLineNumber = lineNumber
	}

	statements := make([]string, 0, len(logicalLines))
	for _, line := range logicalLines {
		for _, statement := range splitStatements(line) {
			if statement != "" {
				statements = append(statements, statement)
			}
		}
	}
	return statements, nil
}

// Format returns a canonical, one-statement-per-line representation of source.
func Format(source string) (string, error) {
	statements, err := CleanStatements(source)
	if err != nil {
		return "", err
	}
	return strings.Join(statements, "\n"), nil
}

func splitStatements(line string) []string {
	statements := []string{}
	start := 0
	quote := byte(0)
	escaped := false
	for index := 0; index < len(line); index++ {
		character := line[index]
		if escaped {
			escaped = false
			continue
		}
		if quote != 0 {
			if character == '\\' {
				escaped = true
			} else if character == quote {
				quote = 0
			}
			continue
		}
		if character == '\'' || character == '"' || character == '`' {
			quote = character
		} else if character == '/' && index+1 < len(line) && line[index+1] == '/' {
			statements = append(statements, strings.TrimSpace(line[start:]))
			return statements
		} else if character == ';' {
			statements = append(statements, strings.TrimSpace(line[start:index]))
			start = index + 1
		}
	}
	return append(statements, strings.TrimSpace(line[start:]))
}

func splitComment(line string) (code, comment string) {
	quote := byte(0)
	escaped := false
	for index := 0; index < len(line)-1; index++ {
		character := line[index]
		if escaped {
			escaped = false
			continue
		}
		if quote != 0 {
			if character == '\\' {
				escaped = true
			} else if character == quote {
				quote = 0
			}
			continue
		}
		if character == '\'' || character == '"' || character == '`' {
			quote = character
			continue
		}
		if character == '/' && line[index+1] == '/' {
			return line[:index], strings.TrimSpace(line[index:])
		}
	}
	return line, ""
}

func parenthesisDepth(line string) int {
	depth := 0
	quote := byte(0)
	escaped := false
	for index := 0; index < len(line); index++ {
		character := line[index]
		if escaped {
			escaped = false
			continue
		}
		if quote != 0 {
			if character == '\\' {
				escaped = true
			} else if character == quote {
				quote = 0
			}
			continue
		}
		if character == '\'' || character == '"' || character == '`' {
			quote = character
		} else if character == '(' {
			depth++
		} else if character == ')' {
			depth--
		}
	}
	return depth
}
