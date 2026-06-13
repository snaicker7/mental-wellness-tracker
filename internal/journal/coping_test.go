package journal

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
)

func TestCrisisKeywordFilter(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		isCrisis bool
	}{
		{"Clean input", "Feeling tired but fine.", false},
		{"Crisis suicide", "I have suicidal thoughts today.", true},
		{"Crisis self harm", "Sometimes I just want to harm myself.", true},
		{"Crisis medical", "Need doctor to prescribe medicine", true},
		{"Crisis depression", "Struggling with severe depression and pills", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isCrisisInput(tt.text)
			if got != tt.isCrisis {
				t.Errorf("expected isCrisisInput(%q) = %v, got %v", tt.text, tt.isCrisis, got)
			}
		})
	}
}

func TestGenerateCoping(t *testing.T) {
	// Normal entry
	entryNormal := Entry{
		UserID:            "student-1",
		EntryText:         "Study was good but sleep was low.",
		MoodLevel:         4,
		EnergyLevel:       5,
		SleepHours:        5.0,
		StudyHours:        7.0,
		ExamCountdownDays: 10,
	}

	resNormal, err := GenerateCoping(context.Background(), entryNormal)
	if err != nil {
		t.Fatalf("GenerateCoping error: %v", err)
	}

	if resNormal.IsCrisis {
		t.Error("expected IsCrisis to be false for normal entry")
	}
	if resNormal.Guidance.MotivationalPrompt == "" {
		t.Error("expected MotivationalPrompt to be populated")
	}
	if resNormal.Guidance.BreathingExercise == "" {
		t.Error("expected BreathingExercise to be populated")
	}
	if resNormal.Guidance.MindfulnessActivity == "" {
		t.Error("expected MindfulnessActivity to be populated")
	}

	// Crisis entry
	entryCrisis := Entry{
		UserID:            "student-1",
		EntryText:         "Feeling hopeless, want to end my life.",
		MoodLevel:         1,
		EnergyLevel:       1,
		SleepHours:        2.0,
		StudyHours:        0.0,
		ExamCountdownDays: 20,
	}

	resCrisis, err := GenerateCoping(context.Background(), entryCrisis)
	if err != nil {
		t.Fatalf("GenerateCoping error: %v", err)
	}

	if !resCrisis.IsCrisis {
		t.Error("expected IsCrisis to be true for crisis entry")
	}
	// Check that helplines are returned in breathing or mindfulness or motivational fields
	if !containsCrisisInfo(resCrisis.Guidance) {
		t.Errorf("expected guidance to contain helpline contacts, got: %+v", resCrisis.Guidance)
	}
}

func containsCrisisInfo(g CopingGuidance) bool {
	info := g.MotivationalPrompt + " " + g.BreathingExercise + " " + g.MindfulnessActivity
	return (containsSub(info, "AASRA") || containsSub(info, "Vandrevala") || containsSub(info, "helpline"))
}

func containsSub(str, sub string) bool {
	return (strings.Contains(str, sub) || strings.Contains(strings.ToLower(str), strings.ToLower(sub)))
}

func TestGetCopingGuidanceEmptyState(t *testing.T) {
	tmpDir := t.TempDir()
	storePath := filepath.Join(tmpDir, "journals.json")
	key := mustKey(t)
	store, err := NewFileStore(storePath, key)
	if err != nil {
		t.Fatalf("NewFileStore() error = %v", err)
	}

	service := NewService(store)
	res, err := service.GetCopingGuidance(context.Background(), "student-empty")
	if err != nil {
		t.Fatalf("GetCopingGuidance error: %v", err)
	}

	if res.IsCrisis {
		t.Error("expected IsCrisis to be false")
	}
	if res.Guidance.MotivationalPrompt == "" {
		t.Error("expected default MotivationalPrompt for empty state")
	}
}

func TestHTTPCopingEndpoint(t *testing.T) {
	tmpDir := t.TempDir()
	storePath := filepath.Join(tmpDir, "journals.json")
	key := mustKey(t)
	store, err := NewFileStore(storePath, key)
	if err != nil {
		t.Fatalf("NewFileStore() error = %v", err)
	}

	server := NewServer(store)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/coping?user_id=student-1", nil)
	server.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	var response CopingResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("decode coping response error: %v", err)
	}

	if response.UserID != "student-1" {
		t.Errorf("expected student-1, got %s", response.UserID)
	}
}
