package command

import "fmt"

type RootCommandInterface interface {
	Run()
}

type RootCommand struct {
	BaseCommand
}

func NewRootCommand() *RootCommand {
	return &RootCommand{}
}

func (c *RootCommand) Run() {
	fmt.Println("コマンドのヘルプを見るには --help を付けて実行してください")
}
