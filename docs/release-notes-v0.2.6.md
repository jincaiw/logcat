# logcat v0.2.6

## Highlights

- Refined the Web UI with a cleaner, more consistent layout and card system.
- Unified the sidebar, top bar, dashboard, tables, forms, and modal styles.
- Improved the Logs, Stats, Settings, Profile, and Login page presentation.
- Standardized the release workflow with version sync, package builds, and GitHub Actions publishing.

## Assets

- `logcat-0.2.6-linux-amd64.tar.gz`
- `logcat-0.2.6-linux-arm64.tar.gz`

## Docker

```bash
docker run -d --name logcat \
  --restart unless-stopped \
  -p 8080:8080 \
  -p 5140:5140/udp \
  -p 5140:5140/tcp \
  -v logcat-data:/app/data \
  qing1205/logcat:0.2.6
```

## Default account

```text
admin / admin123
```

Change the password after first login, or set `LOGCAT_ADMIN_PASSWORD` before the first start.
