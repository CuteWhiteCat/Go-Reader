package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/whitecat/go-reader/internal/service"
	"github.com/whitecat/go-reader/pkg/utils"
)

type CrawlerHandler struct {
	crawler *service.CrawlerService
}

func NewCrawlerHandler(crawler *service.CrawlerService) *CrawlerHandler {
	return &CrawlerHandler{crawler: crawler}
}

// POST /api/crawler/search {query}
func (h *CrawlerHandler) Search(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Query string `json:"query"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Query == "" {
		utils.WriteError(w, http.StatusBadRequest, "query is required")
		return
	}
	res, err := h.crawler.Search(req.Query)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.WriteSuccess(w, res)
}

// POST /api/crawler/import {title,author,url}
func (h *CrawlerHandler) Import(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Title  string `json:"title"`
		Author string `json:"author"`
		URL    string `json:"url"`
		Latest string `json:"latest"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.URL == "" {
		utils.WriteError(w, http.StatusBadRequest, "url is required")
		return
	}
	book, err := h.crawler.Import(serviceToNovel(req))
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.WriteCreated(w, book)
}

// POST /api/crawler/import/start {title,author,url,latest}
func (h *CrawlerHandler) StartImport(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Title  string `json:"title"`
		Author string `json:"author"`
		URL    string `json:"url"`
		Latest string `json:"latest"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.URL == "" {
		utils.WriteError(w, http.StatusBadRequest, "url is required")
		return
	}
	jobID := h.crawler.StartImport(serviceToNovel(req))
	utils.WriteSuccess(w, map[string]string{"job_id": jobID})
}

// GET /api/crawler/import/status?id=xxx
func (h *CrawlerHandler) ImportStatus(w http.ResponseWriter, r *http.Request) {
	jobID := r.URL.Query().Get("id")
	if jobID == "" {
		utils.WriteError(w, http.StatusBadRequest, "id is required")
		return
	}
	job, err := h.crawler.GetJob(jobID)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, err.Error())
		return
	}
	utils.WriteSuccess(w, job)
}

func serviceToNovel(req struct {
	Title  string `json:"title"`
	Author string `json:"author"`
	URL    string `json:"url"`
	Latest string `json:"latest"`
}) service.NovelInput {
	return service.NovelInput{
		Title:  req.Title,
		Author: req.Author,
		Latest: req.Latest,
		URL:    req.URL,
	}
}
