# opork

CLI wrapper for the Porkbun API, optimized for agentic use.

## Install

```bash
npm install -g opork
```

## Configuration

Set credentials via environment variables:

```bash
export PORKBUN_API_KEY=pk1_xxx
export PORKBUN_SECRET_KEY=sk1_xxx
```

Or create a config file:

```bash
opork config init --api-key pk1_xxx --secret-key sk1_xxx
```

Config file location: `~/.config/opork/config.yaml`

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

### Domains

```bash
opork domain list
opork domain get <domain>

opork domain register <domain>
opork domain register example.com --years 2 --ns ns1.example.com --ns ns2.example.com

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
opork ssl get <domain>
opork ssl get <domain> --part cert
opork ssl get <domain> --part key
opork ssl get <domain> --part intermediate
```

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

## JSON Output

Add `--json` to any command for JSON output:

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
