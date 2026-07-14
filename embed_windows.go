//go:build windows

package main

import "embed"

//go:embed vlc-files/*
var vlcFS embed.FS
