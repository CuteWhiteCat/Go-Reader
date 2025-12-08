package repository

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/whitecat/go-reader/internal/models"
)

// ProgressRepository handles database operations for reading progress
type ProgressRepository struct {
	db *sqlx.DB
}

// NewProgressRepository creates a new ProgressRepository
func NewProgressRepository(db *sqlx.DB) *ProgressRepository {
	return &ProgressRepository{db: db}
}

// GetByBookID retrieves reading progress for a book
func (r *ProgressRepository) GetByBookID(bookID string) (*models.ReadingProgress, error) {
	var progress models.ReadingProgress
	query := `SELECT * FROM reading_progress WHERE book_id = ?`
	err := r.db.Get(&progress, query, bookID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("progress not found")
		}
		return nil, fmt.Errorf("failed to get progress: %w", err)
	}
	return &progress, nil
}

// Upsert creates or updates reading progress for a book
func (r *ProgressRepository) Upsert(progress *models.ReadingProgress) error {
	query := `
		INSERT INTO reading_progress (book_id, current_chapter, current_position, progress_percentage, last_read_at)
		VALUES (:book_id, :current_chapter, :current_position, :progress_percentage, :last_read_at)
		ON CONFLICT(book_id) DO UPDATE SET
			current_chapter = :current_chapter,
			current_position = :current_position,
			progress_percentage = :progress_percentage,
			last_read_at = :last_read_at
	`
	_, err := r.db.NamedExec(query, progress)
	if err != nil {
		return fmt.Errorf("failed to upsert progress: %w", err)
	}
	return nil
}

// DeleteByBookID deletes reading progress for a book
func (r *ProgressRepository) DeleteByBookID(bookID string) error {
	query := `DELETE FROM reading_progress WHERE book_id = ?`
	_, err := r.db.Exec(query, bookID)
	if err != nil {
		return fmt.Errorf("failed to delete progress: %w", err)
	}
	return nil
}

// BookmarkRepository handles database operations for bookmarks
type BookmarkRepository struct {
	db *sqlx.DB
}

// NewBookmarkRepository creates a new BookmarkRepository
func NewBookmarkRepository(db *sqlx.DB) *BookmarkRepository {
	return &BookmarkRepository{db: db}
}

// Create creates a new bookmark
func (r *BookmarkRepository) Create(bookmark *models.Bookmark) error {
	query := `
		INSERT INTO bookmarks (id, book_id, chapter_id, position, note, created_at)
		VALUES (:id, :book_id, :chapter_id, :position, :note, :created_at)
	`
	_, err := r.db.NamedExec(query, bookmark)
	if err != nil {
		return fmt.Errorf("failed to create bookmark: %w", err)
	}
	return nil
}

// GetByID retrieves a bookmark by its ID
func (r *BookmarkRepository) GetByID(id string) (*models.Bookmark, error) {
	var bookmark models.Bookmark
	query := `SELECT * FROM bookmarks WHERE id = ?`
	err := r.db.Get(&bookmark, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("bookmark not found")
		}
		return nil, fmt.Errorf("failed to get bookmark: %w", err)
	}
	return &bookmark, nil
}

// GetByBookID retrieves all bookmarks for a book
func (r *BookmarkRepository) GetByBookID(bookID string) ([]models.Bookmark, error) {
	var bookmarks []models.Bookmark
	query := `SELECT * FROM bookmarks WHERE book_id = ? ORDER BY created_at DESC`
	err := r.db.Select(&bookmarks, query, bookID)
	if err != nil {
		return nil, fmt.Errorf("failed to get bookmarks: %w", err)
	}
	return bookmarks, nil
}

// Delete deletes a bookmark by ID
func (r *BookmarkRepository) Delete(id string) error {
	query := `DELETE FROM bookmarks WHERE id = ?`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete bookmark: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("bookmark not found")
	}

	return nil
}
