# logcat

[English](README.md) | 中文

logcat 是一个轻量级安全日志采集与管理平台。它支持通过 UDP/TCP 接收 Syslog，完成日志解析、过滤、入库、查询、告警、统计、审计和系统运维管理。

## 功能特性

- Syslog UDP/TCP 日志接收。
- 设备、设备组、字段映射、解析模板、过滤策略、输出模板管理。
- 日志查询、日志追踪、清理、统计、仪表盘和审计日志。
- 告警规则、告警记录、聚合告警和告警处置。
- HTTP、邮件、Syslog 转发配置。
- 基于 RBAC 的用户和角色管理。
- 默认使用 SQLite，并支持 MySQL 生产部署。

## 默认端口

| 服务 | 端口 |
| --- | --- |
| Web 控制台和 API | `5080/tcp` |
| Syslog UDP | `5140/udp` |
| Syslog TCP | `5140/tcp` |

健康检查：

```bash
curl http://127.0.0.1:5080/healthz
curl http://127.0.0.1:5080/readyz
```

## 下载

从 GitHub Releases 下载 Linux 二进制文件：

```bash
curl -L -o logcat https://github.com/jincaiw/logcat/releases/download/v0.1.0/logcat-linux-amd64
chmod +x logcat
```

如果 release 中提供了 `SHA256SUMS`，请同时校验文件完整性。

## Linux 单文件部署

创建安装目录：

```bash
sudo mkdir -p /opt/logcat/{configs,data,logs}
sudo install -m 0755 logcat /opt/logcat/logcat
sudo cp configs/config.yaml /opt/logcat/configs/config.yaml
```

设置初始化管理员密码并启动：

```bash
cd /opt/logcat
export LOGCAT_ADMIN_PASSWORD='change-this-password'
./logcat --config configs/config.yaml
```

访问控制台：

```text
http://<服务器IP>:5080
```

默认用户名：

```text
admin
```

首次启动时会读取 `LOGCAT_ADMIN_PASSWORD` 作为管理员密码。如果未设置该环境变量，logcat 会生成随机密码并写入 `data/.admin_password`。

## systemd 服务部署

创建专用系统用户：

```bash
sudo useradd --system --home-dir /opt/logcat --shell /usr/sbin/nologin logcat
sudo chown -R logcat:logcat /opt/logcat
```

安装服务文件：

```bash
sudo cp systemd/logcat.service /etc/systemd/system/logcat.service
sudo systemctl daemon-reload
```

编辑 `/etc/systemd/system/logcat.service`，设置强密码：

```ini
Environment="LOGCAT_ADMIN_PASSWORD=change-this-password"
```

启动服务：

```bash
sudo systemctl enable --now logcat
sudo systemctl status logcat
```

查看日志：

```bash
sudo journalctl -u logcat -f
```

验证服务：

```bash
curl http://127.0.0.1:5080/healthz
curl http://127.0.0.1:5080/readyz
```

## 配置

默认配置文件是 `configs/config.yaml`。常用环境变量覆盖项：

| 变量 | 说明 |
| --- | --- |
| `LOGCAT_SERVER_HOST` | HTTP 监听地址 |
| `LOGCAT_SERVER_PORT` | HTTP 端口 |
| `LOGCAT_DATABASE_TYPE` | `sqlite` 或 `mysql` |
| `LOGCAT_SQLITE_PATH` | SQLite 数据库路径 |
| `LOGCAT_MYSQL_HOST` | MySQL 地址 |
| `LOGCAT_MYSQL_PORT` | MySQL 端口 |
| `LOGCAT_MYSQL_DATABASE` | MySQL 数据库名 |
| `LOGCAT_MYSQL_USERNAME` | MySQL 用户名 |
| `LOGCAT_MYSQL_PASSWORD` | MySQL 密码 |
| `LOGCAT_ADMIN_PASSWORD` | 初始化管理员密码 |
| `LOGCAT_SYSLOG_ENABLED` | 是否启用 Syslog 接收 |
| `LOGCAT_SYSLOG_UDP_PORT` | Syslog UDP 端口 |
| `LOGCAT_SYSLOG_TCP_PORT` | Syslog TCP 端口 |

MySQL 示例：

```bash
LOGCAT_DATABASE_TYPE=mysql \
LOGCAT_MYSQL_HOST=127.0.0.1 \
LOGCAT_MYSQL_PORT=3306 \
LOGCAT_MYSQL_DATABASE=logcat \
LOGCAT_MYSQL_USERNAME=logcat \
LOGCAT_MYSQL_PASSWORD='strong-password' \
LOGCAT_ADMIN_PASSWORD='change-this-password' \
/opt/logcat/logcat --config /opt/logcat/configs/config.yaml
```

## 从源码构建

构建前端和本地二进制：

```bash
cd web
npm install
npm run build
cd ..
go build -o logcat .
```

运行测试：

```bash
go test ./...
go vet ./...
cd web && npm run build
```

由于项目使用 CGO 支持 SQLite，Linux release 二进制通过 GitHub Actions 在 Ubuntu runner 中构建。

## License

当前尚未声明许可证。向组织外分发前请补充许可证。
