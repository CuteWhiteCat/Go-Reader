package main

import (
	"fmt"
	"log"

	"github.com/whitecat/go-reader/internal/scraper"
)

func main() {
	// 测试搜索 "在大宋破碎虚空"
	keyword := "在大宋破碎虚空"
	fmt.Printf("搜索关键字: %s\n\n", keyword)

	results, err := scraper.Search(keyword)
	if err != nil {
		log.Fatalf("搜索失败: %v", err)
	}

	if len(results) == 0 {
		fmt.Println("没有找到结果")
		return
	}

	fmt.Printf("找到 %d 个结果\n\n", len(results))

	// 使用第一个结果
	novel := results[0]
	fmt.Printf("书名: %s\n", novel.Title)
	fmt.Printf("作者: %s\n", novel.Author)
	fmt.Printf("URL: %s\n\n", novel.URL)

	// 获取章节列表和封面
	fmt.Println("正在获取章节列表和封面...")
	chapters, coverURL, err := scraper.GetChapterList(novel.URL)
	if err != nil {
		log.Fatalf("获取章节列表失败: %v", err)
	}

	fmt.Printf("\n章节数量: %d\n", len(chapters))
	fmt.Printf("封面URL: %s\n", coverURL)

	if coverURL == "" {
		fmt.Println("\n警告: 未找到封面URL!")
	} else {
		fmt.Println("\n封面URL获取成功!")
	}
}
