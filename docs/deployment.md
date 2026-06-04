# logcat 部署说明

## 默认端口

- Web/API: `5080`
- Syslog UDP: `5140/udp`
- Syslog TCP: `5140/tcp`

## 目录结构

```text
logcat/
├── logcat
├── configs/
│   └── config.yaml
├── data/
│   └── logcat.db
├── logs/
└── systemd/
    └── logcat.service
```

## Linux 单文件部署

1. 构建 Linux 二进制：

```bash
make build-linux
```

2. 将下列文件复制到目标主机：

- `bin/logcat`
- `configs/config.yaml`
- `systemd/logcat.service`

3. 目标主机准备目录：

```bash
sudo mkdir -p /opt/logcat/{configs,data,logs}
sudo cp bin/logcat /opt/logcat/logcat
sudo cp configs/config.yaml /opt/logcat/configs/config.yaml
sudo cp systemd/logcat.service /etc/systemd/system/logcat.service
```

4. 启动服务：

```bash
sudo systemctl daemon-reload
sudo systemctl enable --now logcat
```

5. 验证：

```bash
curl http://127.0.0.1:5080/healthz
```

## Docker Compose 部署

1. 构建并启动：

```bash
docker compose up -d --build
```

2. 验证容器状态：

```bash
docker compose ps
```

3. 验证健康检查：

```bash
curl http://127.0.0.1:5080/healthz
```

## 配置覆盖

支持通过环境变量覆盖关键配置：

- `LOGCAT_SERVER_HOST`
- `LOGCAT_SERVER_PORT`
- `LOGCAT_DATABASE_TYPE`
- `LOGCAT_SQLITE_PATH`
- `LOGCAT_MYSQL_HOST`
- `LOGCAT_MYSQL_PORT`
- `LOGCAT_MYSQL_DATABASE`
- `LOGCAT_MYSQL_USERNAME`
- `LOGCAT_MYSQL_PASSWORD`
- `LOGCAT_SESSION_EXPIRE_HOURS`
- `LOGCAT_MAX_FAILED_LOGIN`
- `LOGCAT_LOCK_DURATION_MINUTES`
- `LOGCAT_RETENTION_DAYS`
- `LOGCAT_UNMATCHED_RETENTION_DAYS`
- `LOGCAT_PARSE_WORKERS`
- `LOGCAT_FILTER_WORKERS`
- `LOGCAT_PUSH_WORKERS`
- `LOGCAT_QUEUE_CAPACITY`
- `LOGCAT_QUEUE_FULL_POLICY`
- `LOGCAT_SYSLOG_ENABLED`
- `LOGCAT_SYSLOG_UDP_PORT`
- `LOGCAT_SYSLOG_TCP_PORT`

## MySQL 模式

使用环境变量切换到 MySQL：

```bash
LOGCAT_DATABASE_TYPE=mysql \
LOGCAT_MYSQL_HOST=127.0.0.1 \
LOGCAT_MYSQL_PORT=3306 \
LOGCAT_MYSQL_DATABASE=logcat \
LOGCAT_MYSQL_USERNAME=logcat \
LOGCAT_MYSQL_PASSWORD=logcat_password \
go run . --config configs/config.yaml
```
