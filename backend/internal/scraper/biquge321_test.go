package scraper

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func TestExtractCoverURL(t *testing.T) {
	path := filepath.Join("testdata", "biquge_cover.html")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read testdata: %v", err)
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(data)))
	if err != nil {
		t.Fatalf("parse html: %v", err)
	}

	got := extractCoverURL(doc, "https://www.example.com/book/123/")
	want := "https://www.example.com/files/article/image/123/real-cover.png"
	if got != want {
		t.Fatalf("cover url mismatch, want %s, got %s", want, got)
	}
}

func TestExtractCoverURL_LivePage(t *testing.T) {
	const novelURL = "https://www.biquge321.com/xiaoshuo/238022/"

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest(http.MethodGet, novelURL, nil)
	if err != nil {
		t.Fatalf("build request: %v", err)
	}
	req.Header.Set("User-Agent", ua())

	resp, err := client.Do(req)
	if err != nil {
		t.Skipf("skip live cover test (network unavailable): %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Skipf("skip live cover test (unexpected status %d)", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		t.Fatalf("parse live page: %v", err)
	}

	got := extractCoverURL(doc, novelURL)
	if got == "" {
		t.Fatalf("live cover url is empty for %s", novelURL)
	}

	expectedFragment := "/files/article/image/222/222887/222887s.jpg"
	if !strings.Contains(got, expectedFragment) {
		t.Fatalf("unexpected cover url, want fragment %s, got %s", expectedFragment, got)
	}
}
