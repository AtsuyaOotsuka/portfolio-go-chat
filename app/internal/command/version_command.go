package command

import (
	"fmt"
)

type VersionCommandInterface interface {
	Run(args []string)
}

type VersionCommand struct {
	BaseCommand
}

func NewVersionCommand() *VersionCommand {
	return &VersionCommand{}
}

func (c *VersionCommand) Run(args []string) {
	fmt.Println("mycli version 1.0.0")
}
