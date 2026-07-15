//go:build windows

package vlc

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

func (m *Manager) Init() {
	m.vlcDir = filepath.Join(os.TempDir(), "online-tv-vlc")
}

func (m *Manager) ExtractVLC(embeddedFS fs.FS) error {
	if _, err := os.Stat(filepath.Join(m.vlcDir, "vlc.exe")); err == nil {
		return nil
	}

	fmt.Println("Extracting VLC player (first time only)...")

	vlcSub, err := fs.Sub(embeddedFS, "vlc-files")
	if err != nil {
		return fmt.Errorf("failed to access embedded VLC directory: %w", err)
	}

	if err := copyFS(vlcSub, m.vlcDir); err != nil {
		return fmt.Errorf("failed to extract VLC: %w", err)
	}

	if _, err := os.Stat(filepath.Join(m.vlcDir, "vlc.exe")); err != nil {
		return fmt.Errorf("vlc.exe not found after extraction")
	}

	fmt.Println("VLC extracted successfully.")
	return nil
}

func copyFS(src fs.FS, destDir string) error {
	return fs.WalkDir(src, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		rel := d.Name()
		if path != "." {
			rel = path
		}

		dest := filepath.Join(destDir, filepath.FromSlash(rel))

		if d.IsDir() {
			return os.MkdirAll(dest, 0755)
		}

		data, err := fs.ReadFile(src, path)
		if err != nil {
			return err
		}

		if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
			return err
		}

		return os.WriteFile(dest, data, 0755)
	})
}

func (m *Manager) Play(url string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.isPlaying {
		m.stopInternal()
	}

	vlcPath := filepath.Join(m.vlcDir, "vlc.exe")
	if _, err := os.Stat(vlcPath); err != nil {
		return fmt.Errorf("vlc.exe not found at %s", vlcPath)
	}

	args := []string{
		"--fullscreen",
		"--network-caching=1500",
		"--no-video-title-show",
		"--no-osd",
		"--no-stats",
		"--volume", fmt.Sprintf("%d", m.volume*512/100),
		url,
	}

	cmd := exec.Command(vlcPath, args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		CreationFlags: 0x08000000,
	}

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

func (m *Manager) cleanupPlatform() {
	if m.vlcDir != "" {
		os.RemoveAll(m.vlcDir)
	}
}
