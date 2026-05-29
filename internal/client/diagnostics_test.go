package client

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestRedactHeadersRedactsSapientAPIKey(t *testing.T) {
	headers := http.Header{
		"X-Sapient-Api-Key": []string{"sap_cli_secret"},
		"X-Example-Api-Key": []string{"example_secret"},
		"Authorization":     []string{"Bearer auth_secret"},
		"Cookie":            []string{"session=cookie_secret"},
		"X-Request-Id":      []string{"request-1"},
	}

	redacted := redactHeaders(headers)

	for _, key := range []string{
		"X-Sapient-Api-Key",
		"X-Example-Api-Key",
		"Authorization",
		"Cookie",
	} {
		if got := redacted.Get(key); got != "[REDACTED]" {
			t.Fatalf("%s = %q, want [REDACTED]", key, got)
		}
	}
	if got := redacted.Get("X-Request-Id"); got != "request-1" {
		t.Fatalf("X-Request-Id = %q, want request-1", got)
	}
}

func TestDryRunClientRedactsSapientAPIKey(t *testing.T) {
	var stderr bytes.Buffer
	req, err := http.NewRequest(
		http.MethodGet,
		"https://api-public.usesapient.com/v1/status",
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("X-Sapient-Api-Key", "sap_cli_test_secret")
	req.Header.Set("X-Request-Id", "request-1")

	resp, err := (&DryRunClient{Stderr: &stderr}).Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	_, _ = io.ReadAll(resp.Body)

	output := stderr.String()
	if strings.Contains(output, "sap_cli_test_secret") {
		t.Fatalf("dry-run output leaked API key: %s", output)
	}
	if !strings.Contains(output, "X-Sapient-Api-Key: [REDACTED]") {
		t.Fatalf("dry-run output did not redact Sapient API key: %s", output)
	}
	if !strings.Contains(output, "X-Request-Id: request-1") {
		t.Fatalf("dry-run output should preserve non-sensitive headers: %s", output)
	}
}
