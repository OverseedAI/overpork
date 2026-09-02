---
name: opork
description: Safely inspect or manage Porkbun domains, DNS records, nameservers, URL forwarding, pricing, SSL certificates, DNSSEC, and glue records with the opork CLI. Use this skill whenever the user asks to operate on Porkbun resources, troubleshoot an opork command, check domain availability or pricing, retrieve Porkbun certificate material, or change Porkbun DNS, DNSSEC, forwarding, nameservers, auto-renewal, or glue records.
---

# opork

Use `opork` as the supported interface to the Porkbun API. Prefer commands that
are easy to audit, keep credentials out of captured output, and make externally
consequential changes only with clear authorization.

## Establish the environment

1. Check whether the CLI is available with `command -v opork`.
2. If it is missing, offer to install it with `npm install -g opork`.
3. Let `opork` read credentials from `PORKBUN_API_KEY` and
   `PORKBUN_SECRET_KEY`. Never print, interpolate, or pass their values as
   command arguments.
4. If the user prefers a config file, use `opork config path` first. Check only
   whether that path exists; do not read or print the file. If it exists,
   confirm before overwriting it because `config init` has no overwrite prompt.
   Then use interactive `opork config init` so the secret is not placed in
   shell history.
5. Verify authenticated access with `opork ping --json` before depending on
   account data.

`opork pricing list` is public. Other account and availability operations,
including `opork pricing check`, require credentials.

## Choose a safe execution mode

- Add `--json` when reading output programmatically. Successful JSON goes to
  stdout; errors go to stderr and return exit code 1.
- Run a read command before a write when current state affects the change. For
  example, list matching DNS records before updating or deleting one.
- Treat an explicit user request containing the exact domain, record, and new
  value as authorization for that stated change. If any target or value is
  inferred, preview the command and ask first.
- Confirm the exact target and effect before deleting data, replacing
  nameservers, changing DNSSEC or glue records, or exposing a private key.
- Do not execute `opork domain register`. The current CLI does not implement
  Porkbun's required price binding, terms acceptance, dry run, or idempotency
  safeguards for registration. Use `pricing check` only for inspection and
  direct the user to Porkbun's website for the purchase.
- Confirm before enabling auto-renew because it creates future billing. Also
  confirm before disabling it unless the user's request is explicit, because a
  domain could expire.
- Do not claim a dry run occurred. `opork` has no dry-run mode and no built-in
  confirmation prompts.

## Read account and DNS state

Use these commands for inspection:

```bash
opork domain list --json
opork domain get example.com --json
opork domain ns-get example.com --json
opork domain forward-list example.com --json

opork dns list example.com --json
opork dns list example.com --type A --json
opork dns list example.com --type A --subdomain www --json

opork pricing list --json
opork pricing check example.com --json
opork dnssec list example.com --json
opork glue list example.com --json
```

For DNS filters, `--subdomain` requires `--type`. Use `@` for the root domain
only with `dns list --subdomain`, `dns set`, and `dns delete-by-name`. For
`dns create --name` and `dns update --name`, omit the name to target the root;
do not pass a literal `@`.

## Manage DNS records

Create a record with the root domain as the default or `--name` for a
subdomain:

```bash
opork dns create example.com A 192.0.2.10 --name www --json
opork dns create example.com MX mail.example.com --prio 10 --json
```

Update by record ID when the ID came from a fresh list response:

```bash
opork dns update example.com 123456 A 192.0.2.20 --name www --json
```

Update by type and subdomain when that pair identifies the intended record:

```bash
opork dns set example.com A www 192.0.2.20 --json
opork dns set example.com A @ 192.0.2.20 --json
```

The by-name/type endpoint can edit every record matching that pair. Use it only
after a fresh list shows exactly the intended match. For multi-value types such
as MX or TXT, prefer record-ID updates unless replacing all matches is explicit.

Delete only after listing the current record and confirming the identifier or
type/subdomain pair:

```bash
opork dns delete example.com 123456 --json
opork dns delete-by-name example.com A www --json
```

`delete-by-name` can delete every record matching the type/subdomain pair. Use
record-ID deletion when more than one match exists or the blast radius is
unclear.

Quote TXT and other content containing spaces or shell metacharacters.

## Manage domains and forwarding

Set nameservers only after reading the current values and confirming the full
replacement set:

```bash
opork domain ns-get example.com --json
opork domain ns-set example.com ns1.example.net ns2.example.net --json
```

Add or delete URL forwarding:

```bash
opork domain forward-add example.com https://www.example.net \
  --subdomain www --type permanent --json
opork domain forward-delete example.com 123456 --json
```

Use only `temporary` or `permanent` for `--type`; the CLI advertises these
values but does not validate them locally.

### Domain registration is temporarily unsupported

Use `pricing check` to inspect availability and the current price signal:

```bash
opork pricing check example.com --json
```

Do not run `opork domain register` with the current CLI. Porkbun's registration
flow requires an exact fresh quote sent as integer cents, explicit acceptance
of terms, a successful dry run, and an idempotency key reused for retries. The
CLI does not currently provide those request-level safeguards. Direct the user
to Porkbun's website until registration support is updated and tested. Do not
provide a copy-pastable registration invocation, even as a do-not-run example.

Manage renewal separately:

```bash
opork domain auto-renew example.com enable --json
opork domain auto-renew example.com disable --json
```

## Retrieve SSL certificate material

Always request the narrowest needed part:

```bash
opork ssl get example.com --part cert --json
opork ssl get example.com --part intermediate --json
opork ssl get example.com --part public --json
```

`opork ssl get example.com` without `--part` includes the private key. Use
`--part key` only when the user explicitly requests the private key and provides
a secure destination. Never paste private-key output into the conversation,
logs, command substitution, or an unprotected file.

## Manage DNSSEC and glue records

Inspect current values before changing them:

```bash
opork dnssec list example.com --json
opork glue list example.com --json
```

Create or update with values supplied by the user or authoritative provider:

```bash
opork dnssec create example.com \
  --keytag 12345 --algorithm 13 --digest-type 2 --digest ABCDEF --json
opork glue create example.com ns1 192.0.2.53 2001:db8::53 --json
opork glue update example.com ns1 192.0.2.54 2001:db8::54 --json
```

`glue update` replaces the complete IP set for that glue record. List the
current record first and pass every IP that should remain, not only the address
being changed.

Delete only after confirmation:

```bash
opork dnssec delete example.com 12345 --json
opork glue delete example.com ns1 --json
```

DNSSEC and glue mistakes can make a domain unreachable. Do not invent key tags,
algorithms, digests, nameserver labels, or IP addresses.

## Report results

- Summarize the operation, exact domain or record affected, and returned status.
- For reads, present the relevant fields rather than dumping large raw payloads.
- For writes, do not imply the resulting state was verified unless a follow-up
  read command confirms it.
- Redact credentials, private keys, and unrelated account data from all output.
