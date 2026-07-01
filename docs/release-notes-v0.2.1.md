# logcat v0.2.1

## Highlights

- Added HTTP notification channel.
- HTTP channel supports configurable URL, timeout, retry count, retry delay, and Notes IDs.
- Updated notification UI and docs for the new channel.
- Bumped release version to v0.2.1.

## Assets

- `logcat-0.2.1-linux-amd64.tar.gz`
- `logcat-0.2.1-linux-arm64.tar.gz`

## Docker

```bash
docker run -d --name logcat \
  --restart unless-stopped \
  -p 8080:8080 \
  -p 5140:5140/udp \
  -p 5140:5140/tcp \
  -v logcat-data:/app/data \
  qing1205/logcat:0.2.1
```

## Default account

```text
admin / admin123
```

Change the password after first login, or set `LOGCAT_ADMIN_PASSWORD` before the first start.
