# logcat v0.2.9

## Highlights

- Refined the web UI into a consistent, cleaner visual system.
- Unified cards, forms, sidebars, dialogs, and login layout.
- Removed inline styles and scattered page-specific styling.
- Improved accessibility labels and i18n coverage.

## Assets

- `logcat-0.2.9-linux-amd64.tar.gz`
- `logcat-0.2.9-linux-arm64.tar.gz`

## Docker

```bash
docker run -d --name logcat \
  --restart unless-stopped \
  -p 8080:8080 \
  -p 5140:5140/udp \
  -p 5140:5140/tcp \
  -v logcat-data:/app/data \
  qing1205/logcat:0.2.9
```

## Default account

```text
admin / admin123
```

Change the password after first start, or set `LOGCAT_ADMIN_PASSWORD` before the first start.
