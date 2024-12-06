package scrape

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/regalias/scry/pkg/throttle"
)

type ScrapeClient struct {
	client *http.Client
}

func NewScrapeClient() (*ScrapeClient, error) {

	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.MaxConnsPerHost = 2
	transport.MaxIdleConnsPerHost = 2

	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create cookiejar: %w", err)
	}

	return &ScrapeClient{
		client: &http.Client{
			Transport: throttle.NewWrappedTransport(transport, 2),
			Timeout:   time.Second * 60,
			Jar:       jar,
		},
	}, nil
}

func (sc *ScrapeClient) newRequest(ctx context.Context, method string, url string) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", USER_AGENT)
	return req, nil
}

func (sc *ScrapeClient) doRequest(ctx context.Context, method string, url string, expectedMIME string) (*http.Response, error) {
	req, err := sc.newRequest(ctx, method, url)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := sc.client.Do(req)
	if err != nil {
		resp.Body.Close()
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("got unexpected HTTP status response: %s", resp.Status)
	}

	// Loose content type check
	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(strings.ToLower(contentType), expectedMIME) {
		resp.Body.Close()
		return nil, fmt.Errorf("got wrong content type: %s", contentType)
	}
	return resp, err
}

func (sc *ScrapeClient) Visit(ctx context.Context, url string) (*goquery.Document, error) {

	resp, err := sc.doRequest(ctx, http.MethodGet, url, "html")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)

	if err != nil {
		return nil, fmt.Errorf("failed to load document: %w", err)
	}

	return doc, nil
}

func (sc *ScrapeClient) VisitJson(ctx context.Context, url string) ([]byte, error) {
	resp, err := sc.doRequest(ctx, http.MethodGet, url, "application/json")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed reading body: %w", err)
	}

	return data, nil
}
