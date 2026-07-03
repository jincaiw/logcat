# logcat

[中文文档](README-zh.md) · [Installation Guide](docs/installation.md) · [User Guide](docs/user-guide.md) · [Project Site](http://logcat.mujizi.com/)

logcat is a lightweight web application for receiving, parsing, filtering, forwarding, and alerting on Syslog messages. It is designed for security operations, blue-team monitoring, and small production environments where a simple self-hosted Syslog alert pipeline is needed.

## Highlights

- Receive Syslog over UDP/TCP.
- Parse JSON, Syslog+JSON, delimiter, key-value, and regex logs.
- Filter logs with flexible rules and deduplication.
- Send alerts to Feishu, Email, HTTP API, or forward to another Syslog server.
- Manage devices, parsing templates, field mappings, and alert rules from a web UI.
- Single binary deployment and Docker deployment.
- SQLite storage with local persistent data.

## Demo

![Dashboard](docs/assets/demo-dashboard.png)

![Notifications](docs/assets/demo-notifications.png)

![Logs](docs/assets/demo-logs.png)

## Quick Start

### Docker Compose

```bash
curl -O https://raw.githubusercontent.com/jincaiw/logcat/v0.2.8/docker-compose.yml
docker compose up -d
```

Open `http://localhost:8080`.

Default account:

```text
Username: admin
Password: admin123
```

Change the password after first login.

### Linux one-line install

```bash
curl -fsSL https://raw.githubusercontent.com/jincaiw/logcat/v0.2.8/scripts/install-linux.sh | sudo bash
```

Open `http://<server-ip>:8080`.

### Binary package

Download `logcat-0.2.8-linux-amd64.tar.gz` or `logcat-0.2.8-linux-arm64.tar.gz` from the release page:

```bash
tar -xzf logcat-0.2.8-linux-amd64.tar.gz
cd logcat-0.2.8-linux-amd64
./start.sh 8080
```

## Docker

```bash
docker run -d --name logcat \
  --restart unless-stopped \
  -p 8080:8080 \
  -p 5140:5140/udp \
  -p 5140:5140/tcp \
  -v logcat-data:/app/data \
  qing1205/logcat:0.2.8
```

Set a custom first admin password before the first start:

```bash
-e LOGCAT_ADMIN_PASSWORD='change-me-now'
```

## Ports

| Port | Protocol | Description |
| --- | --- | --- |
| 8080 | TCP | Web UI and API |
| 5140 | UDP/TCP | Syslog receiver |

## Data and configuration

By default, data is stored next to the executable in `data/`. In Docker, data is stored in `/app/data` and should be mounted as a volume.

Useful environment variables:

| Variable | Description |
| --- | --- |
| `SYSLG_ALERT_DATA_DIR` | Data directory |
| `SYSLG_ALERT_TEMPLATES_DIR` | Built-in templates directory |
| `LOGCAT_OPEN_BROWSER` | Set `1` to open browser automatically |
| `LOGCAT_ADMIN_USERNAME` | Initial admin username |
| `LOGCAT_ADMIN_PASSWORD` | Initial admin password |

## Build from source

```bash
npm -C frontend ci
npm -C frontend run build
go test ./...
go build -o logcat .
```

Build Linux packages:

```bash
APP_VERSION=$(cat VERSION) TARGET_OS=linux TARGET_ARCH=amd64 bash build-web.sh
APP_VERSION=$(cat VERSION) TARGET_OS=linux TARGET_ARCH=arm64 bash build-web.sh
```

Build Docker image:

```bash
docker build -t qing1205/logcat:$(cat VERSION) .
```

## Release process

One-shot:

```bash
bash scripts/publish-release.sh <version>
```

This will auto-commit release changes, create and push the tag, push DockerHub images, and verify GitHub Release / Pages / DockerHub.

Standalone checks:

```bash
bash scripts/release-check.sh <version>
bash scripts/release-check.sh --github-only <version>
bash scripts/release-check.sh --pages-only <version>
bash scripts/release-check.sh --docker-only <version>
```

Dry run:

```bash
bash scripts/publish-release.sh --dry-run <version>
```

Step-by-step:

```bash
bash scripts/release.sh <version>
git push origin <branch>
git push origin v<version>
```

See [docs/release-process.md](docs/release-process.md) for the standard release flow, one-shot publish command, and GitHub Actions release automation.

## Documentation

- [Installation Guide](docs/installation.md)
- [User Guide](docs/user-guide.md)
- [Release Process](docs/release-process.md)
- [Project Site](http://logcat.mujizi.com/)

## License

See repository license.
