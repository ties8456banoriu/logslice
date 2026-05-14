# logslice

Command-line utility to extract and filter structured log ranges by time window across multiple files.

## Installation

```bash
go install github.com/yourusername/logslice@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/logslice.git
cd logslice && go build -o logslice .
```

## Usage

```
logslice [flags] <file> [file...]
```

### Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--from` | Start of time window (RFC3339) | required |
| `--to` | End of time window (RFC3339) | required |
| `--field` | JSON field name containing the timestamp | `"time"` |
| `--format` | Timestamp format layout | RFC3339 |

### Example

Extract logs between two timestamps across multiple files:

```bash
logslice --from 2024-03-01T10:00:00Z --to 2024-03-01T11:00:00Z app.log app.log.1 app.log.2
```

Filter a specific time window and pipe to `jq`:

```bash
logslice --from 2024-03-01T10:00:00Z --to 2024-03-01T10:30:00Z --field timestamp service.log | jq '.level'
```

## How It Works

`logslice` reads structured (JSON) log lines from one or more files, parses the timestamp field in each entry, and outputs only the lines that fall within the specified time window. Files are processed in order and output is written to stdout.

## Requirements

- Go 1.21 or later

## License

MIT © [yourusername](https://github.com/yourusername)