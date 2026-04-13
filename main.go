package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"

	"github.com/huangke/bt-books/book"
)

func defaultDownloadDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, "Documents", "Books")
}

func main() {
	fmt.Println("📚 BT-Books v0.1.0")
	fmt.Println("输入 search <关键词> 搜索电子书，book <关键词> 也可以，quit 退出")
	fmt.Println()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		fmt.Println("\n👋 正在退出...")
		os.Exit(0)
	}()

	scanner := bufio.NewScanner(os.Stdin)
	var lastBooks []book.BookResult
	downloadDir := defaultDownloadDir()

	fmt.Printf("💾 下载目录: %s\n\n", downloadDir)

	for {
		fmt.Print("book> ")
		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}

		switch {
		case strings.ToLower(input) == "quit" || strings.ToLower(input) == "exit" || strings.ToLower(input) == "q":
			fmt.Println("👋 再见!")
			return

		case strings.HasPrefix(strings.ToLower(input), "search ") || strings.HasPrefix(strings.ToLower(input), "book "):
			keyword := strings.TrimSpace(input[strings.IndexByte(input, ' ')+1:])
			if keyword == "" {
				fmt.Println("⚠️  请输入书名或作者")
				continue
			}
			fmt.Printf("📚 搜索电子书: %s\n", keyword)
			books, err := book.SearchBooks(keyword, book.DefaultProviders())
			if err != nil {
				fmt.Printf("❌ 搜索失败: %v\n", err)
				continue
			}
			if len(books) == 0 {
				fmt.Println("未找到相关电子书")
				continue
			}
			lastBooks = books
			fmt.Printf("\n找到 %d 个结果:\n\n", len(books))
			for i, b := range books {
				author := b.Author
				if author == "" {
					author = "未知作者"
				}
				fmt.Printf("  [%d] %s\n      %s | %s | %s | %s\n",
					i+1, b.Title, author, b.Format, b.Size, b.Year)
			}
			fmt.Println("\n输入序号下载（回车跳过）: ")

		default:
			if num, err := strconv.Atoi(input); err == nil {
				if num >= 1 && num <= len(lastBooks) {
					b := lastBooks[num-1]
					if b.ID == "" {
						fmt.Println("⚠️  这个结果没有可下载的 ID")
						continue
					}
					fmt.Printf("⬇️  下载电子书: %s\n", b.Title)
					if err := book.ZlibDownload(b.ID, downloadDir); err != nil {
						fmt.Printf("❌ 下载失败: %v\n", err)
					}
					fmt.Println()
				} else {
					fmt.Println("⚠️  序号超出范围")
				}
			} else {
				fmt.Println("⚠️  未知命令。输入 search <关键词> 搜索电子书，或输入序号下载，或 quit 退出")
			}
		}
	}
}
