package api

import (
	"encoding/json"
	"net/http"
	"sort"
	"strings"
	"time"
)

// Media corresponds to components.schemas.Media in openapi.yaml
type Media struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// Post corresponds to components.schemas.Post in openapi.yaml
type Post struct {
	Content   string  `json:"content"`
	Timestamp string  `json:"timestamp"` // RFC3339 date-time
	Media     []Media `json:"media,omitempty"`
}

// Server holds in-memory data implementing the spec operations.
type Server struct {
	posts map[string]Post // key: YYYY-MM-DD (date part)
}

// NewServer creates a new instance of our server with sample data.
func NewServer() *Server {
	now := time.Now().UTC()
	sample := map[string]Post{}
	// two sample posts on fixed dates for deterministic behavior
	d1 := time.Date(now.Year(), 1, 1, 10, 0, 0, 0, time.UTC)
	d2 := time.Date(now.Year(), 1, 15, 15, 30, 0, 0, time.UTC)
	sample[d1.Format("2006-01-02")] = Post{
		Content:   "Hello, world!",
		Timestamp: d1.Format(time.RFC3339),
		Media:     []Media{{Name: "intro.png", Type: "image/png"}},
	}
	sample[d2.Format("2006-01-02")] = Post{
		Content:   "This is a diary entry.",
		Timestamp: d2.Format(time.RFC3339),
	}
	return &Server{posts: sample}
}

// NewServerWithPosts allows injecting predefined posts (date->Post) for tests.
func NewServerWithPosts(posts map[string]Post) *Server {
	// shallow copy to avoid external mutation side-effects
	cp := make(map[string]Post, len(posts))
	for k, v := range posts {
		cp[k] = v
	}
	return &Server{posts: cp}
}

// GetPosts handles GET /posts returning all posts.
func (s *Server) GetPosts(w http.ResponseWriter, r *http.Request) {
	list := make([]Post, 0, len(s.posts))
	for _, p := range s.posts {
		list = append(list, p)
	}
	// stable order by timestamp
	sort.Slice(list, func(i, j int) bool { return list[i].Timestamp < list[j].Timestamp })
	writeJSON(w, http.StatusOK, list)
}

// GetPostByDate handles GET /posts/{date}
func (s *Server) GetPostByDate(w http.ResponseWriter, r *http.Request, date string) {
	// date expected format YYYY-MM-DD; minimal validation
	if len(date) != 10 || strings.Count(date, "-") != 2 {
		http.NotFound(w, r)
		return
	}
	post, ok := s.posts[date]
	if !ok {
		http.NotFound(w, r)
		return
	}
	writeJSON(w, http.StatusOK, post)
}

// GetPostDates handles GET /posts/dates
func (s *Server) GetPostDates(w http.ResponseWriter, r *http.Request) {
	dates := make([]string, 0, len(s.posts))
	for d := range s.posts {
		dates = append(dates, d)
	}
	sort.Strings(dates)
	writeJSON(w, http.StatusOK, dates)
}

// Helper to write JSON responses.
func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
