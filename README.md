# Online TV

A desktop application for watching live TV channels. Built with Wails, React, and VLC.

## Features

- Browse and search TV channels by category
- Play live TV channels in VLC
- Volume control
- Dark theme UI

## Installation

### Download

Download the latest release from [Releases](https://github.com/keiz7en/online-tv/releases).

### Build from Source

**Prerequisites:**
- [Go](https://golang.org/dl/) 1.21+
- [Node.js](https://nodejs.org/) 18+
- [Wails](https://wails.io/docs/gettingstarted/installation) v2

```bash
# Clone
git clone https://github.com/keiz7en/online-tv.git
cd online-tv

# Download VLC (required for building)
.\setup-vlc.ps1

# Build
wails build
```

## Tech Stack

- **Backend:** Go + Wails
- **Frontend:** React + Bootstrap
- **Player:** VLC (embedded)
