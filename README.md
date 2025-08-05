# Ledger

[![Build Status](https://github.com/askolesov/ledger/workflows/build/badge.svg)](https://github.com/askolesov/ledger/actions)
[![Go Version](https://img.shields.io/badge/go-1.24-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/license-CC--BY--4.0-green.svg)](LICENSE)

A command-line tool for managing personal financial ledgers using the [**Open Ledger Format (OLF)**](Open%20Ledger%20Format.md) — a minimal, human-readable schema for everyday finances that keeps every cent accounted for.

## Features

- 🏗️ **Human-readable** plain-text files (YAML/JSON)
- 📊 **Validate** ledger files against OLF specifications
- 📈 **Generate reports** from financial data
- 🐳 **Cross-platform** with Docker support

## Quick Start

### Installation

**Download binary releases:**
```bash
# Download from GitHub releases
curl -L https://github.com/askolesov/ledger/releases/latest/download/ledger_Linux_x86_64.tar.gz | tar xz
```

**Install with Go:**
```bash
go install github.com/askolesov/ledger/cmd/ledger@latest
```

**Docker:**
```bash
docker pull ghcr.io/askolesov/ledger:latest
```

### Basic Usage

```bash
# Validate OLF v2.0 file
ledger validate ledger.yaml

# Generate report for OLF v2.0
ledger report ledger.yaml

# Show version
ledger version
```

## Open Ledger Format

The Open Ledger Format stores personal finance data in plain-text files with strict validation rules:

```yaml
years:
  2025:
    opening_balance: 1000
    closing_balance: 1050
    months:
      1:
        opening_balance: 1000
        closing_balance: 1050
        accounts:
          Checking:
            opening_balance: 600
            closing_balance: 620
            entries:
              - amount: 50
                note: "Salary (January)"
                date: "2025-01-28"
                tag: "Income"
```

**Key principles:**
- **Human-oriented** — maintain with any text editor
- **Consistency** — validation ensures strict, deterministic balances
- **High-level focus** — summary budgeting, not micro-transactions

For complete specification, see [Open Ledger Format.md](Open%20Ledger%20Format.md).

## Development

### Building from Source

```bash
# Clone repository
git clone https://github.com/askolesov/ledger.git
cd ledger

# Build binary
make build

# Run tests
make test-all

# Install locally
make install
```

### Requirements

- Go 1.24+
- Make (optional, for convenience commands)

## License

Released under [CC-BY-4.0](LICENSE) — use, fork, extend; credit the authors.

## Links

- [Releases](https://github.com/askolesov/ledger/releases)
- [Open Ledger Format Specification](Open%20Ledger%20Format.md)
- [Docker Images](https://github.com/askolesov/ledger/pkgs/container/ledger)
