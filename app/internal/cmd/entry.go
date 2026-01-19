package cmd

import (
	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/command"
	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/usecase"
	"github.com/AtsuyaOotsuka/portfolio-go-lib/atylabmongo"
	"github.com/spf13/cobra"
)

func (c *Cmd) bind() {
	c.rootCmd = command.NewRootCommand()
	c.versionCmd = command.NewVersionCommand()
	c.roomListCmd = command.NewRoomListCommand()
	c.forbiddenWordsCmd = command.NewForbiddenWordsCommand()
}

func (c *Cmd) rootSetUp() {
	c.Cmd = &cobra.Command{
		Use:   "mycli",
		Short: "MyCLI is a sample CLI tool",
		Long:  "MyCLI is a sample CLI tool built with Cobra",
		Run: func(cmd *cobra.Command, args []string) {
			c.rootCmd.Run(args)
		},
	}
}

func (c *Cmd) entry() {
	c.set(
		"version",
		"Print version",
		func(args []string) {
			c.versionCmd.Run(args)
		},
	)
	c.set(
		"room-list",
		"List all chat rooms",
		func(args []string) {
			c.roomListCmd.SetUp(c.initMongo())
			c.roomListCmd.Run(args)
		},
	)
	c.set(
		"forbidden-words",
		"Manage forbidden words",
		func(args []string) {
			c.forbiddenWordsCmd.SetUp(c.initMongo(), 100)
			c.forbiddenWordsCmd.Run(args)
		},
	)
}

func (c *Cmd) set(
	use string,
	short string,
	run func(args []string),
) {
	c.Cmd.AddCommand(
		&cobra.Command{
			Use:   use,
			Short: short,
			Run: func(cmd *cobra.Command, args []string) {
				run(args)
			},
		},
	)
}

func (c *Cmd) initMongo() *usecase.MongoUseCaseStruct {
	mongo := usecase.NewMongo()

	return usecase.NewMongoUseCaseStruct(
		atylabmongo.NewMongoConnectionStruct(),
		mongo,
	)
}
