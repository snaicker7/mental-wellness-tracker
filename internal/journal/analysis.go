package journal

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// AnalyzeStress performs rules-based trigger extraction and optionally queries Gemini API for empathetic interpretation.
func AnalyzeStress(ctx context.Context, entry Entry) (StressAnalysis, error) {
	// 1. Rule-based trigger extraction
	triggers := extractTriggers(entry)

	// 2. Base Stress Score calculation
	stressScore := float64(11 - entry.MoodLevel)

	// 3. Attempt LLM Interpretation if key is present
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey != "" {
		analysis, err := queryGemini(ctx, apiKey, entry, triggers, stressScore)
		if err == nil {
			return analysis, nil
		}
		// Log or handle error, then fallback
	}

	// 4. Fallback generation
	return generateFallback(entry, triggers, stressScore), nil
}

func extractTriggers(entry Entry) []string {
	var triggers []string

	// Metric triggers
	if entry.SleepHours < 6.0 {
		triggers = append(triggers, "Sleep deprivation")
	}
	if entry.StudyHours > 8.0 {
		triggers = append(triggers, "Academic overload")
	}
	if entry.ExamCountdownDays > 0 && entry.ExamCountdownDays <= 7 {
		triggers = append(triggers, "Exam anxiety")
	}
	if entry.MoodLevel <= 3 {
		triggers = append(triggers, "Severe mood dip")
	}
	if entry.EnergyLevel <= 3 {
		triggers = append(triggers, "Physical/mental exhaustion")
	}

	// Keyword-based triggers
	text := strings.ToLower(entry.EntryText)
	if containsAny(text, "exam", "test", "quiz", "revision", "syllabus", "grade", "results", "marks", "fail", "study", "prep") {
		if !contains(triggers, "Academic pressure") {
			triggers = append(triggers, "Academic pressure")
		}
	}
	if containsAny(text, "sleep", "tired", "awake", "insomnia", "exhausted", "nightmare", "restless") {
		if !contains(triggers, "Sleep disruption") {
			triggers = append(triggers, "Sleep disruption")
		}
	}
	if containsAny(text, "friend", "argument", "fight", "lonely", "alone", "social", "parents", "family", "relationship") {
		if !contains(triggers, "Social/interpersonal stress") {
			triggers = append(triggers, "Social/interpersonal stress")
		}
	}
	if containsAny(text, "deadline", "schedule", "time", "hurry", "late", "busy", "overwhelmed") {
		if !contains(triggers, "Time management stress") {
			triggers = append(triggers, "Time management stress")
		}
	}

	return triggers
}

func contains(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

func containsAny(text string, keywords ...string) bool {
	for _, kw := range keywords {
		if strings.Contains(text, kw) {
			return true
		}
	}
	return false
}

type geminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

type geminiParsedResult struct {
	Summary               string `json:"summary"`
	EmpatheticExplanation string `json:"empathetic_explanation"`
}

func queryGemini(ctx context.Context, apiKey string, entry Entry, triggers []string, stressScore float64) (StressAnalysis, error) {
	prompt := fmt.Sprintf(`You are an empathetic student mental wellness assistant.
Analyze the student's journal entry and mood metrics to identify hidden triggers, repeated patterns, and offer supportive guidance.

Input metrics:
- Mood: %d/10
- Energy: %d/10
- Sleep: %.1f hours
- Study: %.1f hours
- Exam Countdown: %d days
- Detected rule triggers: %s

Journal entry text:
"%s"

You must respond with a JSON object containing exactly these fields:
{
  "summary": "1-2 sentences summarizing the overall stress and mood state.",
  "empathetic_explanation": "Empathetic analysis of hidden triggers or patterns and gentle, actionable support."
}
Do not include any other markdown formatting outside of the JSON block.`,
		entry.MoodLevel,
		entry.EnergyLevel,
		entry.SleepHours,
		entry.StudyHours,
		entry.ExamCountdownDays,
		strings.Join(triggers, ", "),
		entry.EntryText,
	)

	reqBody, err := json.Marshal(map[string]any{
		"contents": []any{
			map[string]any{
				"parts": []any{
					map[string]any{
						"text": prompt,
					},
				},
			},
		},
		"generationConfig": map[string]any{
			"responseMimeType": "application/json",
		},
	})
	if err != nil {
		return StressAnalysis{}, err
	}

	url := "https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash:generateContent"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return StressAnalysis{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Goog-Api-Key", apiKey)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return StressAnalysis{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return StressAnalysis{}, fmt.Errorf("gemini api returned status %d", resp.StatusCode)
	}

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return StressAnalysis{}, err
	}

	var geminiResp geminiResponse
	if err := json.Unmarshal(respBytes, &geminiResp); err != nil {
		return StressAnalysis{}, err
	}

	if len(geminiResp.Candidates) == 0 || len(geminiResp.Candidates[0].Content.Parts) == 0 {
		return StressAnalysis{}, fmt.Errorf("empty gemini response candidates")
	}

	rawText := geminiResp.Candidates[0].Content.Parts[0].Text
	// Sometimes LLMs return json wrapped in markdown code blocks, strip them if present
	rawText = strings.TrimPrefix(rawText, "```json")
	rawText = strings.TrimSuffix(rawText, "```")
	rawText = strings.TrimSpace(rawText)

	var parsed geminiParsedResult
	if err := json.Unmarshal([]byte(rawText), &parsed); err != nil {
		return StressAnalysis{}, err
	}

	return StressAnalysis{
		EntryID:               entry.ID,
		StressScore:           stressScore,
		Triggers:              triggers,
		Summary:               parsed.Summary,
		EmpatheticExplanation: parsed.EmpatheticExplanation,
	}, nil
}

func generateFallback(entry Entry, triggers []string, stressScore float64) StressAnalysis {
	var summary string
	var explanation string

	hasSleepIssue := contains(triggers, "Sleep deprivation") || contains(triggers, "Sleep disruption")
	hasAcademicIssue := contains(triggers, "Academic overload") || contains(triggers, "Academic pressure") || contains(triggers, "Exam anxiety")
	hasSocialIssue := contains(triggers, "Social/interpersonal stress")
	hasTimeIssue := contains(triggers, "Time management stress")

	if hasSleepIssue && hasAcademicIssue {
		summary = "Academic demands appear to be affecting your sleep, compounding your overall stress."
		explanation = "Balancing long study hours with less than 6 hours of sleep can lead to a cycle of fatigue. Prioritizing consistent rest can improve both your focus and memory retention."
	} else if hasAcademicIssue {
		summary = "You're experiencing significant academic pressure, likely related to studying or upcoming exams."
		explanation = "Preparation anxiety is very common. Breaking study materials into smaller chunks and incorporating short, active breaks can help keep your stress manageable."
	} else if hasSleepIssue {
		summary = "Your logs show signs of physical exhaustion and disrupted sleep cycles."
		explanation = "A good night's rest is foundational for mental resilience. Try establishing a screen-free wind-down routine 30 minutes before bed to signal your mind to rest."
	} else if hasSocialIssue {
		summary = "Social or family dynamics might be contributing to a dip in your mood."
		explanation = "Interpersonal stress can feel incredibly draining. Reach out to a trusted peer, or give yourself space to process these interactions without self-judgment."
	} else if hasTimeIssue {
		summary = "Time constraints and packed schedules seem to be causing a sense of being overwhelmed."
		explanation = "When tasks pile up, time pressure mounts. Try mapping out just the top three priorities for tomorrow and focus on one single task at a time."
	} else {
		summary = "You are checking in on your wellness, keeping track of your daily emotional patterns."
		explanation = "It's normal to have fluctuations in energy and mood. Taking a moment to write down how you feel is a great habit for long-term emotional awareness."
	}

	// Safe formatting: trim whitespace
	summary = strings.TrimSpace(summary)
	explanation = strings.TrimSpace(explanation)

	return StressAnalysis{
		EntryID:               entry.ID,
		StressScore:           stressScore,
		Triggers:              triggers,
		Summary:               summary,
		EmpatheticExplanation: explanation,
	}
}
