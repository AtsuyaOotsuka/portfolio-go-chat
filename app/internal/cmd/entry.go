package cmd

import (
	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/command"
	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/usecase"
	"github.com/AtsuyaOotsuka/portfolio-go-chat/public_lib/atylabmongo"
	"github.com/spf13/cobra"
)

func (c *Cmd) bind() {
	c.rootCmd = command.NewRootCommand()
	c.versionCmd = command.NewVersionCommand()
	c.roomListCmd = command.NewRoomListCommand()
}

func (c *Cmd) rootSetUp() {
	c.Cmd = &cobra.Command{
		Use:   "mycli",
		Short: "MyCLI is a sample CLI tool",
		Long:  "MyCLI is a sample CLI tool built with Cobra",
		Run: func(cmd *cobra.Command, args []string) {
			c.rootCmd.Run()
		},
	}
}

func (c *Cmd) entry() {
	c.set(
		"version",
		"Print version",
		func() {
			c.versionCmd.Run()
		},
	)
	c.set(
		"room-list",
		"List all chat rooms",
		func() {
			c.roomListCmd.SetUp(c.initMongo())
			c.roomListCmd.Run()
		},
	)
}

func (c *Cmd) set(
	use string,
	short string,
	run func(),
) {
	c.Cmd.AddCommand(
		&cobra.Command{
			Use:   use,
			Short: short,
			Run: func(cmd *cobra.Command, args []string) {
				run()
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
