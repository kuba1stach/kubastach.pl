package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func fixedPosts() map[string]Post {
	t1 := time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC)
	t2 := time.Date(2025, 1, 15, 15, 30, 0, 0, time.UTC)
	return map[string]Post{
		t1.Format("2006-01-02"): {Content: "One", Timestamp: t1.Format(time.RFC3339)},
		t2.Format("2006-01-02"): {Content: "Two", Timestamp: t2.Format(time.RFC3339)},
	}
}

func TestGetPosts(t *testing.T) {
	s := NewServerWithPosts(fixedPosts())
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/posts", nil)

	router := NewRouter(s)
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}
	var posts []Post
	if err := json.Unmarshal(rr.Body.Bytes(), &posts); err != nil {
		t.Fatalf("json decode: %v", err)
	}
	if len(posts) != 2 {
		t.Fatalf("expected 2 posts, got %d", len(posts))
	}
	if posts[0].Timestamp > posts[1].Timestamp {
		t.Fatalf("posts not sorted by timestamp ascending")
	}
}

func TestGetPostByDateFound(t *testing.T) {
	s := NewServerWithPosts(fixedPosts())
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/posts/2025-01-15", nil)
	NewRouter(s).ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", rr.Code)
	}
	var p Post
	if err := json.Unmarshal(rr.Body.Bytes(), &p); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if p.Content != "Two" {
		t.Fatalf("unexpected content: %s", p.Content)
	}
}

func TestGetPostByDateNotFound(t *testing.T) {
	s := NewServerWithPosts(fixedPosts())
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/posts/2025-02-01", nil)
	NewRouter(s).ServeHTTP(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected 404 got %d", rr.Code)
	}
}

func TestGetPostByDateInvalid(t *testing.T) {
	s := NewServerWithPosts(fixedPosts())
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/posts/not-a-date", nil)
	NewRouter(s).ServeHTTP(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected 404 got %d", rr.Code)
	}
}

func TestGetPostDates(t *testing.T) {
	s := NewServerWithPosts(fixedPosts())
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/posts/dates", nil)
	NewRouter(s).ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", rr.Code)
	}
	var dates []string
	if err := json.Unmarshal(rr.Body.Bytes(), &dates); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(dates) != 2 {
		t.Fatalf("expected 2 dates got %d", len(dates))
	}
	if dates[0] != "2025-01-01" || dates[1] != "2025-01-15" {
		t.Fatalf("unexpected dates order/content: %#v", dates)
	}
}
