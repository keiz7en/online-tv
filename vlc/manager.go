package vlc

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"
)

type Manager struct {
	vlcDir     string
	process    *os.Process
	cmd        *exec.Cmd
	mu         sync.Mutex
	onExit     func()
	isPlaying  bool
	volume     int
	currentURL string
}

func NewManager() *Manager {
	return &Manager{
		vlcDir:  filepath.Join(os.TempDir(), "online-tv-vlc"),
		volume:  100,
		isPlaying: false,
	}
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
		"--network-caching=1500",
		"--no-video-title-show",
		"--no-osd",
		"--no-stats",
		"--volume", fmt.Sprintf("%d", m.volume*512/100),
		url,
	}

	cmd := exec.Command(vlcPath, args...)
	applyPlatformOptions(cmd)

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

func (m *Manager) Stop() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.stopInternal()
}

func (m *Manager) stopInternal() error {
	if m.process != nil {
		m.process.Kill()
		m.process = nil
	}
	m.cmd = nil
	m.isPlaying = false
	m.currentURL = ""
	return nil
}

func (m *Manager) Pause() error {
	return nil
}

func (m *Manager) SetVolume(vol int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if vol < 0 {
		vol = 0
	}
	if vol > 100 {
		vol = 100
	}
	m.volume = vol
	return nil
}

func (m *Manager) GetVolume() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.volume
}

func (m *Manager) Seek(seconds int) error {
	return nil
}

func (m *Manager) IsPlaying() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.isPlaying
}

func (m *Manager) GetCurrentURL() string {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.currentURL
}

func (m *Manager) OnExit(fn func()) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.onExit = fn
}

func (m *Manager) monitorProcess() {
	m.mu.Lock()
	p := m.process
	m.mu.Unlock()

	if p == nil {
		return
	}

	p.Wait()

	m.mu.Lock()
	m.isPlaying = false
	m.process = nil
	m.cmd = nil
	fn := m.onExit
	m.mu.Unlock()

	if fn != nil {
		fn()
	}
}

func (m *Manager) Cleanup() {
	m.Stop()
	if runtime.GOOS == "windows" {
		os.RemoveAll(m.vlcDir)
	}
}

// Ensure io is used to prevent import errors
var _ = io.Copy
