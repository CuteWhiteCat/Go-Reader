package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/whitecat/go-reader/internal/api/handlers"
	"github.com/whitecat/go-reader/internal/api/middleware"
)

// Router holds all HTTP handlers
type Router struct {
	BookHandler     *handlers.BookHandler
	TagHandler      *handlers.TagHandler
	ProgressHandler *handlers.ProgressHandler
	CrawlerHandler  *handlers.CrawlerHandler
}

// NewRouter creates a new API router
func NewRouter(
	bookHandler *handlers.BookHandler,
	tagHandler *handlers.TagHandler,
	progressHandler *handlers.ProgressHandler,
	crawlerHandler *handlers.CrawlerHandler,
) *Router {
	return &Router{
		BookHandler:     bookHandler,
		TagHandler:      tagHandler,
		ProgressHandler: progressHandler,
		CrawlerHandler:  crawlerHandler,
	}
}

// SetupRoutes configures all API routes
func (router *Router) SetupRoutes() *chi.Mux {
	r := chi.NewRouter()

	// Middleware
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(middleware.CORS().Handler)

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	// API routes
	r.Route("/api", func(r chi.Router) {
		// Books
		r.Route("/books", func(r chi.Router) {
			r.Get("/", router.BookHandler.GetAllBooks)
			r.Post("/", router.BookHandler.CreateBook)
			r.Get("/{id}", router.BookHandler.GetBook)
			r.Put("/{id}", router.BookHandler.UpdateBook)
			r.Delete("/{id}", router.BookHandler.DeleteBook)
			r.Get("/{id}/content", router.BookHandler.GetBookContent)
			r.Get("/{id}/chapters", router.BookHandler.GetBookChapters)
			r.Get("/{id}/chapters/{number}", router.BookHandler.GetChapter)
		})

		// Tags
		r.Route("/tags", func(r chi.Router) {
			r.Get("/", router.TagHandler.GetAllTags)
			r.Post("/", router.TagHandler.CreateTag)
			r.Put("/{id}", router.TagHandler.UpdateTag)
			r.Delete("/{id}", router.TagHandler.DeleteTag)
		})

		// Progress
		r.Route("/progress", func(r chi.Router) {
			r.Get("/{bookId}", router.ProgressHandler.GetProgress)
			r.Put("/{bookId}", router.ProgressHandler.UpdateProgress)
		})

		// Bookmarks
		r.Route("/bookmarks", func(r chi.Router) {
			r.Get("/{bookId}", router.ProgressHandler.GetBookmarks)
			r.Post("/", router.ProgressHandler.CreateBookmark)
			r.Delete("/{id}", router.ProgressHandler.DeleteBookmark)
		})

		// Crawler
		r.Route("/crawler", func(r chi.Router) {
			r.Post("/search", router.CrawlerHandler.Search)
			r.Post("/import", router.CrawlerHandler.Import)
			r.Post("/import/start", router.CrawlerHandler.StartImport)
			r.Get("/import/status", router.CrawlerHandler.ImportStatus)
		})
	})

	return r
}
