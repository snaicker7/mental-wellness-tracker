# Detect hidden stress triggers

Type: AFK
Blocked by: 1. Secure daily journal and mood capture
User stories covered: Students can understand what is causing stress, recognize recurring emotional patterns, and get a plain-language summary of risk factors.

## What to build

Add the analysis path that turns journal entries and mood logs into stress signals, trigger hypotheses, and recurring themes. The slice should combine rule-based logic with an AI-assisted interpretation layer so the app can explain patterns in a way students can actually use.

## Acceptance criteria

- [x] The system produces a stress analysis summary from a saved journal entry and mood log.
- [x] Hidden triggers and repeated patterns are identified in clear, empathetic language.
- [x] The analysis flow supports a hybrid rules-plus-LLM approach.
- [x] Cached or repeated analysis requests avoid unnecessary reprocessing when inputs have not changed.
- [x] Tests cover trigger extraction, fallback behavior, and safe output formatting.
- [x] The stress analysis summary and hidden triggers are rendered in a clean UI component (e.g. details page or modal) styled consistently with the Stress Indicator widget in [image.png](file:///c:/Users/Surya/Documents/mental-wellness-tracker/image.png).
