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

func setupTagTestDB(t *testing.T) *TagRepository {
	db := config.NewTestDatabase(t)
	return NewTagRepository(db)
}

func TestTagRepository_Create(t *testing.T) {
	repo := setupTagTestDB(t)

	tag := &models.Tag{
		ID:        uuid.NewString(),
		Name:      "Test Tag",
		Color:     "#FF0000",
		CreatedAt: time.Now(),
	}

	err := repo.Create(tag)
	assert.NoError(t, err)

	createdTag, err := repo.GetByID(tag.ID)
	assert.NoError(t, err)
	assert.NotNil(t, createdTag)
	assert.Equal(t, tag.Name, createdTag.Name)
}

func TestTagRepository_GetByID(t *testing.T) {
	repo := setupTagTestDB(t)

	tag := &models.Tag{
		ID:   uuid.NewString(),
		Name: "Test Tag",
	}
	err := repo.Create(tag)
	assert.NoError(t, err)

	foundTag, err := repo.GetByID(tag.ID)
	assert.NoError(t, err)
	assert.NotNil(t, foundTag)
	assert.Equal(t, tag.ID, foundTag.ID)

	_, err = repo.GetByID(uuid.NewString())
	assert.Error(t, err)
}

func TestTagRepository_GetByName(t *testing.T) {
	repo := setupTagTestDB(t)

	tag := &models.Tag{
		ID:   uuid.NewString(),
		Name: "Unique Tag Name",
	}
	err := repo.Create(tag)
	assert.NoError(t, err)

	foundTag, err := repo.GetByName(tag.Name)
	assert.NoError(t, err)
	assert.NotNil(t, foundTag)
	assert.Equal(t, tag.Name, foundTag.Name)

	_, err = repo.GetByName("non-existent-tag")
	assert.Error(t, err)
}

func TestTagRepository_GetAll(t *testing.T) {
	repo := setupTagTestDB(t)

	tag1 := &models.Tag{ID: uuid.NewString(), Name: "A Tag"}
	tag2 := &models.Tag{ID: uuid.NewString(), Name: "B Tag"}

	err := repo.Create(tag1)
	assert.NoError(t, err)
	err = repo.Create(tag2)
	assert.NoError(t, err)

	tags, err := repo.GetAll()
	assert.NoError(t, err)
	assert.Len(t, tags, 2)
	assert.Equal(t, tag1.Name, tags[0].Name, "Tags should be ordered by name ASC")
}

func TestTagRepository_Update(t *testing.T) {
	repo := setupTagTestDB(t)

	tag := &models.Tag{
		ID:   uuid.NewString(),
		Name: "Original Name",
	}
	err := repo.Create(tag)
	assert.NoError(t, err)

	tag.Name = "Updated Name"
	tag.Color = "#0000FF"

	err = repo.Update(tag)
	assert.NoError(t, err)

	updatedTag, err := repo.GetByID(tag.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Name", updatedTag.Name)
	assert.Equal(t, "#0000FF", updatedTag.Color)
}

func TestTagRepository_Delete(t *testing.T) {
	repo := setupTagTestDB(t)

	tag := &models.Tag{
		ID:   uuid.NewString(),
		Name: "To Be Deleted",
	}
	err := repo.Create(tag)
	assert.NoError(t, err)

	err = repo.Delete(tag.ID)
	assert.NoError(t, err)

	_, err = repo.GetByID(tag.ID)
	assert.Error(t, err, "tag not found")
}