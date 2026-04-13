package book

// BookResult 电子书搜索结果
// ID 由 zlib CLI 搜索结果提供，用于后续下载。
type BookResult struct {
	ID        string // zlib 条目 ID
	Title     string // 书名
	Author    string // 作者
	Format    string // 格式，如 epub/pdf/mobi
	Size      string // 文件大小
	Source    string // 来源站点
	DirectURL string // 下载页面 URL（可选）
	Language  string // 语言
	Year      string // 出版年份
}

// BookProvider 电子书搜索源接口
type BookProvider interface {
	Name() string
	SearchBooks(keyword string) ([]BookResult, error)
}

// SearchBooks 使用所有电子书源搜索，合并结果
func SearchBooks(keyword string, providers []BookProvider) ([]BookResult, error) {
	var all []BookResult
	for _, p := range providers {
		results, err := p.SearchBooks(keyword)
		if err != nil {
			continue
		}
		all = append(all, results...)
	}
	return all, nil
}

// DefaultProviders 返回默认的电子书搜索源。
// 目前以 zlib CLI 为主，后续可继续扩展网页源。
func DefaultProviders() []BookProvider {
	return []BookProvider{NewZlibProvider()}
}
