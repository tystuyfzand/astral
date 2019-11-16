package arguments

type SpaceTokenizer struct {
	Input string
}

// NextToken returns the next token, or an empty string if nothing exists.
func (t *SpaceTokenizer) NextToken() string {
	if len(t.Input) == 0 {
		return ""
	}

	ch := t.Input[0]

	if ch == '"' {
		// Scan until closing quote or end of string
		for i := 1; i < len(t.Input); i++ {
			if t.Input[i] == '"' && t.Input[i-1] != '\\' {
				ret := t.Input[1:i]
				if i + 2 < len(t.Input) {
					t.Input = t.Input[i+2:]
				} else {
					t.Input = ""
				}
				return ret
			}
		}
	} else {
		for i := 0; i < len(t.Input); i++ {
			if t.Input[i] == ' ' {
				ret := t.Input[0:i]
				if i + 1 < len(t.Input) {
					t.Input = t.Input[i+1:]
				} else {
					t.Input = ""
				}
				return ret
			}
		}

	}

	ret := t.Input

	t.Input = ""

	return ret
}

// Empty checks if the input string is empty.
func (t *SpaceTokenizer) Empty() bool {
	return t.Input == ""
}

func NewSpaceTokenizer(input string) *SpaceTokenizer {
	return &SpaceTokenizer{Input: input}
}

// Parse parses a command argument string using a Space Tokenizer
func Parse(command string) []string {
	tokenizer := NewSpaceTokenizer(command)

	arguments := make([]string, 0)

	for !tokenizer.Empty() {
		arguments = append(arguments, tokenizer.NextToken())
	}

	return arguments
}