package models

import "time"

// Chapter represents a chapter in a book
type Chapter struct {
	ID                  string    `json:"id" db:"id"`
	BookID              string    `json:"book_id" db:"book_id"`
	ChapterNumber       int       `json:"chapter_number" db:"chapter_number"`
	VolumeNumber        int       `json:"volume_number" db:"volume_number"`
	VolumeChapterNumber int       `json:"volume_chapter_number" db:"volume_chapter_number"`
	Title               string    `json:"title" db:"title"`
	Content             string    `json:"content,omitempty" db:"content"`
	WordCount           int       `json:"word_count" db:"word_count"`
	CreatedAt           time.Time `json:"created_at" db:"created_at"`
}

// ChapterSummary represents a chapter without content (for listing)
type ChapterSummary struct {
	ID                  string    `json:"id" db:"id"`
	BookID              string    `json:"book_id" db:"book_id"`
	ChapterNumber       int       `json:"chapter_number" db:"chapter_number"`
	VolumeNumber        int       `json:"volume_number" db:"volume_number"`
	VolumeChapterNumber int       `json:"volume_chapter_number" db:"volume_chapter_number"`
	Title               string    `json:"title" db:"title"`
	WordCount           int       `json:"word_count" db:"word_count"`
	CreatedAt           time.Time `json:"created_at" db:"created_at"`
}
