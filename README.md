# BT-Books 📚

电子书搜索下载工具，通过 [Z-Library](https://z-library.sk) 搜索，支持 EPUB / PDF / MOBI 等格式。

从 [BT-Spider](https://github.com/huangke19/BT-Spider) 拆分而来的独立轻量版本。

## 功能

- TUI 终端界面，交互方式与 BT-Spider 类似
- 通过 zlib CLI 搜索电子书，每次返回最多 15 个结果
- 固定列宽结果表格，适配中英文混排
- 展示书名、作者、格式、大小、年份
- 方向键选择结果，回车直接下载
- 提供 `run.sh` 一键编译并启动

## 前置依赖

需要先安装并登录 [heartleo/zlib](https://github.com/heartleo/zlib) CLI：

```bash
# 安装（需要 Go 环境）
GOPATH=/tmp/gopath go install github.com/heartleo/zlib/cmd/zlib@latest
cp /tmp/gopath/bin/zlib ~/bin/zlib

# 登录（需要 Z-Library 账号 + 代理）
HTTPS_PROXY=http://127.0.0.1:7890 ~/bin/zlib login --email your@email.com --password yourpass
```

## 快速开始

### 编译

```bash
go build -o bt-books .
```

### 一键运行

```bash
./run.sh
```

如需代理：

```bash
HTTPS_PROXY=http://127.0.0.1:7890 ./run.sh
```

### 手动运行

```bash
HTTPS_PROXY=http://127.0.0.1:7890 ./bt-books
```

### TUI 操作

| 操作 | 说明 |
|------|------|
| 输入关键词后回车 | 搜索电子书 |
| `search <关键词>` / `book <关键词>` 后回车 | 兼容旧命令格式 |
| `↑` / `↓` 或 `j` / `k` | 选择结果 |
| 空输入时按 `Enter` | 下载当前选中项 |
| `esc` | 清空结果 |
| `q` | 退出 |

### 示例

```text
book> Feynman Lectures on Physics

[01] The Feynman Lectures on Physics: Quantum...  Richard P. Feynman  PDF   8.27 MB  2013
[02] The Feynman Lectures on Physics, Vol. II...  Richard P. Feynman  EPUB  28.40 MB  2015
...
```

## 下载目录

默认：`~/Documents/Books/`

## 代理

Z-Library 在国内需要代理：

```bash
HTTPS_PROXY=http://127.0.0.1:7890 ./run.sh
```

## 项目结构

```text
.
├── main.go          # TUI 入口
├── book/
│   ├── book.go      # BookResult 结构体、BookProvider 接口、SearchBooks 聚合
│   └── zlib.go      # zlib CLI 搜索 / 下载封装
├── tui/
│   └── tui.go       # Bubble Tea TUI、结果表格、键盘交互
├── run.sh           # 一键编译并启动
├── go.mod
├── go.sum
└── README.md
```

## 相关项目

- [BT-Spider](https://github.com/huangke19/BT-Spider) — 完整 BT 下载工具，含引擎 + Telegram Bot
- [BT-Music](https://github.com/huangke19/BT-Music) — 同系列音乐下载工具（B站 yt-dlp + BT搜索）

## 许可证

MIT
