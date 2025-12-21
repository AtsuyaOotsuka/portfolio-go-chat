package command

import (
	"fmt"
)

type VersionCommandInterface interface {
	Run()
}

type VersionCommand struct {
	BaseCommand
}

func NewVersionCommand() *VersionCommand {
	return &VersionCommand{}
}

func (c *VersionCommand) Run() {
	fmt.Println("mycli version 1.0.0")
}
