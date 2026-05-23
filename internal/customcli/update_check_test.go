package customcli

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	generatedcli "github.com/usesapient/cli/internal/cli"
	"github.com/usesapient/cli/internal/output"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (fn roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return fn(req)
}

func withUpdateCheckTestEnv(t *testing.T) (*bytes.Buffer, string) {
	t.Helper()

	previousHTTPClient := updateCheckHTTPClient
	previousNow := updateCheckNow
	previousUserCacheDir := updateCheckUserCacheDir
	previousErrWriter := updateCheckErrWriter
	previousIsTerminal := updateCheckIsTerminal
	previousDisable := os.Getenv(updateCheckDisableEnv)
	previousCI := os.Getenv("CI")

	cacheDir := t.TempDir()
	var stderr bytes.Buffer

	updateCheckNow = func() time.Time {
		return time.Date(2026, 5, 16, 0, 0, 0, 0, time.UTC)
	}
	updateCheckUserCacheDir = func() (string, error) {
		return cacheDir, nil
	}
	updateCheckErrWriter = &stderr
	updateCheckIsTerminal = func() bool {
		return true
	}
	t.Setenv(updateCheckDisableEnv, "")
	t.Setenv("CI", "")
	output.ResetAgentMode()

	t.Cleanup(func() {
		updateCheckHTTPClient = previousHTTPClient
		updateCheckNow = previousNow
		updateCheckUserCacheDir = previousUserCacheDir
		updateCheckErrWriter = previousErrWriter
		updateCheckIsTerminal = previousIsTerminal
		restoreEnv(updateCheckDisableEnv, previousDisable)
		restoreEnv("CI", previousCI)
		output.ResetAgentMode()
	})

	return &stderr, cacheDir
}

func TestRunUpdateCheckUsesFreshCacheWithoutNetwork(t *testing.T) {
	stderr, cacheDir := withUpdateCheckTestEnv(t)
	cachePath := filepath.Join(cacheDir, "sapient", "update-check.json")
	writeUpdateCheckCache(cachePath, updateCheckCache{
		CheckedAt:     updateCheckNow().Add(-time.Hour),
		LatestVersion: "9.9.9",
	})
	updateCheckHTTPClient = &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			t.Fatal("unexpected network request for fresh update cache")
			return nil, nil
		}),
	}

	notified := runUpdateCheck([]string{"sapient", "status", "get"})

	require.True(t, notified)
	require.Contains(t, stderr.String(), "A new Sapient CLI version is available: "+generatedcli.Version+" -> 9.9.9")
	require.Contains(t, stderr.String(), "Update with: brew upgrade sapient")
	cached, ok := readUpdateCheckCache(cachePath)
	require.True(t, ok)
	require.Equal(t, updateCheckNow(), cached.NotifiedAt)
}

func TestRunUpdateCheckDoesNotRepeatFreshCacheNotification(t *testing.T) {
	stderr, cacheDir := withUpdateCheckTestEnv(t)
	cachePath := filepath.Join(cacheDir, "sapient", "update-check.json")
	writeUpdateCheckCache(cachePath, updateCheckCache{
		CheckedAt:     updateCheckNow().Add(-time.Hour),
		LatestVersion: "9.9.9",
		NotifiedAt:    updateCheckNow().Add(-time.Hour),
	})
	updateCheckHTTPClient = &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			t.Fatal("unexpected network request for fresh update cache")
			return nil, nil
		}),
	}

	notified := runUpdateCheck([]string{"sapient", "status", "get"})

	require.False(t, notified)
	require.Empty(t, stderr.String())
}

func TestRunUpdateCheckRefreshesStaleCache(t *testing.T) {
	stderr, cacheDir := withUpdateCheckTestEnv(t)
	cachePath := filepath.Join(cacheDir, "sapient", "update-check.json")
	writeUpdateCheckCache(cachePath, updateCheckCache{
		CheckedAt:     updateCheckNow().Add(-25 * time.Hour),
		LatestVersion: "0.0.1",
	})
	requested := false
	updateCheckHTTPClient = &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			requested = true
			require.Equal(t, http.MethodGet, req.Method)
			require.Equal(t, updateCheckLatestReleaseURL, req.URL.String())
			require.Equal(t, "Sapient/CLI "+generatedcli.Version, req.Header.Get("User-Agent"))
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`{"tag_name":"v9.9.9"}`)),
				Header:     make(http.Header),
			}, nil
		}),
	}

	notified := runUpdateCheck([]string{"sapient", "status", "get"})

	require.True(t, requested)
	require.True(t, notified)
	require.Contains(t, stderr.String(), generatedcli.Version+" -> 9.9.9")
	cached, ok := readUpdateCheckCache(cachePath)
	require.True(t, ok)
	require.Equal(t, "9.9.9", cached.LatestVersion)
	require.Equal(t, updateCheckNow(), cached.CheckedAt)
	require.Equal(t, updateCheckNow(), cached.NotifiedAt)
}

func TestRunUpdateCheckSkipsWhenDisabled(t *testing.T) {
	stderr, _ := withUpdateCheckTestEnv(t)
	t.Setenv(updateCheckDisableEnv, "1")
	updateCheckHTTPClient = &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			t.Fatal("unexpected network request when update checks are disabled")
			return nil, nil
		}),
	}

	notified := runUpdateCheck([]string{"sapient", "status", "get"})

	require.False(t, notified)
	require.Empty(t, stderr.String())
}

func TestRunUpdateCheckSkipsInCI(t *testing.T) {
	stderr, _ := withUpdateCheckTestEnv(t)
	t.Setenv("CI", "true")
	updateCheckHTTPClient = &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			t.Fatal("unexpected network request in CI")
			return nil, nil
		}),
	}

	notified := runUpdateCheck([]string{"sapient", "status", "get"})

	require.False(t, notified)
	require.Empty(t, stderr.String())
}

func TestRunUpdateCheckSkipsNonTTYAndUtilityCommands(t *testing.T) {
	stderr, _ := withUpdateCheckTestEnv(t)
	updateCheckHTTPClient = &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			t.Fatal("unexpected network request for skipped update checks")
			return nil, nil
		}),
	}

	updateCheckIsTerminal = func() bool { return false }
	require.False(t, runUpdateCheck([]string{"sapient", "status", "get"}))
	updateCheckIsTerminal = func() bool { return true }
	require.False(t, runUpdateCheck([]string{"sapient", "version"}))
	require.False(t, runUpdateCheck([]string{"sapient", "--version"}))
	require.False(t, runUpdateCheck([]string{"sapient", "completion", "zsh"}))
	require.False(t, runUpdateCheck([]string{"sapient", "__complete", "prompts"}))
	require.False(t, runUpdateCheck([]string{"sapient", "@completion", "zsh"}))
	require.False(t, runUpdateCheck([]string{"sapient", "@manpages"}))
	require.False(t, runUpdateCheck([]string{"sapient", "--usage"}))
	require.False(t, runUpdateCheck([]string{"sapient", "status", "get", "--help"}))
	require.False(t, runUpdateCheck([]string{"sapient", "status", "get", "--agent-mode"}))

	require.Empty(t, stderr.String())
}

func TestRunUpdateCheckDoesNotNotifyForSameOlderOrInvalidVersion(t *testing.T) {
	for _, latestVersion := range []string{generatedcli.Version, "v0.0.1", "not-a-version"} {
		t.Run(latestVersion, func(t *testing.T) {
			stderr, _ := withUpdateCheckTestEnv(t)
			updateCheckHTTPClient = &http.Client{
				Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(strings.NewReader(`{"tag_name":"` + latestVersion + `"}`)),
						Header:     make(http.Header),
					}, nil
				}),
			}

			notified := runUpdateCheck([]string{"sapient", "status", "get"})

			require.False(t, notified)
			require.Empty(t, stderr.String())
		})
	}
}

func TestIsNewerVersion(t *testing.T) {
	require.True(t, isNewerVersion("v1.2.4", "1.2.3"))
	require.True(t, isNewerVersion("1.3.0", "1.2.9"))
	require.True(t, isNewerVersion("2.0.0", "1.9.9"))
	require.False(t, isNewerVersion("1.2.3", "1.2.3"))
	require.False(t, isNewerVersion("1.2.2", "1.2.3"))
	require.False(t, isNewerVersion("not-a-version", "1.2.3"))
	require.False(t, isNewerVersion("1.2.3", "not-a-version"))
}

func restoreEnv(key string, value string) {
	if value == "" {
		_ = os.Unsetenv(key)
		return
	}
	_ = os.Setenv(key, value)
}
