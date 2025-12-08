package models

import "time"

// Book represents a book in the library
type Book struct {
	ID          string    `json:"id" db:"id"`
	Title       string    `json:"title" db:"title" validate:"required"`
	Author      string    `json:"author" db:"author"`
	Description string    `json:"description" db:"description"`
	CoverPath   string    `json:"cover_path" db:"cover_path"`
	FilePath    string    `json:"file_path" db:"file_path" validate:"required"`
	FileFormat  string    `json:"file_format" db:"file_format" validate:"required,oneof=txt md epub web"`
	FileSize    int64     `json:"file_size" db:"file_size"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
	Tags        []Tag     `json:"tags,omitempty" db:"-"`
}

// CreateBookRequest represents the request to create a new book
type CreateBookRequest struct {
	Title       string   `json:"title" validate:"required"`
	Author      string   `json:"author"`
	Description string   `json:"description"`
	FilePath    string   `json:"file_path" validate:"required"`
	FileFormat  string   `json:"file_format" validate:"required,oneof=txt md epub"`
	TagIDs      []string `json:"tag_ids"`
}

// CreateRemoteBookRequest represents creating a book from scraped chapters (no local file)
type CreateRemoteBookRequest struct {
	Title       string `json:"title"`
	Author      string `json:"author"`
	Description string `json:"description"`
	SourceURL   string `json:"source_url"`
}

// UpdateBookRequest represents the request to update a book
type UpdateBookRequest struct {
	Title       *string  `json:"title"`
	Author      *string  `json:"author"`
	Description *string  `json:"description"`
	CoverPath   *string  `json:"cover_path"`
	TagIDs      []string `json:"tag_ids"`
}
