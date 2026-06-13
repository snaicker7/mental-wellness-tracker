package journal

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"
	"unicode"
	"unicode/utf8"
)

type Repository interface {
	Create(ctx context.Context, entry Entry) (Entry, error)
	ListByUser(ctx context.Context, userID string) ([]Entry, error)
}

type Service struct {
	repo          Repository
	cacheMu       sync.RWMutex
	analysisCache map[string]StressAnalysis
}

func NewService(repo Repository) *Service {
	return &Service{
		repo:          repo,
		analysisCache: make(map[string]StressAnalysis),
	}
}

func (s *Service) CreateEntry(ctx context.Context, request CreateEntryRequest) (Entry, error) {
	cleaned, err := validateAndSanitize(request)
	if err != nil {
		return Entry{}, err
	}

	now := time.Now().UTC()
	entry := Entry{
		ID:                newID(),
		UserID:            cleaned.UserID,
		EntryText:         cleaned.EntryText,
		MoodLevel:         cleaned.MoodLevel,
		EnergyLevel:       cleaned.EnergyLevel,
		SleepHours:        cleaned.SleepHours,
		StudyHours:        cleaned.StudyHours,
		ExamCountdownDays: cleaned.ExamCountdownDays,
		CreatedAt:         now,
		UpdatedAt:         now,
	}

	return s.repo.Create(ctx, entry)
}

func (s *Service) ListEntries(ctx context.Context, userID string) ([]Entry, error) {
	userID = sanitizeUserID(userID)
	if userID == "" {
		return nil, errors.New("user_id is required")
	}

	return s.repo.ListByUser(ctx, userID)
}

func (s *Service) GetTrends(ctx context.Context, userID string) (TrendResponse, error) {
	userID = sanitizeUserID(userID)
	if userID == "" {
		return TrendResponse{}, errors.New("user_id is required")
	}

	entries, err := s.repo.ListByUser(ctx, userID)
	if err != nil {
		return TrendResponse{}, err
	}

	trendMap := make(map[string]*trendAccumulator)
	for _, entry := range entries {
		dateKey := entry.CreatedAt.UTC().Format("2006-01-02")
		bucket, exists := trendMap[dateKey]
		if !exists {
			bucket = &trendAccumulator{}
			trendMap[dateKey] = bucket
		}
		bucket.add(entry)
	}

	points := make([]TrendPoint, 0, len(trendMap))
	for dateKey, bucket := range trendMap {
		points = append(points, bucket.point(dateKey))
	}
	sort.Slice(points, func(i, j int) bool {
		return points[i].Date < points[j].Date
	})

	return TrendResponse{UserID: userID, Points: points}, nil
}

type trendAccumulator struct {
	count int
	moodTotal float64
	stressTotal float64
	energyTotal float64
}

func (a *trendAccumulator) add(entry Entry) {
	a.count++
	a.moodTotal += float64(entry.MoodLevel)
	a.energyTotal += float64(entry.EnergyLevel)
	a.stressTotal += derivedStressScore(entry)
}

func (a *trendAccumulator) point(dateKey string) TrendPoint {
	if a.count == 0 {
		return TrendPoint{Date: dateKey}
	}

	count := float64(a.count)
	return TrendPoint{
		Date:          dateKey,
		EntryCount:    a.count,
		MoodAverage:   a.moodTotal / count,
		StressAverage: a.stressTotal / count,
		EnergyAverage: a.energyTotal / count,
	}
}

func derivedStressScore(entry Entry) float64 {
	return float64(11 - entry.MoodLevel)
}

func validateAndSanitize(request CreateEntryRequest) (CreateEntryRequest, error) {
	request.UserID = sanitizeUserID(request.UserID)
	request.EntryText = sanitizeEntryText(request.EntryText)

	if request.UserID == "" {
		return CreateEntryRequest{}, errors.New("user_id is required")
	}
	if request.EntryText == "" {
		return CreateEntryRequest{}, errors.New("entry_text is required")
	}
	if len(request.EntryText) > maxEntryText {
		return CreateEntryRequest{}, fmt.Errorf("entry_text must be at most %d characters", maxEntryText)
	}
	if request.MoodLevel < 1 || request.MoodLevel > 10 {
		return CreateEntryRequest{}, errors.New("mood_level must be between 1 and 10")
	}
	if request.EnergyLevel < 1 || request.EnergyLevel > 10 {
		return CreateEntryRequest{}, errors.New("energy_level must be between 1 and 10")
	}
	if request.SleepHours < 0 || request.SleepHours > 24 {
		return CreateEntryRequest{}, errors.New("sleep_hours must be between 0 and 24")
	}
	if request.StudyHours < 0 || request.StudyHours > 24 {
		return CreateEntryRequest{}, errors.New("study_hours must be between 0 and 24")
	}
	if request.ExamCountdownDays < 0 {
		return CreateEntryRequest{}, errors.New("exam_countdown_days must be zero or positive")
	}

	return request, nil
}

func sanitizeUserID(value string) string {
	value = strings.TrimSpace(value)
	if value == "" || len(value) > maxUserIDLen || !utf8.ValidString(value) {
		return ""
	}

	return strings.Map(func(r rune) rune {
		if unicode.IsControl(r) {
			return -1
		}
		return r
	}, value)
}

func sanitizeEntryText(value string) string {
	value = strings.TrimSpace(value)
	if value == "" || !utf8.ValidString(value) {
		return ""
	}

	return strings.Map(func(r rune) rune {
		if unicode.IsControl(r) && r != '\n' && r != '\r' && r != '\t' {
			return -1
		}
		return r
	}, value)
}

func newID() string {
	var randomBytes [16]byte
	if _, err := rand.Read(randomBytes[:]); err != nil {
		return hex.EncodeToString([]byte(time.Now().UTC().Format(time.RFC3339Nano)))
	}

	return hex.EncodeToString(randomBytes[:])
}

func (s *Service) GetAnalysis(ctx context.Context, entryID string, userID string) (StressAnalysis, error) {
	entryID = strings.TrimSpace(entryID)
	userID = sanitizeUserID(userID)
	if entryID == "" {
		return StressAnalysis{}, errors.New("entry_id is required")
	}
	if userID == "" {
		return StressAnalysis{}, errors.New("user_id is required")
	}

	s.cacheMu.RLock()
	analysis, cached := s.analysisCache[entryID]
	s.cacheMu.RUnlock()
	if cached {
		return analysis, nil
	}

	entries, err := s.repo.ListByUser(ctx, userID)
	if err != nil {
		return StressAnalysis{}, err
	}

	var foundEntry *Entry
	for i := range entries {
		if entries[i].ID == entryID {
			foundEntry = &entries[i]
			break
		}
	}

	if foundEntry == nil {
		return StressAnalysis{}, errors.New("journal entry not found")
	}

	analysis, err = AnalyzeStress(ctx, *foundEntry)
	if err != nil {
		return StressAnalysis{}, err
	}

	s.cacheMu.Lock()
	s.analysisCache[entryID] = analysis
	s.cacheMu.Unlock()

	return analysis, nil
}

func (s *Service) GetCopingGuidance(ctx context.Context, userID string) (CopingResponse, error) {
	userID = sanitizeUserID(userID)
	if userID == "" {
		return CopingResponse{}, errors.New("user_id is required")
	}

	entries, err := s.repo.ListByUser(ctx, userID)
	if err != nil {
		return CopingResponse{}, err
	}

	if len(entries) == 0 {
		return CopingResponse{
			UserID:      userID,
			MoodLevel:   5,
			StressScore: 5.0,
			IsCrisis:    false,
			Guidance: CopingGuidance{
				MotivationalPrompt:  "Welcome to your wellness dashboard! Add a journal entry to get personalized encouragement and coping activities.",
				BreathingExercise:   "Take a slow breath in for 4 seconds, and release it gently for 4 seconds. Focus on the physical sensation.",
				MindfulnessActivity: "Look around you and notice three things: a color that comforts you, a texture you can touch, and a sound you can hear.",
			},
		}, nil
	}

	latestEntry := entries[len(entries)-1]
	return GenerateCoping(ctx, latestEntry)
}

