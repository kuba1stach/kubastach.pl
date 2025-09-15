package api

import (
	"net/http"
	"strings"
)

// NewRouter wires HTTP routes to server handlers. Base path: /api/v1
func NewRouter(s *Server) http.Handler {
	mux := http.NewServeMux()

	// GET /api/v1/posts
	mux.HandleFunc("/api/v1/posts", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.NotFound(w, r)
			return
		}
		// Distinguish /posts vs /posts/dates vs /posts/{date}
		// Exact match handled here; delegate others below.
		if r.URL.Path == "/api/v1/posts" { // list all posts
			s.GetPosts(w, r)
			return
		}
		http.NotFound(w, r)
	})

	// dynamic: /api/v1/posts/{date} and /api/v1/posts/dates
	mux.HandleFunc("/api/v1/posts/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.NotFound(w, r)
			return
		}
		rest := strings.TrimPrefix(r.URL.Path, "/api/v1/posts/")
		if rest == "dates" { // /posts/dates
			s.GetPostDates(w, r)
			return
		}
		// treat rest as date param
		s.GetPostByDate(w, r, rest)
	})

	return mux
}
