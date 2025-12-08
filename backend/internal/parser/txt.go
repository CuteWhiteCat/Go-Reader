package parser

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
	runelib "unicode"

	"github.com/google/uuid"
	"github.com/saintfish/chardet"
	"github.com/whitecat/go-reader/internal/models"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/ianaindex"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/traditionalchinese"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

// TxtParser parses .txt files
type TxtParser struct{}

// NewTxtParser creates a new TxtParser
func NewTxtParser() *TxtParser {
	return &TxtParser{}
}

// Parse parses a txt file and returns chapters
func (p *TxtParser) Parse(filePath string) ([]models.Chapter, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Detect and decode file encoding so Chinese text displays correctly
	enc := detectFileEncoding(filePath)
	reader := transform.NewReader(file, enc.NewDecoder())

	var chapters []models.Chapter
	var currentChapter *models.Chapter
	chapterNumber := 0
	volumeNumber := 1
	volumeChapterNumber := 0
	started := false
	awaitingChapter := false

	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanLines)
	// Increase buffer to handle long lines without scan errors
	scanner.Buffer(make([]byte, 0, 64*1024), 2*1024*1024)

	var contentBuilder strings.Builder

	saveCurrentChapter := func(force bool) {
		if currentChapter == nil {
			return
		}

		content := contentBuilder.String()
		if strings.TrimSpace(content) == "" && !force {
			// Skip empty chapters (e.g., consecutive headers or volume markers)
			if chapterNumber > 0 {
				chapterNumber--
			}
			if volumeChapterNumber > 0 {
				volumeChapterNumber--
			}
			currentChapter = nil
			contentBuilder.Reset()
			return
		}

		currentChapter.Content = content
		currentChapter.WordCount = len(strings.Fields(content))
		chapters = append(chapters, *currentChapter)
		currentChapter = nil
		contentBuilder.Reset()
	}

	for scanner.Scan() {
		line := scanner.Text()

		// Detect volume markers (e.g., "第一卷", "卷一", "Volume 1")
		if isVolumeTitle(line) {
			started = true
			awaitingChapter = true
			// Save previous chapter if exists before switching volume
			saveCurrentChapter(false)

			if parsedVol := parseVolumeNumber(line); parsedVol > 0 {
				volumeNumber = parsedVol
			} else if chapterNumber == 0 && volumeNumber == 1 {
				// First volume marker with no explicit number defaults to 1
				volumeNumber = 1
			} else {
				volumeNumber++
			}
			volumeChapterNumber = 0

			// Add a standalone chapter entry for the volume page
			chapterNumber++
			volumePage := models.Chapter{
				ID:                  uuid.New().String(),
				ChapterNumber:       chapterNumber,
				VolumeNumber:        volumeNumber,
				VolumeChapterNumber: 0,
				Title:               strings.TrimSpace(line),
				Content:             "",
				WordCount:           0,
			}
			chapters = append(chapters, volumePage)
			currentChapter = nil
			contentBuilder.Reset()
			continue
		}

		// Check if line is a chapter title (e.g., "Chapter 1", "第1章", etc.)
		if isChapterTitle(line) {
			started = true
			awaitingChapter = false
			// Save previous chapter if exists
			saveCurrentChapter(false)

			// Start new chapter
			chapterNumber++
			volumeChapterNumber++
			currentChapter = &models.Chapter{
				ID:                  uuid.New().String(),
				ChapterNumber:       chapterNumber,
				VolumeNumber:        volumeNumber,
				VolumeChapterNumber: volumeChapterNumber,
				Title:               strings.TrimSpace(line),
			}
			contentBuilder.Reset()
		} else if currentChapter != nil {
			// Add content to current chapter
			contentBuilder.WriteString(line)
			contentBuilder.WriteString("\n")
		} else if started && currentChapter != nil && !awaitingChapter {
			// Only record content after we've seen a volume/chapter marker
			contentBuilder.WriteString(line)
			contentBuilder.WriteString("\n")
		}
	}

	// Save last chapter
	saveCurrentChapter(false)

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	// If no chapters found, create a single chapter with all content
	if len(chapters) == 0 {
		return []models.Chapter{
			{
				ID:                  uuid.New().String(),
				ChapterNumber:       1,
				VolumeNumber:        volumeNumber,
				VolumeChapterNumber: 1,
				Title:               "Chapter 1",
				Content:             contentBuilder.String(),
				WordCount:           len(strings.Fields(contentBuilder.String())),
			},
		}, nil
	}

	return chapters, nil
}

// isChapterTitle checks if a line is likely a chapter title
func isChapterTitle(line string) bool {
	trimmed := strings.TrimSpace(line)
	lineLower := strings.ToLower(trimmed)

	if trimmed == "" {
		return false
	}

	// Common English chapter patterns
	englishPatterns := []string{
		"chapter ",
		"chapter:",
		"ch.",
		"ch ",
	}

	for _, pattern := range englishPatterns {
		if strings.HasPrefix(lineLower, pattern) {
			return true
		}
	}

	// Chinese chapter pattern: must contain both "第" and "章"
	// Examples: "第一章", "第2章", "第二十三章"
	// This prevents matching "第二天" (Day 2) which only has "第"
	if strings.HasPrefix(trimmed, "第") && strings.Contains(trimmed, "章") {
		return true
	}

	return false
}

// isVolumeTitle checks if a line is likely a volume title (e.g., "第一卷", "卷一", "Volume 1")
func isVolumeTitle(line string) bool {
	trimmed := strings.TrimSpace(line)
	lineLower := strings.ToLower(trimmed)

	if trimmed == "" {
		return false
	}

	// Bound length to avoid matching long sentences
	if len([]rune(trimmed)) > 50 {
		return false
	}

	// English patterns with explicit keywords + digits
	english := regexp.MustCompile(`^(volume|vol\.?)\s*\d+(?:\s+.+)?$`)
	if english.MatchString(lineLower) {
		return true
	}

	// Chinese patterns: 第X卷 ... or 卷X ...
	chPattern := regexp.MustCompile(`^(第[\p{Han}\d]+卷|卷[\p{Han}\d]+)(\s+.+)?$`)
	if chPattern.MatchString(trimmed) {
		return true
	}

	return false
}

// detectFileEncoding attempts to detect the file encoding and returns a decoder-compatible encoding.
// Defaults to UTF-8 if detection fails or the charset is unsupported.
func detectFileEncoding(filePath string) encoding.Encoding {
	file, err := os.Open(filePath)
	if err != nil {
		return unicode.UTF8
	}
	defer file.Close()

	sample := make([]byte, 4096)
	n, err := file.Read(sample)
	if err != nil && err != io.EOF {
		return unicode.UTF8
	}

	detector := chardet.NewTextDetector()
	result, err := detector.DetectBest(sample[:n])
	if err != nil || result == nil {
		return unicode.UTF8
	}

	if enc := encodingFromName(result.Charset); enc != nil {
		return enc
	}

	return unicode.UTF8
}

// encodingFromName maps common charset names to Go encodings.
func encodingFromName(name string) encoding.Encoding {
	switch strings.ToLower(strings.TrimSpace(name)) {
	case "utf-8", "utf8":
		return unicode.UTF8
	case "big5", "big5-hkscs":
		return traditionalchinese.Big5
	case "gbk", "gb2312", "gb-18030", "gb18030", "gb_18030":
		return simplifiedchinese.GB18030
	case "windows-1252", "cp1252":
		return charmap.Windows1252
	default:
		enc, _ := ianaindex.MIME.Encoding(name)
		return enc
	}
}

// parseVolumeNumber tries to extract a volume number from a volume title.
// Returns 0 if no number can be parsed.
func parseVolumeNumber(line string) int {
	// Try Arabic digits first (first contiguous sequence)
	digits := strings.Builder{}
	for _, r := range line {
		if runelib.IsDigit(r) {
			digits.WriteRune(r)
		} else if digits.Len() > 0 {
			break
		}
	}
	if digits.Len() > 0 {
		if n, err := strconv.Atoi(digits.String()); err == nil && n > 0 {
			return n
		}
	}

	// Try Chinese numerals by extracting relevant runes
	numerals := strings.Builder{}
	for _, r := range line {
		if isChineseNumeralRune(r) {
			numerals.WriteRune(r)
		} else if numerals.Len() > 0 {
			break
		}
	}

	if numerals.Len() > 0 {
		if n := chineseNumeralToInt(numerals.String()); n > 0 {
			return n
		}
	}

	return 0
}

// chineseNumeralToInt converts a simple Chinese numeral string (e.g., "一", "十二", "一百零三") to int.
// It is intentionally minimal to cover common volume/chapter numbering.
func chineseNumeralToInt(s string) int {
	digits := map[rune]int{
		'零': 0, '〇': 0,
		'一': 1, '二': 2, '两': 2, '三': 3, '四': 4,
		'五': 5, '六': 6, '七': 7, '八': 8, '九': 9,
	}
	units := map[rune]int{
		'十': 10,
		'百': 100,
		'千': 1000,
	}

	total := 0
	current := 0
	runes := []rune(s)

	for i, r := range runes {
		if val, ok := digits[r]; ok {
			current = val
			// If this is the last rune and it's a digit, add it
			if i == len(runes)-1 {
				total += current
			}
			continue
		}

		if unit, ok := units[r]; ok {
			if current == 0 {
				current = 1
			}
			total += current * unit
			current = 0
		}
	}

	return total
}

func isChineseNumeralRune(r rune) bool {
	switch r {
	case '零', '〇', '一', '二', '两', '三', '四', '五', '六', '七', '八', '九', '十', '百', '千':
		return true
	default:
		return false
	}
}
