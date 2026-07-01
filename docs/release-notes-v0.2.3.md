# logcat v0.2.3

## Highlights

- Refined the notification-channel UI to keep only placeholder hints for HTTP settings.
- Added HTTP notification channel support and aligned the release docs.
- Improved security headers, request-body limits, and login throttling.
- Added in-memory caching for hot read paths, stats TTL caching, and trace cleanup.
- Improved DB indexes and observability for queue drops and trace cache size.

## Assets

- `logcat-0.2.3-linux-amd64.tar.gz`
- `logcat-0.2.3-linux-arm64.tar.gz`

## Docker

```bash
docker run -d --name logcat \
  --restart unless-stopped \
  -p 8080:8080 \
  -p 5140:5140/udp \
  -p 5140:5140/tcp \
  -v logcat-data:/app/data \
  qing1205/logcat:0.2.3
```

## Default account

```text
admin / admin123
```

Change the password after first login, or set `LOGCAT_ADMIN_PASSWORD` before the first start.
