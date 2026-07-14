package playlist

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
)

type Channel struct {
	Name     string `json:"name"`
	URL      string `json:"url"`
	Logo     string `json:"logo"`
	Category string `json:"category"`
}

type Playlist struct {
	Channels []Channel `json:"channels"`
}

var (
	extInfRegex = regexp.MustCompile(`#EXTINF:-?\d*\s+(.*)`)
	logoRegex   = regexp.MustCompile(`tvg-logo="([^"]*)"`)
	groupRegex  = regexp.MustCompile(`group-title="([^"]*)"`)
)

func FetchPlaylist(url string) (*Playlist, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch playlist: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("playlist server returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read playlist: %w", err)
	}

	return ParseM3U8(string(body))
}

func ParseM3U8(content string) (*Playlist, error) {
	lines := strings.Split(content, "\n")
	playlist := &Playlist{}

	var pendingInfo string

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "#EXTM3U") {
			continue
		}

		if strings.HasPrefix(line, "#EXTINF:") {
			pendingInfo = line
			continue
		}

		if strings.HasPrefix(line, "#") || line == "" {
			continue
		}

		if pendingInfo != "" {
			channel := parseExtInf(pendingInfo, line)
			playlist.Channels = append(playlist.Channels, channel)
			pendingInfo = ""
		}
	}

	return playlist, nil
}

func parseExtInf(infoLine, url string) Channel {
	ch := Channel{
		URL: url,
	}

	if match := extInfRegex.FindStringSubmatch(infoLine); len(match) > 1 {
		remaining := match[1]

		if idx := strings.LastIndex(remaining, ","); idx != -1 {
			ch.Name = strings.TrimSpace(remaining[idx+1:])
			remaining = remaining[:idx]
		}
	}

	if match := logoRegex.FindStringSubmatch(infoLine); len(match) > 1 {
		ch.Logo = match[1]
	}

	if match := groupRegex.FindStringSubmatch(infoLine); len(match) > 1 {
		ch.Category = match[1]
	}

	if ch.Category == "" {
		ch.Category = "Uncategorized"
	}

	if ch.Name == "" {
		ch.Name = extractNameFromURL(url)
	}

	return ch
}

func extractNameFromURL(url string) string {
	parts := strings.Split(url, "/")
	if len(parts) > 0 {
		last := parts[len(parts)-1]
		if idx := strings.Index(last, "."); idx > 0 {
			return last[:idx]
		}
		return last
	}
	return "Unknown Channel"
}

func GetCategories(channels []Channel) []string {
	categoryMap := make(map[string]bool)
	var categories []string

	for _, ch := range channels {
		if !categoryMap[ch.Category] {
			categoryMap[ch.Category] = true
			categories = append(categories, ch.Category)
		}
	}

	return categories
}

func SearchChannels(channels []Channel, query string) []Channel {
	query = strings.ToLower(strings.TrimSpace(query))
	if query == "" {
		return channels
	}

	var results []Channel
	for _, ch := range channels {
		if strings.Contains(strings.ToLower(ch.Name), query) ||
			strings.Contains(strings.ToLower(ch.Category), query) {
			results = append(results, ch)
		}
	}

	return results
}
