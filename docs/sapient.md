## sapient

Sapient Public API: Public API for Sapient CLI and external integrations

### Synopsis

Sapient Public API: Public API for Sapient CLI and external integrations.

```
sapient [flags]
```

### Options

```
      --agent-mode                    Enable structured errors and default TOON output for AI coding agents. Automatically enabled when a known agent environment is detected (CLAUDE_CODE, CURSOR_AGENT, etc.). Use --agent-mode=false to disable.
      --color string                  Control colored output: auto (color when output is a TTY), always, or never. Respects NO_COLOR and FORCE_COLOR env vars. (default "auto")
  -d, --debug                         Log request and response diagnostics to stderr
      --dry-run                       Preview the request that would be sent without executing it (output to stderr)
  -H, --header stringArray            Set a custom HTTP request header (format: "Key: Value"). Can be specified multiple times.
  -h, --help                          help for sapient
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

* [sapient api-performance](sapient_api-performance.md)	 - Operations for api-performance
* [sapient auth](sapient_auth.md)	 - Manage authentication credentials
* [sapient configure](sapient_configure.md)	 - Configure authentication credentials and preferences
* [sapient explore](sapient_explore.md)	 - Interactively browse and run commands
* [sapient prompts](sapient_prompts.md)	 - Operations for prompts
* [sapient status](sapient_status.md)	 - Operations for status
* [sapient version](sapient_version.md)	 - Print the CLI version
* [sapient whoami](sapient_whoami.md)	 - Display current authentication configuration
