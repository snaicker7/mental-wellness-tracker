package journal

import "time"

const (
	AES256KeySize = 32
	maxEntryText   = 4000
	maxUserIDLen   = 128
)

type Entry struct {
	ID                string    `json:"id"`
	UserID            string    `json:"user_id"`
	EntryText         string    `json:"entry_text"`
	MoodLevel         int       `json:"mood_level"`
	EnergyLevel       int       `json:"energy_level"`
	SleepHours        float64   `json:"sleep_hours"`
	StudyHours        float64   `json:"study_hours"`
	ExamCountdownDays int       `json:"exam_countdown_days"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type CreateEntryRequest struct {
	UserID            string  `json:"user_id"`
	EntryText         string  `json:"entry_text"`
	MoodLevel         int     `json:"mood_level"`
	EnergyLevel       int     `json:"energy_level"`
	SleepHours        float64 `json:"sleep_hours"`
	StudyHours        float64 `json:"study_hours"`
	ExamCountdownDays int     `json:"exam_countdown_days"`
}

type TrendPoint struct {
	Date        string  `json:"date"`
	EntryCount  int     `json:"entry_count"`
	MoodAverage float64 `json:"mood_average"`
	StressAverage float64 `json:"stress_average"`
	EnergyAverage float64 `json:"energy_average"`
}

type TrendResponse struct {
	UserID string       `json:"user_id"`
	Points []TrendPoint `json:"points"`
}

type StressAnalysis struct {
	EntryID                string   `json:"entry_id"`
	StressScore            float64  `json:"stress_score"`
	Triggers               []string `json:"triggers"`
	Summary                string   `json:"summary"`
	EmpatheticExplanation  string   `json:"empathetic_explanation"`
}

type CopingGuidance struct {
	MotivationalPrompt  string `json:"motivational_prompt"`
	BreathingExercise   string `json:"breathing_exercise"`
	MindfulnessActivity string `json:"mindfulness_activity"`
}

type CopingResponse struct {
	UserID      string         `json:"user_id"`
	MoodLevel   int            `json:"mood_level"`
	StressScore float64        `json:"stress_score"`
	Guidance    CopingGuidance `json:"guidance"`
	IsCrisis    bool           `json:"is_crisis"`
}

