# logcat

[English README](README.md) · [安装指南](docs/installation.md) · [用户手册](docs/user-guide.md) · [项目主页](http://logcat.mujizi.com/)

logcat 是一个轻量级 Syslog 告警处理工具，适用于安全运营、蓝队值守、小型生产环境中的日志接收、解析、过滤、转发与告警通知。

## 功能特性

- 支持 UDP/TCP 接收 Syslog。
- 支持 JSON、Syslog+JSON、分隔符、键值对、正则等解析方式。
- 支持过滤规则、去重、未匹配日志处理。
- 支持飞书、邮件、HTTP 接口、Syslog 转发通知渠道。
- Web 页面管理设备、解析模板、字段映射、过滤策略和通知规则。
- 支持单二进制部署和 Docker 部署。
- 使用 SQLite 本地持久化数据。

## 界面预览

![仪表盘](docs/assets/demo-dashboard.png)

![通知渠道](docs/assets/demo-notifications.png)

![日志查询](docs/assets/demo-logs.png)

## 快速开始

### Docker Compose

```bash
curl -O https://raw.githubusercontent.com/jincaiw/logcat/v0.2.6/docker-compose.yml
docker compose up -d
```

访问：`http://localhost:8080`

默认账号：

```text
用户名：admin
密码：admin123
```

首次登录后请立即修改密码。

### Linux 一键安装

```bash
curl -fsSL https://raw.githubusercontent.com/jincaiw/logcat/v0.2.6/scripts/install-linux.sh | sudo bash
```

访问：`http://<服务器IP>:8080`

### 二进制包安装

从 Release 下载 `logcat-0.2.6-linux-amd64.tar.gz` 或 `logcat-0.2.6-linux-arm64.tar.gz`：

```bash
tar -xzf logcat-0.2.6-linux-amd64.tar.gz
cd logcat-0.2.6-linux-amd64
./start.sh 8080
```

## Docker 运行

```bash
docker run -d --name logcat \
  --restart unless-stopped \
  -p 8080:8080 \
  -p 5140:5140/udp \
  -p 5140:5140/tcp \
  -v logcat-data:/app/data \
  qing1205/logcat:0.2.6
```

首次启动前可设置管理员密码：

```bash
-e LOGCAT_ADMIN_PASSWORD='change-me-now'
```

## 端口说明

| 端口 | 协议 | 说明 |
| --- | --- | --- |
| 8080 | TCP | Web 页面和 API |
| 5140 | UDP/TCP | Syslog 接收端口 |

## 发布流程

```bash
bash scripts/release.sh <version>
git push origin <branch>
git push origin v<version>
```

推送 tag 后，GitHub Actions 会自动构建并发布 Release。详见 [docs/release-process.md](docs/release-process.md)。

## 文档

- [安装指南](docs/installation.md)
- [用户手册](docs/user-guide.md)
- [发布流程](docs/release-process.md)
- [项目主页](http://logcat.mujizi.com/)
