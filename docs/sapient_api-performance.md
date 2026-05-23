## sapient api-performance

Operations for api-performance

### Synopsis

Operations for api-performance

```
sapient api-performance [flags]
```

### Options

```
  -h, --help   help for api-performance
```

### Options inherited from parent commands

```
      --agent-mode                    Enable structured errors and default TOON output for AI coding agents. Automatically enabled when a known agent environment is detected (CLAUDE_CODE, CURSOR_AGENT, etc.). Use --agent-mode=false to disable.
      --color string                  Control colored output: auto (color when output is a TTY), always, or never. Respects NO_COLOR and FORCE_COLOR env vars. (default "auto")
  -d, --debug                         Log request and response diagnostics to stderr
      --dry-run                       Preview the request that would be sent without executing it (output to stderr)
  -H, --header stringArray            Set a custom HTTP request header (format: "Key: Value"). Can be specified multiple times.
      --include-headers               Include HTTP response headers in the output
  -q, --jq string                     Filter and transform output using a jq expression (e.g., '.name', '.items[] | .id')
      --no-interactive                Disable all interactive features (auto-prompting, explorer auto-launch, TUI forms)
  -o, --output-format string          Specify the output format. Options: pretty, json, yaml, table, toon. (default "pretty")
      --sapient-api-key-auth string   API Key
      --server string                 Select a server by index (for indexed servers) or name (for named servers)
      --server-url string             Override the default server URL
      --timeout string                HTTP request timeout (e.g., 30s, 5m, 100ms)
      --usage                         Print the CLI Usage schema in KDL format
```

### SEE ALSO

* [sapient](sapient.md)	 - Sapient Public API: Public API for Sapient CLI and external integrations
* [sapient api-performance evaluation-config-retrieve](sapient_api-performance_evaluation-config-retrieve.md)	 - Retrieve Evaluation Config
* [sapient api-performance evaluation-config-update](sapient_api-performance_evaluation-config-update.md)	 - Update Evaluation Config
* [sapient api-performance interfaces-list](sapient_api-performance_interfaces-list.md)	 - List Interfaces
* [sapient api-performance operation-prompts-create](sapient_api-performance_operation-prompts-create.md)	 - Create Operation Prompt
* [sapient api-performance operation-prompts-list](sapient_api-performance_operation-prompts-list.md)	 - List Operation Prompts
* [sapient api-performance operations-create](sapient_api-performance_operations-create.md)	 - Create Operation
* [sapient api-performance operations-delete](sapient_api-performance_operations-delete.md)	 - Delete Operation
* [sapient api-performance operations-list](sapient_api-performance_operations-list.md)	 - List Operations
* [sapient api-performance operations-retrieve](sapient_api-performance_operations-retrieve.md)	 - Retrieve Operation
* [sapient api-performance operations-update](sapient_api-performance_operations-update.md)	 - Update Operation
* [sapient api-performance platforms-estimate-cost](sapient_api-performance_platforms-estimate-cost.md)	 - Estimate Platform Cost
* [sapient api-performance platforms-list](sapient_api-performance_platforms-list.md)	 - List Platforms
* [sapient api-performance prompts-delete](sapient_api-performance_prompts-delete.md)	 - Delete Operation Prompt
* [sapient api-performance prompts-retrieve](sapient_api-performance_prompts-retrieve.md)	 - Retrieve Operation Prompt
* [sapient api-performance prompts-update](sapient_api-performance_prompts-update.md)	 - Update Operation Prompt
* [sapient api-performance runs-list](sapient_api-performance_runs-list.md)	 - List Operation Runs
* [sapient api-performance runs-retrieve](sapient_api-performance_runs-retrieve.md)	 - Retrieve Run
* [sapient api-performance use-cases-create](sapient_api-performance_use-cases-create.md)	 - Create Use Case
* [sapient api-performance use-cases-delete](sapient_api-performance_use-cases-delete.md)	 - Delete Use Case
* [sapient api-performance use-cases-list](sapient_api-performance_use-cases-list.md)	 - List Use Cases
* [sapient api-performance use-cases-retrieve](sapient_api-performance_use-cases-retrieve.md)	 - Retrieve Use Case
* [sapient api-performance use-cases-update](sapient_api-performance_use-cases-update.md)	 - Update Use Case
