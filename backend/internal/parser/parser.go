package parser

import (
	"fmt"
	"strings"

	"github.com/whitecat/go-reader/internal/models"
)

// Parser interface for different file formats
type Parser interface {
	Parse(filePath string) ([]models.Chapter, error)
}

// GetParser returns the appropriate parser for the given file format
func GetParser(format string) (Parser, error) {
	format = strings.ToLower(strings.TrimSpace(format))

	switch format {
	case "txt":
		return NewTxtParser(), nil
	case "md", "markdown":
		return NewMarkdownParser(), nil
	case "epub":
		return NewEpubParser(), nil
	default:
		return nil, fmt.Errorf("unsupported file format: %s", format)
	}
}
