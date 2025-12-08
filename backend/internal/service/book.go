package service

import (
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/whitecat/go-reader/internal/models"
	"github.com/whitecat/go-reader/internal/parser"
	"github.com/whitecat/go-reader/internal/repository"
)

// BookService handles business logic for books
type BookService struct {
	bookRepo    *repository.BookRepository
	chapterRepo *repository.ChapterRepository
	tagRepo     *repository.TagRepository
}

// NewBookService creates a new BookService
func NewBookService(
	bookRepo *repository.BookRepository,
	chapterRepo *repository.ChapterRepository,
	tagRepo *repository.TagRepository,
) *BookService {
	return &BookService{
		bookRepo:    bookRepo,
		chapterRepo: chapterRepo,
		tagRepo:     tagRepo,
	}
}

// CreateBook creates a new book and parses its content
func (s *BookService) CreateBook(req *models.CreateBookRequest) (*models.Book, error) {
	// Validate file exists
	fileInfo, err := os.Stat(req.FilePath)
	if err != nil {
		return nil, fmt.Errorf("file not found: %w", err)
	}

	// Create book entity
	book := &models.Book{
		ID:          uuid.New().String(),
		Title:       req.Title,
		Author:      req.Author,
		Description: req.Description,
		FilePath:    req.FilePath,
		FileFormat:  req.FileFormat,
		FileSize:    fileInfo.Size(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Save book to database
	if err := s.bookRepo.Create(book); err != nil {
		return nil, fmt.Errorf("failed to create book: %w", err)
	}

	// Parse book content into chapters
	p, err := parser.GetParser(req.FileFormat)
	if err != nil {
		// If parser not available, still create the book
		return book, nil
	}

	chapters, err := p.Parse(req.FilePath)
	if err != nil {
		// If parsing fails, still return the book
		return book, nil
	}

	// Save chapters
	for i := range chapters {
		chapters[i].BookID = book.ID
		chapters[i].CreatedAt = time.Now()
	}

	if err := s.chapterRepo.BatchCreate(chapters); err != nil {
		return nil, fmt.Errorf("failed to create chapters: %w", err)
	}

	// Add tags if provided
	if len(req.TagIDs) > 0 {
		for _, tagID := range req.TagIDs {
			s.bookRepo.AddTag(book.ID, tagID)
		}
	}

	return book, nil
}

// CreateRemoteBook creates a book using scraped chapters (no local file needed).
func (s *BookService) CreateRemoteBook(req *models.CreateRemoteBookRequest, chapters []models.Chapter) (*models.Book, error) {
	book := &models.Book{
		ID:          uuid.New().String(),
		Title:       req.Title,
		Author:      req.Author,
		Description: req.Description,
		FilePath:    req.SourceURL,
		FileFormat:  "web",
		FileSize:    0,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.bookRepo.Create(book); err != nil {
		return nil, fmt.Errorf("failed to create book: %w", err)
	}

	now := time.Now()
	for i := range chapters {
		chapters[i].BookID = book.ID
		if chapters[i].CreatedAt.IsZero() {
			chapters[i].CreatedAt = now
		}
	}

	if len(chapters) > 0 {
		if err := s.chapterRepo.BatchCreate(chapters); err != nil {
			return nil, fmt.Errorf("failed to save chapters: %w", err)
		}
	}

	return book, nil
}

// GetBook retrieves a book by ID with its tags
func (s *BookService) GetBook(id string) (*models.Book, error) {
	book, err := s.bookRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Load tags
	tags, err := s.bookRepo.GetBookTags(book.ID)
	if err == nil {
		book.Tags = tags
	}

	return book, nil
}

// GetAllBooks retrieves all books
func (s *BookService) GetAllBooks() ([]models.Book, error) {
	books, err := s.bookRepo.GetAll()
	if err != nil {
		return nil, err
	}

	// Load tags for each book
	for i := range books {
		tags, err := s.bookRepo.GetBookTags(books[i].ID)
		if err == nil {
			books[i].Tags = tags
		}
	}

	return books, nil
}

// UpdateBook updates a book's metadata
func (s *BookService) UpdateBook(id string, req *models.UpdateBookRequest) (*models.Book, error) {
	book, err := s.bookRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.Title != nil {
		book.Title = *req.Title
	}
	if req.Author != nil {
		book.Author = *req.Author
	}
	if req.Description != nil {
		book.Description = *req.Description
	}
	if req.CoverPath != nil {
		book.CoverPath = *req.CoverPath
	}

	book.UpdatedAt = time.Now()

	if err := s.bookRepo.Update(book); err != nil {
		return nil, err
	}

	// Update tags if provided
	if req.TagIDs != nil {
		// Remove existing tags
		existingTags, _ := s.bookRepo.GetBookTags(book.ID)
		for _, tag := range existingTags {
			s.bookRepo.RemoveTag(book.ID, tag.ID)
		}

		// Add new tags
		for _, tagID := range req.TagIDs {
			s.bookRepo.AddTag(book.ID, tagID)
		}
	}

	return s.GetBook(id)
}

// DeleteBook deletes a book and its associated data
func (s *BookService) DeleteBook(id string) error {
	book, err := s.bookRepo.GetByID(id)
	if err != nil {
		return err
	}

	// Delete chapters
	if err := s.chapterRepo.DeleteByBookID(id); err != nil {
		return fmt.Errorf("failed to delete chapters: %w", err)
	}

	// Delete book
	if err := s.bookRepo.Delete(id); err != nil {
		return err
	}

	// Optionally delete the book file (commented out for safety)
	// os.Remove(book.FilePath)

	_ = book // Use book to avoid unused variable warning

	return nil
}

// GetBookContent retrieves the full content of a book
func (s *BookService) GetBookContent(id string) ([]models.Chapter, error) {
	book, err := s.bookRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Try to get chapters from database first
	chapters, err := s.chapterRepo.GetByBookID(book.ID)
	if err != nil || len(chapters) == 0 {
		// If no chapters in database, parse the file
		p, err := parser.GetParser(book.FileFormat)
		if err != nil {
			return nil, fmt.Errorf("parser not available: %w", err)
		}

		fullChapters, err := p.Parse(book.FilePath)
		if err != nil {
			return nil, fmt.Errorf("failed to parse book: %w", err)
		}

		// Save chapters for future use
		for i := range fullChapters {
			fullChapters[i].BookID = book.ID
			fullChapters[i].CreatedAt = time.Now()
		}
		s.chapterRepo.BatchCreate(fullChapters)

		return fullChapters, nil
	}

	// Convert summaries to full chapters
	var fullChapters []models.Chapter
	for _, summary := range chapters {
		chapter, err := s.chapterRepo.GetByID(summary.ID)
		if err == nil {
			fullChapters = append(fullChapters, *chapter)
		}
	}

	return fullChapters, nil
}

// GetBooksByTag retrieves books by tag
func (s *BookService) GetBooksByTag(tagID string) ([]models.Book, error) {
	return s.bookRepo.GetBooksByTag(tagID)
}

// GetBookChapters retrieves chapters for a book (summaries only)
func (s *BookService) GetBookChapters(id string) ([]models.ChapterSummary, error) {
	return s.chapterRepo.GetByBookID(id)
}

// GetChapter retrieves a specific chapter with full content
func (s *BookService) GetChapter(bookID string, chapterNumber int) (*models.Chapter, error) {
	return s.chapterRepo.GetByNumber(bookID, chapterNumber)
}
