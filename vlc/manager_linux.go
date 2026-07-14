//go:build linux

package vlc

import (
	"fmt"
	"io/fs"
	"os/exec"
)

func (m *Manager) Init() {}

func (m *Manager) ExtractVLC(embeddedFS fs.FS) error {
	if _, err := exec.LookPath("vlc"); err != nil {
		return fmt.Errorf("vlc not found in PATH. Install with: sudo apt install vlc")
	}
	fmt.Println("VLC found in system PATH.")
	return nil
}

func (m *Manager) Play(url string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.isPlaying {
		m.stopInternal()
	}

	vlcPath, err := exec.LookPath("vlc")
	if err != nil {
		return fmt.Errorf("vlc not found: %w", err)
	}

	args := []string{
		"--network-caching=1500",
		"--no-video-title-show",
		"--no-osd",
		"--no-stats",
		"--volume", fmt.Sprintf("%d", m.volume*512/100),
		url,
	}

	cmd := exec.Command(vlcPath, args...)

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start VLC: %w", err)
	}

	m.process = cmd.Process
	m.cmd = cmd
	m.currentURL = url
	m.isPlaying = true

	go m.monitorProcess()

	return nil
}

func (m *Manager) cleanupPlatform() {}
