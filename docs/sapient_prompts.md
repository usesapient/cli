## sapient prompts

Operations for prompts

### Synopsis

Operations for prompts

```
sapient prompts [flags]
```

### Options

```
  -h, --help   help for prompts
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
* [sapient prompts create](sapient_prompts_create.md)	 - Create Prompt
* [sapient prompts delete](sapient_prompts_delete.md)	 - Delete Prompt
* [sapient prompts estimate-cost](sapient_prompts_estimate-cost.md)	 - Estimate Prompt Cost
* [sapient prompts list](sapient_prompts_list.md)	 - List Prompts
* [sapient prompts platforms-list](sapient_prompts_platforms-list.md)	 - List Platforms
* [sapient prompts retrieve](sapient_prompts_retrieve.md)	 - Retrieve Prompt
* [sapient prompts topics-create](sapient_prompts_topics-create.md)	 - Create Topic
* [sapient prompts topics-delete](sapient_prompts_topics-delete.md)	 - Delete Topic
* [sapient prompts topics-list](sapient_prompts_topics-list.md)	 - List Topics
* [sapient prompts topics-update](sapient_prompts_topics-update.md)	 - Update Topic
* [sapient prompts update](sapient_prompts_update.md)	 - Update Prompt
