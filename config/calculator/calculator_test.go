package calculator

import (
	"testing"
)

func TestCalculateExpression(t *testing.T) {
	tests := []struct {
		expr   string
		expect float64
		hasErr bool
	}{
		{"2 + 2 + 2 + 2 + 2 + 2 + (2 + (2 + (2 + 2)))", 20, false},
		{"1+1", 2, false},
		{"(2+2)*2", 8, false},
		{"2+2*2", 6, false},
		{"1+1*", 0, true},
		{"2 / 0", 0, true},
		{"2.5 * 3", 7.5, false},
		{"-1 + 2", 1, false},
		{"2 * (3 + 4)", 14, false},
		{"2 * 3 + 4", 10, false},
		{"2 * (3 + 4 * 2)", 22, false},
	}

	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			res, err := CalculateExpression(tt.expr, 0)
			if tt.hasErr {
				if err == nil {
					t.Error(tt.expr)
				}
			} else {
				if err != nil {
					t.Error(tt.expr, err)
				}
				if res != tt.expect {
					t.Error(tt.expr, res, tt.expect)
				}
			}
		})
	}
}
