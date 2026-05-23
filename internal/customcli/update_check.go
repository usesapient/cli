package customcli

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	generatedcli "github.com/usesapient/cli/internal/cli"
	"github.com/usesapient/cli/internal/output"
)

const (
	updateCheckDisableEnv = "SAPIENT_NO_UPDATE_CHECK"
	updateCheckInterval   = 24 * time.Hour
	updateCheckTimeout    = 500 * time.Millisecond
)

var (
	updateCheckLatestReleaseURL           = "https://api.github.com/repos/usesapient/cli/releases/latest"
	updateCheckHTTPClient                 = &http.Client{Timeout: updateCheckTimeout}
	updateCheckNow                        = time.Now
	updateCheckUserCacheDir               = os.UserCacheDir
	updateCheckErrWriter        io.Writer = os.Stderr
	updateCheckIsTerminal                 = stderrIsTerminal
)

type updateCheckCache struct {
	CheckedAt     time.Time `json:"checked_at"`
	LatestVersion string    `json:"latest_version"`
	NotifiedAt    time.Time `json:"notified_at,omitempty"`
}

func runUpdateCheck(args []string) {
	if shouldSkipUpdateCheck(args) {
		return
	}

	cachePath, err := updateCheckCachePath()
	if err != nil {
		return
	}

	now := updateCheckNow()
	cached, ok := readUpdateCheckCache(cachePath)
	if ok && now.Sub(cached.CheckedAt) < updateCheckInterval {
		writeUpdateNotificationIfNeeded(cachePath, cached, now)
		return
	}

	latestVersion, err := fetchLatestReleaseVersion(context.Background())
	if err != nil {
		return
	}

	nextCache := updateCheckCache{
		CheckedAt:     now,
		LatestVersion: latestVersion,
	}
	if ok && cached.LatestVersion == latestVersion {
		nextCache.NotifiedAt = cached.NotifiedAt
	}
	writeUpdateNotificationIfNeeded(cachePath, nextCache, now)
}

func shouldSkipUpdateCheck(args []string) bool {
	if envTruthy(os.Getenv(updateCheckDisableEnv)) || envTruthy(os.Getenv("CI")) {
		return true
	}
	if output.IsAgentMode() {
		return true
	}
	if !updateCheckIsTerminal() {
		return true
	}
	for _, arg := range args[1:] {
		switch arg {
		case "--version", "-v", "version", "completion", "__complete", "@completion", "@manpages", "help", "--help", "-h", "--usage":
			return true
		}
		if strings.HasPrefix(arg, "--agent-mode") {
			return true
		}
	}
	return false
}

func fetchLatestReleaseVersion(ctx context.Context) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, updateCheckLatestReleaseURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "Sapient/CLI "+generatedcli.Version)

	resp, err := updateCheckHTTPClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected update check status: %d", resp.StatusCode)
	}

	var payload struct {
		TagName string `json:"tag_name"`
	}
	if err := json.NewDecoder(io.LimitReader(resp.Body, 1<<20)).Decode(&payload); err != nil {
		return "", err
	}
	version, ok := normalizeVersion(payload.TagName)
	if !ok {
		return "", errors.New("latest release tag is not a semantic version")
	}
	return version, nil
}

func updateCheckCachePath() (string, error) {
	cacheDir, err := updateCheckUserCacheDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(cacheDir, "sapient", "update-check.json"), nil
}

func readUpdateCheckCache(path string) (updateCheckCache, bool) {
	data, err := os.ReadFile(path)
	if err != nil {
		return updateCheckCache{}, false
	}
	var cached updateCheckCache
	if err := json.Unmarshal(data, &cached); err != nil {
		return updateCheckCache{}, false
	}
	if cached.CheckedAt.IsZero() || cached.LatestVersion == "" {
		return updateCheckCache{}, false
	}
	return cached, true
}

func writeUpdateCheckCache(path string, cached updateCheckCache) {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return
	}
	data, err := json.MarshalIndent(cached, "", "  ")
	if err != nil {
		return
	}
	_ = os.WriteFile(path, data, 0644)
}

func writeUpdateNotificationIfNeeded(path string, cached updateCheckCache, now time.Time) {
	if !cached.NotifiedAt.IsZero() && now.Sub(cached.NotifiedAt) < updateCheckInterval {
		writeUpdateCheckCache(path, cached)
		return
	}
	if writeUpdateNotification(cached.LatestVersion) {
		cached.NotifiedAt = now
	}
	writeUpdateCheckCache(path, cached)
}

func writeUpdateNotification(latestVersion string) bool {
	latestVersion, ok := normalizeVersion(latestVersion)
	if !ok || !isNewerVersion(latestVersion, generatedcli.Version) {
		return false
	}
	fmt.Fprintf(
		updateCheckErrWriter,
		"A new Sapient CLI version is available: %s -> %s\nUpdate with: brew upgrade sapient\n",
		generatedcli.Version,
		latestVersion,
	)
	return true
}

func isNewerVersion(latest string, current string) bool {
	latestVersion, ok := parseSemver(latest)
	if !ok {
		return false
	}
	currentVersion, ok := parseSemver(current)
	if !ok {
		return false
	}
	for i := range latestVersion {
		if latestVersion[i] > currentVersion[i] {
			return true
		}
		if latestVersion[i] < currentVersion[i] {
			return false
		}
	}
	return false
}

func normalizeVersion(version string) (string, bool) {
	parsed, ok := parseSemver(version)
	if !ok {
		return "", false
	}
	return fmt.Sprintf("%d.%d.%d", parsed[0], parsed[1], parsed[2]), true
}

func parseSemver(version string) ([3]int, bool) {
	version = strings.TrimSpace(version)
	version = strings.TrimPrefix(version, "v")
	version = strings.SplitN(version, "+", 2)[0]
	version = strings.SplitN(version, "-", 2)[0]
	parts := strings.Split(version, ".")
	if len(parts) != 3 {
		return [3]int{}, false
	}
	var parsed [3]int
	for i, part := range parts {
		if part == "" {
			return [3]int{}, false
		}
		value, err := strconv.Atoi(part)
		if err != nil || value < 0 {
			return [3]int{}, false
		}
		parsed[i] = value
	}
	return parsed, true
}

func envTruthy(value string) bool {
	value = strings.TrimSpace(strings.ToLower(value))
	return value != "" && value != "0" && value != "false" && value != "no"
}

func stderrIsTerminal() bool {
	info, err := os.Stderr.Stat()
	if err != nil {
		return false
	}
	return info.Mode()&os.ModeCharDevice != 0
}
