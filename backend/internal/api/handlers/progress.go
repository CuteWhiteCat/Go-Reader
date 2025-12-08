package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/whitecat/go-reader/internal/models"
	"github.com/whitecat/go-reader/internal/service"
	"github.com/whitecat/go-reader/pkg/utils"
)

// ProgressHandler handles progress-related HTTP requests
type ProgressHandler struct {
	progressService *service.ProgressService
}

// NewProgressHandler creates a new ProgressHandler
func NewProgressHandler(progressService *service.ProgressService) *ProgressHandler {
	return &ProgressHandler{
		progressService: progressService,
	}
}

// GetProgress handles GET /api/progress/:bookId
func (h *ProgressHandler) GetProgress(w http.ResponseWriter, r *http.Request) {
	bookID := chi.URLParam(r, "bookId")

	progress, err := h.progressService.GetProgress(bookID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.WriteSuccess(w, progress)
}

// UpdateProgress handles PUT /api/progress/:bookId
func (h *ProgressHandler) UpdateProgress(w http.ResponseWriter, r *http.Request) {
	bookID := chi.URLParam(r, "bookId")

	var req models.UpdateProgressRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	progress, err := h.progressService.UpdateProgress(bookID, &req)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.WriteSuccess(w, progress)
}

// GetBookmarks handles GET /api/bookmarks/:bookId
func (h *ProgressHandler) GetBookmarks(w http.ResponseWriter, r *http.Request) {
	bookID := chi.URLParam(r, "bookId")

	bookmarks, err := h.progressService.GetBookmarks(bookID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.WriteSuccess(w, bookmarks)
}

// CreateBookmark handles POST /api/bookmarks
func (h *ProgressHandler) CreateBookmark(w http.ResponseWriter, r *http.Request) {
	var req models.CreateBookmarkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	bookmark, err := h.progressService.CreateBookmark(&req)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.WriteCreated(w, bookmark)
}

// DeleteBookmark handles DELETE /api/bookmarks/:id
func (h *ProgressHandler) DeleteBookmark(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := h.progressService.DeleteBookmark(id); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.WriteSuccess(w, map[string]string{"message": "Bookmark deleted successfully"})
}
