# Validation Commands

Use these commands to validate the backend, frontend, and deployment packaging for Mental Wellness Tracker.

## Quick start for a new session

If you are picking up the project fresh, do this first:

1. Copy [.env.example](.env.example) to [.env](.env) and set `JOURNAL_ENCRYPTION_KEY` plus `DATABASE_URL` if you want PostgreSQL.
2. Run the backend with `& 'C:\Program Files\Go\bin\go.exe' run ./cmd/api` from the repository root.
3. Run the frontend with `Set-Location frontend; npm.cmd install; $env:VITE_API_BASE_URL = 'http://localhost:8080'; npm.cmd run dev`.
4. Open `http://localhost:5173` for the UI and use the API checks below against port `8080`.

## Go backend

Run from the repository root:

```powershell
& 'C:\Program Files\Go\bin\go.exe' test ./...
```

If dependencies or `go.mod` change:

```powershell
& 'C:\Program Files\Go\bin\go.exe' mod tidy
```

Run the API locally:

```powershell
$bytes = New-Object byte[] 32
[System.Security.Cryptography.RandomNumberGenerator]::Fill($bytes)
$env:JOURNAL_ENCRYPTION_KEY = [Convert]::ToBase64String($bytes)
$env:DATABASE_URL = '<postgres-connection-string>'  # optional; omit to use local file storage
& 'C:\Program Files\Go\bin\go.exe' run ./cmd/api
```

## Frontend

Run these commands from the repository root:

```powershell
Set-Location frontend
npm.cmd install
npm.cmd run build
```

Start the frontend dev server:

```powershell
Set-Location frontend
$env:VITE_API_BASE_URL = 'http://localhost:8080'
npm.cmd run dev
```

If you want to preview the production build:

```powershell
Set-Location frontend
npm.cmd run preview
```

## Docker

Build the backend container:

```powershell
docker build -t mental-wellness-tracker-api .
```

Run the container locally:

```powershell
docker run -p 8080:8080 --env-file .env mental-wellness-tracker-api
```

If you want to use PostgreSQL, uncomment `DATABASE_URL` in `.env` and set it to your connection string.

## API checks on 8080

Health check:

```powershell
curl.exe http://localhost:8080/healthz
```

List saved entries:

```powershell
curl.exe "http://localhost:8080/v1/entries?user_id=student-1"
```

Load trend data:

```powershell
curl.exe "http://localhost:8080/v1/trends?user_id=student-1"
```

Create a journal entry:

```powershell
curl.exe -Method Post http://localhost:8080/v1/entries `
	-ContentType application/json `
	-Body '{"user_id":"student-1","entry_text":"Today felt heavy, but I finished my revision plan.","mood_level":4,"energy_level":5,"sleep_hours":6.5,"study_hours":8,"exam_countdown_days":12}'
```

## Quick validation order

1. `& 'C:\Program Files\Go\bin\go.exe' test ./...`
2. `Set-Location frontend; npm.cmd install`
3. `Set-Location frontend; npm.cmd run build`
4. `docker build -t mental-wellness-tracker-api .`
