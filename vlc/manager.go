package vlc

import (
	"os"
	"os/exec"
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
		volume:    100,
		isPlaying: false,
	}
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
	if m.cmd != nil {
		m.cmd.Wait()
		m.cmd = nil
	}
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
	cmd := m.cmd
	m.mu.Unlock()

	if p == nil {
		return
	}

	cmd.Wait()

	m.mu.Lock()
	if m.process == p {
		m.isPlaying = false
		m.process = nil
		m.cmd = nil
	}
	fn := m.onExit
	m.mu.Unlock()

	if fn != nil {
		fn()
	}
}

func (m *Manager) Cleanup() {
	m.Stop()
	m.cleanupPlatform()
}
