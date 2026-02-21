package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMaskEmailForRegistration(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"user@example.com", "u***@example.com"},
		{"a@example.com", "a***@example.com"},
		{"ab@example.com", "a***@example.com"},
		{"invalid", "***"},
		{"@example.com", "***@example.com"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := maskEmail(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
