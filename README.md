# greplive

> A terminal tool that streams and filters log output in real-time with regex patterns and color-coded severity levels.

---

## Installation

```bash
go install github.com/yourusername/greplive@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/greplive.git
cd greplive
go build -o greplive .
```

---

## Usage

Pipe any log stream into `greplive` and optionally provide a regex filter:

```bash
# Stream and colorize all log output
tail -f /var/log/app.log | greplive

# Filter lines matching a pattern
tail -f /var/log/app.log | greplive -p "ERROR|WARN"

# Follow a running service with journalctl
journalctl -fu my-service | greplive -p "timeout"
```

### Flags

| Flag | Description |
|------|-------------|
| `-p` | Regex pattern to filter log lines |
| `-i` | Case-insensitive matching |
| `--no-color` | Disable color-coded output |

### Severity Colors

| Level | Color |
|-------|-------|
| ERROR | 🔴 Red |
| WARN | 🟡 Yellow |
| INFO | 🟢 Green |
| DEBUG | 🔵 Blue |

---

## Requirements

- Go 1.21+

---

## License

[MIT](LICENSE)