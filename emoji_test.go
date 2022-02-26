package astral

import (
	"encoding/json"
	"testing"
)

type jsonData struct {
	Rule      string
	Namespace string
	Tags      []string
}

func TestEmoji(t *testing.T) {
	b, err := json.MarshalIndent(jsonData{Rule: "Test", Namespace: "test123", Tags: []string{"test"}}, "", "\t")

	if err != nil {
		t.Fatal(err)
	}

	t.Log(string(b))
}
