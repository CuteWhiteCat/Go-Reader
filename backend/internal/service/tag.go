package service

import (
	"time"

	"github.com/google/uuid"
	"github.com/whitecat/go-reader/internal/models"
	"github.com/whitecat/go-reader/internal/repository"
)

// TagService handles business logic for tags
type TagService struct {
	tagRepo *repository.TagRepository
}

// NewTagService creates a new TagService
func NewTagService(tagRepo *repository.TagRepository) *TagService {
	return &TagService{
		tagRepo: tagRepo,
	}
}

// CreateTag creates a new tag
func (s *TagService) CreateTag(req *models.CreateTagRequest) (*models.Tag, error) {
	tag := &models.Tag{
		ID:        uuid.New().String(),
		Name:      req.Name,
		Color:     req.Color,
		CreatedAt: time.Now(),
	}

	if err := s.tagRepo.Create(tag); err != nil {
		return nil, err
	}

	return tag, nil
}

// GetTag retrieves a tag by ID
func (s *TagService) GetTag(id string) (*models.Tag, error) {
	return s.tagRepo.GetByID(id)
}

// GetAllTags retrieves all tags
func (s *TagService) GetAllTags() ([]models.Tag, error) {
	return s.tagRepo.GetAll()
}

// UpdateTag updates a tag
func (s *TagService) UpdateTag(id string, req *models.UpdateTagRequest) (*models.Tag, error) {
	tag, err := s.tagRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		tag.Name = *req.Name
	}
	if req.Color != nil {
		tag.Color = *req.Color
	}

	if err := s.tagRepo.Update(tag); err != nil {
		return nil, err
	}

	return tag, nil
}

// DeleteTag deletes a tag
func (s *TagService) DeleteTag(id string) error {
	return s.tagRepo.Delete(id)
}
