# logcat v0.2.0

## Highlights

- Removed DingTalk and WeWork notification channels.
- Supported notification channels are now Feishu, Email, and Syslog forwarding.
- Added public health check endpoints: `/healthz` and `/api/health`.
- Added Dockerfile and Docker Compose support.
- Added Linux binary package build script and one-line installer.
- Added GitHub Pages documentation.
- Updated README to English by default with Chinese README link.
- Added demo screenshots.
- Improved production behavior: browser auto-open is disabled by default on Linux and in containers.
- Added environment variables for initial admin username and password.

## Assets

- `logcat-0.2.0-linux-amd64.tar.gz`
- `logcat-0.2.0-linux-arm64.tar.gz`

## Docker

```bash
docker run -d --name logcat \
  --restart unless-stopped \
  -p 8080:8080 \
  -p 5140:5140/udp \
  -p 5140:5140/tcp \
  -v logcat-data:/app/data \
  jincaiw/logcat:0.2.0
```

## Default account

```text
admin / admin123
```

Change the password after first login, or set `LOGCAT_ADMIN_PASSWORD` before the first start.
