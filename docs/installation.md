# Installation Guide

This guide covers the recommended production installation methods for logcat.

## Recommended: Docker Compose

1. Install Docker and Docker Compose.
2. Download the compose file.

```bash
curl -O https://raw.githubusercontent.com/jincaiw/logcat/v0.2.10/docker-compose.yml
```

3. Optional: edit `docker-compose.yml` and set a strong initial admin password before the first start.

```yaml
environment:
  LOGCAT_ADMIN_PASSWORD: change-me-now
```

4. Start logcat.

```bash
docker compose up -d
```

5. Open the web UI.

```text
http://<server-ip>:8080
```

Default account: `admin / admin123` if no custom password was set.

## One-line Linux install

```bash
curl -fsSL https://raw.githubusercontent.com/jincaiw/logcat/v0.2.10/scripts/install-linux.sh | sudo bash
```

The installer will:

- download the Linux binary package,
- install it to `/opt/logcat`,
- create a systemd service,
- start the service automatically.

Check service status:

```bash
systemctl status logcat
```

View logs:

```bash
journalctl -u logcat -f
```

## Binary package

```bash
tar -xzf logcat-0.2.10-linux-amd64.tar.gz
cd logcat-0.2.10-linux-amd64
./start.sh 8080
```

## Ports

| Port | Protocol | Description |
| --- | --- | --- |
| 8080 | TCP | Web UI and API |
| 5140 | UDP/TCP | Syslog receiver |

## Environment variables

| Variable | Description |
| --- | --- |
| `SYSLG_ALERT_DATA_DIR` | Data directory |
| `SYSLG_ALERT_TEMPLATES_DIR` | Template directory |
| `LOGCAT_OPEN_BROWSER` | Set `1` to open browser automatically |
| `LOGCAT_ADMIN_USERNAME` | Initial admin username |
| `LOGCAT_ADMIN_PASSWORD` | Initial admin password |

## Upgrade

For Docker Compose:

```bash
docker compose pull
docker compose up -d
```

For systemd installation, rerun the installer with the target version:

```bash
curl -fsSL https://raw.githubusercontent.com/jincaiw/logcat/v0.2.10/scripts/install-linux.sh | sudo VERSION=0.2.10 bash
```

## Backup

Back up the data directory or Docker volume.

Docker volume example:

```bash
docker run --rm -v logcat-data:/data -v "$PWD":/backup alpine \
  tar -czf /backup/logcat-data-backup.tar.gz -C /data .
```
