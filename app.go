package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/bmaupin/go-epub"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx           context.Context
	cookieManager *CookieManager
	// 用于通知前端的进度channel
	progressChan chan DownloadProgress
}

// 添加新的类型定义
type ToValues struct {
	Txt  string `json:"txt"`
	Epub string `json:"epub"`
	Html string `json:"html"`
}

type DelayConfig struct {
	Min int `json:"min"`
	Max int `json:"max"`
}

// Config 应用配置
type Config struct {
	SavePath    string      `json:"savePath"`
	SaveFormat  string      `json:"saveFormat"`
	ThreadCount int         `json:"threadCount"`
	ConvertType ToValues    `json:"convertType"`
	Delay       DelayConfig `json:"delay"`
	RetryConfig RetryConfig `json:"retryConfig"`
	LogConfig   LogConfig   `json:"logConfig"`
}

// 新增重试置
type RetryConfig struct {
	MaxRetries int `json:"maxRetries"`
	RetryDelay int `json:"retryDelay"`
}

// 新增日志配置
type LogConfig struct {
	Enabled bool   `json:"enabled"`
	Level   string `json:"level"`
}

// Book 结构体用于表示书籍信息
type Book struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Author    string `json:"author"`
	Cover     string `json:"cover"`
	WordCount int    `json:"wordCount"`
	Status    string `json:"status"`
}

// Chapter 结构体表示章节信息
type Chapter struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	ID      string `json:"id"`
}

// BookDetail 结构体表示书籍详细信息
type BookDetail struct {
	Title    string             `json:"title"`
	Chapters map[string]Chapter `json:"chapters"`
	Status   string             `json:"status"`
}

// DownloadProgress 结构体用于表示下载进度
type DownloadProgress struct {
	Total        int    `json:"total"`
	Current      int    `json:"current"`
	ChapterTitle string `json:"chapterTitle"`
	Status       string `json:"status"`
}

// Cookie管理相关结构体
type CookieManager struct {
	Cookie     string    `json:"cookie"`
	UpdateTime time.Time `json:"updateTime"`
}

var (
	headers = []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/93.0.4577.63 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:91.0) Gecko/20100101 Firefox/91.0",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/93.0.4577.63 Safari/537.36 Edg/93.0.961.47",
	}

	// 从Python代码移植的字符映射表
	codeRanges = [][2]int{
		{58344, 58715},
		{58345, 58716},
	}

	// 字符映射表
	charsets = [][]string{
		// 第一个映射表
		{"D", "在", "主", "特", "家", "军", "然", "表", "场", "4", "要", "只", "v", "和", "?", "6", "别", "还", "g", "现", "儿", "岁", "?", "?", "此", "象", "月", "3", "出", "战", "工", "相", "o", "男", "直", "失", "世", "F", "都", "平", "文", "什", "V", "O", "将", "真", "T", "那", "当", "?", "会", "立", "些", "u", "是", "十", "张", "学", "气", "大", "爱", "两", "命", "全", "后", "东", "性", "通", "被", "1", "它", "乐", "接", "而", "感", "车", "山", "公", "了", "常", "以", "何", "可", "话", "先", "p", "i", "叫", "轻", "M", "士", "w", "着", "变", "尔", "快", "l", "个", "说", "色", "里", "安", "花", "远", "7", "难", "师", "放", "t", "报", "认", "面", "道", "S", "?", "克", "地", "度", "I", "好", "机", "U", "民", "写", "把", "万", "同", "水", "新", "没", "书", "电", "吃", "像", "斯", "5", "为", "y", "白", "几", "日", "教", "看", "但", "第", "加", "候", "作", "上", "拉", "住", "有", "法", "r", "事", "应", "位", "利", "你", "声", "身", "国", "问", "马", "女", "他", "Y", "比", "父", "x", "A", "H", "N", "s", "X", "边", "美", "对", "所", "金", "活", "回", "意", "到", "z", "从", "j", "知", "又", "内", "因", "点", "Q", "三", "定", "8", "R", "b", "正", "或", "夫", "向", "德", "听", "更", "?", "得", "告", "并", "本", "q", "过", "记", "L", "让", "打", "f", "人", "就", "者", "去", "原", "满", "体", "做", "经", "K", "走", "如", "孩", "c", "G", "给", "使", "物", "?", "最", "笑", "部", "?", "员", "等", "受", "k", "行", "一", "条", "果", "动", "光", "门", "头", "见", "往", "自", "解", "成", "处", "天", "能", "于", "名", "其", "发", "总", "母", "的", "死", "手", "入", "路", "进", "心", "来", "h", "时", "力", "多", "开", "已", "许", "d", "至", "由", "很", "界", "n", "小", "与", "Z", "想", "代", "么", "分", "生", "口", "再", "妈", "望", "次", "西", "风", "种", "带", "J", "?", "实", "情", "才", "这", "?", "E", "我", "神", "格", "长", "觉", "间", "年", "眼", "无", "不", "亲", "关", "结", "0", "友", "信", "下", "却", "重", "己", "老", "2", "音", "字", "m", "呢", "明", "之", "前", "高", "P", "B", "目", "太", "e", "9", "起", "稜", "她", "也", "W", "用", "方", "子", "英", "每", "理", "便", "四", "数", "期", "中", "C", "外", "样", "a", "海", "们", "任"},
		// 第二个映射表
		{"s", "?", "作", "口", "在", "他", "能", "并", "B", "士", "4", "U", "克", "才", "正", "们", "字", "声", "高", "全", "尔", "活", "者", "动", "其", "主", "报", "多", "望", "放", "h", "w", "次", "年", "?", "中", "3", "特", "于", "十", "入", "要", "男", "同", "G", "面", "分", "方", "K", "什", "再", "教", "本", "己", "结", "1", "等", "世", "N", "?", "说", "g", "u", "期", "Z", "外", "美", "M", "行", "给", "9", "文", "将", "两", "许", "张", "友", "0", "英", "应", "向", "像", "此", "白", "安", "少", "何", "打", "气", "常", "定", "间", "花", "见", "孩", "它", "直", "风", "数", "使", "道", "第", "水", "已", "女", "山", "解", "d", "P", "的", "通", "关", "性", "叫", "儿", "L", "妈", "问", "回", "神", "来", "S", "", "四", "望", "前", "国", "些", "O", "v", "l", "A", "心", "平", "自", "无", "军", "光", "代", "是", "好", "却", "c", "得", "种", "就", "意", "先", "立", "z", "子", "过", "Y", "j", "表", "", "么", "所", "接", "了", "名", "金", "受", "J", "满", "眼", "没", "部", "那", "m", "每", "车", "度", "可", "R", "斯", "经", "现", "门", "明", "V", "如", "走", "命", "y", "6", "E", "战", "很", "上", "f", "月", "西", "7", "长", "", "想", "话", "变", "海", "机", "x", "到", "W", "一", "成", "生", "信", "笑", "但", "父", "开", "内", "东", "马", "日", "小", "而", "后", "带", "以", "三", "几", "为", "认", "X", "死", "员", "目", "位", "之", "", "远", "人", "音", "呢", "我", "q", "乐", "象", "重", "对", "个", "被", "别", "F", "也", "书", "稜", "D", "写", "还", "因", "家", "发", "时", "i", "或", "住", "德", "当", "o", "l", "比", "觉", "然", "吃", "去", "公", "a", "老", "亲", "情", "体", "太", "b", "万", "C", "电", "理", "?", "失", "力", "更", "拉", "物", "着", "原", "她", "工", "实", "色", "感", "记", "看", "出", "相", "路", "大", "你", "候", "2", "和", "?", "与", "p", "样", "新", "只", "便", "最", "不", "进", "T", "r", "做", "格", "母", "总", "爱", "身", "师", "轻", "知", "往", "加", "从", "?", "天", "e", "H", "?", "听", "场", "由", "快", "边", "让", "把", "任", "8", "条", "头", "事", "至", "起", "点", "真", "手", "这", "难", "", "界", "用", "法", "n", "处", "下", "又", "Q", "告", "地", "5", "k", "t", "岁", "有", "会", "果", "利", "?"},
	}
)

// NewApp 创建新的App实例
func NewApp() *App {
	return &App{
		cookieManager: &CookieManager{},
		progressChan:  make(chan DownloadProgress, 1),
	}
}

// 启动时初始化
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	// 确保配置目录存在
	configDir := filepath.Join("data")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		fmt.Printf("创建配置目录失败: %v\n", err)
	}

	// 需要在这里加载已保存的cookie
	cookiePath := filepath.Join("data", "cookie.json")
	if data, err := os.ReadFile(cookiePath); err == nil {
		json.Unmarshal(data, a.cookieManager)
	}
}

// GetConfig 获取配置
func (a *App) GetConfig() (Config, error) {
	var config Config
	configPath := filepath.Join("data", "config.json")

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			// 如果配置文件不存在,创建默认配置
			config = getDefaultConfig()
			return config, nil
		}
		return config, err
	}

	err = json.Unmarshal(data, &config)
	return config, err
}

// SaveConfig 保存配置
func (a *App) SaveConfig(config Config) error {
	configPath := filepath.Join("data", "config.json")
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(configPath, data, 0644)
}

// SearchBooks 搜索书籍
func (a *App) SearchBooks(keyword string) ([]Book, error) {
	url := fmt.Sprintf("https://api5-normal-lf.fqnovel.com/reading/bookapi/search/page/v/?query=%s&aid=1967&channel=0&os_version=0&device_type=0&device_platform=0&iid=466614321180296", keyword)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Code int `json:"code"`
		Data []struct {
			BookData []struct {
				BookID     string `json:"book_id"`
				BookName   string `json:"book_name"`
				Author     string `json:"author"`
				Cover      string `json:"cover"`
				WordNumber int    `json:"word_number"`
				Status     string `json:"status"`
			} `json:"book_data"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	var books []Book
	for _, item := range result.Data {
		if len(item.BookData) > 0 {
			book := item.BookData[0]
			books = append(books, Book{
				ID:        book.BookID,
				Title:     book.BookName,
				Author:    book.Author,
				Cover:     book.Cover,
				WordCount: book.WordNumber,
				Status:    book.Status,
			})
		}
	}

	return books, nil
}

// DownloadBook 实现书籍下载
func (a *App) DownloadBook(bookID string, format string) error {
	// 1. 获取书籍信息
	bookDetail, err := a.getBookDetail(bookID)
	if err != nil {
		return fmt.Errorf("获取书籍信息失败: %v", err)
	}

	// 2. 下载所有章节
	err = a.downloadChapters(bookDetail)
	if err != nil {
		return fmt.Errorf("下载章节失败: %v", err)
	}

	// 3. 根据格式保存文件
	err = a.saveBook(bookDetail, format)
	if err != nil {
		return fmt.Errorf("保存文件失败: %v", err)
	}

	return nil
}

// getBookDetail 获取书籍详细信息
func (a *App) getBookDetail(bookID string) (*BookDetail, error) {
	url := fmt.Sprintf("https://fanqienovel.com/page/%s", bookID)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", getRandomHeader())

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	title := doc.Find("h1").Text()
	status := doc.Find("span.info-label-yellow").Text()

	chapters := make(map[string]Chapter)
	doc.Find("div.chapter div a").Each(func(i int, s *goquery.Selection) {
		chapterTitle := s.Text()
		href, _ := s.Attr("href")
		chapterID := strings.Split(href, "/")[len(strings.Split(href, "/"))-1]

		chapters[chapterTitle] = Chapter{
			Title: chapterTitle,
			ID:    chapterID,
		}
	})

	return &BookDetail{
		Title:    title,
		Chapters: chapters,
		Status:   status,
	}, nil
}

// downloadChapters 下载所有章节
func (a *App) downloadChapters(book *BookDetail) error {
	config, err := a.GetConfig()
	if err != nil {
		return err
	}

	type downloadResult struct {
		Title      string
		Content    string
		Error      error
		RetryCount int
	}

	jobs := make(chan Chapter, len(book.Chapters))
	results := make(chan downloadResult, len(book.Chapters))

	// 启动工作协程
	for w := 0; w < config.ThreadCount; w++ {
		go func() {
			for chapter := range jobs {
				var result downloadResult
				result.Title = chapter.Title

				// 添加重试逻辑
				for result.RetryCount < config.RetryConfig.MaxRetries {
					content, err := a.downloadChapterContent(chapter.ID)
					if err == nil {
						result.Content = content
						break
					}

					result.Error = err
					result.RetryCount++

					if result.RetryCount < config.RetryConfig.MaxRetries {
						time.Sleep(time.Duration(config.RetryConfig.RetryDelay) * time.Second)
					}
				}

				results <- result
				time.Sleep(time.Duration(rand.Intn(config.Delay.Max-config.Delay.Min)+config.Delay.Min) * time.Millisecond)
			}
		}()
	}

	// 发送任务
	for _, chapter := range book.Chapters {
		jobs <- chapter
	}
	close(jobs)

	// 收集结果
	totalChapters := len(book.Chapters)
	completedChapters := 0

	// 收集结果时添加进度通知
	for i := 0; i < len(book.Chapters); i++ {
		result := <-results
		if result.Error != nil {
			a.NotifyProgress(DownloadProgress{
				Total:        totalChapters,
				Current:      completedChapters,
				ChapterTitle: result.Title,
				Status:       "error",
			})
			return result.Error
		}

		completedChapters++
		chapter := book.Chapters[result.Title]
		chapter.Content = result.Content
		book.Chapters[result.Title] = chapter

		// 发送进度通知
		a.NotifyProgress(DownloadProgress{
			Total:        totalChapters,
			Current:      completedChapters,
			ChapterTitle: result.Title,
			Status:       "downloading",
		})
	}

	// 下载完成通知
	a.NotifyProgress(DownloadProgress{
		Total:        totalChapters,
		Current:      totalChapters,
		ChapterTitle: "完成",
		Status:       "completed",
	})

	return nil
}

// saveBook 根据格式保存书籍
func (a *App) saveBook(book *BookDetail, format string) error {
	switch format {
	case "txt":
		return a.saveAsTxt(book)
	case "epub":
		return a.saveAsEpub(book)
	case "html":
		return a.saveAsHtml(book)
	default:
		return fmt.Errorf("不支持的格式: %s", format)
	}
}

// saveAsTxt 保存为txt格式
func (a *App) saveAsTxt(book *BookDetail) error {
	config, err := a.GetConfig()
	if err != nil {
		return err
	}

	// 确保保存目存在
	if err := os.MkdirAll(config.SavePath, 0755); err != nil {
		return err
	}

	// 创建文件
	filename := filepath.Join(config.SavePath, sanitizeFilename(book.Title+".txt"))
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// 写入内容
	for _, chapter := range sortChapters(book.Chapters) {
		// 写入章节标题
		if _, err := file.WriteString("\n" + chapter.Title + "\n\n"); err != nil {
			return err
		}
		// 写入章节内容
		if _, err := file.WriteString(chapter.Content + "\n"); err != nil {
			return err
		}
	}

	return nil
}

// saveAsEpub 保存为epub格式
func (a *App) saveAsEpub(book *BookDetail) error {
	config, err := a.GetConfig()
	if err != nil {
		return err
	}

	// 创建新的epub书
	e := epub.NewEpub(book.Title)

	// 设置元数据
	e.SetLang("zh")
	e.SetAuthor("Unknown") // TODO: 添加作者信息

	// 添加章节
	for _, chapter := range sortChapters(book.Chapters) {
		// 将章节内容转换为HTML格式
		content := fmt.Sprintf("<h1>%s</h1><div>%s</div>",
			chapter.Title,
			strings.ReplaceAll(chapter.Content, "\n", "<br/>"))

		// 添加章节到epub
		_, err := e.AddSection(content, chapter.Title, "", "")
		if err != nil {
			return err
		}
	}

	// 保存epub文件
	filename := filepath.Join(config.SavePath, sanitizeFilename(book.Title+".epub"))
	return e.Write(filename)
}

// saveAsHtml 保存为html格式
func (a *App) saveAsHtml(book *BookDetail) error {
	config, err := a.GetConfig()
	if err != nil {
		return err
	}

	// 创建保存目录
	bookDir := filepath.Join(config.SavePath, sanitizeFilename(book.Title))
	if err := os.MkdirAll(bookDir, 0755); err != nil {
		return err
	}

	// 创建目录页
	indexContent := `<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <title>%s</title>
    <style>
        body { font-family: Arial, sans-serif; max-width: 800px; margin: 0 auto; padding: 20px; }
        a { color: #333; text-decoration: none; }
        a:hover { text-decoration: underline; }
    </style>
</head>
<body>
    <h1>%s</h1>
    <ul>
`
	indexContent = fmt.Sprintf(indexContent, book.Title, book.Title)

	// 添加章节链接
	for _, chapter := range sortChapters(book.Chapters) {
		indexContent += fmt.Sprintf(`        <li><a href="%s.html">%s</a></li>`,
			sanitizeFilename(chapter.Title), chapter.Title)
	}

	indexContent += `
    </ul>
</body>
</html>`

	// 保存目录页
	if err := os.WriteFile(filepath.Join(bookDir, "index.html"), []byte(indexContent), 0644); err != nil {
		return err
	}

	// 保存章节页面
	chapterTemplate := `<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <title>%s - %s</title>
    <style>
        body { font-family: Arial, sans-serif; max-width: 800px; margin: 0 auto; padding: 20px; }
        .chapter-content { line-height: 1.6; }
        .nav { margin: 20px 0; }
        a { color: #333; text-decoration: none; }
        a:hover { text-decoration: underline; }
    </style>
</head>
<body>
    <div class="nav">
        <a href="index.html">返回目录</a>
    </div>
    <h1>%s</h1>
    <div class="chapter-content">%s</div>
</body>
</html>`

	for _, chapter := range sortChapters(book.Chapters) {
		content := strings.ReplaceAll(chapter.Content, "\n", "<br/>")
		htmlContent := fmt.Sprintf(chapterTemplate,
			chapter.Title, book.Title,
			chapter.Title, content)

		filename := filepath.Join(bookDir, sanitizeFilename(chapter.Title)+".html")
		if err := ioutil.WriteFile(filename, []byte(htmlContent), 0644); err != nil {
			return err
		}
	}

	return nil
}

// 工具函数
func getRandomHeader() string {
	return headers[rand.Intn(len(headers))]
}

// 字符解码函数
func interpreter(uni int, mode int) string {
	bias := uni - codeRanges[mode][0]
	if bias < 0 || bias >= len(charsets[mode]) {
		return string(rune(uni))
	}
	return charsets[mode][bias]
}

// 字符串解码函数
func strInterpreter(text string, mode int) string {
	if text == "" {
		return ""
	}

	var result strings.Builder
	result.Grow(len(text)) // 预分配空间提升性能

	for _, r := range text {
		uni := int(r)
		if uni >= len(codeRanges) || mode >= len(codeRanges) {
			result.WriteRune(r)
			continue
		}

		if codeRanges[mode][0] <= uni && uni <= codeRanges[mode][1] {
			decoded := interpreter(uni, mode)
			result.WriteString(decoded)
		} else {
			result.WriteRune(r)
		}
	}

	return result.String()
}

// 下载章节内容
func (a *App) downloadChapterContent(chapterID string) (string, error) {
	url := fmt.Sprintf("https://fanqienovel.com/reader/%s", chapterID)
	return a.downloadWithRetry(url, 3)
}

// 添加重试机制的下载函数
func (a *App) downloadWithRetry(url string, maxRetries int) (string, error) {
	var lastErr error
	for i := 0; i < maxRetries; i++ {
		content, err := a.downloadOnce(url)
		if err == nil {
			return content, nil
		}

		lastErr = err
		// 如果是Cookie失效，刷新Cookie后继续
		if strings.Contains(err.Error(), "Cookie无效") {
			if err := a.refreshCookie(); err != nil {
				return "", fmt.Errorf("刷新Cookie失败: %v", err)
			}
			continue
		}

		// 其他错误等待后重试
		time.Sleep(time.Duration(i+1) * time.Second)
	}

	return "", fmt.Errorf("重试%d次后仍然失败: %v", maxRetries, lastErr)
}

func (a *App) downloadOnce(url string) (string, error) {
	client := &http.Client{
		Timeout: 30 * time.Second, // 添加超时设置
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %v", err)
	}

	req.Header.Set("User-Agent", getRandomHeader())
	req.Header.Set("Cookie", a.getCookie())

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP错误: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", fmt.Errorf("解析HTML失败: %v", err)
	}

	var content strings.Builder

	// 改进选择器
	doc.Find("div.muye-reader-content.noselect p").Each(func(i int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())
		if text != "" {
			content.WriteString(text)
			content.WriteString("\n")
		}
	})

	if content.Len() < 100 {
		return "", fmt.Errorf("内容长度异常，可能是Cookie无效")
	}

	// 解码内容
	decodedContent := strInterpreter(content.String(), 0)
	return decodedContent, nil
}

// 添加进度通知方法
func (a *App) NotifyProgress(progress DownloadProgress) {
	// 通过Wails运行时发送事件到前端
	runtime.EventsEmit(a.ctx, "download:progress", progress)
}

// Cookie管理相关方法
func (a *App) getCookie() string {
	if a.cookieManager.Cookie == "" || time.Since(a.cookieManager.UpdateTime) > 30*time.Minute {
		for i := 0; i < 3; i++ { // 添加重试逻辑
			if err := a.refreshCookie(); err != nil {
				fmt.Printf("第%d次刷新Cookie失败: %v\n", i+1, err)
				time.Sleep(time.Second * time.Duration(i+1))
				continue
			}
			break
		}
	}
	return a.cookieManager.Cookie
}

func (a *App) refreshCookie() error {
	baseNum := int64(1000000000000000000)
	maxTries := 10

	// 增加随机范围
	for try := 0; try < maxTries; try++ {
		// 参考 Python 版本的范围：random.randint(bas * 6, bas * 8)
		minNum := baseNum * 6
		maxNum := baseNum * 8
		randNum := minNum + rand.Int63n(maxNum-minNum)
		cookie := fmt.Sprintf("novel_web_id=%d", randNum)

		// 添加合适的延迟
		time.Sleep(time.Duration(50+rand.Intn(100)) * time.Millisecond)

		if a.validateCookie(cookie) {
			a.cookieManager.Cookie = cookie
			a.cookieManager.UpdateTime = time.Now()
			return a.saveCookieToFile()
		}
	}
	return fmt.Errorf("无法获取有效的Cookie,请稍后重试")
}

func (a *App) validateCookie(cookie string) bool {
	client := &http.Client{}
	// 应该使用 API 接口而不是网页来验证
	req, err := http.NewRequest("GET", "https://fanqienovel.com/api/reader/full?itemId=7143038691944959011", nil)
	if err != nil {
		return false
	}

	req.Header.Set("Cookie", cookie)
	req.Header.Set("User-Agent", getRandomHeader())

	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false
	}

	// 检查响应长度和内容
	return len(body) > 200
}

func (a *App) saveCookieToFile() error {
	data, err := json.Marshal(a.cookieManager)
	if err != nil {
		return err
	}

	cookiePath := filepath.Join("data", "cookie.json")
	return os.WriteFile(cookiePath, data, 0644)
}

// 工具函数
func sanitizeFilename(filename string) string {
	// 替换Windows文件名中的非法字符
	illegal := []string{"<", ">", ":", "\"", "/", "\\", "|", "?", "*"}
	legal := []string{"＜", "＞", "：", "＂", "／", "＼", "｜", "？", "＊"}

	for i := range illegal {
		filename = strings.ReplaceAll(filename, illegal[i], legal[i])
	}
	return filename
}

// 对章节进行排序
func sortChapters(chapters map[string]Chapter) []Chapter {
	var sorted []Chapter
	for _, chapter := range chapters {
		sorted = append(sorted, chapter)
	}
	// 按照章节ID排序
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].ID < sorted[j].ID
	})
	return sorted
}

// 添加日志功能
func (a *App) logError(format string, args ...interface{}) {
	config, err := a.GetConfig()
	if err != nil {
		return
	}

	if config.LogConfig.Enabled {
		msg := fmt.Sprintf(format, args...)
		// 可以根据需要将日志写入文件或发送到前端
		runtime.EventsEmit(a.ctx, "log:error", msg)
	}
}

func getDefaultConfig() Config {
	return Config{
		SavePath:    "",
		SaveFormat:  "txt",
		ThreadCount: 1,
		ConvertType: ToValues{
			Txt:  "txt",
			Epub: "epub",
			Html: "html",
		},
		Delay: DelayConfig{
			Min: 50,
			Max: 150,
		},
		RetryConfig: RetryConfig{
			MaxRetries: 3,
			RetryDelay: 2,
		},
		LogConfig: LogConfig{
			Enabled: true,
			Level:   "info",
		},
	}
}

// 自定义错误类型
type DownloadError struct {
	ChapterTitle string
	Err          error
	RetryCount   int
}

func (e *DownloadError) Error() string {
	return fmt.Sprintf("下载章节 %s 失败: %v (重试次数: %d)", e.ChapterTitle, e.Err, e.RetryCount)
}

// 添加错误处理中间件
func (a *App) withErrorHandling(operation func() error) error {
	err := operation()
	if err != nil {
		a.logError("操作失败: %v", err)
		// 可以在这里添加错误恢复逻辑
		return err
	}
	return nil
}

// 新增一个导出方法
func (a *App) GetBookDetail(bookID string) (*BookDetail, error) {
	return a.getBookDetail(bookID)
}

func (a *App) GetCookie() string {
	return a.getCookie()
}
