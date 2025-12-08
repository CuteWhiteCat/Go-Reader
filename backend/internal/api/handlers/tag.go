package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/whitecat/go-reader/internal/models"
	"github.com/whitecat/go-reader/internal/service"
	"github.com/whitecat/go-reader/pkg/utils"
)

// TagHandler handles tag-related HTTP requests
type TagHandler struct {
	tagService *service.TagService
}

// NewTagHandler creates a new TagHandler
func NewTagHandler(tagService *service.TagService) *TagHandler {
	return &TagHandler{
		tagService: tagService,
	}
}

// GetAllTags handles GET /api/tags
func (h *TagHandler) GetAllTags(w http.ResponseWriter, r *http.Request) {
	tags, err := h.tagService.GetAllTags()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.WriteSuccess(w, tags)
}

// CreateTag handles POST /api/tags
func (h *TagHandler) CreateTag(w http.ResponseWriter, r *http.Request) {
	var req models.CreateTagRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	tag, err := h.tagService.CreateTag(&req)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.WriteCreated(w, tag)
}

// UpdateTag handles PUT /api/tags/:id
func (h *TagHandler) UpdateTag(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req models.UpdateTagRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	tag, err := h.tagService.UpdateTag(id, &req)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.WriteSuccess(w, tag)
}

// DeleteTag handles DELETE /api/tags/:id
func (h *TagHandler) DeleteTag(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := h.tagService.DeleteTag(id); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.WriteSuccess(w, map[string]string{"message": "Tag deleted successfully"})
}
