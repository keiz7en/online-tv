//go:build linux

package main

import "embed"

//go:embed vlc-linux/*
var vlcFS embed.FS
