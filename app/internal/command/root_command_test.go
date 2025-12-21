package command

import (
	"testing"

	"github.com/AtsuyaOotsuka/portfolio-go-chat/test_helper/funcs"
)

func TestRootCmdRun(t *testing.T) {
	cmd := NewRootCommand()

	outPut := funcs.CaptureStdout(t, func() {
		cmd.Run([]string{})
	})

	expected := "コマンドのヘルプを見るには --help を付けて実行してください\n"
	if outPut != expected {
		t.Errorf("expected %q but got %q", expected, outPut)
	}
}
