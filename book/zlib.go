package book

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

// ZlibBin returns the path to the zlib executable (prefers ~/bin/zlib, falls back to PATH).
func ZlibBin() string {
	home, _ := os.UserHomeDir()
	local := filepath.Join(home, "bin", "zlib")
	if _, err := os.Stat(local); err == nil {
		return local
	}
	if p, err := exec.LookPath("zlib"); err == nil {
		return p
	}
	return "zlib"
}

// zlibSession stores the zlib CLI session (cookies + domain).
type zlibSession struct {
	Cookies map[string]string `json:"cookies"`
	Domain  string            `json:"domain"`
}

// ZlibProvider is a BookProvider backed by the zlib CLI.
type ZlibProvider struct{}

func NewZlibProvider() *ZlibProvider { return &ZlibProvider{} }

func (z *ZlibProvider) Name() string { return "zlib CLI" }

func (z *ZlibProvider) SearchBooks(keyword string) ([]BookResult, error) {
	return ZlibSearch(keyword)
}

// LoadZlibSession reads the zlib CLI saved session.
func LoadZlibSession() (*zlibSession, error) {
	home, _ := os.UserHomeDir()
	data, err := os.ReadFile(filepath.Join(home, ".config", "zlib", "session.json"))
	if err != nil {
		return nil, fmt.Errorf("未找到 zlib session，请先运行: ~/bin/zlib login")
	}
	var s zlibSession
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, err
	}
	return &s, nil
}

// ZlibSearch searches for ebooks via the zlib CLI and parses the table output.
func ZlibSearch(keyword string) ([]BookResult, error) {
	if _, err := LoadZlibSession(); err != nil {
		return nil, err
	}
	cmd := exec.Command(ZlibBin(), "search", keyword, "-n", "15")
	cmd.Env = os.Environ()
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("zlib search 失败: %w", err)
	}
	return ParseZlibTable(string(out)), nil
}

// ParseZlibTable parses the box-drawing table output from zlib search.
// Format: │ # │ ID │ Title │ Authors │ Year │ Format │ Size │
func ParseZlibTable(output string) []BookResult {
	var books []BookResult
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "│") {
			continue
		}
		line = strings.Trim(line, "│")
		cols := strings.Split(line, "│")
		if len(cols) < 7 {
			continue
		}
		num := strings.TrimSpace(cols[0])
		if num == "#" || num == "" {
			continue
		}
		if _, err := strconv.Atoi(num); err != nil {
			continue
		}
		id := strings.TrimSpace(cols[1])
		title := strings.TrimSpace(cols[2])
		author := strings.TrimSpace(cols[3])
		year := strings.TrimSpace(cols[4])
		format := strings.TrimSpace(cols[5])
		size := strings.TrimSpace(cols[6])
		if id == "" || title == "" {
			continue
		}
		books = append(books, BookResult{
			ID:     id,
			Title:  title,
			Author: author,
			Year:   year,
			Format: format,
			Size:   size,
			Source: "zlib CLI",
		})
	}
	return books
}

// ZlibDownload calls zlib download <id> via script to provide a pseudo TTY.
func ZlibDownload(id, destDir string) error {
	if destDir == "" {
		home, _ := os.UserHomeDir()
		destDir = filepath.Join(home, "Documents", "Books")
	}
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return err
	}
	cmd := exec.Command("script", "-q", "/dev/null",
		ZlibBin(), "download", id, "-d", destDir)
	cmd.Env = os.Environ()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
