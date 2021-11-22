package provider

import (
	"math/rand"
	"strconv"
	"time"
)

// SampleContentProvider is an example for a Provider's client
type SampleContentProvider struct {
	Source Provider
}

// GetContent returns content items given a user IP, and the number of content items desired.
func (cp *SampleContentProvider) GetContent(userIP string, count int) ([]*ContentItem, error) {
	resp := make([]*ContentItem, count)
	for i, _ := range resp {
		resp[i] = &ContentItem{
			ID:      strconv.Itoa(rand.Int()),
			Title:   "title",
			Source:  string(cp.Source),
			Summary: "",
			Link:    "",
			Expiry:  time.Now(),
		}

	}
	return resp, nil
}
