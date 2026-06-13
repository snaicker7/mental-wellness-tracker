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

// GenerateCoping returns contextual coping strategies, breathing exercises, and mindfulness prompts.
func GenerateCoping(ctx context.Context, entry Entry) (CopingResponse, error) {
	stressScore := float64(11 - entry.MoodLevel)

	// 1. Crisis / Safety guardrails keyword filtering
	if isCrisisInput(entry.EntryText) {
		return CopingResponse{
			UserID:      entry.UserID,
			MoodLevel:   entry.MoodLevel,
			StressScore: stressScore,
			IsCrisis:    true,
			Guidance: CopingGuidance{
				MotivationalPrompt:  "Please remember that you do not have to carry this weight alone, and professional care makes a real difference.",
				BreathingExercise:   "For immediate assistance, please call the AASRA helpline at 91-9820466726 or the Vandrevala Foundation at 9999 666 555.",
				MindfulnessActivity: "We recommend pausing and reaching out to a professional mental health counselor, doctor, or a trusted loved one.",
			},
		}, nil
	}

	// 2. Hybrid LLM logic if configured
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey != "" {
		guidance, err := queryGeminiCoping(ctx, apiKey, entry, stressScore)
		if err == nil {
			return CopingResponse{
				UserID:      entry.UserID,
				MoodLevel:   entry.MoodLevel,
				StressScore: stressScore,
				IsCrisis:    false,
				Guidance:    guidance,
			}, nil
		}
	}

	// 3. Rules-based fallback engine
	return CopingResponse{
		UserID:      entry.UserID,
		MoodLevel:   entry.MoodLevel,
		StressScore: stressScore,
		IsCrisis:    false,
		Guidance:    generateFallbackCoping(entry.MoodLevel, stressScore),
	}, nil
}

func isCrisisInput(text string) bool {
	normalized := strings.ToLower(text)
	crisisKeywords := []string{
		"suicide", "suicidal", "kill myself", "harm myself", "self-harm",
		"end my life", "depression", "diagnose", "prescribe", "medicine",
		"pill", "doctor", "clinician", "psychiatrist", "therapist", "die",
	}
	for _, kw := range crisisKeywords {
		if strings.Contains(normalized, kw) {
			return true
		}
	}
	return false
}

func queryGeminiCoping(ctx context.Context, apiKey string, entry Entry, stressScore float64) (CopingGuidance, error) {
	prompt := fmt.Sprintf(`You are an empathetic student mental wellness companion. 
Generate safe coping guidance, breathing exercises, and mindfulness activities for a student with the following context:
- Mood Level: %d/10
- Stress Score: %.1f/10
- Recent Journal Entry: "%s"

Your tone must be warm, supportive, and non-judgmental.
You must respond with a JSON object containing exactly these fields:
{
  "motivational_prompt": "1-2 sentences of warm motivational encouragement.",
  "breathing_exercise": "A simple, quick description of a breathing exercise (e.g. Box Breathing).",
  "mindfulness_activity": "A simple grounding or mindfulness activity."
}
Do not include any other markdown formatting outside of the JSON block.`,
		entry.MoodLevel,
		stressScore,
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
		return CopingGuidance{}, err
	}

	url := "https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash:generateContent"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return CopingGuidance{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Goog-Api-Key", apiKey)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return CopingGuidance{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return CopingGuidance{}, fmt.Errorf("gemini api status %d", resp.StatusCode)
	}

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return CopingGuidance{}, err
	}

	var geminiResp geminiResponse
	if err := json.Unmarshal(respBytes, &geminiResp); err != nil {
		return CopingGuidance{}, err
	}

	if len(geminiResp.Candidates) == 0 || len(geminiResp.Candidates[0].Content.Parts) == 0 {
		return CopingGuidance{}, fmt.Errorf("empty gemini candidates response")
	}

	rawText := geminiResp.Candidates[0].Content.Parts[0].Text
	rawText = strings.TrimPrefix(rawText, "```json")
	rawText = strings.TrimSuffix(rawText, "```")
	rawText = strings.TrimSpace(rawText)

	var guidance CopingGuidance
	if err := json.Unmarshal([]byte(rawText), &guidance); err != nil {
		return CopingGuidance{}, err
	}

	return guidance, nil
}

func generateFallbackCoping(moodLevel int, stressScore float64) CopingGuidance {
	if stressScore > 7 { // High Stress / Low Mood
		return CopingGuidance{
			MotivationalPrompt:  "Take things one step at a time today. You are strong, but it's okay to let yourself rest and breathe.",
			BreathingExercise:   "Try Box Breathing: Inhale for 4s, hold for 4s, exhale for 4s, hold for 4s. Repeat this 3-4 times.",
			MindfulnessActivity: "5-4-3-2-1 Grounding: Identify 5 things you see, 4 you feel, 3 you hear, 2 you smell, and 1 you taste around you.",
		}
	} else if stressScore > 4 { // Medium Stress
		return CopingGuidance{
			MotivationalPrompt:  "You are doing your best, and that is absolutely enough. Take a pause to celebrate small progress.",
			BreathingExercise:   "Try 4-7-8 Breathing: Inhale for 4s, hold for 7s, exhale completely for 8s. Helps reduce immediate anxiety.",
			MindfulnessActivity: "Mindful Listening: Close your eyes and focus on the furthest sound you can hear, then the closest one, for 60 seconds.",
		}
	} else { // Low Stress / Good Mood
		return CopingGuidance{
			MotivationalPrompt:  "It is wonderful to see you in a steady space. Share this positive energy or write down what went well today!",
			BreathingExercise:   "Equal Breathing: Inhale for 5 seconds, and exhale for 5 seconds. Feel the rhythm of air entering and leaving.",
			MindfulnessActivity: "Gratitude list: Name three minor things that brought you comfort or a smile today, and hold onto that feeling.",
		}
	}
}
