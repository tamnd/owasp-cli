---
title: "CLI"
description: "Every command and subcommand, with the flags that matter."
weight: 10
---

```
owasp <command> [subcommand] [flags]
```

Run `owasp <command> --help` for the full flag list on any command. This
page is the map; keep it in step with the real command tree as you add to it.

## Commands

| Command | What it does |
|---|---|
| `list` | List all OWASP cheat sheets (one GitHub API request) |
| `search <query>` | Search sheets by topic name (client-side filter) |
| `version` | Print the version and exit |

## Global flags

| Flag | Default | Description |
|---|---|---|
| `-o, --output` | `auto` | Output format: `table`, `json`, `jsonl`, `csv`, `tsv`, `url`, `raw` |
| `--fields` | all | Comma-separated columns to include |
| `--no-header` | false | Omit header row in table/csv/tsv |
| `--template` | | Go `text/template` applied per record |
| `-n, --limit` | 0 (all) | Limit number of records returned |
| `-q, --quiet` | false | Suppress progress messages on stderr |
| `--delay` | 200ms | Minimum gap between HTTP requests |
| `--timeout` | 30s | Per-request timeout |
| `--retries` | 3 | Retry attempts on 429/5xx |
| `--user-agent` | Chrome UA | User-Agent header sent to GitHub |
