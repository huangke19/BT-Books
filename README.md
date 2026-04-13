# BT-Books 📚

电子书搜索下载工具，通过 [Z-Library](https://z-library.sk) 搜索，支持 EPUB / PDF / MOBI 等格式。

从 [BT-Spider](https://github.com/huangke19/BT-Spider) 拆分而来的独立轻量版本。

## 功能

- 通过 zlib CLI 搜索电子书，每次返回最多 15 个结果
- 展示书名、作者、格式、大小、年份
- 输入序号直接下载，自动保存到本地
- 无额外依赖，纯 Go 标准库

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

### 运行

```bash
# 需要代理（用于访问 Z-Library）
HTTPS_PROXY=http://127.0.0.1:7890 timeout 60 ./bt-books
```

### 命令

| 命令 | 说明 |
|------|------|
| `search <关键词>` / `book <关键词>` | 搜索电子书 |
| `<序号>` | 下载对应结果 |
| `quit` / `q` | 退出 |

### 示例

```
book> search Feynman Lectures on Physics
📚 搜索电子书: Feynman Lectures on Physics

找到 15 个结果:

  [1] The Feynman Lectures on Physics: Quantum
      Richard P. Feynman | PDF | 8.27 MB | 2013
  [2] The Feynman Lectures on Physics, Vol. II
      Richard P. Feynman | EPUB | 28.40 MB | 2015
  ...

输入序号下载（回车跳过）:
1
⬇️  下载电子书: The Feynman Lectures on Physics: Quantum Mechanics
✓ Saved to: ~/Documents/Books/...pdf
```

## 下载目录

默认：`~/Documents/Books/`

## 代理

Z-Library 在国内需要代理：

```bash
HTTPS_PROXY=http://127.0.0.1:7890 ./bt-books
```

## 项目结构

```
.
├── main.go          # CLI 入口，交互式 REPL
├── book/
│   ├── book.go      # BookResult 结构体、BookProvider 接口、SearchBooks 聚合
│   └── zlib.go      # zlib CLI 搜索 / 下载封装，表格解析
├── go.mod
└── README.md
```

## 相关项目

- [BT-Spider](https://github.com/huangke19/BT-Spider) — 完整 BT 下载工具，含引擎 + Telegram Bot
- [BT-Music](https://github.com/huangke19/BT-Music) — 同系列音乐下载工具（B站 yt-dlp + BT搜索）

## 许可证

MIT
