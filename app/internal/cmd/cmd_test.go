package cmd

import "testing"

func TestNewCmd(t *testing.T) {
	cmd := NewCmd()
	if cmd == nil {
		t.Error("Expected NewCmd() to return a non-nil Cmd instance")
	}
}

func TestCmdInit(t *testing.T) {
	cmd := NewCmd()
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Cmd.Init() panicked: %v", r)
		}
	}()
	cmd.Init()
}
