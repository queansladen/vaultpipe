package env

// Stage is a function that transforms a slice of KEY=VALUE pairs.
type Stage func([]string) ([]string, error)

// Pipeline chains multiple stages together, passing the output of each stage
// as the input to the next. If any stage returns an error the pipeline halts
// and that error is returned immediately.
type Pipeline struct {
	stages []Stage
}

// NewPipeline returns a Pipeline that will execute the provided stages in
// order. Passing zero stages is valid; Run will return the input unchanged.
func NewPipeline(stages ...Stage) *Pipeline {
	return &Pipeline{stages: stages}
}

// Add appends one or more stages to the end of the pipeline.
func (p *Pipeline) Add(stages ...Stage) {
	p.stages = append(p.stages, stages...)
}

// Len returns the number of stages currently in the pipeline.
func (p *Pipeline) Len() int {
	return len(p.stages)
}

// Run executes all stages in order against env, returning the final
// transformed slice or the first error encountered.
func (p *Pipeline) Run(env []string) ([]string, error) {
	current := env
	for _, stage := range p.stages {
		var err error
		current, err = stage(current)
		if err != nil {
			return nil, err
		}
	}
	return current, nil
}
