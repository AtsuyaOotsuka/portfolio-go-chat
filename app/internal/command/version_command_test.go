package command

import (
	"testing"

	"github.com/AtsuyaOotsuka/portfolio-go-chat/test_helper/funcs"
)

func TestVersionCmdRun(t *testing.T) {
	cmd := NewVersionCommand()

	outPut := funcs.CaptureStdout(t, func() {
		cmd.Run()
	})

	expected := "mycli version 1.0.0\n"
	if outPut != expected {
		t.Errorf("expected %q but got %q", expected, outPut)
	}
}
