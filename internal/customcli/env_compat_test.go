package customcli

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/usesapient/cli/internal/config"
	"github.com/usesapient/cli/internal/output"
)

func TestInstallEnvCompatibilityAliasesCopiesLegacyAPIKey(t *testing.T) {
	t.Setenv(legacyAPIKeyEnv, "legacy-secret")
	t.Setenv(canonicalAPIKeyEnv, "")

	installEnvCompatibilityAliases()

	require.Equal(t, "legacy-secret", os.Getenv(canonicalAPIKeyEnv))
}

func TestInstallEnvCompatibilityAliasesKeepsCanonicalAPIKey(t *testing.T) {
	t.Setenv(legacyAPIKeyEnv, "legacy-secret")
	t.Setenv(canonicalAPIKeyEnv, "canonical-secret")

	installEnvCompatibilityAliases()

	require.Equal(t, "canonical-secret", os.Getenv(canonicalAPIKeyEnv))
}

func TestLegacyAPIKeyEnvResolvesAsEnvCredential(t *testing.T) {
	t.Setenv(legacyAPIKeyEnv, "legacy-secret")
	t.Setenv(canonicalAPIKeyEnv, "")
	config.Reset()
	output.ResetAgentMode()
	t.Cleanup(func() {
		config.Reset()
		output.ResetAgentMode()
	})

	root, err := NewRootCommand()
	require.NoError(t, err)

	var stdout bytes.Buffer
	root.SetOut(&stdout)
	root.SetArgs([]string{"whoami"})

	require.NoError(t, root.Execute())
	require.Contains(t, stdout.String(), "--sapient-api-key-auth")
	require.Contains(t, stdout.String(), "[env")
}
