package application

import (
	"testing"

	"github.com/magiconair/properties/assert"
)

func TestFilterBadWords(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "No bad words",
			input:    "No bad words",
			expected: "No bad words",
		},
		{
			name:     "One bad word kerfuffle",
			input:    "One bad word kerfuffle",
			expected: "One bad word ****",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := filterBadWords(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
