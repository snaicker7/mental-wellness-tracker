# Secure daily journal and mood capture

Type: AFK
Blocked by: None - can start immediately
User stories covered: Students can log a daily journal entry, record a mood level, attach optional study/sleep metadata, and view their saved history.

## What to build

Create the first end-to-end writing path for the app: a student can submit a daily journaling entry plus mood log data, and the system validates, stores, and retrieves that record securely. This slice should establish the durable data model and API contract that later AI and reporting features build on.

## Acceptance criteria

- [x] A student can submit a journal entry with mood, energy, sleep, study hours, and exam countdown metadata.
- [x] Inputs are validated and sanitized before persistence.
- [x] Journaling data is stored encrypted at rest and served over TLS-friendly API boundaries.
- [x] A student can retrieve their prior entries for review.
- [x] Unit and integration tests cover the happy path and invalid payload handling.
- [x] The submission UI matches the playful, rounded card theme of [image.png](file:///c:/Users/Surya/Documents/mental-wellness-tracker/image.png) and supports both quick-mood selection (dashboard) and full entry creation (with text input area for journaling and numeric metadata sliders/inputs).
