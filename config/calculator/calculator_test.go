package calculator

import (
	"fmt"
	"os"
	"testing"
)

func TestCalc(t *testing.T) {
	// Создаем директорию database, если её нет
	if err := os.MkdirAll("database_test", os.ModePerm); err != nil {
		t.Fatalf("Failed to create database directory: %v", err)
	}

	// Создаем файл results.jsonl, если его нет
	file, err := os.OpenFile("database_test/results_test.jsonl", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		t.Fatalf("Failed to create or open results.jsonl: %v", err)
	}
	file.Close()

	fmt.Println("Testing Calc function:")

	testCases := []struct {
		expr     string
		expected float64
		hasError bool
	}{
		{"2 + 2 + 2 + 2 + 2 + 2 + (2 + (2 + (2 + 2)))", 20, false},
		{"1 + 1", 2, false},
		{"(2 + 2) * 2", 8, false},
		{"2 + 2 * 2", 6, false},
		{"1 + 1 *", 0, true},
		{"2 / 0", 0, true},
		{"2.5 * 3", 7.5, false},
		{"2 * (3 + 4)", 14, false},
		{"2 * 3 + 4", 10, false},
		{"2 * (3 + 4 * 2)", 22, false},
	}

	for _, tc := range testCases {
		t.Run(tc.expr, func(t *testing.T) {
			result, err := Calc(tc.expr, 0)
			if tc.hasError {
				if err == nil {
					t.Errorf("Expected error for expression: %s, but got none", tc.expr)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error for expression: %s, error: %v", tc.expr, err)
				}
				if result != tc.expected {
					t.Errorf("For expression: %s, expected: %v, got: %v", tc.expr, tc.expected, result)
				}
			}
		})
	}
}
