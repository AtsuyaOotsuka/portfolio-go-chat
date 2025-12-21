package command

import "fmt"

type RootCommandInterface interface {
	Run(args []string)
}

type RootCommand struct {
	BaseCommand
}

func NewRootCommand() *RootCommand {
	return &RootCommand{}
}

func (c *RootCommand) Run(args []string) {
	fmt.Println("コマンドのヘルプを見るには --help を付けて実行してください")
}
