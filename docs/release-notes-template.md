# logcat v{{VERSION}}

## Highlights

- 
- 
- 

## Assets

- `logcat-{{VERSION}}-linux-amd64.tar.gz`
- `logcat-{{VERSION}}-linux-arm64.tar.gz`

## Docker

```bash
docker run -d --name logcat \
  --restart unless-stopped \
  -p 8080:8080 \
  -p 5140:5140/udp \
  -p 5140:5140/tcp \
  -v logcat-data:/app/data \
  qing1205/logcat:{{VERSION}}
```

## Default account

```text
admin / admin123
```

Change the password after first login, or set `LOGCAT_ADMIN_PASSWORD` before the first start.
