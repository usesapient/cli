package customcli

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewRootCommandShapesCurrentGeneratedCommands(t *testing.T) {
	root, err := NewRootCommand()
	require.NoError(t, err)

	for _, path := range [][]string{
		{"prompts", "config", "retrieve"},
		{"prompts", "topics", "create"},
		{"prompts", "actions", "list"},
		{"api-performance", "targets", "list"},
		{"api-performance", "config", "update"},
		{"api-performance", "evals", "update"},
		{"api-performance", "custom-evals", "create"},
		{"api-performance", "starter-projects", "list"},
		{"api-performance", "actions", "verify"},
	} {
		_, err := findCommand(root, path)
		require.NoError(t, err, "command path %v should exist", path)
	}
}

func TestNewRootCommandRemovesFlattenedGeneratedCommands(t *testing.T) {
	root, err := NewRootCommand()
	require.NoError(t, err)

	for _, path := range [][]string{
		{"prompts", "config-retrieve"},
		{"prompts", "actions-list"},
		{"api-performance", "targets-list"},
		{"api-performance", "custom-evals-create"},
		{"api-performance", "starter-projects-list"},
		{"api-performance", "actions-verify"},
	} {
		_, err := findCommand(root, path)
		require.Error(t, err, "command path %v should have been moved", path)
	}
}
