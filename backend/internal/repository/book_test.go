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

func setupTestDB(t *testing.T) *BookRepository {
	db := config.NewTestDatabase(t)
	return NewBookRepository(db)
}

func TestBookRepository_Create(t *testing.T) {
	repo := setupTestDB(t)

	book := &models.Book{
		ID:         uuid.NewString(),
		Title:      "Test Book",
		Author:     "Test Author",
		FilePath:   "/test.txt",
		FileFormat: "txt",
		FileSize:   12345,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	err := repo.Create(book)
	assert.NoError(t, err)

	createdBook, err := repo.GetByID(book.ID)
	assert.NoError(t, err)
	assert.NotNil(t, createdBook)
	assert.Equal(t, book.Title, createdBook.Title)
}

func TestBookRepository_GetByID(t *testing.T) {
	repo := setupTestDB(t)

	book := &models.Book{
		ID:         uuid.NewString(),
		Title:      "Test Book",
		Author:     "Test Author",
		FilePath:   "/test.txt",
		FileFormat: "txt",
		FileSize:   12345,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	err := repo.Create(book)
	assert.NoError(t, err)

	// Test found
	foundBook, err := repo.GetByID(book.ID)
	assert.NoError(t, err)
	assert.NotNil(t, foundBook)
	assert.Equal(t, book.ID, foundBook.ID)

	// Test not found
	_, err = repo.GetByID(uuid.NewString())
	assert.Error(t, err)
	assert.Equal(t, "book not found", err.Error())
}

func TestBookRepository_GetAll(t *testing.T) {
	repo := setupTestDB(t)

	book1 := &models.Book{
		ID:         uuid.NewString(),
		Title:      "Book 1",
		FilePath:   "/test1.txt",
		FileFormat: "txt",
		UpdatedAt:  time.Now().Add(-time.Hour),
	}
	book2 := &models.Book{
		ID:         uuid.NewString(),
		Title:      "Book 2",
		FilePath:   "/test2.txt",
		FileFormat: "txt",
		UpdatedAt:  time.Now(),
	}

	err := repo.Create(book1)
	assert.NoError(t, err)
	err = repo.Create(book2)
	assert.NoError(t, err)

	books, err := repo.GetAll()
	assert.NoError(t, err)
	assert.Len(t, books, 2)
	assert.Equal(t, book2.Title, books[0].Title, "Books should be ordered by updated_at DESC")
}

func TestBookRepository_Update(t *testing.T) {
	repo := setupTestDB(t)

	book := &models.Book{
		ID:         uuid.NewString(),
		Title:      "Original Title",
		Author:     "Original Author",
		FilePath:   "/test.txt",
		FileFormat: "txt",
		UpdatedAt:  time.Now(),
	}
	err := repo.Create(book)
	assert.NoError(t, err)

	book.Title = "Updated Title"
	book.Author = "Updated Author"
	book.UpdatedAt = time.Now()

	err = repo.Update(book)
	assert.NoError(t, err)

	updatedBook, err := repo.GetByID(book.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Title", updatedBook.Title)
	assert.Equal(t, "Updated Author", updatedBook.Author)
}

func TestBookRepository_Delete(t *testing.T) {
	repo := setupTestDB(t)

	book := &models.Book{
		ID:         uuid.NewString(),
		Title:      "To Be Deleted",
		FilePath:   "/delete.txt",
		FileFormat: "txt",
	}
	err := repo.Create(book)
	assert.NoError(t, err)

	err = repo.Delete(book.ID)
	assert.NoError(t, err)

	_, err = repo.GetByID(book.ID)
	assert.Error(t, err, "book not found")
}

func TestBookRepository_TagManagement(t *testing.T) {
	bookRepo := setupTestDB(t)
	tagRepo := NewTagRepository(bookRepo.db)

	book := &models.Book{ID: uuid.NewString(), Title: "Tagged Book", FilePath: "/tagged.txt", FileFormat: "txt"}
	tag1 := &models.Tag{ID: uuid.NewString(), Name: "Tag 1"}
	tag2 := &models.Tag{ID: uuid.NewString(), Name: "Tag 2"}

	err := bookRepo.Create(book)
	assert.NoError(t, err)
	err = tagRepo.Create(tag1)
	assert.NoError(t, err)
	err = tagRepo.Create(tag2)
	assert.NoError(t, err)

	// Add tags
	err = bookRepo.AddTag(book.ID, tag1.ID)
	assert.NoError(t, err)
	err = bookRepo.AddTag(book.ID, tag2.ID)
	assert.NoError(t, err)

	// Get tags for book
	tags, err := bookRepo.GetBookTags(book.ID)
	assert.NoError(t, err)
	assert.Len(t, tags, 2)

	// Get books by tag
	books, err := bookRepo.GetBooksByTag(tag1.ID)
	assert.NoError(t, err)
	assert.Len(t, books, 1)
	assert.Equal(t, book.ID, books[0].ID)

	// Remove tag
	err = bookRepo.RemoveTag(book.ID, tag1.ID)
	assert.NoError(t, err)
	tags, err = bookRepo.GetBookTags(book.ID)
	assert.NoError(t, err)
	assert.Len(t, tags, 1)
	assert.Equal(t, tag2.Name, tags[0].Name)
}
