package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"kubastach.pl/backend/internal/models"
)

// PostRepository defines data access needed by handlers.
type PostRepository interface {
	ListPosts(ctx context.Context) ([]models.Post, error)
	GetPostByDate(ctx context.Context, date string) (*models.Post, error)
	ListDates(ctx context.Context) ([]string, error)
	GetActivitiesInDateRange(ctx context.Context, minDate, maxDate time.Time) ([]models.Activity, error)
	GetJunkFoodInDateRange(ctx context.Context, minDate, maxDate time.Time) ([]models.JunkFood, error)
}

type ProgressService interface {
	GetProgressInRange(ctx context.Context, minDate, maxDate time.Time) (models.ProgressPerDay, error)
}

// Register registers all HTTP handlers onto the provided mux.
func Register(mux *http.ServeMux, repo PostRepository, service ProgressService) {
	mux.HandleFunc("/api/v1/posts", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		posts, err := repo.ListPosts(r.Context())
		if err != nil {
			internalError(w, err)
			return
		}
		writeJSON(w, posts)
	})

	mux.HandleFunc("/api/v1/posts/dates", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		dates, err := repo.ListDates(r.Context())
		if err != nil {
			internalError(w, err)
			return
		}
		writeJSON(w, dates)
	})

	mux.HandleFunc("/api/v1/posts/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/v1/posts/"), "/")
		if len(parts) != 1 || parts[0] == "" || parts[0] == "dates" {
			http.NotFound(w, r)
			return
		}
		date := parts[0]
		if _, err := time.Parse("2006-01-02", date); err != nil {
			http.Error(w, "invalid date", http.StatusBadRequest)
			return
		}
		post, err := repo.GetPostByDate(r.Context(), date)
		if err != nil {
			internalError(w, err)
			return
		}
		if post == nil {
			http.NotFound(w, r)
			return
		}
		writeJSON(w, post)
	})

	mux.HandleFunc("/api/v1/activities/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/v1/activities/"), "/")
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			http.NotFound(w, r)
			return
		}

		startDateStr := parts[0]
		endDateStr := parts[1]

		startDate, err := time.Parse("2006-01-02", startDateStr)
		if err != nil {
			http.Error(w, "invalid start date", http.StatusBadRequest)
			return
		}

		endDate, err := time.Parse("2006-01-02", endDateStr)
		if err != nil {
			http.Error(w, "invalid end date", http.StatusBadRequest)
			return
		}

		if endDate.Before(startDate) {
			http.Error(w, "end date must be after start date", http.StatusBadRequest)
			return
		}
		activities, err := repo.GetActivitiesInDateRange(r.Context(), startDate, endDate)
		if err != nil {
			internalError(w, err)
			return
		}
		writeJSON(w, activities)
	})

	mux.HandleFunc("/api/v1/junkFood/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/v1/junkFood/"), "/")
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			http.NotFound(w, r)
			return
		}

		startDateStr := parts[0]
		endDateStr := parts[1]

		startDate, err := time.Parse("2006-01-02", startDateStr)
		if err != nil {
			http.Error(w, "invalid start date", http.StatusBadRequest)
			return
		}

		endDate, err := time.Parse("2006-01-02", endDateStr)
		if err != nil {
			http.Error(w, "invalid end date", http.StatusBadRequest)
			return
		}

		if endDate.Before(startDate) {
			http.Error(w, "end date must be after start date", http.StatusBadRequest)
			return
		}
		junkFoods, err := repo.GetJunkFoodInDateRange(r.Context(), startDate, endDate)
		if err != nil {
			internalError(w, err)
			return
		}
		writeJSON(w, junkFoods)
	})

	mux.HandleFunc("/api/v1/progress/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/v1/progress/"), "/")
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			http.NotFound(w, r)
			return
		}

		startDateStr := parts[0]
		endDateStr := parts[1]

		startDate, err := time.Parse("2006-01-02", startDateStr)
		if err != nil {
			http.Error(w, "invalid start date", http.StatusBadRequest)
			return
		}

		endDate, err := time.Parse("2006-01-02", endDateStr)
		if err != nil {
			http.Error(w, "invalid end date", http.StatusBadRequest)
			return
		}

		if endDate.Before(startDate) {
			http.Error(w, "end date must be after start date", http.StatusBadRequest)
			return
		}
		junkFoods, err := service.GetProgressInRange(r.Context(), startDate, endDate)
		if err != nil {
			internalError(w, err)
			return
		}
		writeJSON(w, junkFoods)
	})
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("encode response: %v", err)
	}
}

func internalError(w http.ResponseWriter, err error) {
	log.Printf("internal error: %v", err)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
