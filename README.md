# gh-jwt-fetch

A lightweight Go CLI tool that downloads a file from a GitHub repository using GitHub App authentication. It generates a JWT from the App's private key, obtains an Installation Access Token, and fetches the specified file.

Designed for container environments with all configuration via environment variables. Supports both one-shot and loop (daemon) modes.

## Modes

- **One-shot mode**: Default. Downloads the file once and exits.
- **Loop mode**: Set `GH_INTERVAL` to run as a long-lived process that periodically re-downloads the file. Handles SIGINT/SIGTERM for graceful shutdown.

## Environment Variables

| Variable | Required | Description |
|---|---|---|
| `GH_APP_ID` | Yes | GitHub App ID |
| `GH_APP_PRIVATE_KEY` | * | Private key PEM content (inline) |
| `GH_APP_PRIVATE_KEY_PATH` | * | Path to private key PEM file |
| `GH_APP_INSTALLATION_ID` | Yes | Installation ID |
| `GH_REPO_OWNER` | Yes | Repository owner |
| `GH_REPO_NAME` | Yes | Repository name |
| `GH_FILE_PATH` | Yes | File path to download |
| `GH_OUTPUT_PATH` | Yes | Output file path |
| `GH_API_BASE_URL` | No | API base URL (default: `https://api.github.com`) |
| `GH_REF` | No | Branch, tag, or commit (default: repository's default branch) |
| `GH_INTERVAL` | No | Fetch interval (e.g. `5m`, `1h`). Omit for one-shot mode |

\* Either `GH_APP_PRIVATE_KEY` or `GH_APP_PRIVATE_KEY_PATH` is required. If both are set, `GH_APP_PRIVATE_KEY` takes precedence.

## Usage

### Build

```bash
go build -o gh-jwt-fetch .
```

### One-shot

```bash
GH_APP_ID=123456 \
GH_APP_PRIVATE_KEY_PATH=./private-key.pem \
GH_APP_INSTALLATION_ID=789012 \
GH_REPO_OWNER=myorg \
GH_REPO_NAME=myrepo \
GH_FILE_PATH=config/settings.json \
GH_OUTPUT_PATH=/tmp/settings.json \
./gh-jwt-fetch
```

### Loop (every 5 minutes)

```bash
GH_INTERVAL=5m \
GH_APP_ID=123456 \
GH_APP_PRIVATE_KEY_PATH=./private-key.pem \
GH_APP_INSTALLATION_ID=789012 \
GH_REPO_OWNER=myorg \
GH_REPO_NAME=myrepo \
GH_FILE_PATH=config/settings.json \
GH_OUTPUT_PATH=/tmp/settings.json \
./gh-jwt-fetch
```

### Docker

```bash
docker build -t gh-jwt-fetch .

docker run --rm \
  -e GH_APP_ID=123456 \
  -e GH_APP_PRIVATE_KEY_PATH=/key.pem \
  -e GH_APP_INSTALLATION_ID=789012 \
  -e GH_REPO_OWNER=myorg \
  -e GH_REPO_NAME=myrepo \
  -e GH_FILE_PATH=config/settings.json \
  -e GH_OUTPUT_PATH=/out/settings.json \
  -e GH_INTERVAL=5m \
  -v /path/to/private-key.pem:/key.pem:ro \
  -v /path/to/output:/out \
  gh-jwt-fetch
```

### GitHub Enterprise Server

Set `GH_API_BASE_URL` to your GHE API endpoint:

```bash
GH_API_BASE_URL=https://ghe.example.com/api/v3 \
# ... other variables ...
./gh-jwt-fetch
```

## Exit Codes

| Code | Meaning |
|---|---|
| 0 | Success |
| 1 | Configuration error |
| 2 | Runtime error (JWT generation, token fetch, download, or write) |

In loop mode, individual fetch errors are logged and the process continues.

## Testing

```bash
go test -v -race ./...
```

## Details

- **Zero external dependencies** — JWT RS256 signing uses only the Go standard library
- **PKCS1 / PKCS8 support** — handles both private key formats
- **Raw content download** — uses `Accept: application/vnd.github.raw+json` (up to 100 MB)
- **Distroless container** — `gcr.io/distroless/static-debian12:nonroot` base image

## License

MIT
