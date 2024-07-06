package random

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewString(t *testing.T) {
	tests := []struct {
		name   string
		length int
	}{
		{name: "should generate string with length 1", length: 1},
		{name: "should generate string with length 2", length: 2},
		{name: "should generate string with length 3", length: 3},
		{name: "should generate string with length 4", length: 4},
		{name: "should generate string with length 5", length: 5},
		{name: "should generate string with length 6", length: 6},
		{name: "should generate string with length 7", length: 7},
		{name: "should generate string with length 8", length: 8},
		{name: "should generate string with length 9", length: 9},
		{name: "should generate string with length 10", length: 10},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			string1 := NewString(tt.length)
			string2 := NewString(tt.length)

			assert.Len(t, string1, tt.length)
			assert.Len(t, string2, tt.length)

			assert.NotEqual(t, string1, string2)

			for _, char := range string1 {
				assert.Contains(t, characters, string(char))
			}
		})
	}
}
