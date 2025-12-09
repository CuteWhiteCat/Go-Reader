package parser

import (
	"archive/zip"
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStripHTMLTags(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Simple HTML", "<p>Hello</p> <b>World</b>", "Hello World"},
		{"With attributes", `<a href="test">Link</a>`, "Link"},
		{"No HTML", "Just plain text.", "Just plain text."},
		{"Empty", "", ""},
		{"Multiple spaces", "  Hello   World  ", "Hello World"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, stripHTMLTags(tt.input))
		})
	}
}

func TestExtractHTMLTitle(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Title tag", "<title>My Title</title>", "My Title"},
		{"H1 tag", "<h1>My H1</h1>", "My H1"},
		{"Title with attributes", `<title lang="en">My Title</title>`, "My Title"},
		{"H1 with attributes", `<h1 class="main">My H1</h1>`, "My H1"},
		{"No title", "<p>Some content</p>", ""},
		{"Both present", "<title>Title</title><h1>H1</h1>", "Title"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, extractHTMLTitle(tt.input))
		})
	}
}

func TestNaturalLess(t *testing.T) {
	tests := []struct {
		s1, s2   string
		expected bool
	}{
		{"item1.html", "item2.html", true},
		{"item2.html", "item10.html", true},
		{"item10.html", "item1.html", false},
		{"chapter1_part1.html", "chapter1_part2.html", true},
		{"z1.html", "z10.html", true},
	}

	for _, tt := range tests {
		t.Run(tt.s1+" vs "+tt.s2, func(t *testing.T) {
			assert.Equal(t, tt.expected, naturalLess(tt.s1, tt.s2))
		})
	}
}

// createTestEpub creates a dummy epub file for testing
func createTestEpub(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	filePath := filepath.Join(dir, "test.epub")

	// Create a buffer to write our zip archive to.
	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)

	// Create container.xml
	containerContent := `<?xml version="1.0" encoding="UTF-8"?>
<container version="1.0" xmlns="urn:oasis:names:tc:opendocument:xmlns:container">
  <rootfiles>
    <rootfile full-path="OEBPS/content.opf" media-type="application/oebps-package+xml"/>
  </rootfiles>
</container>`
	f, err := zipWriter.Create("META-INF/container.xml")
	assert.NoError(t, err)
	_, err = f.Write([]byte(containerContent))
	assert.NoError(t, err)

	// Create content.opf
	opfContent := `<?xml version="1.0" encoding="UTF-8"?>
<package xmlns="http://www.idpf.org/2007/opf" unique-identifier="pub-id" version="2.0">
  <metadata xmlns:dc="http://purl.org/dc/elements/1.1/">
    <dc:title>Test Book</dc:title>
  </metadata>
  <manifest>
    <item id="chap1" href="chap1.xhtml" media-type="application/xhtml+xml"/>
    <item id="chap2" href="chap2.xhtml" media-type="application/xhtml+xml"/>
  </manifest>
  <spine toc="ncx">
    <itemref idref="chap1"/>
    <itemref idref="chap2"/>
  </spine>
</package>`
	f, err = zipWriter.Create("OEBPS/content.opf")
	assert.NoError(t, err)
	_, err = f.Write([]byte(opfContent))
	assert.NoError(t, err)

	// Create chapter files
	chap1Content := `<?xml version="1.0" encoding="UTF-8"?><html xmlns="http://www.w3.org/1999/xhtml"><head><title>Chapter 1: The Beginning</title></head><body><p>Hello world.</p></body></html>`
	f, err = zipWriter.Create("OEBPS/chap1.xhtml")
	assert.NoError(t, err)
	_, err = f.Write([]byte(chap1Content))
	assert.NoError(t, err)

	chap2Content := `<?xml version="1.0" encoding="UTF-8"?><html xmlns="http://www.w3.org/1999/xhtml"><head><title>Chapter 2: The End</title></head><body><p>Goodbye world.</p></body></html>`
	f, err = zipWriter.Create("OEBPS/chap2.xhtml")
	assert.NoError(t, err)
	_, err = f.Write([]byte(chap2Content))
	assert.NoError(t, err)

	// Close the zip writer
	err = zipWriter.Close()
	assert.NoError(t, err)

	// Write the buffer to file
	err = os.WriteFile(filePath, buf.Bytes(), 0644)
	assert.NoError(t, err)

	return filePath
}

func TestEpubParser_Parse(t *testing.T) {
	parser := NewEpubParser()
	filePath := createTestEpub(t)

	chapters, err := parser.Parse(filePath)
	assert.NoError(t, err)
	assert.Len(t, chapters, 2)
	assert.Equal(t, "Chapter 1: The Beginning", chapters[0].Title)
	assert.Contains(t, chapters[0].Content, "Hello world")
	assert.Equal(t, "Chapter 2: The End", chapters[1].Title)
	assert.Contains(t, chapters[1].Content, "Goodbye world")
}