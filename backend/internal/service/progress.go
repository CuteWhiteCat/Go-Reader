package service

import (
	"time"

	"github.com/google/uuid"
	"github.com/whitecat/go-reader/internal/models"
	"github.com/whitecat/go-reader/internal/repository"
)

// ProgressService handles business logic for reading progress
type ProgressService struct {
	progressRepo  *repository.ProgressRepository
	bookmarkRepo  *repository.BookmarkRepository
}

// NewProgressService creates a new ProgressService
func NewProgressService(
	progressRepo *repository.ProgressRepository,
	bookmarkRepo *repository.BookmarkRepository,
) *ProgressService {
	return &ProgressService{
		progressRepo:  progressRepo,
		bookmarkRepo:  bookmarkRepo,
	}
}

// GetProgress retrieves reading progress for a book
func (s *ProgressService) GetProgress(bookID string) (*models.ReadingProgress, error) {
	progress, err := s.progressRepo.GetByBookID(bookID)
	if err != nil {
		// Return default progress if not found
		return &models.ReadingProgress{
			BookID:             bookID,
			CurrentChapter:     0,
			CurrentPosition:    0,
			ProgressPercentage: 0.0,
			LastReadAt:         time.Now(),
		}, nil
	}
	return progress, nil
}

// UpdateProgress updates reading progress for a book
func (s *ProgressService) UpdateProgress(bookID string, req *models.UpdateProgressRequest) (*models.ReadingProgress, error) {
	progress := &models.ReadingProgress{
		BookID:             bookID,
		CurrentChapter:     req.CurrentChapter,
		CurrentPosition:    req.CurrentPosition,
		ProgressPercentage: req.ProgressPercentage,
		LastReadAt:         time.Now(),
	}

	if err := s.progressRepo.Upsert(progress); err != nil {
		return nil, err
	}

	return progress, nil
}

// CreateBookmark creates a new bookmark
func (s *ProgressService) CreateBookmark(req *models.CreateBookmarkRequest) (*models.Bookmark, error) {
	bookmark := &models.Bookmark{
		ID:        uuid.New().String(),
		BookID:    req.BookID,
		ChapterID: req.ChapterID,
		Position:  req.Position,
		Note:      req.Note,
		CreatedAt: time.Now(),
	}

	if err := s.bookmarkRepo.Create(bookmark); err != nil {
		return nil, err
	}

	return bookmark, nil
}

// GetBookmarks retrieves all bookmarks for a book
func (s *ProgressService) GetBookmarks(bookID string) ([]models.Bookmark, error) {
	return s.bookmarkRepo.GetByBookID(bookID)
}

// DeleteBookmark deletes a bookmark
func (s *ProgressService) DeleteBookmark(id string) error {
	return s.bookmarkRepo.Delete(id)
}
