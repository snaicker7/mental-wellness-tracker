package journal

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
)

func TestExtractTriggers(t *testing.T) {
	tests := []struct {
		name     string
		entry    Entry
		expected []string
	}{
		{
			name: "Sleep and study overload",
			entry: Entry{
				EntryText:         "Normal day",
				MoodLevel:         5,
				EnergyLevel:       5,
				SleepHours:        5.0,
				StudyHours:        9.0,
				ExamCountdownDays: 10,
			},
			expected: []string{"Sleep deprivation", "Academic overload"},
		},
		{
			name: "Exam anxiety and academic keywords",
			entry: Entry{
				EntryText:         "Preparing for the exam tomorrow",
				MoodLevel:         6,
				EnergyLevel:       6,
				SleepHours:        7.0,
				StudyHours:        5.0,
				ExamCountdownDays: 1,
			},
			expected: []string{"Exam anxiety", "Academic pressure"},
		},
		{
			name: "Social and time keywords with severe mood dip",
			entry: Entry{
				EntryText:         "Had a fight with a friend, schedule is very busy",
				MoodLevel:         2,
				EnergyLevel:       3,
				SleepHours:        8.0,
				StudyHours:        4.0,
				ExamCountdownDays: 15,
			},
			expected: []string{"Severe mood dip", "Physical/mental exhaustion", "Social/interpersonal stress", "Time management stress"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractTriggers(tt.entry)
			if len(got) != len(tt.expected) {
				t.Fatalf("expected %d triggers, got %v", len(tt.expected), got)
			}
			for _, exp := range tt.expected {
				if !contains(got, exp) {
					t.Errorf("missing expected trigger: %s in %v", exp, got)
				}
			}
		})
	}
}

func TestGenerateFallback(t *testing.T) {
	entry := Entry{
		ID:                "test-id",
		EntryText:         "I am so tired and have exams coming up",
		MoodLevel:         4,
		EnergyLevel:       3,
		SleepHours:        4.5,
		StudyHours:        10.0,
		ExamCountdownDays: 3,
	}

	triggers := extractTriggers(entry)
	analysis := generateFallback(entry, triggers, 7.0)

	if analysis.EntryID != "test-id" {
		t.Errorf("expected entry id test-id, got %s", analysis.EntryID)
	}
	if analysis.StressScore != 7.0 {
		t.Errorf("expected stress score 7.0, got %f", analysis.StressScore)
	}
	if !contains(analysis.Triggers, "Sleep deprivation") {
		t.Errorf("expected trigger Sleep deprivation, got %v", analysis.Triggers)
	}
	if analysis.Summary == "" || analysis.EmpatheticExplanation == "" {
		t.Errorf("expected fallback summary and explanation to be non-empty")
	}
}

func TestGetAnalysisCachingAndFallback(t *testing.T) {
	tmpDir := t.TempDir()
	storePath := filepath.Join(tmpDir, "journals.json")
	key := mustKey(t)
	store, err := NewFileStore(storePath, key)
	if err != nil {
		t.Fatalf("NewFileStore() error = %v", err)
	}

	service := NewService(store)
	ctx := context.Background()

	// Create test entry
	entry, err := service.CreateEntry(ctx, CreateEntryRequest{
		UserID:            "student-1",
		EntryText:         "Studying all night for midterm quiz.",
		MoodLevel:         3,
		EnergyLevel:       4,
		SleepHours:        5.0,
		StudyHours:        10.0,
		ExamCountdownDays: 4,
	})
	if err != nil {
		t.Fatalf("failed to create entry: %v", err)
	}

	// 1. First analysis request
	analysis1, err := service.GetAnalysis(ctx, entry.ID, "student-1")
	if err != nil {
		t.Fatalf("GetAnalysis error: %v", err)
	}

	if analysis1.EntryID != entry.ID {
		t.Errorf("expected EntryID %s, got %s", entry.ID, analysis1.EntryID)
	}

	// 2. Second request (should hit cache)
	analysis2, err := service.GetAnalysis(ctx, entry.ID, "student-1")
	if err != nil {
		t.Fatalf("GetAnalysis error on second try: %v", err)
	}

	if analysis1.Summary != analysis2.Summary || analysis1.EmpatheticExplanation != analysis2.EmpatheticExplanation {
		t.Errorf("cached result should match original: %+v vs %+v", analysis1, analysis2)
	}
}

func TestHTTPAnalysisEndpoint(t *testing.T) {
	tmpDir := t.TempDir()
	storePath := filepath.Join(tmpDir, "journals.json")
	key := mustKey(t)
	store, err := NewFileStore(storePath, key)
	if err != nil {
		t.Fatalf("NewFileStore() error = %v", err)
	}

	service := NewService(store)
	server := NewServer(store)
	ctx := context.Background()

	// Create entry
	entry, err := service.CreateEntry(ctx, CreateEntryRequest{
		UserID:            "student-1",
		EntryText:         "Busy week ahead, exam prep.",
		MoodLevel:         4,
		EnergyLevel:       5,
		SleepHours:        6.0,
		StudyHours:        8.0,
		ExamCountdownDays: 5,
	})
	if err != nil {
		t.Fatalf("failed to create entry: %v", err)
	}

	// Request with missing parameters
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/analysis?user_id=student-1", nil)
	server.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for missing entry_id, got %d", rec.Code)
	}

	// Successful request
	rec2 := httptest.NewRecorder()
	req2 := httptest.NewRequest(http.MethodGet, "/v1/analysis?user_id=student-1&entry_id="+entry.ID, nil)
	server.ServeHTTP(rec2, req2)
	if rec2.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec2.Code, rec2.Body.String())
	}

	var analysis StressAnalysis
	if err := json.Unmarshal(rec2.Body.Bytes(), &analysis); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if analysis.EntryID != entry.ID {
		t.Errorf("expected entry ID %s, got %s", entry.ID, analysis.EntryID)
	}
	if len(analysis.Triggers) == 0 {
		t.Errorf("expected triggers to be detected")
	}
}
