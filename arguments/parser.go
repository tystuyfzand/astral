package arguments

import "github.com/diamondburned/arikawa/v3/utils/bot/extras/shellwords"

// Parse parses a command argument string using a Space Tokenizer
func Parse(command string) []string {
	args, err := shellwords.Parse(command)

	if err != nil {
		return []string{command}
	}

	return args
}
