# logcat v0.2.11

## Highlights

- Enabled simultaneous TCP + UDP syslog listening by default.
- Normalized system config protocol to `both` for new and existing installs.
- Fixed empty list responses so `output-templates` and `robots` return `[]` instead of `null`.
- Updated test message copy to use `logcat` branding.
- Cleaned up historical `Syslog2Bot` mentions in visible pages and docs.

## Assets

- `logcat-0.2.11-linux-amd64.tar.gz`
- `logcat-0.2.11-linux-arm64.tar.gz`

## Docker

```bash
docker run -d --name logcat \
  --restart unless-stopped \
  -p 8080:8080 \
  -p 5140:5140/udp \
  -p 5140:5140/tcp \
  -v logcat-data:/app/data \
  qing1205/logcat:0.2.11
```

## Default account

```text
admin / admin123
```

Change the password after first start, or set `LOGCAT_ADMIN_PASSWORD` before the first start.
