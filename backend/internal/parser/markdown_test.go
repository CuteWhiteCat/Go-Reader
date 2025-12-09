package parser

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createTestMarkdownFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	filePath := filepath.Join(dir, "test.md")
	err := os.WriteFile(filePath, []byte(content), 0644)
	assert.NoError(t, err)
	return filePath
}

func TestMarkdownParser_Parse(t *testing.T) {
	parser := NewMarkdownParser()

	t.Run("Multiple Chapters", func(t *testing.T) {
		content := `
# Chapter 1
This is the first chapter.

## Chapter 2
This is the second chapter.
With two lines.
`
		filePath := createTestMarkdownFile(t, content)
		chapters, err := parser.Parse(filePath)
		assert.NoError(t, err)
		assert.Len(t, chapters, 2)
		assert.Equal(t, "Chapter 1", chapters[0].Title)
		assert.Contains(t, chapters[0].Content, "first chapter")
		assert.Equal(t, "Chapter 2", chapters[1].Title)
		assert.Contains(t, chapters[1].Content, "second chapter")
	})

	t.Run("No Chapters", func(t *testing.T) {
		content := "This is a file with no chapter headers."
		filePath := createTestMarkdownFile(t, content)
		chapters, err := parser.Parse(filePath)
		assert.NoError(t, err)
		assert.Len(t, chapters, 1)
		assert.Equal(t, "Chapter 1", chapters[0].Title)
		assert.Equal(t, content, chapters[0].Content)
	})

	t.Run("Empty File", func(t *testing.T) {
		filePath := createTestMarkdownFile(t, "")
		chapters, err := parser.Parse(filePath)
		assert.NoError(t, err)
		assert.Len(t, chapters, 1)
		assert.Equal(t, "Chapter 1", chapters[0].Title)
		assert.Equal(t, "", chapters[0].Content)
	})

	t.Run("File with only headers", func(t *testing.T) {
		content := `
# Chapter 1
# Chapter 2
## Chapter 3
`
		filePath := createTestMarkdownFile(t, content)
		chapters, err := parser.Parse(filePath)
		assert.NoError(t, err)
		assert.Len(t, chapters, 3)
		assert.Equal(t, "Chapter 1", chapters[0].Title)
		assert.Equal(t, "Chapter 2", chapters[1].Title)
		assert.Equal(t, "Chapter 3", chapters[2].Title)
		assert.Equal(t, "", chapters[0].Content)
	})
}