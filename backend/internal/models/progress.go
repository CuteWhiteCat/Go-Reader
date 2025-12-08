package models

import "time"

// ReadingProgress represents a user's reading progress for a book
type ReadingProgress struct {
	BookID             string    `json:"book_id" db:"book_id"`
	CurrentChapter     int       `json:"current_chapter" db:"current_chapter"`
	CurrentPosition    int       `json:"current_position" db:"current_position"`
	ProgressPercentage float64   `json:"progress_percentage" db:"progress_percentage"`
	LastReadAt         time.Time `json:"last_read_at" db:"last_read_at"`
}

// UpdateProgressRequest represents the request to update reading progress
type UpdateProgressRequest struct {
	CurrentChapter  int     `json:"current_chapter" validate:"min=0"`
	CurrentPosition int     `json:"current_position" validate:"min=0"`
	ProgressPercentage float64 `json:"progress_percentage" validate:"min=0,max=100"`
}

// Bookmark represents a bookmark in a book
type Bookmark struct {
	ID         string    `json:"id" db:"id"`
	BookID     string    `json:"book_id" db:"book_id"`
	ChapterID  string    `json:"chapter_id" db:"chapter_id"`
	Position   int       `json:"position" db:"position"`
	Note       string    `json:"note" db:"note"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}

// CreateBookmarkRequest represents the request to create a bookmark
type CreateBookmarkRequest struct {
	BookID    string `json:"book_id" validate:"required"`
	ChapterID string `json:"chapter_id" validate:"required"`
	Position  int    `json:"position" validate:"min=0"`
	Note      string `json:"note"`
}
