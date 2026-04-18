package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-runewidth"

	"github.com/huangke/bt-books/book"
)

const version = "0.2.0"

type searchDoneMsg struct {
	keyword string
	results []book.BookResult
	err     error
}

type downloadDoneMsg struct {
	title  string
	output string
	err    error
}

type statusMsg struct {
	text  string
	isErr bool
}

var (
	titleStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("230"))
	chromeStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62")).
			Padding(0, 1)
	panelStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240")).
			Padding(0, 1)
	panelTitleStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("117"))
	accentStyle     = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("81"))
	dimStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
	okStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
	errStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
	selectedStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("229"))
	markerStyle     = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("86"))
	metaLabelStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
	valueStyle      = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("229"))
	statusBoxStyle  = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("238")).
			Padding(0, 1)
)

type Model struct {
	input       textinput.Model
	results     []book.BookResult
	selected    int
	status      string
	isErr       bool
	width       int
	height      int
	searching   bool
	downloading bool
	downloadDir string
	lastKeyword string
}

func New(downloadDir string) Model {
	ti := textinput.New()
	ti.Placeholder = "输入书名/作者，回车搜索"
	ti.Prompt = "book> "
	ti.CharLimit = 512
	ti.Focus()

	return Model{
		input:       ti,
		status:      "输入关键词回车搜索，↑/↓ 选择结果，Enter 下载，q 退出",
		downloadDir: downloadDir,
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		if msg.Width > 8 {
			m.input.Width = msg.Width - 8
		}
		return m, nil

	case searchDoneMsg:
		m.searching = false
		if msg.err != nil {
			m.results = nil
			m.selected = 0
			m.status = "搜索失败: " + msg.err.Error()
			m.isErr = true
			return m, nil
		}
		m.results = msg.results
		m.selected = 0
		m.lastKeyword = msg.keyword
		if len(msg.results) == 0 {
			m.status = "未找到相关电子书"
			m.isErr = true
			return m, nil
		}
		m.status = fmt.Sprintf("找到 %d 个结果，按 Enter 下载选中项", len(msg.results))
		m.isErr = false
		return m, nil

	case downloadDoneMsg:
		m.downloading = false
		if msg.err != nil {
			m.status = "下载失败: " + compactOutput(msg.output, msg.err.Error())
			m.isErr = true
			return m, nil
		}
		m.status = fmt.Sprintf("下载完成: %s", msg.title)
		if s := compactOutput(msg.output, ""); s != "" {
			m.status += " | " + s
		}
		m.isErr = false
		return m, nil

	case statusMsg:
		m.status = msg.text
		m.isErr = msg.isErr
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			if len(m.results) > 0 && m.selected > 0 {
				m.selected--
			}
			return m, nil

		case "down", "j":
			if len(m.results) > 0 && m.selected < len(m.results)-1 {
				m.selected++
			}
			return m, nil

		case "esc":
			m.results = nil
			m.selected = 0
			m.status = "结果已清空"
			m.isErr = false
			return m, nil

		case "enter":
			if m.searching || m.downloading {
				return m, nil
			}
			query := normalizeQuery(m.input.Value())
			if query != "" {
				m.searching = true
				m.status = "搜索中: " + query
				m.isErr = false
				m.input.SetValue("")
				return m, searchCmd(query)
			}
			if len(m.results) == 0 {
				return m, nil
			}
			selected := m.results[m.selected]
			if selected.ID == "" {
				return m, statusCmd("该结果没有可下载 ID", true)
			}
			m.downloading = true
			m.status = "下载中: " + selected.Title
			m.isErr = false
			return m, downloadCmd(selected, m.downloadDir)
		}
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	header := chromeStyle.Width(m.contentWidth()).Render(strings.Join([]string{
		lipgloss.JoinHorizontal(lipgloss.Top, titleStyle.Render("BT-Books"), dimStyle.Render("  v"+version)),
		dimStyle.Render("Z-Library 电子书搜索下载"),
		metaLabelStyle.Render("下载目录  ") + dimStyle.Render(m.downloadDir),
	}, "\n"))

	inputPanel := panelStyle.Width(m.contentWidth()).Render(
		panelTitleStyle.Render("搜索") + "\n" + m.input.View(),
	)

	resultsPanel := panelStyle.Width(m.contentWidth()).Render(
		m.renderResults(),
	)

	statusPanel := m.renderStatusBar()

	help := dimStyle.Render("回车搜索或下载 | ↑/↓ 或 j/k 移动 | esc 清空结果 | q 退出")

	return lipgloss.JoinVertical(lipgloss.Left, header, inputPanel, resultsPanel, statusPanel, help)
}

func (m Model) renderResults() string {
	var body strings.Builder
	body.WriteString(panelTitleStyle.Render("搜索结果"))
	body.WriteString("\n")
	body.WriteString(dimStyle.Render(m.resultsSummary()))
	body.WriteString("\n\n")

	if m.searching {
		body.WriteString(dimStyle.Render("正在搜索，请稍候..."))
		return body.String()
	}
	if m.downloading {
		body.WriteString(dimStyle.Render("正在下载选中书籍，请稍候..."))
		return body.String()
	}
	if len(m.results) == 0 {
		body.WriteString(dimStyle.Render("暂无结果。输入书名、作者或 `search <关键词>` 后回车开始搜索。"))
		return body.String()
	}

	contentWidth := m.contentWidth() - panelStyle.GetHorizontalFrameSize()
	if contentWidth < 48 {
		contentWidth = 48
	}
	rowWidth := contentWidth - 2
	if rowWidth < 44 {
		rowWidth = 44
	}
	limit := visibleResultLimit(m.height)
	if len(m.results) < limit {
		limit = len(m.results)
	}
	start := resultWindowStart(m.selected, len(m.results), limit)
	end := start + limit
	if end > len(m.results) {
		end = len(m.results)
	}

	indexWidth, formatWidth, sizeWidth, yearWidth := resultColumnWidths(m.results)
	authorWidth := 14
	fixedWidth := indexWidth + authorWidth + formatWidth + sizeWidth + yearWidth + 10
	titleWidth := rowWidth - fixedWidth
	if titleWidth < 16 {
		titleWidth = 16
	}

	body.WriteString(dimStyle.Render("  "))
	body.WriteString(dimStyle.Render(padLeftDisplay("#", indexWidth)))
	body.WriteString(dimStyle.Render(" "))
	body.WriteString(dimStyle.Render(padDisplay("标题", titleWidth)))
	body.WriteString(dimStyle.Render("  "))
	body.WriteString(dimStyle.Render(padDisplay("作者", authorWidth)))
	body.WriteString(dimStyle.Render("  "))
	body.WriteString(dimStyle.Render(padDisplay("格式", formatWidth)))
	body.WriteString(dimStyle.Render("  "))
	body.WriteString(dimStyle.Render(padLeftDisplay("大小", sizeWidth)))
	body.WriteString(dimStyle.Render("  "))
	body.WriteString(dimStyle.Render(padLeftDisplay("年份", yearWidth)))
	body.WriteString("\n")
	body.WriteString(dimStyle.Render("  " + strings.Repeat("─", maxInt(24, rowWidth))))
	body.WriteString("\n")

	for i := start; i < end; i++ {
		result := m.results[i]
		author := result.Author
		if author == "" {
			author = "未知作者"
		}
		year := result.Year
		if year == "" {
			year = "-"
		}
		format := result.Format
		if format == "" {
			format = "-"
		}
		size := result.Size
		if size == "" {
			size = "-"
		}

		indexText := padLeftDisplay(fmt.Sprintf("[%02d]", i+1), indexWidth)
		titleText := padDisplay(truncateDisplay(result.Title, titleWidth), titleWidth)
		authorText := padDisplay(truncateDisplay(author, authorWidth), authorWidth)
		formatText := padDisplay(strings.ToUpper(truncateDisplay(format, formatWidth)), formatWidth)
		sizeText := padLeftDisplay(size, sizeWidth)
		yearText := padLeftDisplay(year, yearWidth)

		if i == m.selected {
			body.WriteString(markerStyle.Render("› "))
			body.WriteString(selectedStyle.Render(indexText))
			body.WriteString(" ")
			body.WriteString(selectedStyle.Render(titleText))
			body.WriteString("  ")
			body.WriteString(selectedStyle.Render(authorText))
			body.WriteString("  ")
			body.WriteString(selectedStyle.Render(formatText))
			body.WriteString("  ")
			body.WriteString(selectedStyle.Render(sizeText))
			body.WriteString("  ")
			body.WriteString(selectedStyle.Render(yearText))
		} else {
			body.WriteString("  ")
			body.WriteString(dimStyle.Render(indexText))
			body.WriteString(" ")
			body.WriteString(titleText)
			body.WriteString("  ")
			body.WriteString(dimStyle.Render(authorText))
			body.WriteString("  ")
			body.WriteString(dimStyle.Render(formatText))
			body.WriteString("  ")
			body.WriteString(dimStyle.Render(sizeText))
			body.WriteString("  ")
			body.WriteString(dimStyle.Render(yearText))
		}
		body.WriteString("\n")
	}
	if len(m.results) > limit {
		body.WriteString("\n")
		body.WriteString(dimStyle.Render(fmt.Sprintf("显示 %d-%d / %d", start+1, end, len(m.results))))
	}
	return strings.TrimRight(body.String(), "\n")
}

func (m Model) renderStatusBar() string {
	line := m.status
	if line == "" {
		line = "准备就绪"
	}
	if m.isErr {
		return statusBoxStyle.Width(m.contentWidth()).Render(errStyle.Render("✖ " + line))
	}
	return statusBoxStyle.Width(m.contentWidth()).Render(okStyle.Render("• " + line))
}

func (m Model) resultsSummary() string {
	if len(m.results) == 0 {
		if m.lastKeyword == "" {
			return "等待输入"
		}
		return fmt.Sprintf("关键词“%s”暂无结果", m.lastKeyword)
	}
	return fmt.Sprintf("关键词“%s”共 %d 条结果", m.lastKeyword, len(m.results))
}

func (m Model) contentWidth() int {
	if m.width <= 0 {
		return 88
	}
	if m.width < 40 {
		return m.width - 2
	}
	return m.width - 4
}

func visibleResultLimit(height int) int {
	switch {
	case height >= 42:
		return 18
	case height >= 34:
		return 14
	default:
		return 10
	}
}

func resultWindowStart(selected, total, limit int) int {
	if total <= limit {
		return 0
	}
	start := selected - limit/2
	if start < 0 {
		start = 0
	}
	maxStart := total - limit
	if start > maxStart {
		start = maxStart
	}
	return start
}

func resultColumnWidths(results []book.BookResult) (indexWidth, formatWidth, sizeWidth, yearWidth int) {
	indexWidth = runewidth.StringWidth(fmt.Sprintf("[%02d]", len(results)))
	if indexWidth < runewidth.StringWidth("[00]") {
		indexWidth = runewidth.StringWidth("[00]")
	}
	formatWidth = runewidth.StringWidth("FORMAT")
	sizeWidth = runewidth.StringWidth("00.00 MB")
	yearWidth = runewidth.StringWidth("0000")

	for _, result := range results {
		formatWidth = maxInt(formatWidth, runewidth.StringWidth(strings.ToUpper(result.Format)))
		sizeWidth = maxInt(sizeWidth, runewidth.StringWidth(result.Size))
		yearWidth = maxInt(yearWidth, runewidth.StringWidth(result.Year))
	}
	if formatWidth > 6 {
		formatWidth = 6
	}
	if sizeWidth > 9 {
		sizeWidth = 9
	}
	if yearWidth > 4 {
		yearWidth = 4
	}
	return indexWidth, formatWidth, sizeWidth, yearWidth
}

func truncateDisplay(s string, maxCols int) string {
	s = strings.TrimSpace(s)
	if runewidth.StringWidth(s) <= maxCols {
		return s
	}
	if maxCols <= 3 {
		w := 0
		var out []rune
		for _, r := range s {
			rw := runewidth.RuneWidth(r)
			if w+rw > maxCols {
				break
			}
			out = append(out, r)
			w += rw
		}
		return string(out)
	}
	target := maxCols - 3
	w := 0
	var out []rune
	for _, r := range s {
		rw := runewidth.RuneWidth(r)
		if w+rw > target {
			break
		}
		out = append(out, r)
		w += rw
	}
	return string(out) + "..."
}

func padDisplay(s string, width int) string {
	if runewidth.StringWidth(s) >= width {
		return s
	}
	return s + strings.Repeat(" ", width-runewidth.StringWidth(s))
}

func padLeftDisplay(s string, width int) string {
	if runewidth.StringWidth(s) >= width {
		return s
	}
	return strings.Repeat(" ", width-runewidth.StringWidth(s)) + s
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func searchCmd(keyword string) tea.Cmd {
	return func() tea.Msg {
		results, err := book.SearchBooks(keyword, book.DefaultProviders())
		return searchDoneMsg{keyword: keyword, results: results, err: err}
	}
}

func downloadCmd(selected book.BookResult, destDir string) tea.Cmd {
	return func() tea.Msg {
		output, err := book.ZlibDownloadQuiet(selected.ID, destDir)
		return downloadDoneMsg{title: selected.Title, output: output, err: err}
	}
}

func statusCmd(text string, isErr bool) tea.Cmd {
	return func() tea.Msg {
		return statusMsg{text: text, isErr: isErr}
	}
}

func normalizeQuery(input string) string {
	query := strings.TrimSpace(input)
	lower := strings.ToLower(query)
	switch {
	case strings.HasPrefix(lower, "search "):
		return strings.TrimSpace(query[len("search "):])
	case strings.HasPrefix(lower, "book "):
		return strings.TrimSpace(query[len("book "):])
	default:
		return query
	}
}

func compactOutput(output, fallback string) string {
	text := strings.TrimSpace(output)
	text = strings.ReplaceAll(text, "\r", "\n")
	lines := strings.FieldsFunc(text, func(r rune) bool {
		return r == '\n'
	})
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if strings.Contains(line, "Saved to:") || strings.Contains(line, "saved to:") {
			return line
		}
	}
	if len(lines) > 0 {
		return lines[len(lines)-1]
	}
	if fallback == "" {
		return ""
	}
	return fallback
}
