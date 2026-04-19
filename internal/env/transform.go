package env

// Transformer applies an ordered pipeline of env slice transformations.
type Transformer struct {
	steps []func([]string) []string
}

// NewTransformer returns a Transformer with the given steps.
func NewTransformer(steps ...func([]string) []string) *Transformer {
	return &Transformer{steps: steps}
}

// Apply runs each step in order, passing the output of one as the input of the next.
func (t *Transformer) Apply(env []string) []string {
	result := make([]string, len(env))
	copy(result, env)
	for _, step := range t.steps {
		result = step(result)
	}
	return result
}

// Add appends additional transformation steps to the pipeline.
func (t *Transformer) Add(steps ...func([]string) []string) {
	t.steps = append(t.steps, steps...)
}

// Len returns the number of steps in the pipeline.
func (t *Transformer) Len() int {
	return len(t.steps)
}
