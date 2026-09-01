# opork

CLI wrapper for the Porkbun API, optimized for agentic use.

## Install

```bash
npm install -g opork
```

### Agent Skill

Install the bundled skill with the [Skills CLI](https://skills.sh/):

```bash
npx skills add OverseedAI/overpork --skill opork
```

The skill teaches compatible agents to use `opork` safely for Porkbun domain,
DNS, pricing, SSL, DNSSEC, and glue-record operations.

## Configuration

Set credentials via environment variables:

```bash
export PORKBUN_API_KEY=pk1_xxx
export PORKBUN_SECRET_KEY=sk1_xxx
```

Or create a config file:

```bash
opork config init
```

Interactive setup keeps the secret key out of shell history and process
arguments. Run `opork config path` first and back up an existing file before
reinitializing it; `config init` overwrites the file without prompting.

Config file location: run `opork config path` to print the exact path for your system. The file lives in an `overpork` subdirectory of the OS-specific user config directory (`os.UserConfigDir()`), so the defaults are:

- **macOS**: `~/Library/Application Support/overpork/config.yaml`
- **Linux**: `~/.config/overpork/config.yaml` (or `$XDG_CONFIG_HOME/overpork/config.yaml` if `$XDG_CONFIG_HOME` is set)
- **Windows**: `%AppData%\overpork\config.yaml` (e.g. `C:\Users\<you>\AppData\Roaming\overpork\config.yaml`)

## Commands

### General

```bash
opork ping                    # Test connectivity
opork version                 # Print version
opork config path             # Show config path
```

### DNS Records

```bash
opork dns list <domain>
opork dns list <domain> --type A
opork dns list <domain> --type A --subdomain www

opork dns create <domain> <type> <content>
opork dns create example.com A 192.168.1.1
opork dns create example.com A 192.168.1.1 --name www
opork dns create example.com MX mail.example.com --prio 10

opork dns update <domain> <id> <type> <content>
opork dns set <domain> <type> <subdomain> <content>   # Update by name
opork dns set example.com A www 192.168.1.1
opork dns set example.com A @ 192.168.1.1             # @ = root

opork dns delete <domain> <id>
opork dns delete-by-name <domain> <type> <subdomain>
```

`dns set` and `dns delete-by-name` can affect every record matching the
type/subdomain pair. List records first and prefer record-ID operations when
more than one match exists.

### Domains

```bash
opork domain list
opork domain get <domain>
```

The `opork domain register` command exists but must not be used with Porkbun's
current registration flow. It does not yet send the required exact price,
explicit terms acceptance, dry-run request, or idempotency key. Use
`opork pricing check <domain> --json` for availability, then register on the
Porkbun website.

```bash
opork domain auto-renew <domain> enable
opork domain auto-renew <domain> disable

opork domain ns-get <domain>
opork domain ns-set <domain> <ns1> [ns2] [ns3]...

opork domain forward-list <domain>
opork domain forward-add <domain> <url> [--subdomain www] [--type permanent]
opork domain forward-delete <domain> <id>
```

### Pricing

```bash
opork pricing list            # List all TLD prices
opork pricing check <domain>  # Check availability and price
```

### SSL Certificates

```bash
opork ssl get <domain> --part cert
opork ssl get <domain> --part intermediate
opork ssl get <domain> --part public
```

Without `--part`, the output includes the private key. Use `--part key` only
when the private key is explicitly required and the output has a secure
destination.

### DNSSEC

```bash
opork dnssec list <domain>
opork dnssec create <domain> --keytag X --algorithm Y --digest-type Z --digest ABC
opork dnssec delete <domain> <keytag>
```

### Glue Records

```bash
opork glue list <domain>
opork glue create <domain> <subdomain> <ip> [ip...]
opork glue update <domain> <subdomain> <ip> [ip...]
opork glue delete <domain> <subdomain>
```

`glue update` replaces the complete IP set. List the current record first and
include every address that should remain.

## JSON Output

Add `--json` to supported data-producing commands for JSON output:

```bash
opork dns list example.com --json
opork pricing check example.com --json
```

## Exit Codes

- `0` - Success
- `1` - Error (message printed to stderr)

## Development

### Local Build

```bash
make build      # Build binary
make test       # Run tests
make lint       # Run linter
make dist       # Build for all platforms
```

### Releasing

Releases are automated via GitHub Actions. To publish a new version:

1. Bump the version in `package.json`:
   ```bash
   npm version patch   # 0.1.0 -> 0.1.1
   npm version minor   # 0.1.0 -> 0.2.0
   npm version major   # 0.1.0 -> 1.0.0
   ```

2. Push to main:
   ```bash
   git push
   ```

The workflow will automatically:
- Build binaries for all platforms (darwin/linux/windows, amd64/arm64)
- Create a GitHub release with the binaries
- Publish to npm

## License

MIT
