package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/huangke/bt-books/tui"
)

func defaultDownloadDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, "Documents", "Books")
}

func main() {
	downloadDir := defaultDownloadDir()
	if err := os.MkdirAll(downloadDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "❌ 无法创建下载目录: %v\n", err)
		os.Exit(1)
	}

	p := tea.NewProgram(tui.New(downloadDir), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "❌ TUI 运行出错: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("👋 已退出，下载目录: %s (%s)\n", downloadDir, time.Now().Format("15:04:05"))
}
