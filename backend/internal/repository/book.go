package repository

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/whitecat/go-reader/internal/models"
)

// BookRepository handles database operations for books
type BookRepository struct {
	db *sqlx.DB
}

// NewBookRepository creates a new BookRepository
func NewBookRepository(db *sqlx.DB) *BookRepository {
	return &BookRepository{db: db}
}

// Create creates a new book in the database
func (r *BookRepository) Create(book *models.Book) error {
	query := `
		INSERT INTO books (id, title, author, description, cover_path, file_path, file_format, file_size, created_at, updated_at)
		VALUES (:id, :title, :author, :description, :cover_path, :file_path, :file_format, :file_size, :created_at, :updated_at)
	`
	_, err := r.db.NamedExec(query, book)
	if err != nil {
		return fmt.Errorf("failed to create book: %w", err)
	}
	return nil
}

// GetByID retrieves a book by its ID
func (r *BookRepository) GetByID(id string) (*models.Book, error) {
	var book models.Book
	query := `SELECT * FROM books WHERE id = ?`
	err := r.db.Get(&book, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("book not found")
		}
		return nil, fmt.Errorf("failed to get book: %w", err)
	}
	return &book, nil
}

// GetAll retrieves all books
func (r *BookRepository) GetAll() ([]models.Book, error) {
	var books []models.Book
	query := `SELECT * FROM books ORDER BY updated_at DESC`
	err := r.db.Select(&books, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get books: %w", err)
	}
	return books, nil
}

// Update updates a book
func (r *BookRepository) Update(book *models.Book) error {
	query := `
		UPDATE books
		SET title = :title, author = :author, description = :description,
		    cover_path = :cover_path, updated_at = :updated_at
		WHERE id = :id
	`
	result, err := r.db.NamedExec(query, book)
	if err != nil {
		return fmt.Errorf("failed to update book: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("book not found")
	}

	return nil
}

// Delete deletes a book by ID
func (r *BookRepository) Delete(id string) error {
	query := `DELETE FROM books WHERE id = ?`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete book: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("book not found")
	}

	return nil
}

// GetBooksByTag retrieves books by tag ID
func (r *BookRepository) GetBooksByTag(tagID string) ([]models.Book, error) {
	var books []models.Book
	query := `
		SELECT b.* FROM books b
		INNER JOIN book_tags bt ON b.id = bt.book_id
		WHERE bt.tag_id = ?
		ORDER BY b.updated_at DESC
	`
	err := r.db.Select(&books, query, tagID)
	if err != nil {
		return nil, fmt.Errorf("failed to get books by tag: %w", err)
	}
	return books, nil
}

// AddTag adds a tag to a book
func (r *BookRepository) AddTag(bookID, tagID string) error {
	query := `INSERT INTO book_tags (book_id, tag_id) VALUES (?, ?)`
	_, err := r.db.Exec(query, bookID, tagID)
	if err != nil {
		return fmt.Errorf("failed to add tag to book: %w", err)
	}
	return nil
}

// RemoveTag removes a tag from a book
func (r *BookRepository) RemoveTag(bookID, tagID string) error {
	query := `DELETE FROM book_tags WHERE book_id = ? AND tag_id = ?`
	_, err := r.db.Exec(query, bookID, tagID)
	if err != nil {
		return fmt.Errorf("failed to remove tag from book: %w", err)
	}
	return nil
}

// GetBookTags retrieves all tags for a book
func (r *BookRepository) GetBookTags(bookID string) ([]models.Tag, error) {
	var tags []models.Tag
	query := `
		SELECT t.* FROM tags t
		INNER JOIN book_tags bt ON t.id = bt.tag_id
		WHERE bt.book_id = ?
		ORDER BY t.name
	`
	err := r.db.Select(&tags, query, bookID)
	if err != nil {
		return nil, fmt.Errorf("failed to get book tags: %w", err)
	}
	return tags, nil
}
