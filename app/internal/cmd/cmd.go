package cmd

import (
	"github.com/AtsuyaOotsuka/portfolio-go-chat/internal/command"
	"github.com/spf13/cobra"
)

type Cmd struct {
	Cmd *cobra.Command

	rootCmd           command.RootCommandInterface
	versionCmd        command.VersionCommandInterface
	roomListCmd       command.RoomListCommandInterface
	forbiddenWordsCmd command.ForbiddenWordsCommandInterface
}

func NewCmd() *Cmd {
	return &Cmd{}
}

func (c *Cmd) Init() {
	c.rootSetUp()
	c.setupFlags()
	c.bind()
	c.entry()
}
