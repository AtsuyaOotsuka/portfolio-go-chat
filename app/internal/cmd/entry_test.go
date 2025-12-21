package cmd

import (
	"testing"

	"github.com/AtsuyaOotsuka/portfolio-go-chat/test_helper/mocks/command_mock"
	"github.com/stretchr/testify/mock"
)

func TestRootSetUp(t *testing.T) {
	c := &Cmd{}
	rootCmd := new(command_mock.RootCommandMock)
	rootCmd.On("Run", mock.Anything).Return()

	c.rootCmd = rootCmd
	c.rootSetUp()

	c.Cmd.Execute()

	rootCmd.AssertCalled(t, "Run", []string{})
}

func TestEntry(t *testing.T) {
	expected := map[string]map[string]any{
		"version":         {"cmd": "version"},
		"room-list":       {"cmd": "room-list"},
		"forbidden-words": {"cmd": "forbidden-words"},
	}

	for name, expect := range expected {
		t.Run(name, func(t *testing.T) {

			c := &Cmd{}
			rootCmd := new(command_mock.RootCommandMock)
			versionCmd := new(command_mock.VersionCommandMock)
			roomListCmd := new(command_mock.RoomListCommandMock)
			forbiddenWordsCmd := new(command_mock.ForbiddenWordsCommandMock)

			versionCmd.On("Run", mock.Anything).Return()
			roomListCmd.On("SetUp", mock.Anything).Return()
			roomListCmd.On("Run", mock.Anything).Return()
			forbiddenWordsCmd.On("SetUp", mock.Anything, mock.Anything).Return()
			forbiddenWordsCmd.On("Run", mock.Anything).Return()

			c.rootCmd = rootCmd
			c.versionCmd = versionCmd
			c.roomListCmd = roomListCmd
			c.forbiddenWordsCmd = forbiddenWordsCmd
			c.rootSetUp()

			c.entry()

			c.Cmd.SetArgs([]string{expect["cmd"].(string)})
			c.Cmd.Execute()

			rootCmd.AssertNotCalled(t, "Run")
			if expect["cmd"] == "version" {
				versionCmd.AssertExpectations(t)
			} else {
				versionCmd.AssertNotCalled(t, "Run")
			}

			if expect["cmd"] == "room-list" {
				roomListCmd.AssertExpectations(t)
			} else {
				roomListCmd.AssertNotCalled(t, "Run")
				roomListCmd.AssertNotCalled(t, "SetUp")
			}

			if expect["cmd"] == "forbidden-words" {
				forbiddenWordsCmd.AssertExpectations(t)
			} else {
				forbiddenWordsCmd.AssertNotCalled(t, "Run")
				forbiddenWordsCmd.AssertNotCalled(t, "SetUp")
			}

		})
	}
}

func TestInitMongo(t *testing.T) {
	c := &Cmd{}
	mongo := c.initMongo()
	if mongo == nil {
		t.Errorf("Expected mongo to be initialized, got nil")
	}
}
