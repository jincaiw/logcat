# logcat 安全日志告警流转平台 Spec

## Why
安全运营团队需要一个轻量化的 Web 平台来统一接收、解析、筛选和推送多源安全设备的 Syslog 日志，替代分散的日志处理脚本，实现告警从接收到处置的完整闭环。

## What Changes
- 从零构建完整 logcat 平台（后端 Go + 前端 Vue 3）
- 实现 Syslog UDP/TCP 接收与异步日志处理流水线
- 实现日志解析引擎（JSON、分隔符、键值对、正则、子模板）
- 实现筛选策略、白名单、去重、聚合、高频 IP 识别
- 实现 HTTP / 邮箱 / Syslog 三种推送通道
- 实现用户认证、RBAC 权限控制、审计日志
- 实现设备管理、模板库、字段映射、数据统计
- 实现告警处置闭环、日志全链路追踪
- 支持 SQLite（默认）和 MySQL 8.0+ 双数据库模式
- 支持 Linux 单文件部署和 Docker Compose 部署
- 内嵌前端静态资源到 Go 二进制

## Impact
- Affected specs: 无（新建项目）
- Affected code: 全部新建

---

## ADDED Requirements

### Requirement: 项目脚手架与基础架构
系统 SHALL 基于 Go 1.22+ + Gin + GORM 构建后端，基于 Vue 3 + TypeScript + Naive UI + Vite 构建前端，采用前后端分离开发、单文件同源部署架构。

#### Scenario: 项目初始化
- **WHEN** 开发者执行项目初始化
- **THEN** 创建 Go module 项目结构和 Vue 3 前端项目结构
- **THEN** 配置 GORM 支持 SQLite 和 MySQL 双数据库驱动
- **THEN** 前端构建产物通过 Go embed 内嵌到二进制

### Requirement: 数据库与存储
系统 SHALL 默认使用 SQLite 数据库（启用 WAL 模式），支持通过配置切换为 MySQL 8.0+（utf8mb4 字符集）。系统启动时自动执行数据库迁移和基础数据初始化。

#### Scenario: SQLite 默认模式启动
- **WHEN** 系统以默认配置启动
- **THEN** 在 `data/logcat.db` 创建 SQLite 数据库并启用 WAL
- **THEN** 自动创建所有数据表、索引、默认管理员账号和内置角色权限

#### Scenario: MySQL 模式启动
- **WHEN** 配置 `database.type: mysql` 并提供有效 MySQL 连接信息
- **THEN** 连接 MySQL 8.0+ 并自动创建所有数据表、索引和基础数据

#### Scenario: 数据库迁移幂等
- **WHEN** 重复执行数据库迁移
- **THEN** 不破坏已有数据，不报错

### Requirement: 用户认证
系统 SHALL 提供基于 HttpOnly Cookie Session 的用户认证，支持登录、退出、密码修改、首次登录强制改密、会话过期和登录失败锁定。

#### Scenario: 用户登录
- **WHEN** 用户提供有效用户名和密码
- **THEN** 创建 Session 并设置 HttpOnly + SameSite Cookie，返回用户信息

#### Scenario: 未登录访问
- **WHEN** 未认证请求访问受保护 API
- **THEN** 返回 401 状态码

#### Scenario: 首次登录强制改密
- **WHEN** 默认管理员首次登录且 `must_change_password` 为 true
- **THEN** 强制跳转到密码修改页面，不允许访问其他功能

#### Scenario: 登录失败锁定
- **WHEN** 用户连续登录失败达到阈值
- **THEN** 锁定账号并记录审计日志

#### Scenario: 初始管理员密码
- **WHEN** 系统首次启动
- **THEN** 通过环境变量 `LOGCAT_ADMIN_PASSWORD` 设置初始密码，或引导进入初始化页面设置密码

### Requirement: RBAC 权限控制
系统 SHALL 实现基于角色的权限控制，包括菜单级、按钮级和 API 级权限校验。内置系统管理员、安全运营人员、值守人员、只读审计人员四种角色。

#### Scenario: 菜单权限控制
- **WHEN** 用户登录后访问管理后台
- **THEN** 只显示其角色授权的菜单项

#### Scenario: API 权限校验
- **WHEN** 用户请求无权访问的 API
- **THEN** 返回 403 状态码

### Requirement: Syslog 日志接收
系统 SHALL 支持 UDP 和 TCP 两种协议的 Syslog 日志接收，默认监听端口 5140，支持启停控制。

#### Scenario: UDP Syslog 接收
- **WHEN** 外部设备发送 UDP Syslog 消息到监听端口
- **THEN** 接收消息并放入日志接收队列

#### Scenario: TCP Syslog 接收
- **WHEN** 外部设备建立 TCP 连接并发送 Syslog 消息
- **THEN** 接收消息并按换行分割，放入日志接收队列

#### Scenario: 服务启停
- **WHEN** 管理员通过 Web 界面停止 Syslog 服务
- **THEN** 关闭监听端口，不再接收新日志
- **WHEN** 管理员启动 Syslog 服务
- **THEN** 重新开始监听并接收日志

### Requirement: 异步日志处理流水线
系统 SHALL 采用异步队列 + Worker 池架构处理日志，包括日志接收队列、解析 Worker 池、筛选 Worker 池、数据库写入队列和推送任务队列。HTTP/邮件/Syslog 推送异常不得阻塞日志接收主流程。

#### Scenario: 队列满处理
- **WHEN** 日志接收队列满载
- **THEN** 默认短时间阻塞后丢弃并记录计数

#### Scenario: Worker 异常恢复
- **WHEN** 某个 Worker 异常退出
- **THEN** 系统自动恢复或重建 Worker

#### Scenario: 优雅停止
- **WHEN** 服务收到停止信号
- **THEN** 处理完已接收但未完成的日志后再退出

### Requirement: 设备管理
系统 SHALL 支持安全设备/日志源的 CRUD 管理，包括设备名称、IP 地址（唯一）、分组、解析模板关联、设备模板关联和启用/禁用。

#### Scenario: 新增设备
- **WHEN** 管理员创建设备并填写 IP 地址
- **THEN** IP 地址必须唯一，否则返回 409

#### Scenario: 设备识别
- **WHEN** 系统接收到 Syslog 日志
- **THEN** 根据来源 IP 匹配已登记设备

#### Scenario: 未登记设备
- **WHEN** 日志来源 IP 未匹配任何设备
- **THEN** 按默认策略处理

### Requirement: 设备分组管理
系统 SHALL 支持设备分组的 CRUD 管理，分组名称唯一，支持颜色配置和排序。

### Requirement: 设备模板库
系统 SHALL 支持设备模板的 CRUD 和导入导出，模板包含设备类型、关联解析模板、字段映射和推荐筛选策略。新增设备时可选择模板快速生成配置。

### Requirement: 字段映射文档库
系统 SHALL 支持字段映射文档的 CRUD 和批量导入，维护设备字段说明和标准字段映射（标准字段名、原始字段名、字段说明、字段类型）。

### Requirement: 日志解析引擎
系统 SHALL 支持 Syslog+JSON、纯 JSON、分隔符、键值对、正则表达式和子模板六种解析方式。解析模板支持 CRUD、启用/禁用和实时测试。

#### Scenario: 解析模板测试
- **WHEN** 用户输入示例日志并选择解析模板
- **THEN** 返回解析成功状态、错误信息、字段列表和解析结果

#### Scenario: 子模板解析
- **WHEN** 设备配置了子模板
- **THEN** 根据事件类型或匹配规则选择对应子模板进行解析

### Requirement: 筛选策略引擎
系统 SHALL 支持多条件组合的筛选策略，条件支持等于/不等于、包含、前缀/后缀、多值匹配、正则、存在判断和数值比较等操作符，支持 AND/OR 逻辑组合。策略支持按设备、设备分组、解析模板关联。

#### Scenario: 多条件 AND 筛选
- **WHEN** 策略配置多个条件且逻辑为 AND
- **THEN** 所有条件均满足时才命中

#### Scenario: 策略动作
- **WHEN** 策略命中且动作为 keep
- **THEN** 日志继续流转到后续处理
- **WHEN** 策略动作为 discard
- **THEN** 丢弃日志

#### Scenario: 策略优先级
- **WHEN** 多条策略可匹配同一条日志
- **THEN** 按优先级顺序匹配，首次命中即停止

#### Scenario: 筛选策略测试
- **WHEN** 用户提供测试日志并选择策略
- **THEN** 返回策略是否命中及各条件匹配结果

### Requirement: 白名单策略
系统 SHALL 支持按字段配置白名单过滤，默认字段为源 IP。白名单内的值自动放行或丢弃（根据配置）。

### Requirement: 告警去重
系统 SHALL 支持按时间窗口对相同内容的告警进行去重，减少重复推送。

#### Scenario: 去重窗口内重复告警
- **WHEN** 相同去重键的告警在去重窗口内再次触发
- **THEN** 不重复推送，记录去重计数

### Requirement: 告警聚合
系统 SHALL 支持按源 IP、目的 IP、事件类型、设备、告警等级等字段进行聚合，配置聚合时间窗口和阈值，生成聚合告警。

#### Scenario: 聚合触发
- **WHEN** 聚合窗口内相同聚合键的日志数量达到阈值
- **THEN** 生成一条聚合告警，记录首次/最近出现时间和数量

#### Scenario: 查看聚合告警关联日志
- **WHEN** 用户查看聚合告警详情
- **THEN** 展示关联的原始日志列表

### Requirement: 高频 IP 扫描识别
系统 SHALL 支持按源 IP 和时间窗口识别高频扫描行为，配置时间窗口和告警阈值，生成聚合告警。

#### Scenario: 高频 IP 识别触发
- **WHEN** 指定时间窗口内某源 IP 日志数量超过阈值
- **THEN** 生成高频 IP 告警，记录出现次数、首次和最后出现时间

### Requirement: 字段脱敏
系统 SHALL 支持配置脱敏字段规则，对 IP、账号、Token、URL 参数、邮箱、手机号等字段进行脱敏处理，外发推送内容按模板控制字段范围。敏感配置默认脱敏展示。

### Requirement: 输出模板
系统 SHALL 支持按推送通道（HTTP/邮箱/Syslog）配置输出模板，支持 `{{field_name}}` 模板变量替换，支持按设备类型配置模板。

### Requirement: HTTP 接口推送
系统 SHALL 支持 HTTP POST 方式推送告警到下游业务系统，支持配置接口地址、超时、重试、接收标识列表、认证方式（Bearer/Basic/Custom Header）、请求头和请求体模板。

#### Scenario: HTTP 推送成功
- **WHEN** 系统按告警规则触发 HTTP 推送
- **THEN** 发送 HTTP POST 请求，记录状态码和响应摘要

#### Scenario: HTTP 推送失败重试
- **WHEN** HTTP 推送返回需重试状态码或网络错误
- **THEN** 按配置的重试次数和间隔进行重试，记录每次重试结果

#### Scenario: HTTP 推送测试
- **WHEN** 用户点击推送测试
- **THEN** 发送测试请求并展示请求地址、状态码、响应摘要和失败原因

### Requirement: 邮箱推送
系统 SHALL 支持通过 SMTP 发送告警邮件，支持配置 SMTP Host/Port/Username/Password、发件人、收件人、邮件标题模板和正文模板。

### Requirement: Syslog 转发
系统 SHALL 支持将处理后的日志通过 UDP/TCP 转发到其他 Syslog 服务器，支持 JSON、RFC3164 和 RFC5424 格式，支持转发字段选择。

### Requirement: 告警规则编排
系统 SHALL 支持告警规则的 CRUD 管理，规则关联筛选策略、推送配置和输出模板，支持启用/禁用。

### Requirement: 日志管理
系统 SHALL 保存原始日志、解析后数据、过滤状态、告警状态等完整信息，支持分页查询、按设备/时间/关键词/来源IP/目的IP/事件类型/告警等级/推送状态/解析字段/日志ID 等多维度查询。

#### Scenario: 日志清理
- **WHEN** 管理员执行日志清理操作
- **THEN** 按配置保留天数清理过期日志，写入审计日志

### Requirement: 告警记录
系统 SHALL 保存告警发送记录，包含关联日志 ID、推送配置、告警规则、设备名称、消息内容、发送状态、错误信息、响应摘要和重试信息。

### Requirement: 告警处置闭环
系统 SHALL 支持告警的确认、忽略、关闭、备注等处置操作，记录处置人、处置时间和处置内容，支持查看处置历史。告警状态包括未处理、处理中、已确认、已忽略、已关闭。

### Requirement: 日志 ID 全链路追踪
系统 SHALL 支持按日志 ID 查询完整处理链路，展示接收状态、解析状态、解析模板、解析结果/错误、过滤状态、命中策略、去重/白名单/聚合判断结果和各推送目标的发送结果。

### Requirement: 数据统计
系统 SHALL 支持按筛选策略、时间范围和字段进行 Top N 统计，返回总日志数、唯一值数量、字段值/数量/占比/最后出现时间。支持可统计字段列表展示、CSV 导出和 IP 列表复制。

### Requirement: 配置导入导出
系统 SHALL 支持解析模板、筛选策略、推送配置和设备模板的导入导出（JSON 格式），导出包含版本号和导出时间。导入时名称已存在则更新，不存在则创建。

### Requirement: 系统配置
系统 SHALL 支持配置监听端口、协议、日志保留天数、最大日志大小、告警开关/间隔、未匹配日志策略、默认过滤动作、深色/浅色主题、数据目录、会话过期时间、登录锁定策略、Worker 数量和队列策略。

### Requirement: 系统仪表盘
系统 SHALL 在首页展示服务运行状态、日志接收速率、TCP 连接数、最近接收时间、今日接收总量、今日告警数量、推送成功率、队列积压数量和最近推送失败信息摘要。

### Requirement: 健康检查与指标监控
系统 SHALL 提供 `/healthz` 健康检查、`/readyz` 就绪检查和 `/metrics` 指标接口，展示运行指标（接收量、解析成功率、推送成功率、队列积压、失败计数）。

### Requirement: 审计日志
系统 SHALL 记录登录、退出、失败登录、配置变更、删除、导入导出、清理和处置等关键操作的审计日志，包含用户、操作类型、资源类型、客户端 IP、操作结果和详情。审计日志不得被普通用户删除。

### Requirement: 部署支持
系统 SHALL 支持 Linux 单文件二进制部署（Go embed 内嵌前端静态资源）和 Docker Compose 部署（含 MySQL 8.0+）。提供 config.yaml 配置文件、systemd 服务示例和环境变量覆盖关键配置能力。

#### Scenario: Linux 单文件部署
- **WHEN** 用户上传二进制文件到 Linux 服务器并启动
- **THEN** 同时提供 Web 管理后台、REST API 和 Syslog UDP/TCP 接收服务

#### Scenario: Docker Compose 部署
- **WHEN** 用户执行 `docker-compose up`
- **THEN** 启动 logcat 服务和 MySQL 8.0+ 数据库，完成数据库初始化