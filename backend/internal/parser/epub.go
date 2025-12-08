package parser

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"io"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"unicode"

	"github.com/google/uuid"
	"github.com/whitecat/go-reader/internal/models"
)

// EpubParser parses .epub files
type EpubParser struct{}

// NewEpubParser creates a new EpubParser
func NewEpubParser() *EpubParser {
	return &EpubParser{}
}

// EPUB structures for XML parsing
const (
	opfNS = "http://www.idpf.org/2007/opf"
	dcNS  = "http://purl.org/dc/elements/1.1/"
)

type Container struct {
	Rootfiles []Rootfile `xml:"rootfiles>rootfile"`
}

type Rootfile struct {
	FullPath string `xml:"full-path,attr"`
}

type Package struct {
	Metadata Metadata  `xml:"{http://www.idpf.org/2007/opf}metadata"`
	Manifest []Item    `xml:"{http://www.idpf.org/2007/opf}manifest>{http://www.idpf.org/2007/opf}item"`
	Spine    SpineData `xml:"{http://www.idpf.org/2007/opf}spine"`
}

type Metadata struct {
	Title []string `xml:"{http://purl.org/dc/elements/1.1/}title"`
}

type Item struct {
	ID        string `xml:"id,attr"`
	Href      string `xml:"href,attr"`
	MediaType string `xml:"media-type,attr"`
}

type SpineData struct {
	ItemRefs []ItemRef `xml:"{http://www.idpf.org/2007/opf}itemref"`
}

type ItemRef struct {
	IDRef string `xml:"idref,attr"`
}

// Parse parses an epub file and returns chapters
func (p *EpubParser) Parse(filePath string) ([]models.Chapter, error) {
	// Open ZIP file (EPUB is a ZIP file)
	reader, err := zip.OpenReader(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open epub: %w", err)
	}
	defer reader.Close()

	// Find container.xml
	containerPath := "META-INF/container.xml"
	containerFile := findFileInZip(&reader.Reader, containerPath)
	if containerFile == nil {
		return nil, fmt.Errorf("container.xml not found")
	}

	// Parse container.xml to find OPF file
	rc, err := containerFile.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open container.xml: %w", err)
	}
	defer rc.Close()

	var container Container
	if err := xml.NewDecoder(rc).Decode(&container); err != nil {
		return nil, fmt.Errorf("failed to parse container.xml: %w", err)
	}

	if len(container.Rootfiles) == 0 {
		return nil, fmt.Errorf("no rootfile found in container")
	}

	// Parse OPF file
	opfPath := container.Rootfiles[0].FullPath
	opfFile := findFileInZip(&reader.Reader, opfPath)
	if opfFile == nil {
		return nil, fmt.Errorf("OPF file not found: %s", opfPath)
	}

	rc, err = opfFile.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open OPF file: %w", err)
	}
	defer rc.Close()

	var pkg Package
	if err := xml.NewDecoder(rc).Decode(&pkg); err != nil {
		return nil, fmt.Errorf("failed to parse OPF file: %w", err)
	}

	// Get base directory from OPF path (normalize slashes first)
	baseDir := path.Dir(filepath.ToSlash(opfPath))

	// Create chapters from spine
	var chapters []models.Chapter
	manifest := make(map[string]Item)
	for _, it := range pkg.Manifest {
		manifest[it.ID] = it
	}

	for i, itemRef := range pkg.Spine.ItemRefs {
		// Find manifest item
		item, ok := manifest[itemRef.IDRef]
		if !ok {
			continue
		}
		if !strings.Contains(item.MediaType, "html") {
			// Skip non-html items such as nav/cover/images
			continue
		}
		href := item.Href

		// Resolve path
		fullPath := path.Clean(path.Join(baseDir, href))
		contentFile := findFileInZip(&reader.Reader, fullPath)
		if contentFile == nil {
			continue
		}

		// Read content
		rc, err := contentFile.Open()
		if err != nil {
			continue
		}

		content, err := io.ReadAll(rc)
		rc.Close()
		if err != nil {
			continue
		}

		// Extract text
		raw := string(content)
		contentStr := stripHTMLTags(raw)
		if strings.TrimSpace(contentStr) == "" {
			continue
		}

		// Create chapter
		title := fmt.Sprintf("Chapter %d", i+1)
		if i == 0 && len(pkg.Metadata.Title) > 0 {
			title = pkg.Metadata.Title[0]
		} else if t := extractHTMLTitle(raw); t != "" {
			title = t
		}

		chapter := models.Chapter{
			ID:                  uuid.New().String(),
			ChapterNumber:       i + 1,
			VolumeNumber:        1,
			VolumeChapterNumber: i + 1,
			Title:               title,
			Content:             contentStr,
			WordCount:           len([]rune(contentStr)),
		}

		chapters = append(chapters, chapter)
	}

	if len(chapters) == 0 {
		// Fallback: grab all HTML/XHTML files in the zip (excluding nav/cover) sorted by name
		var htmlFiles []string
		for _, f := range reader.File {
			name := strings.ReplaceAll(f.Name, "\\", "/")
			ln := strings.ToLower(name)
			if strings.HasSuffix(ln, ".html") || strings.HasSuffix(ln, ".xhtml") {
				if strings.Contains(ln, "nav") || strings.Contains(ln, "cover") || strings.Contains(ln, "toc") {
					continue
				}
				htmlFiles = append(htmlFiles, name)
			}
		}
		sort.Slice(htmlFiles, func(i, j int) bool {
			return naturalLess(htmlFiles[i], htmlFiles[j])
		})

		for i, name := range htmlFiles {
			contentFile := findFileInZip(&reader.Reader, name)
			if contentFile == nil {
				continue
			}
			rc, err := contentFile.Open()
			if err != nil {
				continue
			}
			rawBytes, err := io.ReadAll(rc)
			rc.Close()
			if err != nil {
				continue
			}
			raw := string(rawBytes)
			contentStr := stripHTMLTags(raw)
			if strings.TrimSpace(contentStr) == "" {
				continue
			}
			title := extractHTMLTitle(raw)
			if title == "" {
				title = fmt.Sprintf("Chapter %d", i+1)
			}
			chapters = append(chapters, models.Chapter{
				ID:                  uuid.New().String(),
				ChapterNumber:       len(chapters) + 1,
				VolumeNumber:        1,
				VolumeChapterNumber: len(chapters) + 1,
				Title:               title,
				Content:             contentStr,
				WordCount:           len([]rune(contentStr)),
			})
		}
	}

	if len(chapters) == 0 {
		return nil, fmt.Errorf("no chapters found in epub (manifest=%d, spine=%d)", len(pkg.Manifest), len(pkg.Spine.ItemRefs))
	}

	return chapters, nil
}

// naturalLess compares strings by numeric parts to avoid 1,10,100 ordering issues.
func naturalLess(a, b string) bool {
	ai, bi := 0, 0
	for ai < len(a) && bi < len(b) {
		ra, rb := rune(a[ai]), rune(b[bi])
		if unicode.IsDigit(ra) && unicode.IsDigit(rb) {
			// consume full number
			var sa, sb strings.Builder
			for ai < len(a) && unicode.IsDigit(rune(a[ai])) {
				sa.WriteByte(a[ai])
				ai++
			}
			for bi < len(b) && unicode.IsDigit(rune(b[bi])) {
				sb.WriteByte(b[bi])
				bi++
			}
			na, nb := sa.String(), sb.String()
			if na != nb {
				if len(na) != len(nb) {
					return len(na) < len(nb)
				}
				return na < nb
			}
			continue
		}
		if ra != rb {
			return ra < rb
		}
		ai++
		bi++
	}
	return len(a) < len(b)
}

// findFileInZip finds a file in a ZIP archive
func findFileInZip(r *zip.Reader, name string) *zip.File {
	target := strings.ReplaceAll(name, "\\", "/")
	targetLower := strings.ToLower(target)
	for _, f := range r.File {
		normalized := strings.ReplaceAll(f.Name, "\\", "/")
		if normalized == target || strings.ToLower(normalized) == targetLower {
			return f
		}
	}
	return nil
}

// stripHTMLTags removes HTML tags from content (simple implementation)
func stripHTMLTags(content string) string {
	var result strings.Builder
	inTag := false

	for _, char := range content {
		if char == '<' {
			inTag = true
			continue
		}
		if char == '>' {
			inTag = false
			result.WriteRune(' ') // Add space after tag
			continue
		}
		if !inTag {
			result.WriteRune(char)
		}
	}

	// Clean up multiple spaces
	text := strings.TrimSpace(result.String())
	text = strings.Join(strings.Fields(text), " ")

	return text
}

// extractHTMLTitle returns the first <title>...</title> or <h1>... if available (case-insensitive).
func extractHTMLTitle(content string) string {
	lower := strings.ToLower(content)
	findBetween := func(open, close string) string {
		start := strings.Index(lower, open)
		if start == -1 {
			return ""
		}
		end := strings.Index(lower[start+len(open):], close)
		if end == -1 {
			return ""
		}
		return strings.TrimSpace(stripHTMLTags(content[start+len(open) : start+len(open)+end]))
	}
	if t := findBetween("<title", "</title>"); t != "" {
		// remove any attributes from <title ...>
		if idx := strings.Index(t, ">"); idx != -1 && idx < len(t)-1 {
			return strings.TrimSpace(t[idx+1:])
		}
		return t
	}
	if h := findBetween("<h1", "</h1>"); h != "" {
		if idx := strings.Index(h, ">"); idx != -1 && idx < len(h)-1 {
			return strings.TrimSpace(h[idx+1:])
		}
		return h
	}
	return ""
}
