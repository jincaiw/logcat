# 标准发布流程

适用于 logcat 的统一发版流程。

## 版本源

- 根目录 `VERSION` 是发布版本的主输入。
- `scripts/bump-version.sh` 会同步所有版本引用。
- `scripts/release.sh` 会完成校验、构建、打包和本地 tag 创建。

## 推荐步骤

### 1. 准备版本

```bash
bash scripts/bump-version.sh <version>
```

### 2. 准备发布说明

如果没有对应的发布说明文件，先创建：

```bash
cp docs/release-notes-template.md docs/release-notes-v<version>.md
```

然后补充 Highlights / Assets / Docker 内容。

### 3. 执行标准发布

```bash
bash scripts/release.sh <version>
```

脚本会：

- 同步版本引用
- 运行 `go test ./...`
- 构建 Linux amd64/arm64 发布包
- 创建本地 tag `v<version>`

### 4. 发布到 GitHub

```bash
git push origin <branch>
git push origin v<version>
```

推送 tag 后，GitHub Actions 会自动：

- 运行校验与构建
- 生成 `amd64/arm64` 发布包
- 使用 `docs/release-notes-v<version>.md` 创建 GitHub Release

## 约定

- 所有新版本统一走上述脚本，不再手工散改版本号。
- `build-web.sh` 与 `scripts/install-linux.sh` 默认读取根目录 `VERSION`。
