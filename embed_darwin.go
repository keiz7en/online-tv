//go:build darwin

package main

import "embed"

//go:embed vlc-linux/*
var vlcFS embed.FS
