package arguments

import (
	"testing"
)

func TestParse(t *testing.T) {
	args := Parse("normal \"testing quoted\" and normal \"end closed\"")
	expected := []string{"normal", "testing quoted", "and", "normal", "end closed"}

	for i := 0; i < len(expected); i++ {
		if args[i] != expected[i] {
			t.Errorf("Expected %s, got %s at index %d", expected[i], args[i], i)
		}
	}
}