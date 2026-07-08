# logcat v0.2.10

## Highlights

- Fixed Pages version sync so the root static site and the `docs/` mirror are updated together.
- Standardized the version bump flow to cover every published page.
- Kept the web UI refactor and accessibility improvements from the previous release line.

## Assets

- `logcat-0.2.10-linux-amd64.tar.gz`
- `logcat-0.2.10-linux-arm64.tar.gz`

## Docker

```bash
docker run -d --name logcat \
  --restart unless-stopped \
  -p 8080:8080 \
  -p 5140:5140/udp \
  -p 5140:5140/tcp \
  -v logcat-data:/app/data \
  qing1205/logcat:0.2.10
```

## Default account

```text
admin / admin123
```

Change the password after first start, or set `LOGCAT_ADMIN_PASSWORD` before the first start.
