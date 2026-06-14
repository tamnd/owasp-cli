package owasp_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/tamnd/owasp-cli/owasp"
)

func makeFiles(names []string) []byte {
	type file struct {
		Name        string `json:"name"`
		DownloadURL string `json:"download_url"`
		HTMLURL     string `json:"html_url"`
		Size        int    `json:"size"`
		Type        string `json:"type"`
	}
	files := make([]file, len(names))
	for i, n := range names {
		files[i] = file{
			Name:        n,
			DownloadURL: "https://raw.githubusercontent.com/OWASP/CheatSheetSeries/master/cheatsheets/" + n,
			HTMLURL:     "https://github.com/OWASP/CheatSheetSeries/blob/master/cheatsheets/" + n,
			Size:        1000,
			Type:        "file",
		}
	}
	b, _ := json.Marshal(files)
	return b
}

func TestList(t *testing.T) {
	payload := makeFiles([]string{
		"AJAX_Security_Cheat_Sheet.md",
		"Access_Control_Cheat_Sheet.md",
		"Authentication_Cheat_Sheet.md",
	})

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("User-Agent") == "" {
			t.Error("request carried no User-Agent")
		}
		_, _ = w.Write(payload)
	}))
	defer srv.Close()

	cfg := owasp.DefaultConfig()
	cfg.BaseURL = srv.URL
	cfg.Rate = 0

	c := owasp.NewClient(cfg)
	sheets, err := c.List(context.Background(), 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(sheets) != 3 {
		t.Fatalf("got %d sheets, want 3", len(sheets))
	}
	if sheets[0].Name != "AJAX Security" {
		t.Errorf("name = %q, want %q", sheets[0].Name, "AJAX Security")
	}
	if sheets[0].Rank != 1 {
		t.Errorf("rank = %d, want 1", sheets[0].Rank)
	}
}

func TestListLimit(t *testing.T) {
	names := make([]string, 5)
	for i := range names {
		names[i] = "Sheet_Cheat_Sheet.md"
	}
	payload := makeFiles(names)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write(payload)
	}))
	defer srv.Close()

	cfg := owasp.DefaultConfig()
	cfg.BaseURL = srv.URL
	cfg.Rate = 0

	c := owasp.NewClient(cfg)
	sheets, err := c.List(context.Background(), 3)
	if err != nil {
		t.Fatal(err)
	}
	if len(sheets) != 3 {
		t.Fatalf("got %d sheets, want 3 (limit applied)", len(sheets))
	}
}

func TestSearch(t *testing.T) {
	payload := makeFiles([]string{
		"AJAX_Security_Cheat_Sheet.md",
		"SQL_Injection_Prevention_Cheat_Sheet.md",
		"XSS_Prevention_Cheat_Sheet.md",
	})

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write(payload)
	}))
	defer srv.Close()

	cfg := owasp.DefaultConfig()
	cfg.BaseURL = srv.URL
	cfg.Rate = 0

	c := owasp.NewClient(cfg)
	sheets, err := c.Search(context.Background(), "sql", 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(sheets) != 1 {
		t.Fatalf("got %d sheets, want 1", len(sheets))
	}
	if sheets[0].Name != "SQL Injection Prevention" {
		t.Errorf("name = %q", sheets[0].Name)
	}
}
