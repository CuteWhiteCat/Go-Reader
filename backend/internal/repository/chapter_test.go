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

func setupChapterTestDB(t *testing.T) (*ChapterRepository, *BookRepository) {
	db := config.NewTestDatabase(t)
	return NewChapterRepository(db), NewBookRepository(db)
}

func createTestBook(t *testing.T, bookRepo *BookRepository) *models.Book {
	book := &models.Book{
		ID:         uuid.NewString(),
		Title:      "Test Book for Chapters",
		FilePath:   "/test.txt",
		FileFormat: "txt",
	}
	err := bookRepo.Create(book)
	assert.NoError(t, err)
	return book
}

func TestChapterRepository_Create(t *testing.T) {
	chapterRepo, bookRepo := setupChapterTestDB(t)
	book := createTestBook(t, bookRepo)

	chapter := &models.Chapter{
		ID:            uuid.NewString(),
		BookID:        book.ID,
		ChapterNumber: 1,
		Title:         "Chapter 1",
		Content:       "This is the content.",
		WordCount:     4,
		CreatedAt:     time.Now(),
	}

	err := chapterRepo.Create(chapter)
	assert.NoError(t, err)

	createdChapter, err := chapterRepo.GetByID(chapter.ID)
	assert.NoError(t, err)
	assert.NotNil(t, createdChapter)
	assert.Equal(t, chapter.Title, createdChapter.Title)
}

func TestChapterRepository_GetByID(t *testing.T) {
	chapterRepo, bookRepo := setupChapterTestDB(t)
	book := createTestBook(t, bookRepo)

	chapter := &models.Chapter{
		ID:            uuid.NewString(),
		BookID:        book.ID,
		ChapterNumber: 1,
		Title:         "Chapter 1",
	}
	err := chapterRepo.Create(chapter)
	assert.NoError(t, err)

	foundChapter, err := chapterRepo.GetByID(chapter.ID)
	assert.NoError(t, err)
	assert.NotNil(t, foundChapter)
	assert.Equal(t, chapter.ID, foundChapter.ID)

	_, err = chapterRepo.GetByID(uuid.NewString())
	assert.Error(t, err)
}

func TestChapterRepository_GetByBookID(t *testing.T) {
	chapterRepo, bookRepo := setupChapterTestDB(t)
	book := createTestBook(t, bookRepo)

	ch1 := &models.Chapter{ID: uuid.NewString(), BookID: book.ID, ChapterNumber: 1, Title: "Chapter 1"}
	ch2 := &models.Chapter{ID: uuid.NewString(), BookID: book.ID, ChapterNumber: 2, Title: "Chapter 2"}

	err := chapterRepo.Create(ch1)
	assert.NoError(t, err)
	err = chapterRepo.Create(ch2)
	assert.NoError(t, err)

	chapters, err := chapterRepo.GetByBookID(book.ID)
	assert.NoError(t, err)
	assert.Len(t, chapters, 2)
	assert.Equal(t, ch1.Title, chapters[0].Title)
}

func TestChapterRepository_GetByNumber(t *testing.T) {
	chapterRepo, bookRepo := setupChapterTestDB(t)
	book := createTestBook(t, bookRepo)

	ch1 := &models.Chapter{ID: uuid.NewString(), BookID: book.ID, ChapterNumber: 1, Title: "Chapter 1"}
	err := chapterRepo.Create(ch1)
	assert.NoError(t, err)

	foundChapter, err := chapterRepo.GetByNumber(book.ID, 1)
	assert.NoError(t, err)
	assert.NotNil(t, foundChapter)
	assert.Equal(t, ch1.ID, foundChapter.ID)

	_, err = chapterRepo.GetByNumber(book.ID, 99)
	assert.Error(t, err)
}

func TestChapterRepository_DeleteByBookID(t *testing.T) {
	chapterRepo, bookRepo := setupChapterTestDB(t)
	book := createTestBook(t, bookRepo)

	ch1 := &models.Chapter{ID: uuid.NewString(), BookID: book.ID, ChapterNumber: 1, Title: "Chapter 1"}
	err := chapterRepo.Create(ch1)
	assert.NoError(t, err)

	err = chapterRepo.DeleteByBookID(book.ID)
	assert.NoError(t, err)

	chapters, err := chapterRepo.GetByBookID(book.ID)
	assert.NoError(t, err)
	assert.Len(t, chapters, 0)
}

func TestChapterRepository_BatchCreate(t *testing.T) {
	chapterRepo, bookRepo := setupChapterTestDB(t)
	book := createTestBook(t, bookRepo)

	chapters := []models.Chapter{
		{ID: uuid.NewString(), BookID: book.ID, ChapterNumber: 1, Title: "Batch 1"},
		{ID: uuid.NewString(), BookID: book.ID, ChapterNumber: 2, Title: "Batch 2"},
		{ID: uuid.NewString(), BookID: book.ID, ChapterNumber: 3, Title: "Batch 3"},
	}

	err := chapterRepo.BatchCreate(chapters)
	assert.NoError(t, err)

	retrievedChapters, err := chapterRepo.GetByBookID(book.ID)
	assert.NoError(t, err)
	assert.Len(t, retrievedChapters, 3)
}