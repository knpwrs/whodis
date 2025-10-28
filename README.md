# whodis

[![Test](https://github.com/knpwrs/whodis/actions/workflows/test.yml/badge.svg)](https://github.com/knpwrs/whodis/actions/workflows/test.yml)
[![Release](https://github.com/knpwrs/whodis/actions/workflows/release.yml/badge.svg)](https://github.com/knpwrs/whodis/actions/workflows/release.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/knpwrs/whodis)](https://goreportcard.com/report/github.com/knpwrs/whodis)

> A modern, colorful CLI tool for WHOIS lookups

**whodis** is to `whois` what [doggo](https://github.com/mr-karan/doggo) is to `dig` - a beautiful, user-friendly command-line interface that makes WHOIS queries more readable and enjoyable.

## Features

- 🎨 **Colorful terminal output** with clearly organized sections
- 📊 **Multiple output formats**: pretty terminal, JSON, short table, and raw
- ⚡ **Concurrent lookups** for multiple domains
- 🌍 **Automatic server discovery** via IANA with referral following
- 📝 **Parsed data** for 600+ TLDs
- ⏱️ **Query timing** display

## Installation

### Homebrew

```bash
brew install knpwrs/tap/whodis
```

### eget

```bash
eget knpwrs/whodis
```

### Pre-built binaries

Download the latest release from the [releases page](https://github.com/knpwrs/whodis/releases).

### Using Go

```bash
go install github.com/knpwrs/whodis/cmd/whodis@latest
```

### Build from source

```bash
git clone https://github.com/knpwrs/whodis.git
cd whodis
go build -o whodis ./cmd/whodis
```

## Usage

### Basic lookup

```bash
whodis example.com
```

Output:

```
DOMAIN INFORMATION
├─ Name: example.com
├─ Status: clientDeleteProhibited, clientTransferProhibited, clientUpdateProhibited
├─ Created: 1995-08-14T04:00:00Z (30 years ago)
├─ Updated: 2025-08-14T07:01:39Z
└─ Expires: 2026-08-13T04:00:00Z (in 288 days)

REGISTRAR
├─ Name: RESERVED-Internet Assigned Numbers Authority
└─ WHOIS Server: http://res-dom.iana.org

NAMESERVERS
├─ a.iana-servers.net
└─ b.iana-servers.net

Query time: 517ms
```

### Multiple domains

```bash
whodis example.com github.com google.com
```

### JSON output

```bash
whodis --json example.com
```

### Short table format

```bash
whodis --short example.com github.com
```

Output:

```
DOMAIN       REGISTRAR                                     EXPIRES               STATUS
example.com  RESERVED-Internet Assigned Numbers Authority  2026-08-13T04:00:00Z  clientDeleteProhibited
github.com   MarkMonitor Inc.                              2026-10-09T18:20:50Z  clientDeleteProhibited
```

### Raw WHOIS output

```bash
whodis --raw example.com
```

## Flags

- `-j, --json` - Output results in JSON format
- `-s, --short` - Output minimal information in table format
- `-r, --raw` - Output raw WHOIS data
- `-t, --timeout` - Query timeout in seconds (default: 30)
- `-d, --debug` - Show debug information
- `-h, --help` - Show help message
- `-v, --version` - Show version information

## Libraries Used

- [github.com/likexian/whois](https://github.com/likexian/whois) - WHOIS client with automatic server discovery
- [github.com/likexian/whois-parser](https://github.com/likexian/whois-parser) - WHOIS data parser supporting 600+ TLDs
- [github.com/spf13/pflag](https://github.com/spf13/pflag) - POSIX-style flag parsing
- [github.com/fatih/color](https://github.com/fatih/color) - Terminal colorization

## Inspiration

This project is inspired by [doggo](https://github.com/mr-karan/doggo), a fantastic DNS client with beautiful output formatting.

## License

This project is dual-licensed under [CC0 1.0 Universal](LICENSE) and [The Unlicense](UNLICENSE). You can use, modify, and distribute this software without any restrictions.
