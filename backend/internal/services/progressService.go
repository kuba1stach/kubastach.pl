package services

import (
	"context"
	"time"

	"kubastach.pl/backend/internal/models"
	"kubastach.pl/backend/internal/repositories"
)

// ProgressService provides methods to track and report progress of long-running operations.
type ProgressService struct {
	cosmosRepository *repositories.CosmosRepository
}

// NewProgressService creates a new instance of ProgressService.
func NewProgressService(ctx context.Context, cosmosRepository *repositories.CosmosRepository) (*ProgressService, error) {
	return &ProgressService{cosmosRepository: cosmosRepository}, nil
}

func (s *ProgressService) GetProgressInRange(ctx context.Context, minDate, maxDate time.Time) (models.ProgressPerDay, error) {
	activities, err := s.cosmosRepository.GetActivitiesInDateRange(ctx, minDate, maxDate)
	if err != nil {
		return nil, err
	}
	junkFoods, err := s.cosmosRepository.GetJunkFoodInDateRange(ctx, minDate, maxDate)
	if err != nil {
		return nil, err
	}

	progress := make(models.ProgressPerDay)
	for _, activity := range activities {
		dayProgress, exists := progress[activity.Date]
		if !exists {
			dayProgress = models.DailyProgress{}
		}

		activityRespone := models.ActivityResponse{
			Type:        activity.Type,
			ElapsedTime: activity.ElapsedTime,
		}

		switch activity.Type {
		case "Walk":
			dayProgress.Walks = append(dayProgress.Walks, activityRespone)
		default:
			dayProgress.Workouts = append(dayProgress.Workouts, activityRespone)
		}

		progress[activity.Date] = dayProgress
	}

	for _, junkFood := range junkFoods {
		dayProgress, exists := progress[junkFood.Date]
		if !exists {
			dayProgress = models.DailyProgress{}
		}

		dayProgress.JunkFoods = append(dayProgress.JunkFoods, models.JunkFoodResponse{
			Type: junkFood.Type,
		})

		progress[junkFood.Date] = dayProgress
	}

	// Ensure every day in range has a key (even if empty)
	start := time.Date(minDate.Year(), minDate.Month(), minDate.Day(), 0, 0, 0, 0, minDate.Location())
	end := time.Date(maxDate.Year(), maxDate.Month(), maxDate.Day(), 0, 0, 0, 0, maxDate.Location())
	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		key := d.Format("2006-01-02")
		if _, exists := progress[key]; !exists {
			progress[key] = models.DailyProgress{}
		}
	}

	return progress, nil
}
