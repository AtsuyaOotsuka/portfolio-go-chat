package funcs

import (
	"bytes"
	"io"
	"os"
	"testing"
)

func CaptureStdout(t *testing.T, f func()) string {
	t.Helper()

	// 元の stdout を退避
	old := os.Stdout

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe error: %v", err)
	}
	os.Stdout = w

	// 実行
	f()

	// 後始末
	_ = w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)

	return buf.String()
}
