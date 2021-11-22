package app

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dmitriivoitovich/test-assignment-sliide/app/config"
	"github.com/dmitriivoitovich/test-assignment-sliide/app/provider"
)

var (
	defaultHandler = App{
		ContentClients: map[provider.Provider]provider.Client{
			provider.Provider1: &provider.ContentProviderMock{Source: provider.Provider1},
			provider.Provider2: &provider.ContentProviderMock{Source: provider.Provider2},
			provider.Provider3: &provider.ContentProviderMock{Source: provider.Provider3},
		},
		Config: config.DefaultContentMix,
	}
)

func TestDefaultRequest(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	content := runRequest(t, defaultHandler, req)

	if len(content) != 5 {
		t.Fatalf("Got %d items back, want 5", len(content))
	}
}

func TestResponseCount(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/?count=10", nil)
	content := runRequest(t, defaultHandler, req)

	if len(content) != 10 {
		t.Fatalf("Got %d items back, want 5", len(content))
	}
}

func TestResponseOrder(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/?count=5", nil)
	content := runRequest(t, defaultHandler, req)

	if len(content) != 5 {
		t.Fatalf("Got %d items back, want 5", len(content))
	}

	for i, item := range content {
		if provider.Provider(item.Source) != defaultHandler.Config[i%len(defaultHandler.Config)].Type {
			t.Errorf(
				"Position %d: Got Provider %v instead of Provider %v",
				i, item.Source, defaultHandler.Config[i].Type,
			)
		}
	}
}

func TestResponseOrder_ClientReturnsError(t *testing.T) {
	provider1Client := &provider.ContentProviderMock{Source: provider.Provider1}
	provider2Client := &provider.ContentProviderMock{Source: provider.Provider2}
	provider3Client := &provider.ContentProviderMock{Source: provider.Provider3}

	provider3Client.SetError(errors.New("expected error"))

	handler := App{
		ContentClients: map[provider.Provider]provider.Client{
			provider.Provider1: provider1Client,
			provider.Provider2: provider2Client,
			provider.Provider3: provider3Client,
		},
		Config: config.ContentMix{
			config.ContentConfig{Type: provider.Provider1},
			config.ContentConfig{Type: provider.Provider2},
			config.ContentConfig{Type: provider.Provider3},
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	content := runRequest(t, handler, req)

	if len(content) != 2 {
		t.Fatalf("Got %d items back, want 2", len(content))
	}

	if provider.Provider(content[0].Source) != provider.Provider1 {
		t.Errorf("Got Provider %v instead of Provider %v", content[0].Source, provider.Provider1)
	}
	if provider.Provider(content[1].Source) != provider.Provider2 {
		t.Errorf("Got Provider %v instead of Provider %v", content[1].Source, provider.Provider2)
	}
}

func TestResponseOrder_NotEnoughResults(t *testing.T) {
	providerClient := &provider.ContentProviderMock{Source: provider.Provider1}

	responseMock := []*provider.ContentItem{
		{ID: "1"},
		{ID: "2"},
	}
	providerClient.SetResponse(responseMock)

	handler := App{
		ContentClients: map[provider.Provider]provider.Client{provider.Provider1: providerClient},
		Config:         config.ContentMix{config.ContentConfig{Type: provider.Provider1}},
	}

	req := httptest.NewRequest(http.MethodGet, "/?count=3", nil)
	content := runRequest(t, handler, req)

	if len(content) != 2 {
		t.Fatalf("Got %d items back, want 2", len(content))
	}

	for i := range content {
		if content[i].ID != responseMock[i].ID {
			t.Errorf("Got ID %v instead of ID %v", content[i].ID, responseMock[i].ID)
		}
	}
}

func TestOffsetResponseOrder(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/?count=5&offset=5", nil)
	content := runRequest(t, defaultHandler, req)

	if len(content) != 5 {
		t.Fatalf("Got %d items back, want 5", len(content))
	}

	for j, item := range content {
		i := j + 5
		if provider.Provider(item.Source) != defaultHandler.Config[i%len(defaultHandler.Config)].Type {
			t.Errorf(
				"Position %d: Got Provider %v instead of Provider %v",
				i, item.Source, defaultHandler.Config[i].Type,
			)
		}
	}
}

func TestFallback_ClientDoesNotRespond(t *testing.T) {
	provider1Client := &provider.ContentProviderMock{Source: provider.Provider1}
	provider1Client.SetDelay(loadContentTimeout + time.Second)
	provider2Client := &provider.ContentProviderMock{Source: provider.Provider2}

	handler := App{
		ContentClients: map[provider.Provider]provider.Client{
			provider.Provider1: provider1Client,
			provider.Provider2: provider2Client,
		},
		Config: config.ContentMix{
			config.ContentConfig{
				Type:     provider.Provider1,
				Fallback: &provider.Provider2,
			},
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/?count=1", nil)
	content := runRequest(t, handler, req)

	if len(content) != 1 {
		t.Fatalf("Got %d items back, want 1", len(content))
	}

	if provider.Provider(content[0].Source) != provider.Provider2 {
		t.Errorf("Got Provider %v instead of Provider %v", content[0].Source, provider.Provider1)
	}
}

func TestFallback_ClientReturnsError(t *testing.T) {
	provider1Client := &provider.ContentProviderMock{Source: provider.Provider1}
	provider1Client.SetError(errors.New("expected error"))
	provider2Client := &provider.ContentProviderMock{Source: provider.Provider2}

	handler := App{
		ContentClients: map[provider.Provider]provider.Client{
			provider.Provider1: provider1Client,
			provider.Provider2: provider2Client,
		},
		Config: config.ContentMix{
			config.ContentConfig{
				Type:     provider.Provider1,
				Fallback: &provider.Provider2,
			},
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/?count=1", nil)
	content := runRequest(t, handler, req)

	if len(content) != 1 {
		t.Fatalf("Got %d items back, want 1", len(content))
	}

	if provider.Provider(content[0].Source) != provider.Provider2 {
		t.Errorf("Got Provider %v instead of Provider %v", content[0].Source, provider.Provider2)
	}
}

func runRequest(t *testing.T, srv http.Handler, r *http.Request) (content []*provider.ContentItem) {
	response := httptest.NewRecorder()
	srv.ServeHTTP(response, r)

	if response.Code != 200 {
		t.Fatalf("Response code is %d, want 200", response.Code)
		return
	}

	decoder := json.NewDecoder(response.Body)
	if err := decoder.Decode(&content); err != nil {
		t.Fatalf("couldn't decode response json: %v", err)
	}

	return content
}
