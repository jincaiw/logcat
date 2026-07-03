# logcat v0.2.8

## Highlights

- Fixed the GitHub Actions release pipeline to build the frontend before running backend tests.
- Standardized the release automation and version sync flow.
- Removed the promotional `公众号.md` document and cleaned up brand-specific wording.
- No functional changes.

## Assets

- `logcat-0.2.8-linux-amd64.tar.gz`
- `logcat-0.2.8-linux-arm64.tar.gz`

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

## Default account

```text
admin / admin123
```

Change the password after first login, or set `LOGCAT_ADMIN_PASSWORD` before the first start.
