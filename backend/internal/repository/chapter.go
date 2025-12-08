package repository

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/whitecat/go-reader/internal/models"
)

// ChapterRepository handles database operations for chapters
type ChapterRepository struct {
	db *sqlx.DB
}

// NewChapterRepository creates a new ChapterRepository
func NewChapterRepository(db *sqlx.DB) *ChapterRepository {
	return &ChapterRepository{db: db}
}

// Create creates a new chapter in the database
func (r *ChapterRepository) Create(chapter *models.Chapter) error {
	query := `
		INSERT INTO chapters (id, book_id, chapter_number, volume_number, volume_chapter_number, title, content, word_count, created_at)
		VALUES (:id, :book_id, :chapter_number, :volume_number, :volume_chapter_number, :title, :content, :word_count, :created_at)
	`
	_, err := r.db.NamedExec(query, chapter)
	if err != nil {
		return fmt.Errorf("failed to create chapter: %w", err)
	}
	return nil
}

// GetByID retrieves a chapter by its ID
func (r *ChapterRepository) GetByID(id string) (*models.Chapter, error) {
	var chapter models.Chapter
	query := `SELECT * FROM chapters WHERE id = ?`
	err := r.db.Get(&chapter, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("chapter not found")
		}
		return nil, fmt.Errorf("failed to get chapter: %w", err)
	}
	return &chapter, nil
}

// GetByBookID retrieves all chapters for a book
func (r *ChapterRepository) GetByBookID(bookID string) ([]models.ChapterSummary, error) {
	var chapters []models.ChapterSummary
	query := `
		SELECT id, book_id, chapter_number, volume_number, volume_chapter_number, title, word_count, created_at
		FROM chapters
		WHERE book_id = ?
		ORDER BY chapter_number ASC
	`
	err := r.db.Select(&chapters, query, bookID)
	if err != nil {
		return nil, fmt.Errorf("failed to get chapters: %w", err)
	}
	return chapters, nil
}

// GetByNumber retrieves a chapter by book ID and chapter number
func (r *ChapterRepository) GetByNumber(bookID string, chapterNumber int) (*models.Chapter, error) {
	var chapter models.Chapter
	query := `SELECT * FROM chapters WHERE book_id = ? AND chapter_number = ?`
	err := r.db.Get(&chapter, query, bookID, chapterNumber)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("chapter not found")
		}
		return nil, fmt.Errorf("failed to get chapter: %w", err)
	}
	return &chapter, nil
}

// DeleteByBookID deletes all chapters for a book
func (r *ChapterRepository) DeleteByBookID(bookID string) error {
	query := `DELETE FROM chapters WHERE book_id = ?`
	_, err := r.db.Exec(query, bookID)
	if err != nil {
		return fmt.Errorf("failed to delete chapters: %w", err)
	}
	return nil
}

// BatchCreate creates multiple chapters in a single transaction
func (r *ChapterRepository) BatchCreate(chapters []models.Chapter) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `
		INSERT INTO chapters (id, book_id, chapter_number, volume_number, volume_chapter_number, title, content, word_count, created_at)
		VALUES (:id, :book_id, :chapter_number, :volume_number, :volume_chapter_number, :title, :content, :word_count, :created_at)
	`

	for _, chapter := range chapters {
		if _, err := tx.NamedExec(query, &chapter); err != nil {
			return fmt.Errorf("failed to create chapter: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
