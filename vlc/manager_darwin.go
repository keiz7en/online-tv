//go:build darwin

package vlc

import (
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

const (
	vlcMacURL    = "https://downloads.videolan.org/pub/videolan/vlc/3.0.21/macosx/vlc-3.0.21-universal.dmg"
	vlcExtractDir = "online-tv-vlc"
)

func (m *Manager) Init() {
	m.vlcDir = filepath.Join(os.TempDir(), vlcExtractDir)
}

func (m *Manager) ExtractVLC(embeddedFS fs.FS) error {
	vlcBin := filepath.Join(m.vlcDir, "VLC.app", "Contents", "MacOS", "VLC")
	if _, err := os.Stat(vlcBin); err == nil {
		return nil
	}

	if _, err := exec.LookPath("vlc"); err == nil {
		return nil
	}

	fmt.Println("VLC not found. Downloading VLC for macOS (first time only)...")
	fmt.Println("This may take a minute...")

	if err := m.downloadVLC(); err != nil {
		return fmt.Errorf("failed to setup VLC: %w\nInstall manually: brew install vlc", err)
	}

	fmt.Println("VLC downloaded successfully.")
	return nil
}

func (m *Manager) downloadVLC() error {
	dmgPath := filepath.Join(os.TempDir(), "vlc-macos.dmg")
	defer os.Remove(dmgPath)

	resp, err := http.Get(vlcMacURL)
	if err != nil {
		return fmt.Errorf("download failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed: status %d", resp.StatusCode)
	}

	out, err := os.Create(dmgPath)
	if err != nil {
		return err
	}
	_, err = io.Copy(out, resp.Body)
	out.Close()
	if err != nil {
		return err
	}

	cmd := exec.Command("hdiutil", "attach", dmgPath, "-nobrowse", "-quiet")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("mount failed: %w", err)
	}
	defer exec.Command("hdiutil", "detach", "/Volumes/VLC media player", "-quiet").Run()

	src := "/Volumes/VLC media player/VLC.app"
	dst := filepath.Join(m.vlcDir, "VLC.app")

	if err := os.MkdirAll(m.vlcDir, 0755); err != nil {
		return err
	}

	cmd = exec.Command("cp", "-R", src, dst)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("copy failed: %w", err)
	}

	return nil
}

func (m *Manager) findVLCPath() string {
	vlcBin := filepath.Join(m.vlcDir, "VLC.app", "Contents", "MacOS", "VLC")
	if _, err := os.Stat(vlcBin); err == nil {
		return vlcBin
	}

	if path, err := exec.LookPath("vlc"); err == nil {
		return path
	}

	return ""
}

func (m *Manager) Play(url string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.isPlaying {
		m.stopInternal()
	}

	vlcPath := m.findVLCPath()
	if vlcPath == "" {
		return fmt.Errorf("vlc not found. Install with: brew install vlc")
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
	os.RemoveAll(m.vlcDir)
}
