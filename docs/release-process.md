# 标准发布流程

适用于 logcat 的统一发版流程。

## 最近一次发布记录

### v0.2.8
- 统一 `VERSION` 为单一版本源。
- 标准化 `scripts/bump-version.sh` / `scripts/release.sh`。
- 整理并固化 GitHub Release、Pages、DockerHub 的发布路径。
- Pages 首页、安装页、用户手册恢复为简洁版式，并补齐截图资源。
- `main` / `master` 保持同步。

### v0.2.7
- 清理 `公众号.md`。
- 去除品牌化措辞。
- 保持发布流程统一。

### v0.2.6
- 前端 UI 重构完成。
- 发布流程开始标准化。

## 发布原则

- `VERSION` 是唯一版本源。
- 所有版本引用必须通过脚本同步，不手工散改。
- 每个版本都要有对应的 `docs/release-notes-v<version>.md`。
- 先本地校验，再打 tag，再推送到远端。
- 如果 Pages 仍使用 `master` 作为源分支，`main` 与 `master` 必须保持一致。

## 一次完整发布

> 下面命令按顺序执行即可。

```bash
VERSION=0.2.9

# 1. 同步版本号
bash scripts/bump-version.sh "$VERSION"

# 2. 准备 release notes
cp docs/release-notes-template.md "docs/release-notes-v${VERSION}.md"
# 编辑 docs/release-notes-v${VERSION}.md

# 3. 本地校验 + 构建 + 打 tag
bash scripts/release.sh "$VERSION"

# 4. 推送代码和 tag
# 如 Pages 仍以 master 为源，请两条都推
git push origin main
git push origin master
git push origin "v${VERSION}"

# 5. DockerHub 发布
# 如需更新镜像标签，执行：
# docker build -t qing1205/logcat:${VERSION} -t qing1205/logcat:latest .
# docker push qing1205/logcat:${VERSION}
# docker push qing1205/logcat:latest
```

## `scripts/release.sh` 做什么

- 同步版本引用
- 自动生成缺失的 release notes 骨架
- 执行 `go test ./...`
- 构建 Linux `amd64` / `arm64` 发布包
- 创建本地 tag `v<version>`

## 发布后检查

- GitHub Release 已生成，且正文来自 `docs/release-notes-v<version>.md`
- `logcat.mujizi.com` 已刷新
- DockerHub `qing1205/logcat:<version>` 和 `latest` 已更新
- 下载链接、安装脚本、截图链接都可访问

## Pages 注意事项

当前站点页面以仓库根目录静态文件为准：

- `index.html`
- `installation.html`
- `user-guide.html`
- `assets/`

如果首页 / 安装页 / 用户手册内容变更，需要同时维护 `docs/` 下的对应文件，避免文档与站点不同步。

## 约定

- 所有新版本统一走上述脚本，不再手工散改版本号。
- `build-web.sh` 与 `scripts/install-linux.sh` 默认读取根目录 `VERSION`。
- 版本号、tag、Release notes、Docker tag、Pages 内容都必须互相对应。
