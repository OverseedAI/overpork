# overpork

CLI wrapper for the Porkbun API, optimized for agentic use.

## Install

```bash
npm install -g overpork
```

## Configuration

Set credentials via environment variables:

```bash
export PORKBUN_API_KEY=pk1_xxx
export PORKBUN_SECRET_KEY=sk1_xxx
```

Or create a config file:

```bash
overpork config init --api-key pk1_xxx --secret-key sk1_xxx
```

Config file location: `~/.config/overpork/config.yaml`

## Commands

### General

```bash
overpork ping                    # Test connectivity
overpork version                 # Print version
overpork config path             # Show config path
```

### DNS Records

```bash
overpork dns list <domain>
overpork dns list <domain> --type A
overpork dns list <domain> --type A --subdomain www

overpork dns create <domain> <type> <content>
overpork dns create example.com A 192.168.1.1
overpork dns create example.com A 192.168.1.1 --name www
overpork dns create example.com MX mail.example.com --prio 10

overpork dns update <domain> <id> <type> <content>
overpork dns set <domain> <type> <subdomain> <content>   # Update by name
overpork dns set example.com A www 192.168.1.1
overpork dns set example.com A @ 192.168.1.1             # @ = root

overpork dns delete <domain> <id>
overpork dns delete-by-name <domain> <type> <subdomain>
```

### Domains

```bash
overpork domain list
overpork domain get <domain>

overpork domain register <domain>
overpork domain register example.com --years 2 --ns ns1.example.com --ns ns2.example.com

overpork domain auto-renew <domain> enable
overpork domain auto-renew <domain> disable

overpork domain ns-get <domain>
overpork domain ns-set <domain> <ns1> [ns2] [ns3]...

overpork domain forward-list <domain>
overpork domain forward-add <domain> <url> [--subdomain www] [--type permanent]
overpork domain forward-delete <domain> <id>
```

### Pricing

```bash
overpork pricing list            # List all TLD prices
overpork pricing check <domain>  # Check availability and price
```

### SSL Certificates

```bash
overpork ssl get <domain>
overpork ssl get <domain> --part cert
overpork ssl get <domain> --part key
overpork ssl get <domain> --part intermediate
```

### DNSSEC

```bash
overpork dnssec list <domain>
overpork dnssec create <domain> --keytag X --algorithm Y --digest-type Z --digest ABC
overpork dnssec delete <domain> <keytag>
```

### Glue Records

```bash
overpork glue list <domain>
overpork glue create <domain> <subdomain> <ip> [ip...]
overpork glue update <domain> <subdomain> <ip> [ip...]
overpork glue delete <domain> <subdomain>
```

## JSON Output

Add `--json` to any command for JSON output:

```bash
overpork dns list example.com --json
overpork pricing check example.com --json
```

## Exit Codes

- `0` - Success
- `1` - Error (message printed to stderr)

## License

MIT
