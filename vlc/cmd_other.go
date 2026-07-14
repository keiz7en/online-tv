//go:build !windows

package vlc

import "os/exec"

func applyPlatformOptions(cmd *exec.Cmd) {
}
