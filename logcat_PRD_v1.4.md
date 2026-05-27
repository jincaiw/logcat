# logcat 产品需求文档 PRD

> 文档版本：v1.4  
> 产品名称：logcat  
> 产品形态：Web 版安全日志接收、解析、筛选与推送平台  
> 技术方向：Go 1.22+ + Gin + GORM + Vue 3 + TypeScript + Naive UI + 默认 SQLite，支持 MySQL 8.0+

---

## 1. 产品概述

### 1.1 产品名称

**logcat 安全日志告警流转平台**

### 1.2 产品定位

logcat 是一款面向安全运营、蓝队重保、现场值守和轻量化安全告警流转场景的 Web 平台。系统接收安全设备、服务器、应用系统产生的 Syslog 日志，对日志进行结构化解析、字段映射、清洗、筛选、白名单过滤、去重、聚合、脱敏、统计分析、告警推送和轻量处置，实现多源安全日志到 HTTP 接口、邮箱、Syslog 转发目标等下游系统的统一流转。

产品核心定位为：**轻量化安全日志网关 + 告警编排平台 + HTTP/Syslog 转发中枢 + 轻量告警处置平台**。

### 1.3 核心目标

1. 快速接入防火墙、WAF、EDR、主机安全、态势感知、云安全、应用系统等日志源。
2. 将原始 Syslog 日志解析为结构化字段，提升日志可读性和可处理性。
3. 通过筛选策略、白名单、去重、聚合和高频 IP 识别机制降低无效告警、重复告警和低价值告警。
4. 将高价值告警推送到 HTTP 接口、邮箱或其他 Syslog 服务器。
5. 支持日志追踪、数据统计、历史查询、字段脱敏、配置导入导出和轻量告警处置。
6. 支持 Linux 服务器单文件部署、Docker Compose 部署和内网集中访问。
7. 增加用户认证能力，支持登录、会话管理、角色权限和操作审计。
8. 默认使用 SQLite 数据库，支持通过配置切换为 MySQL 8.0+。
9. 所有功能整体开发、整体联调、整体验收、整体交付，不按优先级或阶段拆分。

---

## 2. 产品边界

### 2.1 当前建设范围

| 模块 | 是否纳入 | 说明 |
|---|---:|---|
| Web 管理平台 | 是 | Web 管理后台、REST API、静态资源同源部署 |
| 用户认证 | 是 | 登录、退出、密码修改、会话管理、角色权限 |
| RBAC 权限控制 | 是 | 菜单、按钮、API 权限控制 |
| 审计日志 | 是 | 登录、退出、失败登录、配置变更、删除、导入导出、清理、处置记录 |
| Syslog 日志接收 | 是 | 支持 UDP/TCP 监听 |
| 设备管理 | 是 | 日志源设备维护、启用/禁用、分组 |
| 设备分组 | 是 | 按区域、系统、设备类型组织日志源 |
| 设备模板库 | 是 | 内置和维护常见设备解析模板、字段映射和推荐策略 |
| 解析模板管理 | 是 | 支持 JSON、分隔符、键值对、正则、子模板等解析方式 |
| 字段映射文档库 | 是 | 维护设备字段说明和标准字段映射 |
| 筛选策略 | 是 | 多条件组合、AND/OR、优先级、动作控制 |
| 白名单策略 | 是 | 按字段配置白名单过滤 |
| 告警去重 | 是 | 按窗口期降低重复推送 |
| 告警聚合 | 是 | 支持按字段、时间窗口和阈值生成聚合告警 |
| 高频 IP 扫描识别 | 是 | 支持按源 IP 和时间窗口识别高频扫描行为 |
| 字段脱敏 | 是 | 支持外发字段、展示字段和敏感配置脱敏 |
| 输出模板 | 是 | 支持不同推送通道的消息模板 |
| HTTP 接口推送 | 是 | 支持超时、重试、接收标识列表、认证、请求头和请求体模板 |
| 邮箱推送 | 是 | 支持 SMTP 推送 |
| Syslog 转发 | 是 | 支持 TCP/UDP 转发，支持 JSON/RFC3164/RFC5424 格式 |
| 日志查询 | 是 | 原始日志、解析字段、筛选状态、告警状态查询 |
| 告警记录 | 是 | 保存推送状态、错误信息、发送时间、重试结果 |
| 告警处置闭环 | 是 | 支持确认、忽略、关闭、备注、处置人、处置时间 |
| 日志 ID 追踪 | 是 | 查看接收、解析、过滤、去重、聚合、推送全链路 |
| 数据统计 | 是 | 字段 Top N、IP 统计、CSV 导出 |
| 健康检查与指标监控 | 是 | 支持健康检查、就绪检查、运行指标和队列积压监控 |
| 配置导入导出 | 是 | 解析模板、筛选策略、推送配置备份迁移 |
| 系统配置 | 是 | 端口、协议、数据库、日志保留、认证策略、主题配置 |
| Linux 单文件部署 | 是 | Go 二进制内嵌前端静态资源，默认 SQLite 数据库 |
| Docker Compose 部署 | 是 | 可选部署方式，支持 logcat 与 MySQL 8.0+ |

### 2.2 产品边界声明

logcat 聚焦安全日志接收、解析、筛选、转发、统计和轻量处置，不替代完整 SIEM 关联分析平台、完整 CMDB、完整 DCIM、网络设备自动配置下发系统、多租户 SaaS 平台、工单 SLA 管理平台和长期海量日志全文检索平台。

---

## 3. 用户角色与权限

### 3.1 用户角色

| 角色 | 权限范围 |
|---|---|
| 系统管理员 | 用户管理、角色权限、系统配置、设备管理、模板管理、策略管理、推送配置、日志查询、数据统计、审计日志、部署配置 |
| 安全运营人员 | 设备查看、模板查看、策略配置、推送配置、日志查询、告警记录、日志追踪、数据统计、告警处置 |
| 值守人员 | 日志查询、告警记录、日志追踪、数据统计、告警确认和备注 |
| 只读审计人员 | 只读查看日志、告警记录、系统配置、审计日志和统计数据 |

### 3.2 用户认证需求

| 编号 | 需求 |
|---|---|
| AUTH-001 | 系统必须提供登录页面 |
| AUTH-002 | 用户必须通过用户名和密码登录后才能访问管理后台 |
| AUTH-003 | 系统必须支持退出登录 |
| AUTH-004 | 系统必须支持密码修改 |
| AUTH-005 | 系统必须支持管理员重置用户密码 |
| AUTH-006 | 系统必须支持用户启用、禁用和锁定 |
| AUTH-007 | 系统必须支持基于角色的权限控制 RBAC |
| AUTH-008 | 系统必须支持会话过期时间配置 |
| AUTH-009 | 系统必须支持登录失败次数限制 |
| AUTH-010 | 系统必须记录登录、退出、失败登录、关键配置变更等审计日志 |
| AUTH-011 | 密码必须使用不可逆哈希算法存储，禁止明文存储 |
| AUTH-012 | Web API 必须进行认证校验，未认证请求返回 401 |
| AUTH-013 | 无权限访问必须返回 403 |
| AUTH-014 | 密码哈希算法应使用 bcrypt 或 Argon2id |
| AUTH-015 | 默认管理员首次登录必须修改初始密码 |
| AUTH-016 | Web 管理后台默认使用基于 HttpOnly Cookie 的 Session 认证 |
| AUTH-017 | Session Cookie 必须设置 HttpOnly 和 SameSite |
| AUTH-018 | 系统应支持配置会话空闲超时时间 |
| AUTH-019 | 连续登录失败达到阈值后应锁定账号或临时禁止登录 |
| AUTH-020 | 审计日志不得被普通用户删除 |
| AUTH-021 | 首次启动时系统应自动创建默认管理员账号 |
| AUTH-022 | 默认管理员账号用户名为 `admin` |
| AUTH-023 | 默认管理员初始密码不得在文档或代码中硬编码固定明文 |
| AUTH-024 | 系统应支持通过环境变量 `LOGCAT_ADMIN_PASSWORD` 设置初始管理员密码 |
| AUTH-025 | 未设置初始管理员密码时，系统应在首次初始化页面引导设置密码 |

---

## 4. 典型使用场景

### 4.1 重保期间快速告警推送

安全设备将 Syslog 日志发送到 logcat。logcat 根据来源 IP 识别设备，使用解析模板完成字段解析，再通过筛选策略识别高危告警，最终推送到 HTTP 接口、邮箱或 Syslog 目标系统。

### 4.2 多厂商日志标准化

不同厂商设备日志格式不统一。logcat 通过字段映射文档库、解析模板和值转换规则，将不同格式日志转换为统一字段，便于筛选、统计和下游系统消费。

### 4.3 HTTP 接口对接业务系统

用户配置 HTTP 接口推送地址、超时时间、重试次数、重试间隔、接收标识列表、认证方式、请求头和请求体模板。系统筛选出有效告警后，按照指定请求格式推送到业务消息接口、告警平台或内部通知系统。

### 4.4 高频攻击 IP 统计与识别

系统根据历史日志中的源 IP、目的 IP、事件类型、威胁等级等字段进行 Top N 聚合统计，并在指定时间窗口内识别高频扫描 IP，生成聚合告警，支持 CSV 导出和 IP 列表复制。

### 4.5 日志全链路追踪

用户输入日志 ID，系统展示该日志从接收、解析、过滤、白名单判断、去重、聚合到推送的全过程，用于定位日志未推送、解析失败、策略未命中、接口发送失败等问题。

### 4.6 告警轻量处置

值守人员在告警记录中对告警进行确认、忽略、备注、关闭等操作，系统记录处置人、处置时间和处置内容，并写入审计日志。

---

## 5. 技术栈与架构

### 5.1 技术栈要求

| 层级 | 技术 | 说明 |
|---|---|---|
| 后端语言 | Go 1.22+ | 用于实现 Web API、Syslog 接收、日志处理、推送编排和单文件部署 |
| Web 框架 | Gin | 提供 REST API、静态资源服务、认证中间件、权限中间件和健康检查接口 |
| 数据库 ORM | GORM | 用于数据模型映射、迁移和查询 |
| 数据库 | 默认 SQLite，支持 MySQL 8.0+ | SQLite 用于默认单文件部署；MySQL 8.0+ 用于生产扩展部署 |
| 缓存/限流 | 内存缓存 | 用于登录失败限制、短期去重窗口和基础限流 |
| 前端框架 | Vue 3 | 用于构建 Web 管理后台 |
| 前端语言 | TypeScript | 用于提升前端类型安全和可维护性 |
| UI 组件 | Naive UI | 用于管理后台 UI 组件 |
| 状态管理 | Pinia | 用于前端全局状态管理 |
| 路由 | Vue Router | 用于前端页面路由管理 |
| 构建工具 | Vite | 用于前端构建 |
| 认证方案 | HttpOnly Cookie Session，API 可扩展 JWT | Web 管理后台默认使用 HttpOnly Cookie Session |
| 部署方式 | Linux 单文件部署 / Docker Compose | Linux 单文件部署为默认方式；Docker Compose 为可选方式 |

### 5.2 系统架构

logcat 采用 **前后端分离开发、单文件同源部署** 的 Web 架构：

```text
浏览器
  ↓ HTTP/HTTPS
logcat Web 服务
  ├─ 静态资源服务
  ├─ REST API 服务
  ├─ 用户认证与权限模块
  ├─ Syslog UDP/TCP 接收服务
  ├─ 日志接收队列
  ├─ 日志解析引擎
  ├─ 筛选策略引擎
  ├─ 去重与白名单引擎
  ├─ 高频 IP 识别模块
  ├─ 告警聚合模块
  ├─ 字段脱敏模块
  ├─ 数据库写入模块
  ├─ 推送编排模块
  │   ├─ HTTP 接口推送
  │   ├─ 邮箱推送
  │   └─ Syslog 转发
  ├─ 告警处置模块
  ├─ 日志查询与统计模块
  ├─ 配置导入导出模块
  ├─ 健康检查与指标监控模块
  └─ 审计日志模块
  ↓
默认 SQLite / 可选 MySQL 8.0+
```

### 5.3 日志异步处理架构

Syslog 接收、日志解析、筛选、数据库写入和告警推送必须解耦处理，HTTP 推送、邮件发送、Syslog 转发和数据库慢写入不得阻塞 Syslog 接收主流程。

```text
Syslog UDP/TCP 接收服务
  ↓
日志接收队列
  ↓
解析 Worker 池
  ↓
筛选/白名单/去重 Worker 池
  ↓
高频 IP 识别 / 告警聚合 / 字段脱敏
  ↓
数据库写入队列
  ↓
推送任务队列
  ↓
HTTP / 邮箱 / Syslog 推送 Worker
  ↓
告警记录 / 日志追踪 / 统计数据 / 审计日志
```

### 5.4 日志处理运行机制要求

| 编号 | 需求 |
|---|---|
| RUN-001 | 日志接收队列、解析队列、数据库写入队列、推送队列均应支持容量配置 |
| RUN-002 | 队列满时系统应支持阻塞接收、丢弃新日志、仅保存原始日志三种策略 |
| RUN-003 | 队列满时默认策略为短时间阻塞，超时后记录丢弃计数和错误日志 |
| RUN-004 | 数据库写入应支持批量写入 |
| RUN-005 | Worker 异常退出后系统应自动恢复或重建 Worker |
| RUN-006 | 服务停止时应尝试优雅关闭，处理已接收但未完成的日志 |
| RUN-007 | 系统应统计接收失败、解析失败、筛选失败、推送失败和数据库写入失败次数 |
| RUN-008 | HTTP 推送、邮件推送、Syslog 转发异常不得阻塞 Syslog 接收主流程 |

---

## 6. 总体业务流程

```text
安全设备 / 应用系统 / 主机日志
  ↓ Syslog UDP/TCP
logcat 日志接收服务
  ↓
日志接收队列
  ↓
设备识别 / 原始日志入库
  ↓
解析模板处理
  ↓
字段映射 / 值转换
  ↓
筛选策略匹配
  ↓
白名单判断 / 去重判断
  ↓
高频 IP 识别 / 告警聚合 / 字段脱敏
  ↓
告警规则匹配
  ↓
输出模板渲染
  ↓
HTTP 接口推送 / 邮箱推送 / Syslog 转发
  ↓
告警记录 / 告警处置 / 日志追踪 / 数据统计 / 审计日志
```

---

## 7. 功能需求

### 7.1 系统状态模块

| 编号 | 需求 |
|---|---|
| FR-001 | 系统应支持启动 Syslog 接收服务 |
| FR-002 | 系统应支持停止 Syslog 接收服务 |
| FR-003 | 系统应支持配置监听端口，默认建议为 5140 |
| FR-004 | 系统应支持 UDP 和 TCP 两种协议 |
| FR-005 | 系统应展示服务运行状态 |
| FR-006 | 系统应展示日志接收速率 |
| FR-007 | 系统应展示当前 TCP 连接数 |
| FR-008 | 系统应展示最近接收日志时间 |
| FR-009 | 系统应展示今日接收总量、今日告警数量、推送成功率 |
| FR-010 | 系统应展示日志接收队列积压数量 |
| FR-011 | 系统应展示推送队列积压数量 |
| FR-012 | 系统应展示最近一次推送失败信息摘要 |

### 7.2 用户管理模块

| 编号 | 需求 |
|---|---|
| FR-013 | 系统应支持新增用户 |
| FR-014 | 系统应支持编辑用户 |
| FR-015 | 系统应支持禁用用户 |
| FR-016 | 系统应支持删除用户 |
| FR-017 | 系统应支持重置用户密码 |
| FR-018 | 系统应支持分配角色 |
| FR-019 | 系统应支持查看用户最近登录时间 |
| FR-020 | 系统应支持查看用户状态 |
| FR-021 | 系统初始化时应创建默认管理员账号，并要求首次登录修改密码 |
| FR-022 | 系统应支持解除账号锁定 |
| FR-023 | 系统应支持强制用户下次登录修改密码 |

### 7.3 角色权限模块

| 编号 | 需求 |
|---|---|
| FR-024 | 系统应内置系统管理员、安全运营人员、值守人员、只读审计人员角色 |
| FR-025 | 系统应支持菜单级权限控制 |
| FR-026 | 系统应支持按钮级权限控制 |
| FR-027 | 系统应支持 API 级权限校验 |
| FR-028 | 系统应支持角色权限查看 |
| FR-029 | 关键配置操作必须校验管理员或安全运营权限 |

### 7.4 设备管理模块

| 编号 | 需求 |
|---|---|
| FR-030 | 系统应支持新增设备 |
| FR-031 | 系统应支持编辑设备 |
| FR-032 | 系统应支持删除设备 |
| FR-033 | 系统应支持设备启用/禁用 |
| FR-034 | 系统应支持设备名称配置 |
| FR-035 | 系统应支持设备 IP 地址配置 |
| FR-036 | 系统应支持设备描述配置 |
| FR-037 | 系统应支持设备分组 |
| FR-038 | 系统应支持设备与解析模板关联 |
| FR-039 | 设备 IP 地址必须唯一 |
| FR-040 | 未登记设备日志应支持按默认策略处理 |

### 7.5 设备分组模块

| 编号 | 需求 |
|---|---|
| FR-041 | 系统应支持新增设备分组 |
| FR-042 | 系统应支持编辑设备分组 |
| FR-043 | 系统应支持删除设备分组 |
| FR-044 | 设备分组名称必须唯一 |
| FR-045 | 系统应支持分组颜色配置 |
| FR-046 | 系统应支持分组排序 |
| FR-047 | 筛选策略和告警策略应支持按设备分组关联 |

### 7.6 映射文档库模块

| 编号 | 需求 |
|---|---|
| FR-048 | 系统应支持新增字段映射文档 |
| FR-049 | 系统应支持编辑字段映射文档 |
| FR-050 | 系统应支持删除字段映射文档 |
| FR-051 | 系统应支持按设备类型查询字段映射文档 |
| FR-052 | 系统应支持批量导入字段映射关系 |
| FR-053 | 解析模板应支持引用字段映射文档 |
| FR-054 | 字段映射文档应支持标准字段名、原始字段名、字段说明、字段类型 |

### 7.7 日志解析模块

#### 7.7.1 支持解析类型

| 解析类型 | 说明 |
|---|---|
| Syslog + JSON | 适用于 Syslog 头部 + JSON 内容日志 |
| 纯 JSON | 适用于直接发送 JSON 日志的系统 |
| 分隔符 | 适用于固定分隔符日志 |
| 键值对分隔 | 适用于 `key:value` 或 `key=value` 格式 |
| 正则表达式 | 适用于非结构化日志 |
| 子模板 | 适用于同一设备按事件类型使用不同解析规则 |

#### 7.7.2 功能要求

| 编号 | 需求 |
|---|---|
| FR-055 | 系统应支持新增解析模板 |
| FR-056 | 系统应支持编辑解析模板 |
| FR-057 | 系统应支持删除解析模板 |
| FR-058 | 系统应支持启用/禁用解析模板 |
| FR-059 | 系统应支持解析类型选择 |
| FR-060 | 系统应支持 Header 正则配置 |
| FR-061 | 系统应支持字段映射配置 |
| FR-062 | 系统应支持值转换规则配置 |
| FR-063 | 系统应支持示例日志配置 |
| FR-064 | 系统应支持设备类型配置 |
| FR-065 | 系统应支持分隔符配置 |
| FR-066 | 系统应支持解析模板实时测试 |
| FR-067 | 系统应返回解析成功状态、错误信息、字段列表和解析结果 |

### 7.8 筛选策略模块

#### 7.8.1 支持操作符

| 类型 | 操作符 |
|---|---|
| 等于/不等于 | `==`、`!=` |
| 字符串包含 | `contains`、`not_contains` |
| 前缀/后缀 | `starts_with`、`ends_with` |
| 多值匹配 | `in`、`not_in` |
| 正则 | `regex` |
| 存在判断 | `exists`、`not_exists` |
| 数值比较 | `>`、`>=`、`<`、`<=` |

#### 7.8.2 功能要求

| 编号 | 需求 |
|---|---|
| FR-068 | 系统应支持新增筛选策略 |
| FR-069 | 系统应支持编辑筛选策略 |
| FR-070 | 系统应支持删除筛选策略 |
| FR-071 | 系统应支持启用/禁用筛选策略 |
| FR-072 | 系统应支持按设备关联筛选策略 |
| FR-073 | 系统应支持按设备分组关联筛选策略 |
| FR-074 | 系统应支持按解析模板关联筛选策略 |
| FR-075 | 系统应支持多条件组合 |
| FR-076 | 系统应支持 AND/OR 条件逻辑 |
| FR-077 | 系统应支持白名单配置 |
| FR-078 | 系统应支持白名单字段配置，默认字段建议为源 IP 字段 |
| FR-079 | 系统应支持 keep/discard 动作 |
| FR-080 | 系统应支持策略优先级 |
| FR-081 | 系统应支持告警去重开关 |
| FR-082 | 系统应支持配置去重时间窗口 |
| FR-083 | 系统应支持未匹配日志丢弃配置 |
| FR-084 | 系统应支持筛选策略测试 |

### 7.9 输出模板模块

| 编号 | 需求 |
|---|---|
| FR-085 | 系统应支持新增输出模板 |
| FR-086 | 系统应支持编辑输出模板 |
| FR-087 | 系统应支持删除输出模板 |
| FR-088 | 系统应支持启用/禁用输出模板 |
| FR-089 | 系统应支持按推送通道配置输出模板 |
| FR-090 | 系统应支持模板变量替换 |
| FR-091 | 系统应展示可用字段列表 |
| FR-092 | 系统应支持按设备类型配置模板 |
| FR-093 | 告警规则应可关联输出模板 |

### 7.10 推送配置模块

#### 7.10.1 支持推送通道

| 推送通道 | 说明 |
|---|---|
| HTTP 接口推送 | 通过 HTTP POST 将告警推送到业务系统或内部消息平台 |
| 邮箱推送 | 通过 SMTP 发送告警邮件 |
| Syslog 转发 | 将处理后的日志继续转发到其他 Syslog 服务端 |

#### 7.10.2 通用功能要求

| 编号 | 需求 |
|---|---|
| FR-094 | 系统应支持新增推送配置 |
| FR-095 | 系统应支持编辑推送配置 |
| FR-096 | 系统应支持删除推送配置 |
| FR-097 | 系统应支持启用/禁用推送配置 |
| FR-098 | 系统应支持推送配置连通性测试 |
| FR-099 | 一个推送配置应支持被多个告警规则引用 |
| FR-100 | 推送失败应记录失败原因 |
| FR-101 | 推送失败应按配置重试 |
| FR-102 | 推送敏感参数应加密存储或脱敏展示 |

### 7.11 HTTP 接口推送模块

#### 7.11.1 配置项

| 配置项 | 类型 | 默认值/示例 | 说明 |
|---|---|---|---|
| 接口名称 | string | 安全告警 HTTP 推送 | 推送配置名称 |
| 接口地址 | string | `http://168.63.6.81:8080/cib-message/public/service/sendwbg.do` | HTTP 推送目标地址 |
| 请求方法 | string | POST | 默认 POST |
| 超时时间 | int | 3 | 单次请求超时时间，单位秒 |
| 重试次数 | int | 3 | 发送失败后的重试次数 |
| 重试间隔 | int | 2 | 重试间隔，单位秒 |
| 接收标识列表 | string[] | `420102,420809` | 下游系统接收人、机构、群组或消息路由标识 |
| 请求头 | JSON | `{}` | 可选，支持配置 Content-Type、认证 Token 等 |
| 请求体模板 | JSON/Text | 模板内容 | 支持变量替换 |
| 成功状态码 | int[] | `[200]` | 用于判断 HTTP 推送是否成功 |
| 成功响应关键字 | string | 空 | 可选，用于根据响应体关键字判断成功 |
| 认证方式 | string | none | 支持 none、bearer、basic、custom_header |
| Token | string | 空 | Bearer Token 或自定义认证 Token |
| Content-Type | string | application/json | HTTP 请求内容类型 |
| 需重试状态码 | int[] | `[500,502,503,504]` | 命中后触发重试 |
| 响应记录最大长度 | int | 2048 | 响应内容按最大长度保存摘要 |
| 启用状态 | bool | true | 是否启用该 HTTP 推送配置 |

> HTTP 接口推送参考配置：`SENDWBG_URL`、`SENDWBG_TIMEOUT`、`SENDWBG_RETRY_COUNT`、`SENDWBG_RETRY_DELAY`、`NOTESIDS`。

#### 7.11.2 功能要求

| 编号 | 需求 |
|---|---|
| FR-103 | 系统应支持配置 HTTP 接口地址 |
| FR-104 | 系统应支持配置 HTTP 请求超时时间 |
| FR-105 | 系统应支持配置 HTTP 失败重试次数 |
| FR-106 | 系统应支持配置 HTTP 重试间隔 |
| FR-107 | 系统应支持配置接收标识列表 |
| FR-108 | 系统应支持配置 HTTP 请求头 |
| FR-109 | 系统应支持配置请求体模板 |
| FR-110 | 系统应支持 JSON 请求体 |
| FR-111 | 系统应支持模板变量替换日志字段 |
| FR-112 | 系统应支持 HTTP 推送测试 |
| FR-113 | 系统应记录 HTTP 响应状态码 |
| FR-114 | 系统应记录 HTTP 响应内容摘要 |
| FR-115 | 系统应记录每次重试结果 |
| FR-116 | HTTP 推送失败不应阻塞日志接收主流程 |
| FR-117 | HTTP 接口地址中的敏感参数应支持脱敏展示 |
| FR-118 | HTTP 推送应支持配置成功状态码 |
| FR-119 | HTTP 推送应支持配置响应体成功关键字 |
| FR-120 | HTTP 推送应支持 Bearer Token、Basic Auth 和自定义 Header |
| FR-121 | HTTP 推送应支持配置 Content-Type |
| FR-122 | HTTP 推送响应内容应按长度限制保存摘要 |
| FR-123 | HTTP 推送重试应支持按网络错误、超时和指定状态码触发 |

#### 7.11.3 推荐请求体格式

```json
{
  "notesIds": ["420102", "420809"],
  "title": "安全告警通知",
  "content": "{{message}}",
  "level": "{{severity}}",
  "sourceIp": "{{src_ip}}",
  "destinationIp": "{{dst_ip}}",
  "eventType": "{{event_type}}",
  "deviceName": "{{device_name}}",
  "occurTime": "{{event_time}}",
  "logId": "{{log_id}}"
}
```

### 7.12 邮箱推送模块

| 编号 | 需求 |
|---|---|
| FR-124 | 系统应支持 SMTP Host 配置 |
| FR-125 | 系统应支持 SMTP Port 配置 |
| FR-126 | 系统应支持 SMTP Username 配置 |
| FR-127 | 系统应支持 SMTP Password 配置 |
| FR-128 | 系统应支持发件人配置 |
| FR-129 | 系统应支持收件人配置 |
| FR-130 | 系统应支持邮件标题模板 |
| FR-131 | 系统应支持邮件正文模板 |
| FR-132 | 系统应支持邮件推送测试 |

### 7.13 Syslog 转发模块

| 编号 | 需求 |
|---|---|
| FR-133 | 系统应支持配置转发目标 Host |
| FR-134 | 系统应支持配置转发目标 Port |
| FR-135 | 系统应支持 UDP/TCP 转发 |
| FR-136 | 系统应支持 JSON 格式转发 |
| FR-137 | 系统应支持 RFC3164 格式转发 |
| FR-138 | 系统应支持 RFC5424 格式转发 |
| FR-139 | 系统应支持转发字段选择 |
| FR-140 | 系统应支持转发测试 |
| FR-141 | UDP 报文超限时应给出明确提示或自动截断 |

### 7.14 告警规则模块

| 编号 | 需求 |
|---|---|
| FR-142 | 系统应支持新增告警规则 |
| FR-143 | 系统应支持编辑告警规则 |
| FR-144 | 系统应支持删除告警规则 |
| FR-145 | 系统应支持启用/禁用告警规则 |
| FR-146 | 告警规则必须关联推送配置 |
| FR-147 | 告警规则必须关联筛选策略 |
| FR-148 | 告警规则可关联输出模板 |
| FR-149 | 告警规则应支持推送通道类型选择 |
| FR-150 | 系统应支持按推送配置查询关联规则 |

### 7.15 日志管理模块

| 编号 | 需求 |
|---|---|
| FR-151 | 系统应保存原始日志 |
| FR-152 | 系统应保存解析后数据 |
| FR-153 | 系统应保存关键解析字段 |
| FR-154 | 系统应保存过滤状态 |
| FR-155 | 系统应保存匹配策略 ID |
| FR-156 | 系统应保存告警状态 |
| FR-157 | 系统应保存告警规则 ID |
| FR-158 | 系统应保存日志来源 IP |
| FR-159 | 系统应保存设备名称 |
| FR-160 | 系统应保存 Syslog facility 和 severity |
| FR-161 | 系统应支持分页查询日志 |
| FR-162 | 系统应支持按设备查询日志 |
| FR-163 | 系统应支持按时间范围查询日志 |
| FR-164 | 系统应支持按关键词查询日志 |
| FR-165 | 系统应支持清理指定天数前日志 |
| FR-166 | 系统应支持清空全部日志 |
| FR-167 | 系统应支持未匹配日志数量统计 |
| FR-168 | 系统应支持清理未匹配日志 |
| FR-169 | 系统应支持按来源 IP 查询日志 |
| FR-170 | 系统应支持按目的 IP 查询日志 |
| FR-171 | 系统应支持按事件类型查询日志 |
| FR-172 | 系统应支持按告警等级查询日志 |
| FR-173 | 系统应支持按推送状态查询日志 |
| FR-174 | 系统应支持按解析字段键值查询日志 |
| FR-175 | 系统应支持按日志 ID 精确查询日志 |

### 7.16 告警记录模块

| 编号 | 需求 |
|---|---|
| FR-176 | 系统应保存告警发送记录 |
| FR-177 | 系统应记录关联日志 ID |
| FR-178 | 系统应记录推送配置 ID |
| FR-179 | 系统应记录告警规则 ID |
| FR-180 | 系统应记录设备名称 |
| FR-181 | 系统应记录告警消息内容 |
| FR-182 | 系统应记录发送状态 |
| FR-183 | 系统应记录错误信息 |
| FR-184 | 系统应记录发送时间 |
| FR-185 | 系统应记录推送通道类型 |
| FR-186 | 系统应记录请求摘要、响应状态码、响应摘要和重试次数 |
| FR-187 | 系统应支持分页查看告警记录 |

### 7.17 日志 ID 追踪模块

| 编号 | 需求 |
|---|---|
| FR-188 | 系统应支持按日志 ID 查询处理链路 |
| FR-189 | 系统应展示日志接收状态 |
| FR-190 | 系统应展示日志解析状态 |
| FR-191 | 系统应展示使用的解析模板 |
| FR-192 | 系统应展示解析结果 |
| FR-193 | 系统应展示解析错误 |
| FR-194 | 系统应展示过滤状态 |
| FR-195 | 系统应展示命中的筛选策略 |
| FR-196 | 系统应展示去重判断结果 |
| FR-197 | 系统应展示白名单判断结果 |
| FR-198 | 系统应展示告警聚合结果 |
| FR-199 | 系统应展示各推送目标的发送结果 |
| FR-200 | 系统应展示 HTTP 推送状态码、重试次数和失败原因 |

### 7.18 数据统计模块

| 编号 | 需求 |
|---|---|
| FR-201 | 系统应支持按筛选策略统计 |
| FR-202 | 系统应支持按开始时间和结束时间统计 |
| FR-203 | 系统应支持选择统计字段 |
| FR-204 | 系统应支持配置 Top N |
| FR-205 | 系统应返回总日志数 |
| FR-206 | 系统应返回唯一值数量 |
| FR-207 | 系统应返回字段值、数量、占比和最后出现时间 |
| FR-208 | 系统应支持可统计字段列表 |
| FR-209 | 系统应支持导出 CSV |
| FR-210 | 系统应支持复制 IP 列表 |
| FR-211 | 数据统计查询必须选择时间范围 |

### 7.19 配置导入导出模块

| 编号 | 需求 |
|---|---|
| FR-212 | 系统应支持导出解析模板 |
| FR-213 | 系统应支持导出筛选策略 |
| FR-214 | 系统应支持导出推送配置 |
| FR-215 | 系统应支持导入解析模板 |
| FR-216 | 系统应支持导入筛选策略 |
| FR-217 | 系统应支持导入推送配置 |
| FR-218 | 导出内容应包含版本号和导出时间 |
| FR-219 | 导入时如名称已存在，应更新已有配置 |
| FR-220 | 导入时如名称不存在，应创建新配置 |
| FR-221 | 导入结果应返回成功状态、消息、数量和错误列表 |
| FR-222 | 导入 JSON 配置时必须校验格式、版本和数据合法性 |

### 7.20 系统配置模块

| 编号 | 需求 |
|---|---|
| FR-223 | 系统应支持配置监听端口 |
| FR-224 | 系统应支持配置监听协议 |
| FR-225 | 系统应支持配置日志保留天数 |
| FR-226 | 系统应支持配置最大日志大小 |
| FR-227 | 系统应支持配置告警开关 |
| FR-228 | 系统应支持配置告警间隔 |
| FR-229 | 系统应支持配置未匹配日志保留天数 |
| FR-230 | 系统应支持配置未匹配日志告警开关 |
| FR-231 | 系统应支持配置默认过滤动作 |
| FR-232 | 系统应支持深色/浅色主题 |
| FR-233 | 系统应支持数据目录和配置目录配置 |
| FR-234 | 系统应支持认证会话过期时间配置 |
| FR-235 | 系统应支持登录失败锁定策略配置 |
| FR-236 | 系统应支持配置日志处理 Worker 数量 |
| FR-237 | 系统应支持配置推送 Worker 数量 |
| FR-238 | 系统应支持配置队列满时的处理策略 |

### 7.21 高频 IP 扫描识别模块

| 编号 | 需求 |
|---|---|
| FR-239 | 系统应支持按源 IP 统计单位时间内的日志数量 |
| FR-240 | 系统应支持配置高频 IP 识别时间窗口 |
| FR-241 | 系统应支持配置高频 IP 告警阈值 |
| FR-242 | 系统应支持识别高频扫描 IP 并生成聚合告警 |
| FR-243 | 系统应支持查看高频 IP 明细、出现次数、首次出现时间、最后出现时间 |
| FR-244 | 系统应支持按设备、设备分组、事件类型限制高频 IP 识别范围 |

### 7.22 告警聚合模块

| 编号 | 需求 |
|---|---|
| FR-245 | 系统应支持按源 IP、目的 IP、事件类型、设备、告警等级进行聚合 |
| FR-246 | 系统应支持配置聚合时间窗口 |
| FR-247 | 系统应支持配置聚合阈值 |
| FR-248 | 系统应支持生成聚合告警摘要 |
| FR-249 | 系统应支持查看聚合告警关联的原始日志列表 |
| FR-250 | 系统应支持聚合告警推送和记录 |

### 7.23 字段脱敏模块

| 编号 | 需求 |
|---|---|
| FR-251 | 系统应支持配置脱敏字段 |
| FR-252 | 系统应支持对 IP、账号、Token、URL 参数、邮箱、手机号等字段进行脱敏 |
| FR-253 | 系统应支持外发推送内容按模板控制字段范围 |
| FR-254 | 系统应支持敏感配置列表和详情默认脱敏展示 |
| FR-255 | 系统应支持管理员查看或重置敏感配置 |
| FR-256 | 字段脱敏配置变更必须写入审计日志 |

### 7.24 告警处置闭环模块

| 编号 | 需求 |
|---|---|
| FR-257 | 系统应支持告警状态管理 |
| FR-258 | 告警状态应包括未处理、处理中、已确认、已忽略、已关闭 |
| FR-259 | 系统应支持为告警添加处置备注 |
| FR-260 | 系统应记录处置人和处置时间 |
| FR-261 | 系统应支持查看告警处置历史 |
| FR-262 | 告警处置操作必须写入审计日志 |

### 7.25 设备模板库模块

| 编号 | 需求 |
|---|---|
| FR-263 | 系统应支持维护设备模板 |
| FR-264 | 设备模板应包含设备类型、解析模板、字段映射和推荐筛选策略 |
| FR-265 | 系统应支持新增、编辑、删除、启用、禁用设备模板 |
| FR-266 | 新增设备时应支持选择设备模板快速生成配置 |
| FR-267 | 系统应支持导入和导出设备模板 |

### 7.26 健康检查与指标监控模块

| 编号 | 需求 |
|---|---|
| FR-268 | 系统应提供 `/healthz` 健康检查接口 |
| FR-269 | 系统应提供 `/readyz` 就绪检查接口 |
| FR-270 | 系统应支持展示日志接收量、解析成功率、推送成功率等运行指标 |
| FR-271 | 系统应支持展示日志接收队列、写入队列和推送队列积压数量 |
| FR-272 | 系统应支持记录数据库写入失败次数和推送失败次数 |
| FR-273 | 系统应支持指标接口 `/metrics`，用于输出运行指标 |

---

## 8. 页面需求

| 页面 | 主要功能 |
|---|---|
| 登录页 | 用户名密码登录、错误提示、首次登录修改密码 |
| 首页仪表盘 | 服务状态、接收速率、今日日志、今日告警、推送成功率、队列积压、健康状态 |
| 系统状态 | 启停 Syslog 服务、查看监听协议、端口、连接状态、队列状态、Worker 状态 |
| 用户管理 | 用户 CRUD、重置密码、启用禁用、解除锁定、角色分配 |
| 角色权限 | 角色列表、权限查看、菜单权限、按钮权限 |
| 设备管理 | 设备 CRUD、设备分组、启用禁用、关联解析模板、选择设备模板 |
| 设备模板库 | 设备模板 CRUD、导入导出、模板应用 |
| 映射文档库 | 字段映射文档 CRUD、批量导入 |
| 日志解析 | 解析模板 CRUD、解析测试、预设模板 |
| 筛选策略 | 条件配置、白名单、去重、动作、优先级、策略测试 |
| 告警聚合 | 聚合规则配置、聚合告警列表、关联日志查看 |
| 高频 IP 识别 | 识别规则配置、高频 IP 列表、明细查看 |
| 字段脱敏 | 脱敏字段配置、脱敏规则管理、敏感配置展示策略 |
| 推送配置 | HTTP、邮箱、Syslog 推送配置和连通性测试 |
| 告警规则 | 筛选策略、输出模板和推送配置编排 |
| 日志查询 | 原始日志、解析日志、筛选状态、告警状态、字段键值查询 |
| 告警记录 | 推送记录、状态、错误信息、重试信息、响应摘要 |
| 告警处置 | 告警确认、忽略、关闭、备注、处置历史 |
| 日志追踪 | 按日志 ID 查看接收、解析、过滤、聚合、推送链路 |
| 数据统计 | 字段 Top N、IP 分布、CSV 导出、复制 IP |
| 审计日志 | 登录、退出、失败登录、配置变更、删除、导入导出、清理、处置记录 |
| 系统配置 | 端口、协议、数据库、保留策略、认证策略、主题、目录配置 |

### 8.1 前端交互要求

1. 页面使用左侧导航 + 顶部状态栏 + 右侧内容区布局。
2. 表格页面必须支持搜索、分页、刷新、编辑、删除。
3. 表单页面必须支持必填校验、格式校验和错误提示。
4. 策略配置应图形化，降低手写 JSON 成本。
5. 解析模板测试应实时展示字段识别结果。
6. HTTP 推送测试应明确显示请求地址、状态码、响应摘要和失败原因。
7. 深色/浅色主题应可切换。
8. 删除、清空、导入覆盖等关键操作必须二次确认。
9. 无权限菜单不显示，无权限按钮置灰或隐藏。
10. 敏感字段在列表和详情中默认脱敏展示。
11. 告警处置操作应提供快捷按钮和备注弹窗。
12. 指标监控页面应展示运行状态、队列积压、错误计数和最近失败记录。

---

## 9. 数据模型需求

### 9.1 核心实体

| 实体 | 用途 |
|---|---|
| User | 用户账号 |
| Role | 角色 |
| Permission | 权限 |
| UserRole | 用户角色关系 |
| RolePermission | 角色权限关系 |
| AuditLog | 审计日志 |
| DeviceGroup | 设备分组 |
| Device | 安全设备/日志源 |
| DeviceTemplate | 设备模板 |
| ParseTemplate | 解析模板 |
| OutputTemplate | 输出模板 |
| FilterPolicy | 筛选策略 |
| PushConfig | 推送配置，统一承载 HTTP、邮箱、Syslog |
| AlertRule | 告警规则 |
| SyslogLog | Syslog 日志 |
| AlertRecord | 告警记录 |
| AggregatedAlert | 聚合告警 |
| AlertDisposition | 告警处置记录 |
| DesensitizeRule | 字段脱敏规则 |
| SystemConfig | 系统配置 |
| FieldMappingDoc | 字段映射文档 |
| LogTraceInfo | 日志追踪信息 |
| MetricSnapshot | 指标快照 |

### 9.2 数据关系

```text
User N ── N Role
Role N ── N Permission
DeviceGroup 1 ── N Device
DeviceTemplate 1 ── N Device
Device 1 ── N SyslogLog
ParseTemplate 1 ── N Device
FilterPolicy 1 ── N AlertRule
OutputTemplate 1 ── N AlertRule
PushConfig 1 ── N AlertRule
SyslogLog 1 ── N AlertRecord
AggregatedAlert 1 ── N SyslogLog
AlertRecord 1 ── N AlertDisposition
User 1 ── N AuditLog
```

### 9.3 User 字段要求

| 字段 | 说明 |
|---|---|
| id | 主键 |
| username | 用户名，唯一 |
| password_hash | 密码哈希 |
| display_name | 显示名称 |
| email | 邮箱 |
| status | enabled/disabled/locked |
| failed_login_count | 连续登录失败次数 |
| locked_until | 锁定截止时间 |
| last_login_at | 最近登录时间 |
| must_change_password | 是否必须修改密码 |
| created_at | 创建时间 |
| updated_at | 更新时间 |

### 9.4 Role 字段要求

| 字段 | 说明 |
|---|---|
| id | 主键 |
| name | 角色名称，唯一 |
| code | 角色编码，唯一 |
| description | 角色说明 |
| built_in | 是否内置角色 |
| created_at | 创建时间 |
| updated_at | 更新时间 |

### 9.5 Permission 字段要求

| 字段 | 说明 |
|---|---|
| id | 主键 |
| name | 权限名称 |
| code | 权限编码，唯一 |
| type | 权限类型：menu/button/api |
| resource | 资源标识 |
| action | 操作标识 |
| created_at | 创建时间 |
| updated_at | 更新时间 |

### 9.6 DeviceGroup 字段要求

| 字段 | 说明 |
|---|---|
| id | 主键 |
| name | 分组名称，唯一 |
| description | 描述 |
| color | 分组颜色 |
| sort_order | 排序 |
| created_at | 创建时间 |
| updated_at | 更新时间 |

### 9.7 Device 字段要求

| 字段 | 说明 |
|---|---|
| id | 主键 |
| name | 设备名称 |
| ip_address | 设备 IP 地址，唯一 |
| group_id | 设备分组 ID |
| template_id | 设备模板 ID |
| parse_template_id | 解析模板 ID |
| device_type | 设备类型 |
| description | 描述 |
| enabled | 是否启用 |
| created_at | 创建时间 |
| updated_at | 更新时间 |

### 9.8 ParseTemplate 字段要求

| 字段 | 说明 |
|---|---|
| id | 主键 |
| name | 解析模板名称 |
| device_type | 设备类型 |
| parse_type | 解析类型 |
| header_regex | Header 正则 |
| delimiter | 分隔符 |
| field_mapping | 字段映射 JSON |
| value_transform | 值转换 JSON |
| sample_log | 示例日志 |
| sub_templates | 子模板配置 JSON |
| enabled | 是否启用 |
| created_at | 创建时间 |
| updated_at | 更新时间 |

### 9.9 FilterPolicy 字段要求

| 字段 | 说明 |
|---|---|
| id | 主键 |
| name | 策略名称 |
| device_id | 关联设备 ID，可为空 |
| device_group_id | 关联设备分组 ID，可为空 |
| parse_template_id | 关联解析模板 ID，可为空 |
| conditions | 筛选条件 JSON |
| condition_logic | AND/OR |
| whitelist_enabled | 是否启用白名单 |
| whitelist_field | 白名单字段 |
| whitelist_values | 白名单值列表 |
| action | keep/discard |
| priority | 优先级 |
| dedup_enabled | 是否启用去重 |
| dedup_window | 去重窗口 |
| enabled | 是否启用 |
| created_at | 创建时间 |
| updated_at | 更新时间 |

### 9.10 OutputTemplate 字段要求

| 字段 | 说明 |
|---|---|
| id | 主键 |
| name | 模板名称 |
| channel_type | http/email/syslog |
| content | 模板内容 |
| fields | 可用字段 JSON |
| device_type | 设备类型 |
| enabled | 是否启用 |
| created_at | 创建时间 |
| updated_at | 更新时间 |

### 9.11 PushConfig 字段要求

| 字段 | 说明 |
|---|---|
| id | 主键 |
| name | 推送配置名称 |
| type | 推送类型：http、email、syslog |
| enabled | 是否启用 |
| url | HTTP 接口地址 |
| method | HTTP 请求方法，默认 POST |
| timeout | HTTP 超时时间 |
| retry_count | 重试次数 |
| retry_delay | 重试间隔 |
| notes_ids | 接收标识列表，JSON 数组 |
| headers | HTTP 请求头，JSON |
| body_template | HTTP 请求体模板 |
| success_status_codes | 成功状态码，JSON 数组 |
| success_body_keyword | 成功响应关键字 |
| auth_type | HTTP 认证方式 |
| token | 认证 Token，加密存储 |
| content_type | HTTP Content-Type |
| retry_on_status_codes | 触发重试的状态码，JSON 数组 |
| max_response_log_size | 响应摘要最大长度 |
| smtp_host | SMTP 地址 |
| smtp_port | SMTP 端口 |
| smtp_username | SMTP 用户名 |
| smtp_password | SMTP 密码，加密存储 |
| syslog_host | Syslog 转发地址 |
| syslog_port | Syslog 转发端口 |
| syslog_protocol | tcp/udp |
| syslog_format | json/rfc3164/rfc5424 |
| created_at | 创建时间 |
| updated_at | 更新时间 |

### 9.12 AlertRule 字段要求

| 字段 | 说明 |
|---|---|
| id | 主键 |
| name | 规则名称 |
| filter_policy_id | 筛选策略 ID |
| push_config_id | 推送配置 ID |
| output_template_id | 输出模板 ID |
| channel_type | http/email/syslog |
| enabled | 是否启用 |
| created_at | 创建时间 |
| updated_at | 更新时间 |

### 9.13 SyslogLog 字段要求

| 字段 | 说明 |
|---|---|
| id | 主键 |
| log_id | 日志唯一 ID |
| source_ip | 来源 IP |
| destination_ip | 目的 IP |
| event_type | 事件类型 |
| severity | Syslog severity |
| facility | Syslog facility |
| device_id | 设备 ID |
| device_name | 设备名称 |
| raw_message | 原始日志 |
| parsed_data | 解析后 JSON |
| filter_status | 过滤状态 |
| matched_filter_policy_id | 命中筛选策略 ID |
| alert_status | 告警状态 |
| alert_rule_id | 告警规则 ID |
| aggregated_alert_id | 聚合告警 ID |
| received_at | 接收时间 |
| created_at | 创建时间 |

### 9.14 AlertRecord 字段要求

| 字段 | 说明 |
|---|---|
| id | 主键 |
| log_id | 日志 ID |
| alert_rule_id | 告警规则 ID |
| push_config_id | 推送配置 ID |
| channel_type | http/email/syslog |
| status | success/failed/retrying |
| retry_count | 已重试次数 |
| request_summary | 请求摘要 |
| response_status_code | HTTP 状态码 |
| response_summary | 响应摘要 |
| error_message | 错误信息 |
| disposition_status | 处置状态 |
| sent_at | 发送时间 |
| created_at | 创建时间 |

### 9.15 AggregatedAlert 字段要求

| 字段 | 说明 |
|---|---|
| id | 主键 |
| aggregate_key | 聚合键 |
| aggregate_type | 聚合类型 |
| source_ip | 来源 IP |
| destination_ip | 目的 IP |
| event_type | 事件类型 |
| device_id | 设备 ID |
| severity | 告警等级 |
| count | 聚合数量 |
| first_seen_at | 首次出现时间 |
| last_seen_at | 最近出现时间 |
| status | 状态 |
| created_at | 创建时间 |
| updated_at | 更新时间 |

### 9.16 AlertDisposition 字段要求

| 字段 | 说明 |
|---|---|
| id | 主键 |
| alert_record_id | 告警记录 ID |
| aggregated_alert_id | 聚合告警 ID，可为空 |
| status | 未处理/处理中/已确认/已忽略/已关闭 |
| note | 处置备注 |
| operator_id | 处置人 ID |
| operator_name | 处置人名称 |
| operated_at | 处置时间 |
| created_at | 创建时间 |

### 9.17 DeviceTemplate 字段要求

| 字段 | 说明 |
|---|---|
| id | 主键 |
| name | 模板名称 |
| device_type | 设备类型 |
| parse_template_id | 解析模板 ID |
| field_mapping_doc_id | 字段映射文档 ID |
| recommended_policy | 推荐策略 JSON |
| enabled | 是否启用 |
| created_at | 创建时间 |
| updated_at | 更新时间 |

### 9.18 AuditLog 字段要求

| 字段 | 说明 |
|---|---|
| id | 主键 |
| user_id | 用户 ID |
| username | 用户名 |
| action | 操作类型 |
| resource_type | 资源类型 |
| resource_id | 资源 ID |
| client_ip | 客户端 IP |
| user_agent | User-Agent |
| result | success/failed |
| detail | 操作详情 |
| created_at | 创建时间 |

---

## 10. API 需求

### 10.1 API 通用规范

#### 10.1.1 统一响应格式

```json
{
  "code": 0,
  "message": "success",
  "data": {},
  "requestId": "202605101000000001"
}
```

#### 10.1.2 分页响应格式

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "items": [],
    "total": 100,
    "page": 1,
    "pageSize": 20
  },
  "requestId": "202605101000000001"
}
```

#### 10.1.3 通用错误码

| 错误码 | 含义 |
|---|---|
| 0 | 成功 |
| 400 | 请求参数错误 |
| 401 | 未登录或会话过期 |
| 403 | 无权限 |
| 404 | 资源不存在 |
| 409 | 数据冲突 |
| 500 | 系统错误 |

#### 10.1.4 通用分页参数

| 参数 | 说明 |
|---|---|
| page | 页码，从 1 开始 |
| pageSize | 每页数量 |
| keyword | 关键词 |
| startTime | 开始时间 |
| endTime | 结束时间 |

### 10.2 认证 API

| 方法 | 路径 | 说明 |
|---|---|---|
| POST | `/api/auth/login` | 用户登录并创建 Session |
| POST | `/api/auth/logout` | 用户退出并销毁 Session |
| GET | `/api/auth/me` | 获取当前用户信息 |
| POST | `/api/auth/change-password` | 修改密码 |
| POST | `/api/auth/init-admin` | 首次初始化管理员密码 |
| POST | `/api/auth/refresh` | 仅在启用 JWT API Token 模式时使用 |

### 10.3 用户与权限 API

| 方法 | 路径 | 说明 |
|---|---|---|
| GET | `/api/users` | 用户列表 |
| POST | `/api/users` | 新增用户 |
| PUT | `/api/users/{id}` | 编辑用户 |
| DELETE | `/api/users/{id}` | 删除用户 |
| POST | `/api/users/{id}/reset-password` | 重置密码 |
| POST | `/api/users/{id}/unlock` | 解除用户锁定 |
| GET | `/api/roles` | 角色列表 |
| POST | `/api/roles` | 新增角色 |
| PUT | `/api/roles/{id}` | 编辑角色 |
| DELETE | `/api/roles/{id}` | 删除角色 |
| GET | `/api/roles/{id}/permissions` | 角色权限详情 |
| POST | `/api/roles/{id}/permissions` | 角色权限分配 |
| GET | `/api/permissions` | 权限列表 |

### 10.4 设备与设备模板 API

| 方法 | 路径 | 说明 |
|---|---|---|
| GET | `/api/devices` | 设备列表 |
| POST | `/api/devices` | 新增设备 |
| PUT | `/api/devices/{id}` | 编辑设备 |
| DELETE | `/api/devices/{id}` | 删除设备 |
| GET | `/api/device-groups` | 设备分组列表 |
| POST | `/api/device-groups` | 新增设备分组 |
| PUT | `/api/device-groups/{id}` | 编辑设备分组 |
| DELETE | `/api/device-groups/{id}` | 删除设备分组 |
| GET | `/api/device-templates` | 设备模板列表 |
| POST | `/api/device-templates` | 新增设备模板 |
| PUT | `/api/device-templates/{id}` | 编辑设备模板 |
| DELETE | `/api/device-templates/{id}` | 删除设备模板 |
| POST | `/api/device-templates/{id}/apply` | 应用设备模板 |
| GET | `/api/field-mappings` | 字段映射文档列表 |
| POST | `/api/field-mappings` | 新增字段映射文档 |
| PUT | `/api/field-mappings/{id}` | 编辑字段映射文档 |
| DELETE | `/api/field-mappings/{id}` | 删除字段映射文档 |

### 10.5 解析模板 API

| 方法 | 路径 | 说明 |
|---|---|---|
| GET | `/api/parse-templates` | 解析模板列表 |
| POST | `/api/parse-templates` | 新增解析模板 |
| PUT | `/api/parse-templates/{id}` | 编辑解析模板 |
| DELETE | `/api/parse-templates/{id}` | 删除解析模板 |
| POST | `/api/parse-templates/test` | 解析模板测试 |

### 10.6 筛选策略 API

| 方法 | 路径 | 说明 |
|---|---|---|
| GET | `/api/filter-policies` | 筛选策略列表 |
| POST | `/api/filter-policies` | 新增筛选策略 |
| PUT | `/api/filter-policies/{id}` | 编辑筛选策略 |
| DELETE | `/api/filter-policies/{id}` | 删除筛选策略 |
| POST | `/api/filter-policies/test` | 筛选策略测试 |

### 10.7 推送配置 API

| 方法 | 路径 | 说明 |
|---|---|---|
| GET | `/api/push-configs` | 推送配置列表 |
| POST | `/api/push-configs` | 新增推送配置 |
| PUT | `/api/push-configs/{id}` | 编辑推送配置 |
| DELETE | `/api/push-configs/{id}` | 删除推送配置 |
| POST | `/api/push-configs/{id}/test` | 推送测试 |

### 10.8 输出模板与告警规则 API

| 方法 | 路径 | 说明 |
|---|---|---|
| GET | `/api/output-templates` | 输出模板列表 |
| POST | `/api/output-templates` | 新增输出模板 |
| PUT | `/api/output-templates/{id}` | 编辑输出模板 |
| DELETE | `/api/output-templates/{id}` | 删除输出模板 |
| GET | `/api/alert-rules` | 告警规则列表 |
| POST | `/api/alert-rules` | 新增告警规则 |
| PUT | `/api/alert-rules/{id}` | 编辑告警规则 |
| DELETE | `/api/alert-rules/{id}` | 删除告警规则 |

### 10.9 日志、告警与处置 API

| 方法 | 路径 | 说明 |
|---|---|---|
| GET | `/api/logs` | 日志查询 |
| DELETE | `/api/logs/cleanup` | 日志清理 |
| GET | `/api/alerts` | 告警记录查询 |
| GET | `/api/logs/{id}/trace` | 日志链路追踪 |
| GET | `/api/stats/fields` | 字段统计 |
| GET | `/api/stats/available-fields` | 可统计字段列表 |
| GET | `/api/aggregated-alerts` | 聚合告警列表 |
| GET | `/api/aggregated-alerts/{id}/logs` | 聚合告警关联日志 |
| POST | `/api/alerts/{id}/dispositions` | 新增告警处置记录 |
| GET | `/api/alerts/{id}/dispositions` | 查看告警处置历史 |

### 10.10 高频 IP、脱敏与指标 API

| 方法 | 路径 | 说明 |
|---|---|---|
| GET | `/api/high-frequency-ips` | 高频 IP 列表 |
| PUT | `/api/high-frequency-ips/config` | 更新高频 IP 识别配置 |
| GET | `/api/desensitize-rules` | 脱敏规则列表 |
| POST | `/api/desensitize-rules` | 新增脱敏规则 |
| PUT | `/api/desensitize-rules/{id}` | 编辑脱敏规则 |
| DELETE | `/api/desensitize-rules/{id}` | 删除脱敏规则 |
| GET | `/api/metrics/runtime` | 运行指标 |
| GET | `/` | Web 管理后台入口页面 |
| GET | `/healthz` | 健康检查 |
| GET | `/readyz` | 就绪检查 |
| GET | `/metrics` | 指标输出 |

### 10.11 系统配置与审计 API

| 方法 | 路径 | 说明 |
|---|---|---|
| GET | `/api/system/config` | 获取系统配置 |
| PUT | `/api/system/config` | 更新系统配置 |
| GET | `/api/system/status` | 获取系统状态 |
| POST | `/api/system/syslog/start` | 启动 Syslog 服务 |
| POST | `/api/system/syslog/stop` | 停止 Syslog 服务 |
| GET | `/api/audit-logs` | 审计日志查询 |

### 10.12 导入导出 API

| 方法 | 路径 | 说明 |
|---|---|---|
| POST | `/api/import/parse-templates` | 导入解析模板 |
| POST | `/api/import/filter-policies` | 导入筛选策略 |
| POST | `/api/import/push-configs` | 导入推送配置 |
| POST | `/api/import/device-templates` | 导入设备模板 |
| GET | `/api/export/parse-templates` | 导出解析模板 |
| GET | `/api/export/filter-policies` | 导出筛选策略 |
| GET | `/api/export/push-configs` | 导出推送配置 |
| GET | `/api/export/device-templates` | 导出设备模板 |

---

## 11. 数据库与存储要求

### 11.1 数据库模式要求

系统默认使用 SQLite，支持通过配置切换为 MySQL 8.0+。数据库类型通过 `config.yaml` 中的 `database.type` 指定：

- `sqlite`：默认值，使用本地数据库文件；
- `mysql`：使用 MySQL 8.0+ 数据库。

系统启动时必须根据数据库类型自动初始化数据表。首次启动时必须自动创建默认管理员账号、内置角色和基础权限数据。

### 11.2 SQLite 要求

| 编号 | 需求 |
|---|---|
| DB-001 | SQLite 为默认数据库 |
| DB-002 | SQLite 数据库文件默认路径为 `data/logcat.db` |
| DB-003 | SQLite 模式必须启用 WAL |
| DB-004 | SQLite 模式必须支持自动建表和基础数据初始化 |
| DB-005 | SQLite 模式必须支持日志定期清理 |

### 11.3 MySQL 8.0+ 要求

| 编号 | 需求 |
|---|---|
| DB-006 | 系统必须支持 MySQL 8.0+ |
| DB-007 | MySQL 字符集必须使用 `utf8mb4` |
| DB-008 | MySQL 必须支持连接池配置 |
| DB-009 | MySQL 必须支持自动建表和基础数据初始化 |
| DB-010 | MySQL 必须支持通过配置文件切换启用 |
| DB-011 | MySQL 时区配置必须支持 `Asia/Shanghai` 或业务指定时区 |

### 11.4 数据库配置示例

```yaml
database:
  type: sqlite
  auto_migrate: true
  sqlite:
    path: data/logcat.db
    wal: true
  mysql:
    host: 127.0.0.1
    port: 3306
    database: logcat
    username: logcat
    password: logcat_password
    charset: utf8mb4
    timezone: Asia/Shanghai
    max_open_conns: 50
    max_idle_conns: 10
```

### 11.5 索引要求

| 表 | 索引字段 |
|---|---|
| users | username、status |
| roles | code、name |
| permissions | code、type |
| devices | ip_address、group_id、enabled |
| device_templates | device_type、enabled |
| parse_templates | name、device_type、enabled |
| filter_policies | enabled、priority、device_id、device_group_id |
| syslog_logs | log_id、source_ip、destination_ip、device_id、received_at、filter_status、alert_status、severity、event_type |
| alert_records | log_id、alert_rule_id、push_config_id、channel_type、status、sent_at |
| aggregated_alerts | aggregate_key、source_ip、event_type、first_seen_at、last_seen_at、status |
| alert_dispositions | alert_record_id、aggregated_alert_id、operator_id、operated_at |
| audit_logs | user_id、username、action、created_at |
| push_configs | type、enabled |

### 11.6 日志保留要求

| 编号 | 需求 |
|---|---|
| DB-012 | 系统应支持配置日志保留天数 |
| DB-013 | 系统应支持手动清理指定天数前日志 |
| DB-014 | 系统应支持清理未匹配日志 |
| DB-015 | 系统应支持清理历史告警记录 |
| DB-016 | 清理操作必须写入审计日志 |
| DB-017 | SQLite 和 MySQL 模式下数据模型字段应保持一致 |
| DB-018 | 数据库迁移脚本应支持重复执行且不破坏已有数据 |

---

## 12. 安全与审计要求

| 编号 | 需求 |
|---|---|
| SEC-001 | 系统必须启用用户认证 |
| SEC-002 | 密码必须哈希存储，禁止明文存储 |
| SEC-003 | 敏感配置必须加密存储或脱敏展示 |
| SEC-004 | API 必须进行认证和权限校验 |
| SEC-005 | 系统必须记录关键操作审计日志 |
| SEC-006 | HTTP 推送外发字段应支持按模板控制，避免敏感数据泄露 |
| SEC-007 | 导入 JSON 配置时必须校验格式和数据合法性 |
| SEC-008 | 审计日志不得被普通用户删除 |
| SEC-009 | 告警处置、配置变更、日志清理、导入导出必须写入审计日志 |
| SEC-010 | 敏感字段包括密码、Token、Secret、认证头、SMTP 密码等 |

---

## 13. 非功能需求

### 13.1 性能

| 编号 | 需求 |
|---|---|
| NFR-001 | 系统应支持持续接收 Syslog 日志 |
| NFR-002 | 系统应避免高频重复告警造成推送风暴 |
| NFR-003 | SQLite 部署应启用 WAL 模式 |
| NFR-004 | 系统应支持日志定期清理 |
| NFR-005 | HTTP 推送失败不应阻塞日志接收线程 |
| NFR-006 | 推送任务应采用异步队列或工作池处理 |
| NFR-007 | 数据统计查询应限制时间范围，避免全表扫描 |
| NFR-008 | 系统应支持配置日志处理 Worker 和推送 Worker 数量 |

### 13.2 可用性

| 编号 | 需求 |
|---|---|
| NFR-009 | 系统应支持 Docker Compose 一键部署 |
| NFR-010 | 系统应支持 Linux 单文件二进制部署 |
| NFR-011 | 系统应支持配置备份和恢复 |
| NFR-012 | 系统应支持健康检查接口 `/healthz` |
| NFR-013 | 系统应支持就绪检查接口 `/readyz` |
| NFR-014 | 服务异常时应输出清晰日志 |
| NFR-015 | 系统应支持优雅停止 |

### 13.3 兼容性

| 编号 | 需求 |
|---|---|
| NFR-016 | Web 管理后台应支持 Chrome、Edge、Firefox 主流现代浏览器 |
| NFR-017 | 后端服务应支持 Linux x86_64 和 ARM64 |
| NFR-018 | 系统应支持 UDP/TCP Syslog 输入和输出 |
| NFR-019 | 系统应支持内网无公网环境部署 |

---

## 14. 部署需求

### 14.1 部署方式

| 部署方式 | 是否支持 | 说明 |
|---|---:|---|
| Linux 单文件部署 | 是 | 默认部署方式，后端 Go 二进制内嵌前端静态资源，启动后直接提供 Web 管理后台和 API 服务 |
| Docker Compose 部署 | 是 | 可选部署方式，用于快速部署 logcat 服务和 MySQL 8.0+ 数据库 |

### 14.2 Linux 单文件部署要求

| 编号 | 需求 |
|---|---|
| DEP-001 | 系统应编译为 Linux 可执行二进制文件 |
| DEP-002 | 前端 Vue 构建产物必须通过 Go Embed 内嵌到后端二进制文件 |
| DEP-003 | 用户上传二进制文件到 Linux 服务器后，应可直接启动服务 |
| DEP-004 | 默认数据库使用 SQLite，数据库文件默认存放在 `data/logcat.db` |
| DEP-005 | 系统应支持通过 `config.yaml` 配置监听地址、Web 端口、Syslog 端口、数据库类型和日志保留策略 |
| DEP-006 | 系统应支持通过环境变量覆盖关键配置 |
| DEP-007 | 系统应提供 systemd 服务配置示例 |
| DEP-008 | 系统启动后应同时提供 Web 管理后台、REST API、Syslog UDP/TCP 接收服务 |

### 14.3 Docker Compose 部署要求

| 编号 | 需求 |
|---|---|
| DEP-009 | 系统应提供 `docker-compose.yml` 示例 |
| DEP-010 | Docker Compose 应支持 logcat 服务和 MySQL 8.0+ 服务 |
| DEP-011 | Docker Compose 应支持数据卷挂载 |
| DEP-012 | Docker Compose 应支持配置文件挂载 |
| DEP-013 | Docker Compose 应支持 Web 端口和 Syslog 端口映射 |
| DEP-014 | Docker Compose 应支持通过环境变量配置数据库连接 |
| DEP-015 | Docker Compose 启动后应能完成数据库初始化 |

### 14.4 目录结构要求

```text
logcat/
├── logcat                 # Linux 可执行文件
├── config.yaml            # 配置文件
├── data/
│   └── logcat.db          # SQLite 数据库
├── logs/
│   └── logcat.log         # 服务运行日志
└── backups/
    └── configs/           # 配置导出备份
```

### 14.5 启动命令

```bash
./logcat --config config.yaml
```

### 14.6 systemd 要求

| 编号 | 需求 |
|---|---|
| DEP-016 | 系统应提供 `/etc/systemd/system/logcat.service` 示例 |
| DEP-017 | systemd 服务应支持开机自启动 |
| DEP-018 | systemd 服务应支持异常退出后自动重启 |

### 14.7 端口规划

| 端口 | 协议 | 用途 |
|---|---|---|
| 8080 | TCP | Web 管理后台/API |
| 5140 | UDP | Syslog UDP 接收 |
| 5140 | TCP | Syslog TCP 接收 |

---

## 15. 验收标准

### 15.1 基础验收

| 编号 | 验收项 |
|---|---|
| AC-001 | Web 版可正常启动并通过浏览器访问 |
| AC-002 | 未登录访问后台自动跳转登录页 |
| AC-003 | 登录成功后可进入系统 |
| AC-004 | 不同角色只能访问授权菜单和功能 |
| AC-005 | 默认监听端口可接收 UDP Syslog |
| AC-006 | TCP Syslog 接收可正常工作 |
| AC-007 | 新增设备后可根据来源 IP 识别设备 |
| AC-008 | 原始日志可落库 |
| AC-009 | 日志解析后可生成结构化字段 |
| AC-010 | 筛选策略可正确命中 |
| AC-011 | 白名单可正确排除指定字段值 |
| AC-012 | 告警去重可在指定窗口内生效 |
| AC-013 | HTTP 接口推送可成功发送请求 |
| AC-014 | HTTP 推送失败后可按配置重试 |
| AC-015 | 邮箱推送可成功发送邮件 |
| AC-016 | Syslog 转发可成功转发日志 |
| AC-017 | 推送失败时可记录错误信息 |
| AC-018 | 日志 ID 追踪可展示完整链路 |
| AC-019 | 数据统计可展示 Top N 结果 |
| AC-020 | 模板和策略可导入导出 |
| AC-021 | 日志清理功能可正常执行 |
| AC-022 | 审计日志可记录登录、退出和关键配置变更 |
| AC-023 | 默认管理员首次登录必须修改初始密码 |
| AC-024 | 无权限用户访问受限 API 返回 403 |

### 15.2 增强功能验收

| 编号 | 验收项 |
|---|---|
| AC-025 | 高频 IP 扫描识别可按时间窗口和阈值生成结果 |
| AC-026 | 告警聚合可按配置字段生成聚合告警 |
| AC-027 | 聚合告警可查看关联原始日志 |
| AC-028 | 字段脱敏规则可对展示和推送内容生效 |
| AC-029 | 告警处置支持确认、忽略、关闭和备注 |
| AC-030 | 告警处置操作可记录处置人、处置时间和审计日志 |
| AC-031 | 设备模板库可新增、编辑、删除、导入、导出和应用 |
| AC-032 | `/healthz`、`/readyz` 和 `/metrics` 接口可正常访问 |
| AC-033 | 仪表盘可展示日志接收量、推送成功率、队列积压和失败计数 |

### 15.3 数据库验收

| 编号 | 验收项 |
|---|---|
| AC-034 | 默认 SQLite 模式可正常启动并自动创建 `data/logcat.db` |
| AC-035 | SQLite 模式启用 WAL |
| AC-036 | 切换为 MySQL 8.0+ 后系统可正常启动 |
| AC-037 | MySQL 模式可自动初始化数据表、默认管理员、角色和权限 |
| AC-038 | SQLite 和 MySQL 模式下核心功能均可正常运行 |
| AC-039 | 数据库迁移重复执行不破坏已有数据 |

### 15.4 部署验收

| 编号 | 验收项 |
|---|---|
| AC-040 | Linux 单文件部署可直接启动 Web、API 和 Syslog 服务 |
| AC-041 | 前端静态资源已内嵌到 Go 二进制文件 |
| AC-042 | `config.yaml` 可正确控制 Web 端口、Syslog 端口和数据库类型 |
| AC-043 | 环境变量可覆盖关键配置 |
| AC-044 | systemd 服务可启动、停止、重启和开机自启 |
| AC-045 | Docker Compose 可启动 logcat 和 MySQL 8.0+ |
| AC-046 | Docker Compose 数据卷、配置挂载和端口映射可正常工作 |

### 15.5 性能验收

| 编号 | 验收项 |
|---|---|
| AC-047 | SQLite 模式下单节点每秒接收 100 条 Syslog 日志时，界面和服务保持可用 |
| AC-048 | MySQL 模式下单节点每秒接收 300 条 Syslog 日志时，界面和服务保持可用 |
| AC-049 | 连续运行 24 小时无服务崩溃 |
| AC-050 | 10 万条以内本地日志普通分页查询 3 秒内返回 |
| AC-051 | 10 万条以内指定时间范围统计 5 秒内返回 |
| AC-052 | 高频重复告警不应持续刷屏 |
| AC-053 | HTTP 推送目标异常时不影响日志接收 |
| AC-054 | 推送队列积压时仪表盘应展示积压数量 |

---

## 16. 整体开发、联调与交付要求

### 16.1 开发原则

logcat 所有功能均纳入当前版本整体开发范围，不按优先级拆分，不按阶段交付。开发过程可按模块组织实施，最终必须整体联调、整体验收、整体交付。

系统开发必须覆盖本文档中定义的全部功能需求、页面需求、数据模型需求、API 需求、数据库与存储要求、安全与审计要求、非功能需求、部署需求和验收标准。任何功能不得以“后续优化”“二期建设”“暂不开发”的方式从当前版本范围中移除。

### 16.2 整体功能范围

本版本必须整体完成以下功能：

1. Linux Web 单文件部署。
2. Docker Compose 部署。
3. 默认 SQLite 数据库支持。
4. MySQL 8.0+ 数据库支持。
5. 用户登录、退出、密码修改、首次登录强制改密。
6. RBAC 权限控制。
7. 审计日志。
8. Syslog UDP/TCP 接收。
9. 日志接收队列、解析 Worker、筛选 Worker、数据库写入队列和推送 Worker。
10. 设备管理。
11. 设备分组。
12. 设备模板库。
13. 字段映射文档库。
14. 解析模板管理。
15. 解析模板测试。
16. 筛选策略管理。
17. 筛选策略测试。
18. 白名单策略。
19. 告警去重。
20. 告警聚合。
21. 高频 IP 扫描识别。
22. 字段脱敏。
23. 输出模板管理。
24. HTTP 接口推送。
25. 邮箱推送。
26. Syslog 转发。
27. 推送配置连通性测试。
28. 告警规则管理。
29. 日志查询。
30. 告警记录。
31. 告警处置闭环。
32. 日志 ID 全链路追踪。
33. 数据统计。
34. 配置导入导出。
35. 系统配置。
36. 队列积压监控。
37. 系统健康检查和指标监控。

### 16.3 联调要求

所有模块必须在同一个版本中完成集成联调，联调范围包括：

1. 前端页面与后端 REST API 联调。
2. 用户认证、RBAC 权限与所有 API 的权限校验联调。
3. Syslog UDP/TCP 接收与设备识别联调。
4. 日志解析、字段映射、筛选策略、白名单、去重、聚合、告警规则的完整链路联调。
5. HTTP 接口推送、邮箱推送、Syslog 转发与告警记录联调。
6. 告警处置、日志查询、日志追踪、数据统计、审计日志联调。
7. SQLite 默认数据库模式下的完整功能联调。
8. MySQL 8.0+ 数据库模式下的完整功能联调。
9. Linux 单文件部署模式下的完整功能联调。
10. Docker Compose 部署模式下的完整功能联调。

### 16.4 验收要求

验收必须覆盖本文档第 15 章定义的全部验收标准。所有验收项必须一次性纳入整体验收范围，不再按优先级拆分验收。

系统只有在以下条件全部满足后，才视为整体完成：

1. 所有功能模块开发完成。
2. 所有页面可正常访问和操作。
3. 所有 REST API 完成认证和权限校验。
4. SQLite 和 MySQL 两种数据库模式均通过验证。
5. Linux 单文件部署和 Docker Compose 部署均通过验证。
6. Syslog 接收、解析、筛选、去重、聚合、推送、记录、追踪、统计、处置形成完整闭环。
7. HTTP 推送异常、邮箱异常、Syslog 转发异常和数据库慢写入不影响 Syslog 接收主流程。
8. 审计日志能够记录登录、退出、失败登录、配置变更、删除、导入导出、清理、处置等关键操作。
9. 性能验收、部署验收、数据库验收、增强功能验收和基础功能验收全部通过。

---

## 17. 文档结论

logcat 是 Web 版安全日志接收、解析、筛选、统计与告警转发平台。系统默认使用 SQLite 数据库，支持 MySQL 8.0+；支持 Linux 服务器单文件部署；内置用户认证、RBAC 权限控制、审计日志、Syslog 接收、日志解析、筛选策略、HTTP 接口推送、邮箱推送、Syslog 转发、日志查询、日志追踪、告警聚合、高频 IP 识别、字段脱敏、告警处置和数据统计能力。

本版本 PRD 固定技术栈为 **Go 1.22+ + Gin + GORM + Vue 3 + TypeScript + Naive UI**。系统采用前后端分离开发、单文件同源部署方式，日志接收、解析、筛选、数据库写入和推送任务采用异步队列与 Worker 池处理，确保 HTTP 推送异常、邮件发送异常或数据库慢写入不阻塞 Syslog 日志接收主流程。
