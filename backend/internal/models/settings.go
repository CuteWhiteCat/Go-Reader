package models

import "time"

// Settings represents application settings
type Settings struct {
	Key       string    `json:"key" db:"key"`
	Value     string    `json:"value" db:"value"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// AppSettings represents the complete application settings
type AppSettings struct {
	Theme           string  `json:"theme"` // light, dark
	FontSize        int     `json:"font_size"`
	LineSpacing     float64 `json:"line_spacing"`
	PageMargin      int     `json:"page_margin"`
	DefaultLibraryPath string `json:"default_library_path"`
}
