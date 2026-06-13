package journal

import (
	"context"
	"path/filepath"
	"strings"
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
	valid := CreateEntryRequest{
		UserID:            "student-1",
		EntryText:         "I took a focused break.",
		MoodLevel:         5,
		EnergyLevel:       5,
		SleepHours:        7,
		StudyHours:        6,
		ExamCountdownDays: 10,
	}

	tests := []struct {
		name   string
		mutate func(*CreateEntryRequest)
	}{
		{"missing user", func(r *CreateEntryRequest) { r.UserID = "" }},
		{"invalid user utf8", func(r *CreateEntryRequest) { r.UserID = string([]byte{0xff}) }},
		{"user too long", func(r *CreateEntryRequest) { r.UserID = strings.Repeat("a", maxUserIDLen+1) }},
		{"missing entry", func(r *CreateEntryRequest) { r.EntryText = "" }},
		{"invalid entry utf8", func(r *CreateEntryRequest) { r.EntryText = string([]byte{0xff}) }},
		{"entry too long", func(r *CreateEntryRequest) { r.EntryText = strings.Repeat("a", maxEntryText+1) }},
		{"low mood", func(r *CreateEntryRequest) { r.MoodLevel = 0 }},
		{"high mood", func(r *CreateEntryRequest) { r.MoodLevel = 11 }},
		{"low energy", func(r *CreateEntryRequest) { r.EnergyLevel = 0 }},
		{"high energy", func(r *CreateEntryRequest) { r.EnergyLevel = 11 }},
		{"low sleep", func(r *CreateEntryRequest) { r.SleepHours = -0.1 }},
		{"high sleep", func(r *CreateEntryRequest) { r.SleepHours = 24.1 }},
		{"low study", func(r *CreateEntryRequest) { r.StudyHours = -0.1 }},
		{"high study", func(r *CreateEntryRequest) { r.StudyHours = 24.1 }},
		{"negative countdown", func(r *CreateEntryRequest) { r.ExamCountdownDays = -1 }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := valid
			tt.mutate(&request)
			if _, err := validateAndSanitize(request); err == nil {
				t.Fatal("expected validation error")
			}
		})
	}
}

func TestServiceRejectsInvalidUserQueries(t *testing.T) {
	store, err := NewFileStore(filepath.Join(t.TempDir(), "journals.json"), mustKey(t))
	if err != nil {
		t.Fatalf("NewFileStore() error = %v", err)
	}

	service := NewService(store)
	ctx := context.Background()

	if _, err := service.ListEntries(ctx, " "); err == nil {
		t.Fatal("expected ListEntries to reject empty user id")
	}
	if _, err := service.GetTrends(ctx, " "); err == nil {
		t.Fatal("expected GetTrends to reject empty user id")
	}
	if _, err := service.GetAnalysis(ctx, "", "student-1"); err == nil {
		t.Fatal("expected GetAnalysis to reject empty entry id")
	}
	if _, err := service.GetAnalysis(ctx, "entry-1", " "); err == nil {
		t.Fatal("expected GetAnalysis to reject empty user id")
	}
	if _, err := service.GetAnalysis(ctx, "missing-entry", "student-1"); err == nil {
		t.Fatal("expected GetAnalysis to reject missing entry")
	}
	if _, err := service.GetCopingGuidance(ctx, " "); err == nil {
		t.Fatal("expected GetCopingGuidance to reject empty user id")
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
