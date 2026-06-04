# logcat

English | [中文](README.zh-CN.md)

logcat is a lightweight security log collection and management platform. It receives Syslog over UDP/TCP, parses and filters events, stores logs, manages devices and templates, and provides a web console for search, alerts, statistics, audit logs, and system operations.

## Features

- Syslog UDP/TCP ingestion.
- Device, group, field mapping, parse template, filter policy, and output template management.
- Log query, trace, cleanup, statistics, dashboard, and audit log views.
- Alert rules, alert records, aggregated alerts, and dispositions.
- HTTP, email, and Syslog forwarding configuration.
- RBAC-based user and role management.
- SQLite by default, with MySQL support for production deployments.

## Default Ports

| Service | Port |
| --- | --- |
| Web console and API | `5080/tcp` |
| Syslog UDP | `5140/udp` |
| Syslog TCP | `5140/tcp` |

Health endpoints:

```bash
curl http://127.0.0.1:5080/healthz
curl http://127.0.0.1:5080/readyz
```

## Download

Download the Linux binary from the GitHub Releases page:

```bash
curl -L -o logcat https://github.com/jincaiw/logcat/releases/download/v0.1.0/logcat-linux-amd64
chmod +x logcat
```

Verify the checksum if `SHA256SUMS` is provided with the release.

## Linux Single-Binary Deployment

Create the installation directories:

```bash
sudo mkdir -p /opt/logcat/{configs,data,logs}
sudo install -m 0755 logcat /opt/logcat/logcat
sudo cp configs/config.yaml /opt/logcat/configs/config.yaml
```

Set the initial admin password and start logcat:

```bash
cd /opt/logcat
export LOGCAT_ADMIN_PASSWORD='change-this-password'
./logcat --config configs/config.yaml
```

Open the web console:

```text
http://<server-ip>:5080
```

Default username:

```text
admin
```

The password is read from `LOGCAT_ADMIN_PASSWORD` during first startup. If the variable is not set, logcat generates a random password and writes it to `data/.admin_password`.

## systemd Service

Create a dedicated user:

```bash
sudo useradd --system --home-dir /opt/logcat --shell /usr/sbin/nologin logcat
sudo chown -R logcat:logcat /opt/logcat
```

Install the service file:

```bash
sudo cp systemd/logcat.service /etc/systemd/system/logcat.service
sudo systemctl daemon-reload
```

Edit `/etc/systemd/system/logcat.service` and set a strong password:

```ini
Environment="LOGCAT_ADMIN_PASSWORD=change-this-password"
```

Start logcat:

```bash
sudo systemctl enable --now logcat
sudo systemctl status logcat
```

View logs:

```bash
sudo journalctl -u logcat -f
```

Verify the service:

```bash
curl http://127.0.0.1:5080/healthz
curl http://127.0.0.1:5080/readyz
```

## Configuration

The default configuration file is `configs/config.yaml`. Common environment overrides:

| Variable | Description |
| --- | --- |
| `LOGCAT_SERVER_HOST` | HTTP bind address |
| `LOGCAT_SERVER_PORT` | HTTP port |
| `LOGCAT_DATABASE_TYPE` | `sqlite` or `mysql` |
| `LOGCAT_SQLITE_PATH` | SQLite database path |
| `LOGCAT_MYSQL_HOST` | MySQL host |
| `LOGCAT_MYSQL_PORT` | MySQL port |
| `LOGCAT_MYSQL_DATABASE` | MySQL database name |
| `LOGCAT_MYSQL_USERNAME` | MySQL username |
| `LOGCAT_MYSQL_PASSWORD` | MySQL password |
| `LOGCAT_ADMIN_PASSWORD` | Initial admin password |
| `LOGCAT_SYSLOG_ENABLED` | Enable Syslog receiver |
| `LOGCAT_SYSLOG_UDP_PORT` | Syslog UDP port |
| `LOGCAT_SYSLOG_TCP_PORT` | Syslog TCP port |

MySQL example:

```bash
LOGCAT_DATABASE_TYPE=mysql \
LOGCAT_MYSQL_HOST=127.0.0.1 \
LOGCAT_MYSQL_PORT=3306 \
LOGCAT_MYSQL_DATABASE=logcat \
LOGCAT_MYSQL_USERNAME=logcat \
LOGCAT_MYSQL_PASSWORD='strong-password' \
LOGCAT_ADMIN_PASSWORD='change-this-password' \
/opt/logcat/logcat --config /opt/logcat/configs/config.yaml
```

## Build From Source

Build the web UI and local binary:

```bash
cd web
npm install
npm run build
cd ..
go build -o logcat .
```

Run tests:

```bash
go test ./...
go vet ./...
cd web && npm run build
```

Linux release binaries are built by the GitHub Actions release workflow on Ubuntu because the project uses CGO for SQLite support.

## License

This project is licensed under the MIT License. See [LICENSE](LICENSE) for details.
