---
title: "Introduction"
description: "What owasp is and how it is put together."
weight: 10
---

Browse the OWASP Cheat Sheet Series from the command line.

owasp is a single binary. It fetches the cheat sheet catalog from the
OWASP GitHub repository over plain HTTPS, shapes the responses into clean
records, and gets out of your way. There is nothing to sign up for and nothing
to run alongside it.

## How it is built

- A **library package** (`owasp`) holds the HTTP client, the GitHub Contents
  API call, and the typed data models. It paces requests, sets an honest
  User-Agent, and retries the transient failures that any public API throws
  under load.
- A **command tree** (`cli`) wraps the library in subcommands with shared
  output formats and flags.
- One **`cmd/owasp`** entry point ties them together.

## Scope

owasp is a read-only client over data the OWASP GitHub repository already
serves publicly. The `list` command makes exactly one HTTP request. That
narrow scope keeps it a single small binary with no database, no daemon, and
no setup.

Next: [install it](/getting-started/installation/), then take the
[quick start](/getting-started/quick-start/).
