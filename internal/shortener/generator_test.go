package shortener

import "testing"

func TestIsValidURL(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "valid https url",
			input:    "https://google.com",
			expected: true,
		},
		{
			name:     "valid http url",
			input:    "http://example.com",
			expected: true,
		},
		{
			name:     "invalid url",
			input:    "hello",
			expected: false,
		},
		{
			name:     "invalid scheme",
			input:    "javascript:alert(1)",
			expected: false,
		},
	}
	for _, test := range tests {
		result := IsValidURL(test.input)
		if result != test.expected {
			t.Errorf("IsValidURL(%q) = %v; want %v",
				test.input,
				result,
				test.expected)
		}
	}
}
