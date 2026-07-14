package main

import (
	"context"
	"embed"
	"fmt"
	"online-tv/playlist"
	"online-tv/vlc"
)

//go:embed vlc-files/*
var vlcFS embed.FS

type App struct {
	ctx      context.Context
	vlc      *vlc.Manager
	playlist *playlist.Playlist
}

func NewApp() *App {
	return &App{
		vlc: vlc.NewManager(),
	}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	if err := a.vlc.ExtractVLC(vlcFS); err != nil {
		fmt.Printf("Warning: %v\n", err)
	}

	go a.fetchPlaylist()
}

func (a *App) fetchPlaylist() {
	p, err := playlist.FetchPlaylist("https://raw.githubusercontent.com/imShakil/tvlink/refs/heads/main/iptv.m3u8")
	if err != nil {
		fmt.Printf("Failed to fetch playlist: %v\n", err)
		return
	}
	a.playlist = p
	fmt.Printf("Loaded %d channels\n", len(p.Channels))
}

func (a *App) GetChannels() []playlist.Channel {
	if a.playlist == nil {
		return nil
	}
	return a.playlist.Channels
}

func (a *App) GetCategories() []string {
	if a.playlist == nil {
		return nil
	}
	return playlist.GetCategories(a.playlist.Channels)
}

func (a *App) SearchChannels(query string) []playlist.Channel {
	if a.playlist == nil {
		return nil
	}
	return playlist.SearchChannels(a.playlist.Channels, query)
}

func (a *App) PlayChannel(url string) error {
	fmt.Printf("Playing channel: %s\n", url)
	err := a.vlc.Play(url)
	if err != nil {
		fmt.Printf("VLC play error: %v\n", err)
	}
	return err
}

func (a *App) StopPlayback() error {
	return a.vlc.Stop()
}

func (a *App) TogglePause() error {
	return a.vlc.Pause()
}

func (a *App) SetVolume(vol int) error {
	return a.vlc.SetVolume(vol)
}

func (a *App) GetVolume() int {
	return a.vlc.GetVolume()
}

func (a *App) IsPlaying() bool {
	return a.vlc.IsPlaying()
}

func (a *App) GetCurrentURL() string {
	return a.vlc.GetCurrentURL()
}

func (a *App) ReloadPlaylist() error {
	p, err := playlist.FetchPlaylist("https://raw.githubusercontent.com/imShakil/tvlink/refs/heads/main/iptv.m3u8")
	if err != nil {
		return err
	}
	a.playlist = p
	return nil
}
