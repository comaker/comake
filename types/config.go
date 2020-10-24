package types

// Step is a build step
type Step struct {
	Name   string   `yaml:"name"`
	Image  string   `yaml:"image"`
	Script []string `yaml:"script"`
}

// Build holds build steps
type Build struct {
	Steps []Step `yaml:"steps"`
}
