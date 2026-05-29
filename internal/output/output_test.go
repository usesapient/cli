package output

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/usesapient/cli/internal/sdk/optionalnullable"
)

type testResponseEnvelope struct {
	Payload interface{}
}

type testOutputPayload struct {
	ID       string                                    `json:"id"`
	SourceID optionalnullable.OptionalNullable[string] `json:"source_id,omitempty"`
	Tags     []string                                  `json:"tags,omitempty"`
	Nested   map[string]string                         `json:"nested,omitempty"`
}

type testListPayload struct {
	Data []testOutputPayload `json:"data"`
	Meta map[string]any      `json:"meta,omitempty"`
}

func testCommand(t *testing.T, format string) (*cobra.Command, *bytes.Buffer) {
	t.Helper()
	var out bytes.Buffer
	cmd := &cobra.Command{}
	cmd.SetOut(&out)
	cmd.Flags().String("output-format", "pretty", "")
	cmd.Flags().String("jq", "", "")
	cmd.Flags().String("color", "never", "")
	cmd.Flags().Bool("include-headers", false, "")
	if format != "" {
		if err := cmd.Flags().Set("output-format", format); err != nil {
			t.Fatal(err)
		}
	}
	return cmd, &out
}

func testPayload() testOutputPayload {
	sourceID := "source-1"
	return testOutputPayload{
		ID:       "target-1",
		SourceID: optionalnullable.From(&sourceID),
		Tags:     []string{"agent", "daily"},
		Nested:   map[string]string{"hidden": "value"},
	}
}

func TestWantsRawJSONKeepsGeneratedCommandsOnTypedResponses(t *testing.T) {
	cmd, _ := testCommand(t, "json")
	if WantsRawJSON(cmd) {
		t.Fatal("json output should not force SDK skip-deserialization")
	}
	if err := cmd.Flags().Set("jq", ".source_id"); err != nil {
		t.Fatal(err)
	}
	if WantsRawJSON(cmd) {
		t.Fatal("jq output should not force SDK skip-deserialization")
	}
}

func TestResultJSONUsesCanonicalTypedSDKShape(t *testing.T) {
	cmd, out := testCommand(t, "json")
	if err := Result(cmd, testResponseEnvelope{Payload: testPayload()}); err != nil {
		t.Fatal(err)
	}

	var got map[string]any
	if err := json.Unmarshal(out.Bytes(), &got); err != nil {
		t.Fatalf("output was not JSON: %v\n%s", err, out.String())
	}
	if got["source_id"] != "source-1" {
		t.Fatalf("source_id was not canonicalized: %#v", got)
	}
	if strings.Contains(out.String(), "true") {
		t.Fatalf("output leaked OptionalNullable internals: %s", out.String())
	}
}

func TestResultYAMLUsesCanonicalTypedSDKShape(t *testing.T) {
	cmd, out := testCommand(t, "yaml")
	if err := Result(cmd, testResponseEnvelope{Payload: testPayload()}); err != nil {
		t.Fatal(err)
	}

	rendered := out.String()
	if !strings.Contains(rendered, "source_id: source-1") {
		t.Fatalf("source_id was not rendered in public YAML shape:\n%s", rendered)
	}
	if strings.Contains(rendered, "true: source-1") {
		t.Fatalf("YAML leaked OptionalNullable internals:\n%s", rendered)
	}
}

func TestResultTableUsesDataEnvelopeRows(t *testing.T) {
	cmd, out := testCommand(t, "table")
	res := testResponseEnvelope{
		Payload: testListPayload{
			Data: []testOutputPayload{testPayload()},
			Meta: map[string]any{"count": 1},
		},
	}
	if err := Result(cmd, res); err != nil {
		t.Fatal(err)
	}

	rendered := out.String()
	for _, expected := range []string{"ID", "SOURCE_ID", "TAGS", "target-1", "source-1", "agent, daily"} {
		if !strings.Contains(rendered, expected) {
			t.Fatalf("table output missing %q:\n%s", expected, rendered)
		}
	}
	if strings.Contains(rendered, "NESTED") {
		t.Fatalf("table output should skip complex nested fields:\n%s", rendered)
	}
}
