package provider

import (
	"time"
)

// ContentItem represent one piece of content fetched from a provider
type ContentItem struct {
	ID      string    `json:"id"`
	Title   string    `json:"title"`
	Source  string    `json:"source"`
	Summary string    `json:"summary"`
	Link    string    `json:"link"`
	Expiry  time.Time `json:"expiry"`
}

// Client represents a provider's client or SDK
type Client interface {
	GetContent(userIP string, count int) ([]*ContentItem, error)
}
