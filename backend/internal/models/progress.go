package models

type DailyProgress struct {
	Walks     []ActivityResponse `json:"walks"`
	Workouts  []ActivityResponse `json:"workouts"`
	JunkFoods []JunkFoodResponse `json:"junkFoods"`
}

type ProgressPerDay map[string]DailyProgress
