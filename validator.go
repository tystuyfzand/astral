package astral

import (
	"errors"
	"fmt"
	"regexp"
)

var (
	UsageError  = errors.New("usage")
	emojiRegexp = regexp.MustCompile("<(a?):(.+?):(\\d+)>")
)

// InvalidValueError is an error type thrown when a value is invalid/unknown
type InvalidValueError struct {
	Argument string
	Value    interface{}
}

// Error constructs a string for the error with the argument and value
func (i InvalidValueError) Error() string {
	return fmt.Sprintf("unknown argument value for %s: %v", i.Argument, i.Value)
}

// Validate checks the context against the Route's defined arguments and ensures all required arguments
// and types are satisfied.
func (r *Route) Validate(ctx *Context) error {
	var err error

	for _, arg := range r.Arguments {
		argValue, exists := ctx.Arguments[arg.Name]

		if arg.Required && (argValue == "" || argValue == nil || !exists) {
			return fmt.Errorf("The %s argument is required.", arg.Name)
		}

		switch arg.Type {
		case ArgumentTypeInt:
			err = validateInt(ctx, arg, argValue.(int64))
		case ArgumentTypeFloat:
			err = validateFloat(ctx, arg, argValue.(float64))
		}

		if err != nil {
			return err
		}

		if len(arg.Choices) > 0 {
			// Ensure options contains value
			found := false

			for _, value := range arg.Choices {
				if value.Value == argValue {
					found = true
					break
				}
			}

			if !found {
				return InvalidValueError{Argument: arg.Name, Value: argValue}
			}
		}
	}

	return nil
}

// validateInt checks an integer argument to ensure it's a valid integer
func validateInt(ctx *Context, arg *Argument, v int64) error {
	if arg.Min != nil && v < arg.Min.(int64) {
		return fmt.Errorf("%s must be larger than %d.", arg.Name, arg.Min)
	}

	if arg.Max != nil && v < arg.Max.(int64) {
		return fmt.Errorf("%s must be smaller than %d.", arg.Name, arg.Max)
	}

	return nil
}

// validateFloat checks an integer argument to ensure it's a valid float
func validateFloat(ctx *Context, arg *Argument, v float64) error {
	if arg.Min != nil && v < arg.Min.(float64) {
		return fmt.Errorf("%s must be larger than %f.", arg.Name, arg.Min)
	}

	if arg.Max != nil && v < arg.Max.(float64) {
		return fmt.Errorf("%s must be smaller than %f.", arg.Name, arg.Max)
	}

	return nil
}
