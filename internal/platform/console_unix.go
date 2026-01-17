//go:build !windows

package platform

import "os/exec"

func HideConsoleWindow(cmd *exec.Cmd) {
	// No-op on Unix-like systems
}
