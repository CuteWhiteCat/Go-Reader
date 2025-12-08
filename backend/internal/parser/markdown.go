package parser

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/whitecat/go-reader/internal/models"
)

// MarkdownParser parses .md files
type MarkdownParser struct{}

// NewMarkdownParser creates a new MarkdownParser
func NewMarkdownParser() *MarkdownParser {
	return &MarkdownParser{}
}

// Parse parses a markdown file and returns chapters
// Chapters are determined by ## or # headers
func (p *MarkdownParser) Parse(filePath string) ([]models.Chapter, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	var chapters []models.Chapter
	var currentChapter *models.Chapter
	chapterNumber := 0

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	var contentBuilder strings.Builder

	for scanner.Scan() {
		line := scanner.Text()

		// Check if line is a markdown header
		if strings.HasPrefix(strings.TrimSpace(line), "#") {
			// Save previous chapter if exists
			if currentChapter != nil {
				currentChapter.Content = contentBuilder.String()
				currentChapter.WordCount = len(strings.Fields(currentChapter.Content))
				chapters = append(chapters, *currentChapter)
			}

			// Start new chapter
			chapterNumber++
			currentChapter = &models.Chapter{
				ID:                  uuid.New().String(),
				ChapterNumber:       chapterNumber,
				VolumeNumber:        1,
				VolumeChapterNumber: chapterNumber,
				Title:               strings.TrimSpace(strings.TrimLeft(line, "#")),
			}
			contentBuilder.Reset()
		} else if currentChapter != nil {
			// Add content to current chapter
			contentBuilder.WriteString(line)
			contentBuilder.WriteString("\n")
		}
	}

	// Save last chapter
	if currentChapter != nil {
		currentChapter.Content = contentBuilder.String()
		currentChapter.WordCount = len(strings.Fields(currentChapter.Content))
		chapters = append(chapters, *currentChapter)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	// If no chapters found, create a single chapter with all content
	if len(chapters) == 0 {
		file.Seek(0, 0)
		content, err := os.ReadFile(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to read file: %w", err)
		}

		return []models.Chapter{
			{
				ID:                  uuid.New().String(),
				ChapterNumber:       1,
				VolumeNumber:        1,
				VolumeChapterNumber: 1,
				Title:               "Chapter 1",
				Content:             string(content),
				WordCount:           len(strings.Fields(string(content))),
			},
		}, nil
	}

	return chapters, nil
}
