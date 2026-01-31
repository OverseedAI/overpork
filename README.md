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

Or create a config file at `~/.config/overpork/config.yaml`:

```yaml
api_key: pk1_xxx
secret_key: sk1_xxx
```

## Usage

```bash
# Test connectivity
overpork ping

# DNS
overpork dns list example.com
overpork dns create example.com A 192.168.1.1 --name www
overpork dns set example.com A www 192.168.1.1
overpork dns delete example.com 123456

# Domains
overpork domain list
overpork domain get example.com
overpork domain ns-get example.com
overpork domain ns-set example.com ns1.porkbun.com ns2.porkbun.com

# Pricing
overpork pricing list
overpork pricing check example.com

# SSL
overpork ssl get example.com
overpork ssl get example.com --part cert
overpork ssl get example.com --part key
```

## JSON Output

Add `--json` to any command for JSON output:

```bash
overpork dns list example.com --json
```

## License

MIT
