package service

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/whitecat/go-reader/internal/models"
	"github.com/whitecat/go-reader/internal/repository"
	"github.com/whitecat/go-reader/internal/scraper"
)

// CrawlerService wraps the custom scraper and persists results to DB.
type CrawlerService struct {
	bookRepo    *repository.BookRepository
	chapterRepo *repository.ChapterRepository

	coversDir string

	mu   sync.Mutex
	jobs map[string]*CrawlerJob
}

type NovelInput struct {
	Title  string
	Author string
	Latest string
	URL    string
}

type CrawlerJob struct {
	ID        string    `json:"id"`
	Status    string    `json:"status"` // pending, running, success, error
	Error     string    `json:"error,omitempty"`
	Total     int       `json:"total"`
	Done      int       `json:"done"`
	BookID    string    `json:"book_id,omitempty"`
	StartedAt time.Time `json:"started_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewCrawlerService(bookRepo *repository.BookRepository, chapterRepo *repository.ChapterRepository) *CrawlerService {
	return &CrawlerService{
		bookRepo:    bookRepo,
		chapterRepo: chapterRepo,
		jobs:        make(map[string]*CrawlerJob),
		coversDir:   "./data/covers",
	}
}

// NewCrawlerServiceWithCoverDir allows specifying custom covers directory.
func NewCrawlerServiceWithCoverDir(bookRepo *repository.BookRepository, chapterRepo *repository.ChapterRepository, coversDir string) *CrawlerService {
	s := NewCrawlerService(bookRepo, chapterRepo)
	if coversDir != "" {
		s.coversDir = coversDir
	}
	return s
}

func (s *CrawlerService) Search(keyword string) ([]scraper.NovelResult, error) {
	return scraper.Search(keyword)
}

// Import downloads chapters synchronously (legacy).
func (s *CrawlerService) Import(novel NovelInput) (*models.Book, error) {
	chaptersInfo, coverURL, err := scraper.GetChapterList(novel.URL)
	if err != nil {
		return nil, fmt.Errorf("get chapter list: %w", err)
	}
	if len(chaptersInfo) == 0 {
		return nil, fmt.Errorf("no chapters found")
	}

	contents := scraper.FetchChapters(chaptersInfo)
	var chapters []models.Chapter
	now := time.Now()
	for i, info := range chaptersInfo {
		chapters = append(chapters, models.Chapter{
			ID:                  uuid.New().String(),
			ChapterNumber:       i + 1,
			VolumeNumber:        1,
			VolumeChapterNumber: i + 1,
			Title:               info.Title,
			Content:             contents[i],
			WordCount:           len([]rune(contents[i])),
			CreatedAt:           now,
		})
	}

	book := &models.Book{
		ID:          uuid.New().String(),
		Title:       novel.Title,
		Author:      novel.Author,
		Description: novel.Latest,
		CoverPath:   s.downloadCover(coverURL),
		FilePath:    novel.URL,
		FileFormat:  "web",
		FileSize:    0,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.bookRepo.Create(book); err != nil {
		return nil, fmt.Errorf("create book: %w", err)
	}
	for i := range chapters {
		chapters[i].BookID = book.ID
	}
	if err := s.chapterRepo.BatchCreate(chapters); err != nil {
		return nil, fmt.Errorf("create chapters: %w", err)
	}

	return book, nil
}

// StartImport kicks off an async import job and returns job id for polling.
func (s *CrawlerService) StartImport(novel NovelInput) string {
	job := &CrawlerJob{
		ID:        uuid.New().String(),
		Status:    "pending",
		StartedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	s.mu.Lock()
	s.jobs[job.ID] = job
	s.mu.Unlock()

	go func() {
		s.setJobStatus(job.ID, "running", "")
		book, err := s.importWithProgress(novel, job.ID)
		if err != nil {
			s.setJobStatus(job.ID, "error", err.Error())
			return
		}
		s.mu.Lock()
		if j, ok := s.jobs[job.ID]; ok {
			j.BookID = book.ID
		}
		s.mu.Unlock()
		s.setJobStatus(job.ID, "success", "")
	}()

	return job.ID
}

func (s *CrawlerService) importWithProgress(novel NovelInput, jobID string) (*models.Book, error) {
	chaptersInfo, coverURL, err := scraper.GetChapterList(novel.URL)
	if err != nil {
		return nil, fmt.Errorf("get chapter list: %w", err)
	}
	if len(chaptersInfo) == 0 {
		return nil, fmt.Errorf("no chapters found")
	}

	s.mu.Lock()
	if job, ok := s.jobs[jobID]; ok {
		job.Total = len(chaptersInfo)
		job.Done = 0
		job.UpdatedAt = time.Now()
	}
	s.mu.Unlock()

	contents := scraper.FetchChaptersWithProgress(chaptersInfo, func(done, total int) {
		s.mu.Lock()
		if job, ok := s.jobs[jobID]; ok {
			job.Done = done
			job.Total = total
			job.UpdatedAt = time.Now()
		}
		s.mu.Unlock()
	})

	var chapters []models.Chapter
	now := time.Now()
	for i, info := range chaptersInfo {
		chapters = append(chapters, models.Chapter{
			ID:                  uuid.New().String(),
			ChapterNumber:       i + 1,
			VolumeNumber:        1,
			VolumeChapterNumber: i + 1,
			Title:               info.Title,
			Content:             contents[i],
			WordCount:           len([]rune(contents[i])),
			CreatedAt:           now,
		})
	}

	book := &models.Book{
		ID:          uuid.New().String(),
		Title:       novel.Title,
		Author:      novel.Author,
		Description: novel.Latest,
		CoverPath:   s.downloadCover(coverURL),
		FilePath:    novel.URL,
		FileFormat:  "web",
		FileSize:    0,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.bookRepo.Create(book); err != nil {
		return nil, fmt.Errorf("create book: %w", err)
	}
	for i := range chapters {
		chapters[i].BookID = book.ID
	}
	if err := s.chapterRepo.BatchCreate(chapters); err != nil {
		return nil, fmt.Errorf("create chapters: %w", err)
	}

	return book, nil
}

func (s *CrawlerService) setJobStatus(id, status, errMsg string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if job, ok := s.jobs[id]; ok {
		job.Status = status
		job.Error = errMsg
		job.UpdatedAt = time.Now()
	}
}

func (s *CrawlerService) GetJob(id string) (*CrawlerJob, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	job, ok := s.jobs[id]
	if !ok {
		return nil, fmt.Errorf("job not found")
	}
	return job, nil
}

// downloadCover saves the cover locally and returns a path accessible by frontend ("/covers/xxx").
// On failure, returns empty string to let frontend fall back to default cover.
func (s *CrawlerService) downloadCover(coverURL string) string {
	if coverURL == "" {
		return ""
	}

	fallback := func() string {
		if u, err := url.Parse(coverURL); err == nil && u.Scheme != "" {
			return coverURL
		}
		return ""
	}

	req, err := http.NewRequest(http.MethodGet, coverURL, nil)
	if err != nil {
		logrus.Warnf("cover: build request failed: %v", err)
		return fallback()
	}
	// Some hosts block default Go UA; mimic browser and set referer to bypass hotlink limits.
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115 Safari/537.36")
	ref := "https://www.biquge321.com/"
	if u, err := url.Parse(coverURL); err == nil && u.Scheme != "" && u.Host != "" {
		ref = u.Scheme + "://" + u.Host + "/"
	}
	req.Header.Set("Referer", ref)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		if err != nil {
			logrus.Warnf("cover: request failed url=%s err=%v", coverURL, err)
		} else {
			logrus.Warnf("cover: non-200 status url=%s status=%d", coverURL, resp.StatusCode)
		}
		return fallback()
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil || len(data) == 0 {
		logrus.Warnf("cover: read body failed url=%s err=%v len=%d", coverURL, err, len(data))
		return fallback()
	}

	ext := ".jpg"
	ct := resp.Header.Get("Content-Type")
	if strings.Contains(ct, "png") {
		ext = ".png"
	} else if strings.Contains(ct, "jpeg") || strings.Contains(ct, "jpg") {
		ext = ".jpg"
	}

	filename := uuid.New().String() + ext
	if s.coversDir == "" {
		s.coversDir = "./data/covers"
	}
	if err := os.MkdirAll(s.coversDir, 0755); err != nil {
		logrus.Warnf("cover: mkdir failed dir=%s err=%v", s.coversDir, err)
		return fallback()
	}
	path := filepath.Join(s.coversDir, filename)
	if err := os.WriteFile(path, data, 0644); err != nil {
		logrus.Warnf("cover: write failed path=%s err=%v", path, err)
		return fallback()
	}
	logrus.Infof("cover: saved url=%s to %s", coverURL, path)
	return "/covers/" + filename
}
