package provider

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

type ContentProviderMock struct {
	Source Provider

	delay    time.Duration
	err      error
	response []*ContentItem
}

func (cp *ContentProviderMock) SetDelay(delay time.Duration) {
	cp.delay = delay
}

func (cp *ContentProviderMock) SetError(err error) {
	cp.err = err
}

func (cp *ContentProviderMock) SetResponse(response []*ContentItem) {
	cp.response = response
}

func (cp *ContentProviderMock) GetContent(_ string, count int) ([]*ContentItem, error) {
	if cp.delay > 0 {
		time.Sleep(cp.delay)
	}

	if cp.err != nil {
		return nil, cp.err
	}

	if cp.response != nil {
		return cp.response, nil
	}

	resp := make([]*ContentItem, count)
	for i := range resp {
		resp[i] = &ContentItem{
			ID:      strconv.Itoa(rand.Int()),
			Title:   fmt.Sprintf("Item #%d", i),
			Source:  string(cp.Source),
			Summary: fmt.Sprintf("Item summary #%d", i),
			Link:    fmt.Sprintf("https://%s.com/%d", cp.Source, i),
			Expiry:  time.Now(),
		}
	}

	return resp, nil
}
