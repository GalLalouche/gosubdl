package common

import (
	"os"
	"os/exec"
)

// This is such a shitty language jesus christ
func Ptr[T any](v T) *T {
    return &v
}

func ReadChar() rune {
	var b []byte = make([]byte, 1)
	os.Stdin.Read(b)
	return rune(b[0])
}

func AllowReadingSingleChar() {
	// Copied from https://stackoverflow.com/a/17278775/736508
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
}
