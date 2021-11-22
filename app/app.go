package app

import (
	"container/list"
	"context"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/dmitriivoitovich/test-assignment-sliide/app/config"
	"github.com/dmitriivoitovich/test-assignment-sliide/app/provider"
	"github.com/dmitriivoitovich/test-assignment-sliide/app/request"
	"github.com/dmitriivoitovich/test-assignment-sliide/app/response"
)

const (
	loadContentTimeout = time.Second * 2
)

type App struct {
	ContentClients map[provider.Provider]provider.Client
	Config         config.ContentMix
}

func (a App) ServeHTTP(w http.ResponseWriter, httpReq *http.Request) {
	// parse request parameters
	req := &request.Request{}
	if err := req.Parse(httpReq); err != nil {
		handleError(w, httpReq, err)

		return
	}

	// load all results from providers
	// and prepare response
	resultsPerProvider := a.loadResults(*req)
	resp := a.prepareResponse(*req, resultsPerProvider)

	handleSuccess(w, httpReq, resp)
}

func (a App) loadResults(req request.Request) map[provider.Provider]*list.List {
	// count how many results we need from each provider
	// including extra results for fallback cases
	resPerProvider := make(map[provider.Provider]int)
	for i := int(req.Offset); i < int(req.Count+req.Offset); i++ {
		providerType := a.Config[i%len(a.Config)].Type
		resPerProvider[providerType]++

		fallbackProviderType := a.Config[i%len(a.Config)].Fallback
		if fallbackProviderType != nil {
			resPerProvider[*fallbackProviderType]++
		}
	}

	// fetch results from all providers simultaneously
	wg := &sync.WaitGroup{}
	contentPerProvider := make(map[provider.Provider]*list.List)

	for i := range resPerProvider {
		providerType := i

		if contentPerProvider[providerType] == nil {
			contentPerProvider[providerType] = list.New()
		}

		wg.Add(1)
		go func() {
			defer wg.Done()

			ctx, cancel := context.WithTimeout(context.Background(), loadContentTimeout)
			defer cancel()

			// ignore errors, will rely on empty result list
			res, _ := getContentWithTimeout(ctx, a.ContentClients[providerType], req.UserIP, resPerProvider[providerType])
			for j := range res {
				contentPerProvider[providerType].PushFront(res[j])
			}
		}()
	}

	wg.Wait()

	return contentPerProvider
}

func (a App) prepareResponse(req request.Request, resultsPerProvider map[provider.Provider]*list.List) response.Response {
	resp := make(response.Response, 0, req.Count)

	for i := int(req.Offset); i < int(req.Count+req.Offset); i++ {
		providerType := a.Config[i%len(a.Config)].Type
		fallbackProviderType := a.Config[i%len(a.Config)].Fallback

		el := resultsPerProvider[providerType].Back()
		if el != nil {
			resultsPerProvider[providerType].Remove(el)
		}

		if el == nil && fallbackProviderType != nil {
			el = resultsPerProvider[*fallbackProviderType].Back()

			if el != nil {
				resultsPerProvider[*fallbackProviderType].Remove(el)
			}
		}

		if el == nil {
			break
		}

		resp = append(resp, *el.Value.(*provider.ContentItem))
	}

	return resp
}

func getContentWithTimeout(ctx context.Context, client provider.Client, userIP net.IP, count int) ([]*provider.ContentItem, error) {
	done := make(chan struct{}, 1)

	var res []*provider.ContentItem
	var err error

	go func() {
		res, err = client.GetContent(userIP.String(), count)
		done <- struct{}{}
	}()

	select {
	case <-done:
		return res, err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}
