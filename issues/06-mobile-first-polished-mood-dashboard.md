# Mobile-first polished mood dashboard

Type: AFK
Blocked by: 1. Secure daily journal and mood capture; 2. Detect hidden stress triggers; 3. Generate safe coping guidance; 4. Show mood and stress progress trends; 5. Production hardening: accessibility, bilingual UI, deployment, and CI/CD
User stories covered: Students can open the app on mobile or desktop and immediately understand their mood state through a polished, playful, high-clarity dashboard.

## What to build

Design and implement the primary user-facing dashboard so the finished product feels like a polished consumer wellness app rather than a plain admin screen. The visual language should be mobile-first, friendly, and expressive, with bold typography, rounded cards, soft backgrounds, emoji-driven mood cues, summary tiles, and a calendar or grid view for mood history.

## Acceptance criteria

- [ ] The landing and dashboard screens use a mobile-first card layout with strong visual hierarchy, matching the design in [image.png](file:///c:/Users/Surya/Documents/mental-wellness-tracker/image.png).
- [ ] The interface includes mood summary components, progress/stat cards (Sleep Duration, Stress Indicator, Activity, Therapy, Discipline), and an at-a-glance mood history grid or calendar matching the visual layout.
- [ ] The visual treatment follows the approved playful wellness style: bold headings, rounded surfaces, soft gradients, and friendly iconography or emoji cues.
- [ ] The layout remains readable and usable on small screens without losing the overall design language.
- [ ] The design stays accessible, with sufficient contrast, keyboard support, and sensible semantics for the visual components.
- [ ] The final frontend feels cohesive across the main user flows after the other slices are complete.

## UI Design Guidelines (from [image.png](file:///c:/Users/Surya/Documents/mental-wellness-tracker/image.png))
- **Landing Screen**: Features friendly, playful colored shapes with expressions and a bold call-to-action ("Not Sure About Your Mood? Let Us Help!").
- **Home Dashboard**:
  - Welcome greeting for the student.
  - Interactive mood selection row (Happy, Angry, Sleepy, Bored).
  - Metrics grid including:
    - **Sleep Duration** card (hours and minutes with daily bar graphs).
    - **Stress Indicator** card (level indicator with sparklines/trends).
    - **Quiz Card** (e.g. Yes/No questions to help assess state).
  - Bottom navigation bar with icons for quick navigation.
- **Mood Calendar**:
  - A grid monthly view where each day displays the corresponding color-coded mood emoji.
  - A summary card ("Monthly Mood Summary") showcasing the dominant mood with a description.
  - Supplementary tracking metrics (Activity/Steps, Therapy sessions, Discipline/Focus score).