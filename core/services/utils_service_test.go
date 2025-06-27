package services

import (
	"testing"
	"time"
)

func TestIsValidLatitude(t *testing.T) {
	tests := []struct {
		lat      float64
		expected bool
	}{
		{0, true},
		{90, true},
		{-90, true},
		{91, false},
		{-91, false},
	}
	for _, tt := range tests {
		if got := isValidLatitude(tt.lat); got != tt.expected {
			t.Errorf("isValidLatitude(%v) = %v; want %v", tt.lat, got, tt.expected)
		}
	}
}

func TestIsValidLongitude(t *testing.T) {
	tests := []struct {
		lon      float64
		expected bool
	}{
		{0, true},
		{180, true},
		{-180, true},
		{181, false},
		{-181, false},
	}
	for _, tt := range tests {
		if got := isValidLongitude(tt.lon); got != tt.expected {
			t.Errorf("isValidLongitude(%v) = %v; want %v", tt.lon, got, tt.expected)
		}
	}
}

func TestGenerateRandomPassword(t *testing.T) {
	pass, err := GenerateRandomPassword(MinPasswordLength)
	if err != nil {
		t.Errorf("GenerateRandomPassword returned error: %v", err)
	}
	if len(pass) != MinPasswordLength {
		t.Errorf("Expected password length %d, got %d", MinPasswordLength, len(pass))
	}
}

func TestGenerateRandomToken(t *testing.T) {
	token := GenerateRandomToken()
	if len(token) != 64 { // 32 bytes hex encoded
		t.Errorf("Expected token length 64, got %d", len(token))
	}
}

func TestParseTime(t *testing.T) {
	cases := []struct {
		input    string
		expected time.Time
	}{
		{"2023-01-01T12:00:00Z", time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)},
		{"2023-01-01 12:00:00", time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)},
		{"invalid", time.Time{}},
	}
	for _, c := range cases {
		got := ParseTime(c.input)
		if !got.Equal(c.expected) {
			t.Errorf("ParseTime(%q) = %v; want %v", c.input, got, c.expected)
		}
	}
}

// MockPaginationParams is a mock for PaginationParams interface

type MockPaginationParams struct {
	limit  int32
	offset int32
}

func (m *MockPaginationParams) GetLimit() int32   { return m.limit }
func (m *MockPaginationParams) GetOffset() int32  { return m.offset }
func (m *MockPaginationParams) SetLimit(l int32)  { m.limit = l }
func (m *MockPaginationParams) SetOffset(o int32) { m.offset = o }

func TestSetDefaultPagination(t *testing.T) {
	m := &MockPaginationParams{limit: 0, offset: -1}
	SetDefaultPagination(m)
	if m.limit != 10 {
		t.Errorf("Expected default limit 10, got %d", m.limit)
	}
	if m.offset != 0 {
		t.Errorf("Expected default offset 0, got %d", m.offset)
	}
}
