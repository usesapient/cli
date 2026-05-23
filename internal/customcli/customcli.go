package customcli

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	generatedcli "github.com/usesapient/cli/internal/cli"
	"github.com/usesapient/cli/internal/explorer"
	"github.com/usesapient/cli/internal/output"
	"golang.org/x/term"
)

const (
	simpleCLIAPIKeyEnv = "CLI_API_KEY"
	legacyAPIKeyEnv    = "SAPIENT_API_KEY"
	canonicalAPIKeyEnv = "CLI_SAPIENT_API_KEY_AUTH"
)

type commandMove struct {
	From []string
	To   []string
}

var commandMoves = []commandMove{
	{From: []string{"prompts", "topics-list"}, To: []string{"prompts", "topics", "list"}},
	{From: []string{"prompts", "topics-create"}, To: []string{"prompts", "topics", "create"}},
	{From: []string{"prompts", "topics-update"}, To: []string{"prompts", "topics", "update"}},
	{From: []string{"prompts", "topics-delete"}, To: []string{"prompts", "topics", "delete"}},
	{From: []string{"prompts", "platforms-list"}, To: []string{"prompts", "platforms", "list"}},

	{From: []string{"api-performance", "platforms-list"}, To: []string{"api-performance", "platforms", "list"}},
	{From: []string{"api-performance", "platforms-estimate-cost"}, To: []string{"api-performance", "platforms", "estimate-cost"}},
	{From: []string{"api-performance", "interfaces-list"}, To: []string{"api-performance", "interfaces", "list"}},
	{From: []string{"api-performance", "evaluation-config-retrieve"}, To: []string{"api-performance", "evaluation-config", "retrieve"}},
	{From: []string{"api-performance", "evaluation-config-update"}, To: []string{"api-performance", "evaluation-config", "update"}},
	{From: []string{"api-performance", "operations-list"}, To: []string{"api-performance", "operations", "list"}},
	{From: []string{"api-performance", "operations-create"}, To: []string{"api-performance", "operations", "create"}},
	{From: []string{"api-performance", "operations-retrieve"}, To: []string{"api-performance", "operations", "retrieve"}},
	{From: []string{"api-performance", "operations-update"}, To: []string{"api-performance", "operations", "update"}},
	{From: []string{"api-performance", "operations-delete"}, To: []string{"api-performance", "operations", "delete"}},
	{From: []string{"api-performance", "operation-prompts-list"}, To: []string{"api-performance", "operations", "prompts", "list"}},
	{From: []string{"api-performance", "operation-prompts-create"}, To: []string{"api-performance", "operations", "prompts", "create"}},
	{From: []string{"api-performance", "runs-list"}, To: []string{"api-performance", "runs", "list"}},
	{From: []string{"api-performance", "runs-retrieve"}, To: []string{"api-performance", "runs", "retrieve"}},
	{From: []string{"api-performance", "prompts-retrieve"}, To: []string{"api-performance", "prompts", "retrieve"}},
	{From: []string{"api-performance", "prompts-update"}, To: []string{"api-performance", "prompts", "update"}},
	{From: []string{"api-performance", "prompts-delete"}, To: []string{"api-performance", "prompts", "delete"}},
	{From: []string{"api-performance", "use-cases-list"}, To: []string{"api-performance", "use-cases", "list"}},
	{From: []string{"api-performance", "use-cases-create"}, To: []string{"api-performance", "use-cases", "create"}},
	{From: []string{"api-performance", "use-cases-retrieve"}, To: []string{"api-performance", "use-cases", "retrieve"}},
	{From: []string{"api-performance", "use-cases-update"}, To: []string{"api-performance", "use-cases", "update"}},
	{From: []string{"api-performance", "use-cases-delete"}, To: []string{"api-performance", "use-cases", "delete"}},
}

// NewRootCommand builds the generated Speakeasy command tree and applies the
// Sapient-specific resource grouping without editing generated files.
func NewRootCommand() (*cobra.Command, error) {
	installEnvCompatibilityAliases()

	root, err := generatedcli.NewRootCommand()
	if err != nil {
		return nil, err
	}
	if err := Shape(root); err != nil {
		return nil, err
	}
	RejectPositionalArgs(root)
	InstallRuntimeUsage(root)
	return root, nil
}

func Execute() error {
	root, err := NewRootCommand()
	if err != nil {
		return err
	}

	output.InitAgentMode(root)
	if shouldAutoExplore() {
		if runUpdateCheck(os.Args) {
			return nil
		}
		return runExplorer(root)
	}

	if err := root.Execute(); err != nil {
		return err
	}
	runUpdateCheck(os.Args)
	return nil
}

func installEnvCompatibilityAliases() {
	if os.Getenv(canonicalAPIKeyEnv) != "" {
		return
	}
	if value := os.Getenv(simpleCLIAPIKeyEnv); value != "" {
		_ = os.Setenv(canonicalAPIKeyEnv, value)
		return
	}
	if value := os.Getenv(legacyAPIKeyEnv); value != "" {
		_ = os.Setenv(canonicalAPIKeyEnv, value)
	}
}

func Shape(root *cobra.Command) error {
	for _, move := range commandMoves {
		cmd, err := findCommand(root, move.From)
		if err != nil {
			return err
		}
		parent := cmd.Parent()
		if parent == nil {
			return fmt.Errorf("command %q has no parent", strings.Join(move.From, " "))
		}

		parent.RemoveCommand(cmd)
		cmd.Use = replaceUseName(cmd.Use, move.To[len(move.To)-1])
		cmd.Example = rewriteExample(cmd.Example, move.From, move.To)

		destParent, err := ensureGroupPath(root, move.To[:len(move.To)-1])
		if err != nil {
			return err
		}
		destParent.AddCommand(cmd)
	}
	return nil
}

func RejectPositionalArgs(root *cobra.Command) {
	walkCommands(root, func(cmd *cobra.Command) {
		if cmd.Args == nil {
			cmd.Args = cobra.NoArgs
		}
	})
}

func InstallRuntimeUsage(root *cobra.Command) {
	walkCommands(root, func(cmd *cobra.Command) {
		originalRunE := cmd.RunE
		cmd.RunE = func(cmd *cobra.Command, args []string) error {
			if usageRequested(cmd) {
				return EmitUsage(cmd.Root(), cmd, cmd.OutOrStdout())
			}
			if originalRunE != nil {
				return originalRunE(cmd, args)
			}
			return cmd.Help()
		}
	})
}

func EmitUsage(root *cobra.Command, selected *cobra.Command, w io.Writer) error {
	if selected == root {
		fmt.Fprintf(w, "name %q\n", root.Name())
		fmt.Fprintf(w, "bin %q\n", root.Name())
		fmt.Fprintf(w, "about %q\n", root.Short)
		fmt.Fprintf(w, "version %q\n", generatedcli.Version)
		fmt.Fprintln(w, "config {")
		fmt.Fprintln(w, "  file \"~/.config/sapient/config.yaml\"")
		fmt.Fprintln(w, "}")
		writeFlags(w, root.PersistentFlags(), "flag", true, 0)
		for _, child := range root.Commands() {
			if child.IsAvailableCommand() && !child.IsAdditionalHelpTopicCommand() {
				writeCommand(w, child, 0)
			}
		}
		return nil
	}

	writeCommand(w, selected, 0)
	return nil
}

func writeCommand(w io.Writer, cmd *cobra.Command, indent int) {
	prefix := strings.Repeat("  ", indent)
	children := availableChildren(cmd)
	hasBody := len(children) > 0 || hasVisibleFlags(cmd.Flags()) || len(cmd.Aliases) > 0

	if hasBody {
		fmt.Fprintf(w, "%scmd %q help=%q {\n", prefix, cmd.Name(), cmd.Short)
		for _, alias := range cmd.Aliases {
			fmt.Fprintf(w, "%s  alias %q\n", prefix, alias)
		}
		writeFlags(w, cmd.LocalFlags(), "flag", false, indent+1)
		for _, child := range children {
			writeCommand(w, child, indent+1)
		}
		fmt.Fprintf(w, "%s}\n", prefix)
		return
	}

	fmt.Fprintf(w, "%scmd %q help=%q\n", prefix, cmd.Name(), cmd.Short)
}

func writeFlags(w io.Writer, flags *pflag.FlagSet, keyword string, global bool, indent int) {
	if flags == nil {
		return
	}
	prefix := strings.Repeat("  ", indent)
	flags.VisitAll(func(flag *pflag.Flag) {
		if flag.Hidden || flag.Name == "help" {
			return
		}
		name := "--" + flag.Name
		if flag.Shorthand != "" {
			name = "-" + flag.Shorthand + " " + name
		}
		if flag.Value.Type() != "bool" {
			name += " <" + strings.ReplaceAll(flag.Name, "-", "_") + ">"
		}
		fmt.Fprintf(w, "%s%s %q help=%q", prefix, keyword, name, flag.Usage)
		if global {
			fmt.Fprint(w, " global=#true")
		}
		if flag.DefValue != "" && flag.DefValue != "[]" {
			fmt.Fprintf(w, " default=%q", flag.DefValue)
		}
		fmt.Fprintln(w)
	})
}

func hasVisibleFlags(flags *pflag.FlagSet) bool {
	if flags == nil {
		return false
	}
	found := false
	flags.VisitAll(func(flag *pflag.Flag) {
		if !flag.Hidden && flag.Name != "help" {
			found = true
		}
	})
	return found
}

func availableChildren(cmd *cobra.Command) []*cobra.Command {
	children := []*cobra.Command{}
	for _, child := range cmd.Commands() {
		if child.IsAvailableCommand() && !child.IsAdditionalHelpTopicCommand() {
			children = append(children, child)
		}
	}
	return children
}

func findCommand(root *cobra.Command, path []string) (*cobra.Command, error) {
	cur := root
	for _, name := range path {
		next := findChild(cur, name)
		if next == nil {
			return nil, fmt.Errorf("command path %q not found", strings.Join(path, " "))
		}
		cur = next
	}
	return cur, nil
}

func findChild(parent *cobra.Command, name string) *cobra.Command {
	for _, child := range parent.Commands() {
		if child.Name() == name {
			return child
		}
	}
	return nil
}

func ensureGroupPath(root *cobra.Command, path []string) (*cobra.Command, error) {
	cur := root
	for _, name := range path {
		child := findChild(cur, name)
		if child == nil {
			child = newGroupCommand(name)
			cur.AddCommand(child)
		}
		cur = child
	}
	return cur, nil
}

func newGroupCommand(name string) *cobra.Command {
	return &cobra.Command{
		Use:           name,
		Short:         "Operations for " + strings.ReplaceAll(name, "-", " "),
		Long:          "Operations for " + strings.ReplaceAll(name, "-", " "),
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if usageRequested(cmd) {
				return EmitUsage(cmd.Root(), cmd, cmd.OutOrStdout())
			}
			return cmd.Help()
		},
	}
}

func replaceUseName(use string, name string) string {
	parts := strings.SplitN(use, " ", 2)
	if len(parts) == 1 {
		return name
	}
	return name + " " + parts[1]
}

func rewriteExample(example string, from []string, to []string) string {
	if example == "" {
		return ""
	}
	oldPath := "sapient " + strings.Join(from, " ")
	newPath := "sapient " + strings.Join(to, " ")
	return strings.ReplaceAll(example, oldPath, newPath)
}

func walkCommands(cmd *cobra.Command, fn func(*cobra.Command)) {
	fn(cmd)
	for _, child := range cmd.Commands() {
		walkCommands(child, fn)
	}
}

func usageRequested(cmd *cobra.Command) bool {
	if cmd == nil {
		return false
	}
	for _, flags := range []*pflag.FlagSet{cmd.Flags(), cmd.InheritedFlags()} {
		if flags == nil {
			continue
		}
		flag := flags.Lookup("usage")
		if flag == nil {
			continue
		}
		value, err := strconv.ParseBool(flag.Value.String())
		if err == nil && value {
			return true
		}
	}
	return false
}

func shouldAutoExplore() bool {
	if len(os.Args) > 1 {
		return false
	}
	if !term.IsTerminal(int(os.Stdin.Fd())) || !term.IsTerminal(int(os.Stdout.Fd())) {
		return false
	}
	return !output.IsAgentMode()
}

func runExplorer(root *cobra.Command) error {
	selectedArgs, err := explorer.Run(root, generatedcli.Version)
	if err != nil {
		return err
	}
	if selectedArgs == nil {
		return nil
	}

	freshRoot, err := NewRootCommand()
	if err != nil {
		return err
	}
	freshRoot.SetArgs(selectedArgs)
	return freshRoot.Execute()
}
