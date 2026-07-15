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

**Windows:** Extract `onlinetv.zip` and run `onlinetv.exe`. VLC is embedded.

**Linux (.deb):**
```bash
sudo apt install vlc
sudo dpkg -i onlinetv_1.0.0_amd64.deb
```

**Linux (manual):**
```bash
sudo apt install vlc
chmod +x onlinetv-linux
./onlinetv-linux
```

**macOS:**
```bash
brew install vlc
unzip onlinetv-macos.zip
cd onlinetv-macos
chmod +x sign-macos.sh
./sign-macos.sh
open "Online TV.app"
```

### Build from Source

**Prerequisites:**
- [Go](https://golang.org/dl/) 1.21+
- [Node.js](https://nodejs.org/) 18+
- [Wails](https://wails.io/docs/gettingstarted/installation) v2

```bash
# Clone
git clone https://github.com/keiz7en/online-tv.git
cd online-tv

# Download VLC (Windows only, for building)
.\setup-vlc.ps1

# Build
wails build
```

## Tech Stack

- **Backend:** Go + Wails
- **Frontend:** React + Bootstrap
- **Player:** VLC
