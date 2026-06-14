---
title: "Quick start"
description: "Run your first owasp command."
weight: 30
---

Once `owasp` is on your `PATH`:

```bash
owasp --help            # see the command tree
owasp version           # build info

owasp list              # list all OWASP cheat sheets
owasp list -o json      # as a JSON array
owasp list -o csv       # as CSV (pipe into spreadsheet tools)

owasp search sql        # find sheets matching "sql"
owasp search xss -o jsonl | jq .url
```

The `list` command fetches the full catalog from the OWASP GitHub repository
in a single request. `search` filters the result client-side, so it also uses
exactly one request.
