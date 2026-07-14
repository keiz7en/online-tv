//go:build windows

package vlc

import (
	"os/exec"
	"syscall"
)

func applyPlatformOptions(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{
		CreationFlags: 0x08000000,
	}
}
