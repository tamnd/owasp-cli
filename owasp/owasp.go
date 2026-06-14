// Package owasp is the library behind the owasp command.
package owasp

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

const DefaultUserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"

// Config holds constructor parameters.
type Config struct {
	BaseURL   string
	UserAgent string
	Rate      time.Duration
	Retries   int
	Timeout   time.Duration
}

// DefaultConfig returns sensible defaults.
func DefaultConfig() Config {
	return Config{
		BaseURL:   "https://api.github.com",
		UserAgent: DefaultUserAgent,
		Rate:      200 * time.Millisecond,
		Retries:   3,
		Timeout:   30 * time.Second,
	}
}

// Client talks to the GitHub Contents API for the OWASP CheatSheetSeries repo.
type Client struct {
	cfg        Config
	httpClient *http.Client
	mu         sync.Mutex
	last       time.Time
}

// NewClient returns a Client with the given config.
func NewClient(cfg Config) *Client {
	return &Client{
		cfg:        cfg,
		httpClient: &http.Client{Timeout: cfg.Timeout},
	}
}

func (c *Client) get(ctx context.Context, url string) ([]byte, error) {
	var lastErr error
	for attempt := 0; attempt <= c.cfg.Retries; attempt++ {
		if attempt > 0 {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(backoff(attempt)):
			}
		}
		b, retry, err := c.do(ctx, url)
		if err == nil {
			return b, nil
		}
		lastErr = err
		if !retry {
			return nil, err
		}
	}
	return nil, fmt.Errorf("get: %w", lastErr)
}

func (c *Client) do(ctx context.Context, url string) ([]byte, bool, error) {
	c.pace()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, false, err
	}
	req.Header.Set("User-Agent", c.cfg.UserAgent)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, true, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode >= 500 {
		return nil, true, fmt.Errorf("http %d", resp.StatusCode)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, false, fmt.Errorf("http %d", resp.StatusCode)
	}

	b, err := io.ReadAll(io.LimitReader(resp.Body, 16<<20))
	if err != nil {
		return nil, true, err
	}
	return b, false, nil
}

func (c *Client) pace() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.cfg.Rate <= 0 {
		return
	}
	if wait := c.cfg.Rate - time.Since(c.last); wait > 0 {
		time.Sleep(wait)
	}
	c.last = time.Now()
}

func backoff(attempt int) time.Duration {
	d := time.Duration(attempt) * 500 * time.Millisecond
	if d > 5*time.Second {
		d = 5 * time.Second
	}
	return d
}

// List fetches all OWASP cheat sheets.
func (c *Client) List(ctx context.Context, limit int) ([]Sheet, error) {
	url := c.cfg.BaseURL + "/repos/OWASP/CheatSheetSeries/contents/cheatsheets"
	raw, err := c.get(ctx, url)
	if err != nil {
		return nil, err
	}
	var files []wireFile
	if err := json.Unmarshal(raw, &files); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	sheets := make([]Sheet, 0, len(files))
	rank := 1
	for _, f := range files {
		if f.Type != "file" || !strings.HasSuffix(f.Name, ".md") {
			continue
		}
		if limit > 0 && len(sheets) >= limit {
			break
		}
		sheets = append(sheets, wireToSheet(f, rank))
		rank++
	}
	return sheets, nil
}

// Search returns sheets whose name contains query (case-insensitive).
func (c *Client) Search(ctx context.Context, query string, limit int) ([]Sheet, error) {
	all, err := c.List(ctx, 0)
	if err != nil {
		return nil, err
	}
	q := strings.ToLower(query)
	var out []Sheet
	for _, s := range all {
		if strings.Contains(strings.ToLower(s.Name), q) {
			out = append(out, s)
			if limit > 0 && len(out) >= limit {
				break
			}
		}
	}
	return out, nil
}

// sheetName converts a filename like "AJAX_Security_Cheat_Sheet.md" to "AJAX Security".
func sheetName(filename string) string {
	name := strings.TrimSuffix(filename, ".md")
	name = strings.TrimSuffix(name, "_Cheat_Sheet")
	return strings.ReplaceAll(name, "_", " ")
}

func wireToSheet(f wireFile, rank int) Sheet {
	return Sheet{
		Rank: rank,
		Name: sheetName(f.Name),
		URL:  f.HTMLURL,
		Raw:  f.DownloadURL,
		Size: f.Size,
	}
}
