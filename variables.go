package astral

func NewVariableBag() *VariableBag {
	return &VariableBag{
		vars: make(map[string]interface{}),
	}
}

type VariableBag struct {
	vars map[string]interface{}
}

// Set sets a variable on the context
func (v *VariableBag) Set(key string, d interface{}) {
	v.vars[key] = d
}

// Get retrieves a variable from the context
func (v *VariableBag) Get(key string) interface{} {
	if c, ok := v.vars[key]; ok {
		return c
	}
	return nil
}
