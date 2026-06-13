# Show mood and stress progress trends

Type: AFK
Blocked by: 1. Secure daily journal and mood capture
User stories covered: Students can see their mood and stress over time and understand whether their well-being is improving or worsening.

## What to build

Build the progress-tracking slice that aggregates saved journal and mood data into visual trends. The slice should present a simple, mobile-friendly view of how stress and mood change across days or weeks so students can spot patterns at a glance.

## Acceptance criteria

- [x] The UI renders mood and stress trends over time from stored entries, aligning with the calendar format and sparklines shown in [image.png](file:///c:/Users/Surya/Documents/mental-wellness-tracker/image.png).
- [x] Trend views work on mobile screens and low-bandwidth connections.
- [x] The visualization is accessible with keyboard navigation and screen-reader labels.
- [x] The data layer exposes the aggregated values needed by the charts.
- [x] Tests verify trend aggregation and chart rendering behavior.
