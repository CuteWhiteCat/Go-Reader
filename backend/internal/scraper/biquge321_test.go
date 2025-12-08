package scraper

import (
	"os"
	"path/filepath"
	"testing"

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
