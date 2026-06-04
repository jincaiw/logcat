# logcat 测试说明

## 基础测试

```bash
go test ./...
cd web && npm run build
```

## SQLite 模式验收

```bash
go run . --config configs/config.yaml
```

验证项：

- 登录、退出、首次改密
- 用户、角色、设备、模板、策略、推送配置 CRUD
- UDP/TCP Syslog 接收
- 设备识别、解析、筛选、日志入库、日志追踪
- 审计日志、数据统计、指标监控、导入导出

## MySQL 模式验收

```bash
docker compose up -d mysql
LOGCAT_DATABASE_TYPE=mysql \
LOGCAT_MYSQL_HOST=127.0.0.1 \
LOGCAT_MYSQL_PORT=3306 \
LOGCAT_MYSQL_DATABASE=logcat \
LOGCAT_MYSQL_USERNAME=logcat \
LOGCAT_MYSQL_PASSWORD=logcat_password \
go run . --config configs/config.yaml
```

## Docker Compose 验收

```bash
docker compose up -d --build
curl http://127.0.0.1:5080/healthz
```

## 推送链路本地模拟

- HTTP：本地回显服务监听一个测试端口，验证状态码和响应摘要
- Email：本地 SMTP 调试服务接收邮件，验证标题与正文模板
- Syslog：本地 UDP/TCP 服务监听指定端口，验证转发格式
