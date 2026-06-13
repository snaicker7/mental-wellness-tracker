package journal

import (
	"context"
	"path/filepath"
	"testing"
	"time"
)

func TestValidateAndSanitize(t *testing.T) {
	request := CreateEntryRequest{
		UserID:            "  student-1  ",
		EntryText:         "  I felt calmer after a walk\x00\n",
		MoodLevel:         7,
		EnergyLevel:       6,
		SleepHours:        7.5,
		StudyHours:        4,
		ExamCountdownDays: 18,
	}

	cleaned, err := validateAndSanitize(request)
	if err != nil {
		t.Fatalf("validateAndSanitize() error = %v", err)
	}
	if cleaned.UserID != "student-1" {
		t.Fatalf("expected trimmed user id, got %q", cleaned.UserID)
	}
	if cleaned.EntryText != "I felt calmer after a walk" {
		t.Fatalf("expected sanitized entry text, got %q", cleaned.EntryText)
	}
}

func TestValidateAndSanitizeRejectsInvalidPayload(t *testing.T) {
	_, err := validateAndSanitize(CreateEntryRequest{
		UserID:    "student-1",
		EntryText: "",
		MoodLevel: 11,
	})
	if err == nil {
		t.Fatal("expected validation error")
	}
}

func TestGetTrendsAggregation(t *testing.T) {
	tmpDir := t.TempDir()
	storePath := filepath.Join(tmpDir, "journals.json")
	key := mustKey(t)
	store, err := NewFileStore(storePath, key)
	if err != nil {
		t.Fatalf("NewFileStore() error = %v", err)
	}

	service := NewService(store)
	ctx := context.Background()

	// 1. Day 1: two entries
	day1 := time.Date(2026, 6, 10, 10, 0, 0, 0, time.UTC)
	_, err = store.Create(ctx, Entry{
		ID:          "e1",
		UserID:      "student-1",
		EntryText:   "Entry 1",
		MoodLevel:   4, // Stress: 7
		EnergyLevel: 6,
		CreatedAt:   day1,
	})
	if err != nil {
		t.Fatalf("store.Create error: %v", err)
	}

	_, err = store.Create(ctx, Entry{
		ID:          "e2",
		UserID:      "student-1",
		EntryText:   "Entry 2",
		MoodLevel:   6, // Stress: 5
		EnergyLevel: 8,
		CreatedAt:   day1.Add(2 * time.Hour),
	})
	if err != nil {
		t.Fatalf("store.Create error: %v", err)
	}

	// 2. Day 2: one entry
	day2 := time.Date(2026, 6, 11, 14, 0, 0, 0, time.UTC)
	_, err = store.Create(ctx, Entry{
		ID:          "e3",
		UserID:      "student-1",
		EntryText:   "Entry 3",
		MoodLevel:   8, // Stress: 3
		EnergyLevel: 9,
		CreatedAt:   day2,
	})
	if err != nil {
		t.Fatalf("store.Create error: %v", err)
	}

	res, err := service.GetTrends(ctx, "student-1")
	if err != nil {
		t.Fatalf("GetTrends failed: %v", err)
	}

	if len(res.Points) != 2 {
		t.Fatalf("expected 2 trend points, got %d", len(res.Points))
	}

	// Validate Day 1 averages: Mood = (4+6)/2 = 5, Stress = (7+5)/2 = 6, Energy = (6+8)/2 = 7
	pt1 := res.Points[0]
	if pt1.Date != "2026-06-10" {
		t.Errorf("expected date 2026-06-10, got %s", pt1.Date)
	}
	if pt1.MoodAverage != 5.0 {
		t.Errorf("expected mood avg 5.0, got %f", pt1.MoodAverage)
	}
	if pt1.StressAverage != 6.0 {
		t.Errorf("expected stress avg 6.0, got %f", pt1.StressAverage)
	}
	if pt1.EnergyAverage != 7.0 {
		t.Errorf("expected energy avg 7.0, got %f", pt1.EnergyAverage)
	}

	// Validate Day 2 averages: Mood = 8, Stress = 3, Energy = 9
	pt2 := res.Points[1]
	if pt2.Date != "2026-06-11" {
		t.Errorf("expected date 2026-06-11, got %s", pt2.Date)
	}
	if pt2.MoodAverage != 8.0 {
		t.Errorf("expected mood avg 8.0, got %f", pt2.MoodAverage)
	}
	if pt2.StressAverage != 3.0 {
		t.Errorf("expected stress avg 3.0, got %f", pt2.StressAverage)
	}
	if pt2.EnergyAverage != 9.0 {
		t.Errorf("expected energy avg 9.0, got %f", pt2.EnergyAverage)
	}
}
