package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/whitecat/go-reader/internal/scraper"
)

func main() {
	// 使用正确的URL测试
	novelURL := "https://www.biquge321.com/xiaoshuo/238022/"
	fmt.Printf("测试URL: %s\n\n", novelURL)

	// 先直接获取HTML看看结构
	client := &http.Client{Timeout: 15 * time.Second}
	req, _ := http.NewRequest("GET", novelURL, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatalf("解析HTML失败: %v", err)
	}

	fmt.Println("=== 寻找封面图片 ===\n")

	// 查找所有img标签
	doc.Find("img").Each(func(i int, s *goquery.Selection) {
		src, exists := s.Attr("src")
		alt, _ := s.Attr("alt")
		class, _ := s.Attr("class")
		if exists {
			fmt.Printf("图片 %d:\n", i+1)
			fmt.Printf("  src: %s\n", src)
			fmt.Printf("  alt: %s\n", alt)
			fmt.Printf("  class: %s\n", class)
			fmt.Printf("  父元素: %s\n\n", goquery.NodeName(s.Parent()))
		}
	})

	fmt.Println("\n=== 使用 scraper.GetChapterList ===\n")

	// 获取章节列表和封面
	chapters, coverURL, err := scraper.GetChapterList(novelURL)
	if err != nil {
		log.Fatalf("获取章节列表失败: %v", err)
	}

	fmt.Printf("章节数量: %d\n", len(chapters))
	fmt.Printf("封面URL: %s\n", coverURL)

	if coverURL == "" {
		fmt.Println("\n警告: 未找到封面URL!")
	} else {
		fmt.Println("\n封面URL获取成功!")
	}

	// 显示前3个章节
	if len(chapters) > 0 {
		fmt.Println("\n前3个章节:")
		for i := 0; i < 3 && i < len(chapters); i++ {
			fmt.Printf("%d. %s\n", i+1, chapters[i].Title)
		}
	}
}
