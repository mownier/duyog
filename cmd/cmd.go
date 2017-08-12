package cmd

// Command interface
type Command interface {
	Parse() bool
	PrintHint()
}

// Parse method
func Parse(c Command) bool {
	return c.Parse()
}

// PrintHint method
func PrintHint(cmd ...Command) {
	for _, c := range cmd {
		c.PrintHint()
	}
}
