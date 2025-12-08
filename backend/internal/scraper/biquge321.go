package scraper

import (
    "fmt"
    "net/http"
    "net/url"
    "regexp"
    "strings"
    "sync"
    "time"

    "github.com/PuerkitoBio/goquery"
    "github.com/liuzl/gocc"
)

// BiQuGe321 provides search and chapter scraping for biquge321.com style sites.
// It works without any account or API; HTML is parsed with goquery.

const (
    bqBaseURL   = "https://www.biquge321.com"
    bqSearchURL = "https://www.biquge321.com/s.php"
)

// NovelResult represents one search result.
type NovelResult struct {
    Title  string `json:"title"`
    Author string `json:"author"`
    Latest string `json:"latest"`
    URL    string `json:"url"`
}

// ChapterInfo is a lightweight chapter DTO.
type ChapterInfo struct {
    Title string
    URL   string
}

var (
    httpClient = &http.Client{Timeout: 15 * time.Second}
    lineClean  = regexp.MustCompile(`(?m)^[=\-_*#]+\s*$`)
    multiNL    = regexp.MustCompile(`\n{3,}`)
    converter  *gocc.OpenCC
)

func init() {
    // Initialize Traditional to Simplified converter
    var err error
    converter, err = gocc.New("t2s") // Traditional to Simplified
    if err != nil {
        fmt.Printf("Warning: Failed to initialize converter: %v\n", err)
    }
}

// traditionalToSimplified converts Traditional Chinese to Simplified Chinese
func traditionalToSimplified(text string) string {
    if converter == nil {
        return text
    }
    result, err := converter.Convert(text)
    if err != nil {
        return text
    }
    return result
}

// Search performs a keyword search and returns novel metadata.
func Search(keyword string) ([]NovelResult, error) {
    // Convert Traditional Chinese to Simplified for better search results
    searchKeyword := traditionalToSimplified(keyword)

    data := url.Values{}
    data.Set("s", searchKeyword)
    data.Set("submit", "")

    req, err := http.NewRequest(http.MethodPost, bqSearchURL, strings.NewReader(data.Encode()))
    if err != nil {
        return nil, fmt.Errorf("build request: %w", err)
    }
    req.Header.Set("User-Agent", ua())
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

    resp, err := httpClient.Do(req)
    if err != nil {
        return nil, fmt.Errorf("do request: %w", err)
    }
    defer resp.Body.Close()

    doc, err := goquery.NewDocumentFromReader(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("parse html: %w", err)
    }

	// Check for rate limit message
	pageText := doc.Text()
	rateLimitSignals := []string{
		"搜索次数已耗尽",
		"搜索过于频繁",
		"搜索次數已耗盡",
		"搜索過於頻繁",
		"一分钟只提供10次搜索机会",
		"一分鐘只提供10次搜索機會",
		"提供10次搜索机会",
		"提供10次搜索機會",
		"防止恶意搜索",
		"防止惡意搜索",
		"防止惡意搜尋",
		"防止恶意搜尋",
		"请稍后再试",
	}
	for _, signal := range rateLimitSignals {
		if strings.Contains(pageText, signal) {
			return nil, fmt.Errorf("rate_limit: 搜索次数已耗尽，请稍后再试")
		}
	}

    var novels []NovelResult
    doc.Find("div.lastupdate ul li").Each(func(idx int, s *goquery.Selection) {
        nameSpan := s.Find("span.name a")
        href, ok := nameSpan.Attr("href")
        if !ok {
            return
        }

        author := strings.TrimSpace(s.Find("span.zuo").Text())
        latestSel := s.Find("span.jie a")
        latest := strings.TrimSpace(latestSel.Text())
        if latest == "" {
            latest = strings.TrimSpace(s.Find("span.jie").Text())
        }

        title := strings.TrimSpace(nameSpan.Text())
        novelURL := joinURL(bqBaseURL, href)

        novels = append(novels, NovelResult{
            Title:  title,
            Author: author,
            Latest: latest,
            URL:    novelURL,
        })
    })

    return novels, nil
}

// GetChapterList fetches the chapter directory for a novel page.
func GetChapterList(novelURL string) ([]ChapterInfo, string, error) {
	req, err := http.NewRequest(http.MethodGet, novelURL, nil)
	if err != nil {
		return nil, "", fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("User-Agent", ua())

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, "", fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, "", fmt.Errorf("parse html: %w", err)
	}

	var chapters []ChapterInfo
	doc.Find("ul.fen_4 a").Each(func(_ int, s *goquery.Selection) {
		href, ok := s.Attr("href")
        if !ok {
            return
        }
        chapters = append(chapters, ChapterInfo{
			Title: strings.TrimSpace(s.Text()),
			URL:   joinURL(novelURL, href),
		})
	})

	coverURL := extractCoverURL(doc, novelURL)

	return chapters, coverURL, nil
}

// FetchChapterContent gets a single chapter text with basic cleanup.
func FetchChapterContent(chapterURL string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, chapterURL, nil)
	if err != nil {
		return "", fmt.Errorf("build request: %w", err)
	}
    req.Header.Set("User-Agent", ua())

    resp, err := httpClient.Do(req)
    if err != nil {
        return "", fmt.Errorf("do request: %w", err)
    }
    defer resp.Body.Close()

    doc, err := goquery.NewDocumentFromReader(resp.Body)
    if err != nil {
        return "", fmt.Errorf("parse html: %w", err)
    }

    content := doc.Find("div#txt")
    content.Find("a, script, style").Remove()
    content.Find("br").ReplaceWithHtml("\n")
    content.Find("p").Each(func(_ int, s *goquery.Selection) {
        s.AfterHtml("\n\n")
    })

    text := strings.TrimSpace(content.Text())
    text = lineClean.ReplaceAllString(text, "")
    text = multiNL.ReplaceAllString(text, "\n\n")
    return text, nil
}

// FetchChapters gets full chapters concurrently (bounded workers).
func FetchChapters(chapterInfos []ChapterInfo) []string {
	return FetchChaptersWithProgress(chapterInfos, nil)
}

// FetchChaptersWithProgress gets full chapters concurrently and reports progress via callback.
func FetchChaptersWithProgress(chapterInfos []ChapterInfo, onProgress func(done int, total int)) []string {
	results := make([]string, len(chapterInfos))
	wg := sync.WaitGroup{}
	sem := make(chan struct{}, 5)
	doneCount := 0
	report := func() {
		if onProgress != nil {
			onProgress(doneCount, len(chapterInfos))
		}
	}
	for i := range chapterInfos {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			sem <- struct{}{}
            defer func() { <-sem }()
			txt, err := FetchChapterContent(chapterInfos[idx].URL)
			if err != nil {
				txt = ""
			}
			results[idx] = txt
			doneCount++
			report()
		}(i)
	}
	wg.Wait()
	report()
	return results
}

func joinURL(base, ref string) string {
	b, err := url.Parse(base)
	if err != nil {
		return ref
	}
    r, err := url.Parse(ref)
    if err != nil {
        return ref
    }
    return b.ResolveReference(r).String()
}

func ua() string {
	return "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115 Safari/537.36"
}

func extractCoverURL(doc *goquery.Document, novelURL string) string {
	// Try common selectors used by various templates
	candidates := []string{
		"a.border_left_a img",
		"#fmimg img",
		".pic img",
		".cover img",
		".bookimg img",
	}
	for _, sel := range candidates {
		if img := doc.Find(sel); img.Length() > 0 {
			if src, ok := img.Attr("src"); ok {
				return joinURL(novelURL, strings.TrimSpace(src))
			}
		}
	}

	// Try og:image
	if meta := doc.Find("meta[property='og:image']"); meta.Length() > 0 {
		if content, ok := meta.Attr("content"); ok {
			return joinURL(novelURL, strings.TrimSpace(content))
		}
	}

	return ""
}
