package cmd

var (
	verbose bool
)

func (c *Cmd) setupFlags() {
	c.Cmd.PersistentFlags().BoolVarP(
		&verbose,
		"verbose",
		"v",
		false,
		"verbose output",
	)
}
