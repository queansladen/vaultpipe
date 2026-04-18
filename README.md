# vaultpipe

Stream secrets from HashiCorp Vault into process environments without writing to disk.

## Installation

```bash
go install github.com/yourusername/vaultpipe@latest
```

Or download a pre-built binary from the [releases page](https://github.com/yourusername/vaultpipe/releases).

## Usage

Inject secrets from a Vault path into a process environment:

```bash
vaultpipe run --path secret/data/myapp -- ./myapp
```

The process inherits your current environment plus the secrets fetched from Vault. Secrets are never written to disk.

### Options

| Flag | Description |
|------|-------------|
| `--path` | Vault secret path to read from |
| `--addr` | Vault server address (default: `$VAULT_ADDR`) |
| `--token` | Vault token (default: `$VAULT_TOKEN`) |
| `--prefix` | Optional env var prefix for injected secrets |

### Example

```bash
export VAULT_ADDR="https://vault.example.com"
export VAULT_TOKEN="s.xxxxxxxx"

vaultpipe run --path secret/data/db --prefix DB_ -- python server.py
```

This fetches all key/value pairs at `secret/data/db` and exposes them as `DB_HOST`, `DB_PASSWORD`, etc. to the child process.

## Requirements

- Go 1.21+
- HashiCorp Vault with KV v2 secrets engine

## License

MIT © [yourusername](https://github.com/yourusername)