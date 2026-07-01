# logcat v0.2.2

## Highlights

- Cleaned up legacy changelog wording for old DingTalk/WeWork references.
- Kept the HTTP notification channel introduced in v0.2.1.
- Bumped release version to v0.2.2.

## Assets

- `logcat-0.2.2-linux-amd64.tar.gz`
- `logcat-0.2.2-linux-arm64.tar.gz`

## Docker

```bash
docker run -d --name logcat \
  --restart unless-stopped \
  -p 8080:8080 \
  -p 5140:5140/udp \
  -p 5140:5140/tcp \
  -v logcat-data:/app/data \
  qing1205/logcat:0.2.2
```

## Default account

```text
admin / admin123
```

Change the password after first login, or set `LOGCAT_ADMIN_PASSWORD` before the first start.
