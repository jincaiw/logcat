# logcat v0.2.7

## Highlights

- Removed the promotional `公众号.md` document.
- Cleaned up brand-specific wording from the codebase.
- Kept the standardized Web UI and release workflow from the previous release.
- No functional changes.

## Assets

- `logcat-0.2.7-linux-amd64.tar.gz`
- `logcat-0.2.7-linux-arm64.tar.gz`

## Docker

```bash
docker run -d --name logcat \
  --restart unless-stopped \
  -p 8080:8080 \
  -p 5140:5140/udp \
  -p 5140:5140/tcp \
  -v logcat-data:/app/data \
  qing1205/logcat:0.2.7
```

## Default account

```text
admin / admin123
```

Change the password after first login, or set `LOGCAT_ADMIN_PASSWORD` before the first start.
