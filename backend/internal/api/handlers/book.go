package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/whitecat/go-reader/internal/models"
	"github.com/whitecat/go-reader/internal/service"
	"github.com/whitecat/go-reader/pkg/utils"
)

// BookHandler handles book-related HTTP requests
type BookHandler struct {
	bookService *service.BookService
}

// NewBookHandler creates a new BookHandler
func NewBookHandler(bookService *service.BookService) *BookHandler {
	return &BookHandler{
		bookService: bookService,
	}
}

// GetAllBooks handles GET /api/books
func (h *BookHandler) GetAllBooks(w http.ResponseWriter, r *http.Request) {
	books, err := h.bookService.GetAllBooks()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.WriteSuccess(w, books)
}

// GetBook handles GET /api/books/:id
func (h *BookHandler) GetBook(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	book, err := h.bookService.GetBook(id)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, err.Error())
		return
	}

	utils.WriteSuccess(w, book)
}

// CreateBook handles POST /api/books
func (h *BookHandler) CreateBook(w http.ResponseWriter, r *http.Request) {
	var req models.CreateBookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	book, err := h.bookService.CreateBook(&req)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.WriteCreated(w, book)
}

// UpdateBook handles PUT /api/books/:id
func (h *BookHandler) UpdateBook(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req models.UpdateBookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	book, err := h.bookService.UpdateBook(id, &req)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.WriteSuccess(w, book)
}

// DeleteBook handles DELETE /api/books/:id
func (h *BookHandler) DeleteBook(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := h.bookService.DeleteBook(id); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.WriteSuccess(w, map[string]string{"message": "Book deleted successfully"})
}

// GetBookContent handles GET /api/books/:id/content
func (h *BookHandler) GetBookContent(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	chapters, err := h.bookService.GetBookContent(id)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.WriteSuccess(w, chapters)
}

// GetBookChapters handles GET /api/books/:id/chapters
func (h *BookHandler) GetBookChapters(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	chapters, err := h.bookService.GetBookChapters(id)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.WriteSuccess(w, chapters)
}

// GetChapter handles GET /api/books/:id/chapters/:number
func (h *BookHandler) GetChapter(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var chapterNumber int
	if _, err := fmt.Sscanf(chi.URLParam(r, "number"), "%d", &chapterNumber); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid chapter number")
		return
	}

	chapter, err := h.bookService.GetChapter(id, chapterNumber)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, err.Error())
		return
	}

	utils.WriteSuccess(w, chapter)
}
