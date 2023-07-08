# shush
Shush is a simple program for redacting secrets from a provided block of text.

Example:

```bash
shush < ./example.txt
{
    token: "[[REDACTED]]
}
```

Shush is built ontop of Gitleak's detection engine so it supports all the same secrets and keys that gitleaks does.

## Installation

If you're using Go:

```bash
go install github.com/bradcypert/shush@latest
```