# sapient

Command-line interface for Sapient.

[![Built by Speakeasy](https://img.shields.io/badge/Built_by-SPEAKEASY-374151?style=for-the-badge&labelColor=f3f4f6)](https://www.speakeasy.com/?utm_source=github-com/usesapient/cli&utm_campaign=cli)
[![License: MIT](https://img.shields.io/badge/LICENSE_//_MIT-3b5bdb?style=for-the-badge&labelColor=eff6ff)](https://opensource.org/licenses/MIT)

<!-- Start Summary [summary] -->
## Summary

Manage Sapient prompts, custom evals, and API Performance.
<!-- End Summary [summary] -->

<!-- Start Table of Contents [toc] -->
## Table of Contents
<!-- $toc-max-depth=2 -->
* [sapient](#sapient)
  * [CLI Installation](#cli-installation)
  * [Shell Completion](#shell-completion)
  * [CLI Example Usage](#cli-example-usage)
  * [Authentication](#authentication)
  * [Available Commands](#available-commands)
  * [Request Body Input](#request-body-input)
  * [Server Selection](#server-selection)
  * [Output Formats](#output-formats)
  * [Error Handling](#error-handling)
  * [Diagnostics](#diagnostics)
* [Development](#development)
  * [Maturity](#maturity)
  * [Contributions](#contributions)

<!-- End Table of Contents [toc] -->

<!-- Start CLI Installation [installation] -->
## CLI Installation

### Quick Install (Linux/macOS)

```bash
curl -fsSL https://raw.githubusercontent.com/usesapient/cli/main/scripts/install.sh | bash
```

### Quick Install (Windows PowerShell)

```powershell
iwr -useb https://raw.githubusercontent.com/usesapient/cli/main/scripts/install.ps1 | iex
```
### Homebrew (macOS/Linux)

```bash
brew install usesapient/tools/sapient
```

### Go Install

Alternatively, install directly via Go:

```bash
go install github.com/usesapient/cli/cmd/sapient@latest
```

### Manual Download

Download pre-built binaries for your platform from the [releases page](https://github.com/usesapient/cli/releases).
<!-- End CLI Installation [installation] -->

<!-- Start Shell Completion [completion] -->
## Shell Completion

Shell completions are available for Bash, Zsh, Fish, and PowerShell.

### Bash

```bash
# Add to ~/.bashrc:
source <(sapient completion bash)

# Or install permanently:
sapient completion bash > /etc/bash_completion.d/sapient
```

### Zsh

```zsh
# Add to ~/.zshrc:
source <(sapient completion zsh)

# Or install permanently:
sapient completion zsh > "${fpath[1]}/_sapient"
```

### Fish

```fish
sapient completion fish | source

# Or install permanently:
sapient completion fish > ~/.config/fish/completions/sapient.fish
```

### PowerShell

```powershell
sapient completion powershell | Out-String | Invoke-Expression
```
<!-- End Shell Completion [completion] -->

<!-- Start CLI Example Usage [usage] -->
## CLI Example Usage

### Example

```bash
sapient status get --sapient-api-key-auth test_api_key

```
<!-- End CLI Example Usage [usage] -->

<!-- Start Authentication [security] -->
## Authentication

Authentication credentials can be configured in four ways (in order of priority):

### 1. Command-line flags

Pass credentials directly as flags to any command:

```bash
sapient --sapient-api-key-auth <value> <command> [arguments]
```

### 2. Environment variables

Set credentials via environment variables:

| Variable | Description |
|----------|-------------|
| `CLI_SAPIENT_API_KEY_AUTH` | API Key |

### 3. OS Keychain (recommended for workstations)

Credentials are stored securely in your operating system's keychain when you run:

```bash
sapient configure
```

Secret credentials (tokens, API keys, passwords) are automatically stored in:
- **macOS**: Keychain
- **Linux**: GNOME Keyring / KWallet (via D-Bus Secret Service)
- **Windows**: Windows Credential Locker

If no keychain is available (e.g., in CI environments), credentials fall back to the config file.

### 4. Configuration file

Run the interactive `configure` command to store non-secret settings:

```bash
sapient configure
```

Configuration is stored in `~/.config/sapient/config.yaml`.
<!-- End Authentication [security] -->

<!-- Start Available Commands [operations] -->
## Available Commands

<details open>
<summary>Available commands</summary>

### [status](docs/sapient_status.md)

* [`get`](docs/sapient_status_get.md) - Get Status
* [`auth-status-get`](docs/sapient_status_auth-status-get.md) - Get Auth Status

### [prompts](docs/sapient_prompts.md)

* [`list`](docs/sapient_prompts_list.md) - List Prompts
* [`create`](docs/sapient_prompts_create.md) - Create Prompts
* [`config-retrieve`](docs/sapient_prompts_config-retrieve.md) - Retrieve Prompt Config
* [`topics-create`](docs/sapient_prompts_topics-create.md) - Create Topic
* [`topics-update`](docs/sapient_prompts_topics-update.md) - Update Topic
* [`topics-delete`](docs/sapient_prompts_topics-delete.md) - Delete Topic
* [`actions-list`](docs/sapient_prompts_actions-list.md) - List Prompt Actionables
* [`retrieve`](docs/sapient_prompts_retrieve.md) - Retrieve Prompt
* [`update`](docs/sapient_prompts_update.md) - Update Prompt
* [`delete`](docs/sapient_prompts_delete.md) - Delete Prompt

### [api-performance](docs/sapient_api-performance.md)

* [`targets-list`](docs/sapient_api-performance_targets-list.md) - List Targets
* [`runs-list`](docs/sapient_api-performance_runs-list.md) - List Eval Runs
* [`runs-retrieve`](docs/sapient_api-performance_runs-retrieve.md) - Retrieve Eval Run
* [`diagnose`](docs/sapient_api-performance_diagnose.md) - Diagnose Eval Runs
* [`config-retrieve`](docs/sapient_api-performance_config-retrieve.md) - Retrieve Eval Config
* [`config-update`](docs/sapient_api-performance_config-update.md) - Update Eval Config
* [`sources-list`](docs/sapient_api-performance_sources-list.md) - List Sources
* [`skills-list`](docs/sapient_api-performance_skills-list.md) - List Skills
* [`evals-list`](docs/sapient_api-performance_evals-list.md) - List Eval Definitions
* [`evals-retrieve`](docs/sapient_api-performance_evals-retrieve.md) - Retrieve Eval Definition
* [`evals-update`](docs/sapient_api-performance_evals-update.md) - Update Eval Definition
* [`custom-evals-list`](docs/sapient_api-performance_custom-evals-list.md) - List Custom Evals
* [`custom-evals-create`](docs/sapient_api-performance_custom-evals-create.md) - Create Custom Eval
* [`custom-evals-retrieve`](docs/sapient_api-performance_custom-evals-retrieve.md) - Retrieve Custom Eval
* [`custom-evals-update`](docs/sapient_api-performance_custom-evals-update.md) - Update Custom Eval
* [`custom-evals-delete`](docs/sapient_api-performance_custom-evals-delete.md) - Delete Custom Eval
* [`starter-projects-list`](docs/sapient_api-performance_starter-projects-list.md) - List Starter Projects
* [`starter-projects-create`](docs/sapient_api-performance_starter-projects-create.md) - Create Starter Project
* [`starter-projects-retrieve`](docs/sapient_api-performance_starter-projects-retrieve.md) - Retrieve Starter Project
* [`starter-projects-update`](docs/sapient_api-performance_starter-projects-update.md) - Update Starter Project
* [`starter-projects-delete`](docs/sapient_api-performance_starter-projects-delete.md) - Delete Starter Project
* [`actions-list`](docs/sapient_api-performance_actions-list.md) - List Actions
* [`actions-refresh`](docs/sapient_api-performance_actions-refresh.md) - Refresh Actions
* [`actions-retrieve`](docs/sapient_api-performance_actions-retrieve.md) - Retrieve Action
* [`actions-update`](docs/sapient_api-performance_actions-update.md) - Update Action
* [`actions-verify`](docs/sapient_api-performance_actions-verify.md) - Verify Action

</details>
<!-- End Available Commands [operations] -->

<!-- Start Request Body Input [stdinpiping] -->
## Request Body Input

Operations that accept a request body support three input methods, with a clear priority chain:

### Individual flags (highest priority)

```bash
sapient <command> --name "Jane" --age 30
```

### `--body` flag

Provide the entire request body as a JSON string:

```bash
sapient <command> --body '{"name": "John", "age": 30}'
```

Individual flags override `--body` values:

```bash
# Result: {name: "Jane", age: 30}
sapient <command> --body '{"name": "John", "age": 30}' --name "Jane"
```

### Stdin piping (lowest priority)

Pipe JSON into any command that accepts a request body:

```bash
echo '{"name": "John", "age": 30}' | sapient <command>
```

Individual flags override stdin values:

```bash
# Result: {name: "Jane", age: 30}
echo '{"name": "John", "age": 30}' | sapient <command> --name "Jane"
```

This is useful for chaining commands, reading from files, or scripting:

```bash
# Read body from a file
sapient <command> < request.json

# Pipe from another command
curl -s https://example.com/data.json | sapient <command>
```

### Priority

When multiple input methods are used, the priority is:

| Priority | Source | Description |
|----------|--------|-------------|
| 1 (highest) | Individual flags | `--name "Jane"` always wins |
| 2 | `--body` flag | Whole-body JSON via flag |
| 3 (lowest) | Stdin | Piped JSON input |
<!-- End Request Body Input [stdinpiping] -->

<!-- Start Server Selection [server] -->
## Server Selection

### Override Server URL

Use `--server-url` to override the server URL entirely, bypassing any named or indexed server selection:

```bash
sapient --server-url https://custom-api.example.com <command> [arguments]
```

**Precedence**: `--server-url` > `--server` > default
<!-- End Server Selection [server] -->

<!-- Start Output Formats [output-formats] -->
## Output Formats

Every command supports a `--output-format` flag that controls how the response is rendered to stdout.

### Available formats

| Format | Flag | Description |
|--------|------|-------------|
| Pretty | `--output-format pretty` (default) | Aligned key-value pairs with color, nested indentation. Human-readable at a glance. |
| JSON | `--output-format json` | JSON output. Passthrough when the response is already JSON (preserves original field order and numeric precision). Falls back to typed marshaling otherwise. |
| YAML | `--output-format yaml` | YAML output via standard marshaling. |
| Table | `--output-format table` | Tabular output for array responses. |
| TOON | `--output-format toon` | [Token-Oriented Object Notation](https://github.com/toon-format/spec) — a compact, line-oriented format that typically uses 30–60% fewer tokens than JSON. Well-suited for piping responses into LLM prompts. |

```bash
# Default pretty output
sapient <command>

# Machine-readable JSON
sapient <command> --output-format json

# TOON for LLM-friendly compact output
sapient <command> --output-format toon

# Pipe JSON to jq without using --output-format
sapient <command> --output-format json | jq '.fieldName'
```

### jq filtering

Use `--jq` to filter or transform the response inline using a [jq](https://jqlang.org) expression. This always outputs JSON and overrides `--output-format`:

```bash
# Extract a single field
sapient <command> --jq '.name'

# Filter an array
sapient <command> --jq '.items[] | select(.active == true)'
```

### Color control

Use `--color` to control terminal colors:

| Value | Behavior |
|-------|----------|
| `auto` (default) | Color when stdout is a TTY, plain text otherwise |
| `always` | Always colorize |
| `never` | Never colorize |

The `NO_COLOR` and `FORCE_COLOR` environment variables are also respected.

### Streaming and pagination

When using `--all` (pagination) or streaming operations, output is written incrementally as items arrive:

| Format | Streaming behavior |
|--------|-------------------|
| `json` | One compact JSON object per line ([NDJSON](https://github.com/ndjson/ndjson-spec)) |
| `yaml` | YAML documents separated by `---` |
| `toon` | One TOON-encoded object per block, separated by blank lines |
| `pretty` (default) | Pretty-printed items separated by blank lines |
<!-- End Output Formats [output-formats] -->

<!-- Start Error Handling [errors] -->
## Error Handling

The CLI uses standard exit codes to indicate success or failure:

| Exit Code | Meaning |
|-----------|---------|
| `0` | Success |
| `1` | Error (API error, invalid input, etc.) |

On success, the response data is printed to **stdout** as JSON. On failure, error details are printed to **stderr**.

```bash
# Capture output and handle errors
sapient ... > output.json 2> error.log
if [ $? -ne 0 ]; then
  echo "Error occurred, see error.log"
fi
```
<!-- End Error Handling [errors] -->

<!-- Start Diagnostics [diagnostics] -->
## Diagnostics

The CLI includes two diagnostic flags available on all commands:

### Dry Run

Preview what would be sent without making any network calls:

```bash
sapient <command> --dry-run
```

Output goes to stderr and includes:
- HTTP method and URL
- Request headers (sensitive values redacted)
- Request body preview (sensitive fields redacted)

The command exits successfully without contacting the API. This is useful for verifying request construction before executing.

### Debug

Log request and response diagnostics while running normally:

```bash
sapient <command> --debug
```

Debug output goes to stderr and includes:
- Request method, URL, headers, and body preview
- Response status, headers, and body preview
- Transport errors (if any)

The command still executes normally and produces its regular output on stdout.

### Flag Precedence

If both `--dry-run` and `--debug` are set, `--dry-run` takes precedence and no network calls are made.

### Security

Sensitive information is automatically redacted in diagnostic output:
- **Headers**: `Authorization`, `Cookie`, `Set-Cookie`, `X-API-Key`, and other security headers show `[REDACTED]`
- **Body**: JSON fields named `password`, `secret`, `token`, `api_key`, `client_secret`, etc. show `[REDACTED]`

Diagnostic output should still be treated as potentially sensitive operational data.
<!-- End Diagnostics [diagnostics] -->

<!-- Placeholder for Future Speakeasy SDK Sections -->

# Development

## Maturity

This CLI is in beta, and there may be breaking changes between versions without a major version update. Therefore, we recommend pinning usage
to a specific package version. This way, you can install the same version each time without breaking changes unless you are intentionally
looking for the latest version.

## Contributions

While we value open-source contributions to this CLI, this library is generated programmatically. Any manual changes added to internal files will be overwritten on the next generation. 
We look forward to hearing your feedback. Feel free to open a PR or an issue with a proof of concept and we'll do our best to include it in a future release. 

### CLI Created by [Speakeasy](https://www.speakeasy.com/?utm_source=github-com/usesapient/cli&utm_campaign=cli)
