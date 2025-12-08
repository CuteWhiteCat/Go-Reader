package repository

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/whitecat/go-reader/internal/models"
)

// TagRepository handles database operations for tags
type TagRepository struct {
	db *sqlx.DB
}

// NewTagRepository creates a new TagRepository
func NewTagRepository(db *sqlx.DB) *TagRepository {
	return &TagRepository{db: db}
}

// Create creates a new tag in the database
func (r *TagRepository) Create(tag *models.Tag) error {
	query := `
		INSERT INTO tags (id, name, color, created_at)
		VALUES (:id, :name, :color, :created_at)
	`
	_, err := r.db.NamedExec(query, tag)
	if err != nil {
		return fmt.Errorf("failed to create tag: %w", err)
	}
	return nil
}

// GetByID retrieves a tag by its ID
func (r *TagRepository) GetByID(id string) (*models.Tag, error) {
	var tag models.Tag
	query := `SELECT * FROM tags WHERE id = ?`
	err := r.db.Get(&tag, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("tag not found")
		}
		return nil, fmt.Errorf("failed to get tag: %w", err)
	}
	return &tag, nil
}

// GetByName retrieves a tag by its name
func (r *TagRepository) GetByName(name string) (*models.Tag, error) {
	var tag models.Tag
	query := `SELECT * FROM tags WHERE name = ?`
	err := r.db.Get(&tag, query, name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("tag not found")
		}
		return nil, fmt.Errorf("failed to get tag: %w", err)
	}
	return &tag, nil
}

// GetAll retrieves all tags
func (r *TagRepository) GetAll() ([]models.Tag, error) {
	var tags []models.Tag
	query := `SELECT * FROM tags ORDER BY name ASC`
	err := r.db.Select(&tags, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get tags: %w", err)
	}
	return tags, nil
}

// Update updates a tag
func (r *TagRepository) Update(tag *models.Tag) error {
	query := `
		UPDATE tags
		SET name = :name, color = :color
		WHERE id = :id
	`
	result, err := r.db.NamedExec(query, tag)
	if err != nil {
		return fmt.Errorf("failed to update tag: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("tag not found")
	}

	return nil
}

// Delete deletes a tag by ID
func (r *TagRepository) Delete(id string) error {
	query := `DELETE FROM tags WHERE id = ?`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete tag: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("tag not found")
	}

	return nil
}
