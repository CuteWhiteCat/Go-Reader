package repository

import (
	"testing"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/whitecat/go-reader/internal/config"
	"github.com/whitecat/go-reader/internal/models"
)

func setupProgressTestDB(t *testing.T) (*ProgressRepository, *BookmarkRepository, *BookRepository, *ChapterRepository) {
	db := config.NewTestDatabase(t)
	return NewProgressRepository(db), NewBookmarkRepository(db), NewBookRepository(db), NewChapterRepository(db)
}

func TestProgressRepository_UpsertAndGet(t *testing.T) {
	progressRepo, _, bookRepo, _ := setupProgressTestDB(t)
	book := createTestBook(t, bookRepo)

	// First time progress
	progress := &models.ReadingProgress{
		BookID:             book.ID,
		CurrentChapter:     1,
		CurrentPosition:    100,
		ProgressPercentage: 50.5,
		LastReadAt:         time.Now(),
	}

	err := progressRepo.Upsert(progress)
	assert.NoError(t, err)

	retrievedProgress, err := progressRepo.GetByBookID(book.ID)
	assert.NoError(t, err)
	assert.NotNil(t, retrievedProgress)
	assert.Equal(t, progress.CurrentChapter, retrievedProgress.CurrentChapter)

	// Update progress
	progress.CurrentChapter = 2
	err = progressRepo.Upsert(progress)
	assert.NoError(t, err)

	updatedProgress, err := progressRepo.GetByBookID(book.ID)
	assert.NoError(t, err)
	assert.Equal(t, 2, updatedProgress.CurrentChapter)
}

func TestProgressRepository_DeleteByBookID(t *testing.T) {
	progressRepo, _, bookRepo, _ := setupProgressTestDB(t)
	book := createTestBook(t, bookRepo)

	progress := &models.ReadingProgress{BookID: book.ID, CurrentChapter: 1}
	err := progressRepo.Upsert(progress)
	assert.NoError(t, err)

	err = progressRepo.DeleteByBookID(book.ID)
	assert.NoError(t, err)

	_, err = progressRepo.GetByBookID(book.ID)
	assert.Error(t, err, "progress not found")
}

func TestBookmarkRepository_CreateAndGet(t *testing.T) {
	_, bookmarkRepo, bookRepo, chapterRepo := setupProgressTestDB(t)
	book := createTestBook(t, bookRepo)
	chapter := &models.Chapter{ID: uuid.NewString(), BookID: book.ID, ChapterNumber: 1, Title: "Chapter 1"}
	err := chapterRepo.Create(chapter)
	assert.NoError(t, err)

	bookmark := &models.Bookmark{
		ID:        uuid.NewString(),
		BookID:    book.ID,
		ChapterID: chapter.ID,
		Position:  123,
		Note:      "A test bookmark",
		CreatedAt: time.Now(),
	}

	err = bookmarkRepo.Create(bookmark)
	assert.NoError(t, err)

	// Get by ID
	found, err := bookmarkRepo.GetByID(bookmark.ID)
	assert.NoError(t, err)
	assert.Equal(t, bookmark.Note, found.Note)

	// Get by Book ID
	bookmarks, err := bookmarkRepo.GetByBookID(book.ID)
	assert.NoError(t, err)
	assert.Len(t, bookmarks, 1)
	assert.Equal(t, bookmark.ID, bookmarks[0].ID)
}

func TestBookmarkRepository_Delete(t *testing.T) {
	_, bookmarkRepo, bookRepo, chapterRepo := setupProgressTestDB(t)
	book := createTestBook(t, bookRepo)
	chapter := &models.Chapter{ID: uuid.NewString(), BookID: book.ID, ChapterNumber: 1, Title: "Chapter 1"}
	err := chapterRepo.Create(chapter)
	assert.NoError(t, err)
	
	bookmark := &models.Bookmark{ID: uuid.NewString(), BookID: book.ID, ChapterID: chapter.ID}
	err = bookmarkRepo.Create(bookmark)
	assert.NoError(t, err)

	err = bookmarkRepo.Delete(bookmark.ID)
	assert.NoError(t, err)

	_, err = bookmarkRepo.GetByID(bookmark.ID)
	assert.Error(t, err)
}