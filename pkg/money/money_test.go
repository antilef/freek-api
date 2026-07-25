package money

import (
	"testing"
)

func TestCLPString(t *testing.T) {
	tests := []struct {
		input    CLP
		expected string
	}{
		{0, "$0"},
		{1, "$1"},
		{100, "$100"},
		{999, "$999"},
		{1000, "$1.000"},
		{1234, "$1.234"},
		{12345, "$12.345"},
		{123456, "$123.456"},
		{1234567, "$1.234.567"},
		{120000, "$120.000"},
		{1200000, "$1.200.000"},
	}

	for _, tt := range tests {
		result := tt.input.String()
		if result != tt.expected {
			t.Errorf("CLP(%d).String() = %s, want %s", int64(tt.input), result, tt.expected)
		}
	}
}

func TestCLPInt64(t *testing.T) {
	tests := []struct {
		input    CLP
		expected int64
	}{
		{0, 0},
		{100, 100},
		{120000, 120000},
	}

	for _, tt := range tests {
		result := tt.input.Int64()
		if result != tt.expected {
			t.Errorf("CLP(%d).Int64() = %d, want %d", int64(tt.input), result, tt.expected)
		}
	}
}

func TestFormatDots(t *testing.T) {
	tests := []struct {
		input    int64
		expected string
	}{
		{0, "0"},
		{1, "1"},
		{100, "100"},
		{999, "999"},
		{1000, "1.000"},
		{1234567, "1.234.567"},
	}

	for _, tt := range tests {
		result := FormatDots(tt.input)
		if result != tt.expected {
			t.Errorf("FormatDots(%d) = %s, want %s", tt.input, result, tt.expected)
		}
	}
}

func TestParse(t *testing.T) {
	tests := []struct {
		input    string
		expected CLP
		hasError bool
	}{
		{"$0", 0, false},
		{"$100", 100, false},
		{"$120.000", 120000, false},
		{"120.000", 120000, false},
		{"$1.234.567", 1234567, false},
		{"100", 100, false},
		{"abc", 0, true},
		{"$abc", 0, true},
	}

	for _, tt := range tests {
		result, err := Parse(tt.input)
		if (err != nil) != tt.hasError {
			t.Errorf("Parse(%q) error = %v, wantError %v", tt.input, err, tt.hasError)
			continue
		}
		if result != tt.expected {
			t.Errorf("Parse(%q) = %d, want %d", tt.input, int64(result), int64(tt.expected))
		}
	}
}

func TestDivideRoundUp(t *testing.T) {
	tests := []struct {
		amount      CLP
		installments int
		expected    CLP
	}{
		{100000, 1, 100000},
		{100000, 2, 50000},
		{100000, 3, 33334},  // 33333.33 → rounds up
		{100000, 6, 16667},  // 16666.66 → rounds up
		{120000, 6, 20000},  // exact
		{100000, 0, 100000}, // no installments
		{100000, -1, 100000},
	}

	for _, tt := range tests {
		result := DivideRoundUp(tt.amount, tt.installments)
		if result != tt.expected {
			t.Errorf("DivideRoundUp(%d, %d) = %d, want %d", int64(tt.amount), tt.installments, int64(result), int64(tt.expected))
		}
	}
}
