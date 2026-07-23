//go:build !windows

package util

import "os"

func IsRunningFromCLI() bool {
	fileInfo, _ := os.Stdin.Stat()
	return (fileInfo.Mode() & os.ModeCharDevice) != 0
}
