package models

import "time"

// Tag represents a tag for organizing books
type Tag struct {
	ID        string    `json:"id" db:"id"`
	Name      string    `json:"name" db:"name" validate:"required"`
	Color     string    `json:"color" db:"color"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// CreateTagRequest represents the request to create a new tag
type CreateTagRequest struct {
	Name  string `json:"name" validate:"required"`
	Color string `json:"color"`
}

// UpdateTagRequest represents the request to update a tag
type UpdateTagRequest struct {
	Name  *string `json:"name"`
	Color *string `json:"color"`
}
