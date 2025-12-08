package parser

import (
	"testing"
)

func TestIsChapterTitle(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		// Valid Chinese chapter titles
		{"Chinese chapter 1", "第一章", true},
		{"Chinese chapter 2", "第二章", true},
		{"Chinese chapter with spaces", "  第三章  ", true},
		{"Chinese chapter with number", "第123章", true},
		{"Chinese chapter with text", "第一章 开始", true},

		// Invalid - contains 第 but not 章
		{"Day (not chapter)", "第二天", false},
		{"First time", "第一次", false},
		{"Just 第", "第", false},

		// Valid English chapter titles
		{"English chapter", "Chapter 1", true},
		{"English chapter lowercase", "chapter 2", true},
		{"English chapter with colon", "Chapter: Introduction", true},
		{"Ch. abbreviation", "Ch. 5", true},
		{"Ch space", "Ch 10", true},

		// Invalid cases
		{"Empty string", "", false},
		{"Just spaces", "   ", false},
		{"Random text", "This is some text", false},
		{"Contains chapter in middle", "The chapter is here", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isChapterTitle(tt.input)
			if result != tt.expected {
				t.Errorf("isChapterTitle(%q) = %v, expected %v", tt.input, result, tt.expected)
			}
		})
	}
}
