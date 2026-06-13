package journal

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCreateAndListEntries(t *testing.T) {
	tmpDir := t.TempDir()
	storePath := filepath.Join(tmpDir, "journals.json")
	key := mustKey(t)
	store, err := NewFileStore(storePath, key)
	if err != nil {
		t.Fatalf("NewFileStore() error = %v", err)
	}

	server := NewServer(store)

	createBody := `{
		"user_id": "student-1",
		"entry_text": "  The mock exam felt overwhelming, but I took a short walk.  ",
		"mood_level": 4,
		"energy_level": 5,
		"sleep_hours": 6.5,
		"study_hours": 8,
		"exam_countdown_days": 12
	}`
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/v1/entries", strings.NewReader(createBody))
	server.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", recorder.Code, recorder.Body.String())
	}

	var created Entry
	if err := json.Unmarshal(recorder.Body.Bytes(), &created); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if created.EntryText != "The mock exam felt overwhelming, but I took a short walk." {
		t.Fatalf("expected trimmed entry text, got %q", created.EntryText)
	}

	data, err := os.ReadFile(storePath)
	if err != nil {
		t.Fatalf("read encrypted file: %v", err)
	}
	if bytes.Contains(data, []byte("overwhelming")) || bytes.Contains(data, []byte("student-1")) {
		t.Fatalf("expected encrypted at-rest storage, got plaintext in %s", storePath)
	}

	listRecorder := httptest.NewRecorder()
	listRequest := httptest.NewRequest(http.MethodGet, "/v1/entries?user_id=student-1", nil)
	server.ServeHTTP(listRecorder, listRequest)

	if listRecorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", listRecorder.Code, listRecorder.Body.String())
	}

	var response struct {
		Entries []Entry `json:"entries"`
	}
	if err := json.Unmarshal(listRecorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("decode list response: %v", err)
	}
	if len(response.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(response.Entries))
	}
	if response.Entries[0].UserID != "student-1" {
		t.Fatalf("unexpected user id: %q", response.Entries[0].UserID)
	}

	trendRecorder := httptest.NewRecorder()
	trendRequest := httptest.NewRequest(http.MethodGet, "/v1/trends?user_id=student-1", nil)
	server.ServeHTTP(trendRecorder, trendRequest)

	if trendRecorder.Code != http.StatusOK {
		t.Fatalf("expected 200 for trends, got %d: %s", trendRecorder.Code, trendRecorder.Body.String())
	}

	var trendResponse TrendResponse
	if err := json.Unmarshal(trendRecorder.Body.Bytes(), &trendResponse); err != nil {
		t.Fatalf("decode trend response: %v", err)
	}
	if trendResponse.UserID != "student-1" {
		t.Fatalf("unexpected trend user id: %q", trendResponse.UserID)
	}
	if len(trendResponse.Points) != 1 {
		t.Fatalf("expected 1 trend point, got %d", len(trendResponse.Points))
	}
	if trendResponse.Points[0].MoodAverage != 4 || trendResponse.Points[0].StressAverage != 7 {
		t.Fatalf("unexpected trend values: %+v", trendResponse.Points[0])
	}
}

func TestCreateEntryRejectsInvalidPayload(t *testing.T) {
	store, err := NewFileStore(filepath.Join(t.TempDir(), "journals.json"), mustKey(t))
	if err != nil {
		t.Fatalf("NewFileStore() error = %v", err)
	}

	server := NewServer(store)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/v1/entries", strings.NewReader(`{"user_id":"student-1","entry_text":"","mood_level":0}`))
	server.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", recorder.Code)
	}
}

func TestTrendsRejectMissingUserID(t *testing.T) {
	store, err := NewFileStore(filepath.Join(t.TempDir(), "journals.json"), mustKey(t))
	if err != nil {
		t.Fatalf("NewFileStore() error = %v", err)
	}

	server := NewServer(store)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/v1/trends", nil)
	server.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", recorder.Code)
	}
}

func TestServerAppliesSecurityHeadersAndRestrictedCORS(t *testing.T) {
	store, err := NewFileStore(filepath.Join(t.TempDir(), "journals.json"), mustKey(t))
	if err != nil {
		t.Fatalf("NewFileStore() error = %v", err)
	}

	server := NewServerWithConfig(store, ServerConfig{
		AllowedOrigins: []string{"https://wellness.example"},
	})

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	request.Header.Set("Origin", "https://wellness.example")
	server.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", recorder.Code)
	}
	if got := recorder.Header().Get("Access-Control-Allow-Origin"); got != "https://wellness.example" {
		t.Fatalf("expected restricted CORS origin, got %q", got)
	}
	if got := recorder.Header().Get("X-Content-Type-Options"); got != "nosniff" {
		t.Fatalf("expected nosniff header, got %q", got)
	}
	if got := recorder.Header().Get("Cache-Control"); got != "no-store" {
		t.Fatalf("expected no-store cache header, got %q", got)
	}

	blockedRecorder := httptest.NewRecorder()
	blockedRequest := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	blockedRequest.Header.Set("Origin", "https://attacker.example")
	server.ServeHTTP(blockedRecorder, blockedRequest)

	if got := blockedRecorder.Header().Get("Access-Control-Allow-Origin"); got != "" {
		t.Fatalf("expected disallowed origin to receive no CORS header, got %q", got)
	}
}

func TestCreateEntryRejectsOversizedPayload(t *testing.T) {
	store, err := NewFileStore(filepath.Join(t.TempDir(), "journals.json"), mustKey(t))
	if err != nil {
		t.Fatalf("NewFileStore() error = %v", err)
	}

	server := NewServerWithConfig(store, ServerConfig{MaxJSONBodyBytes: 32})
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/v1/entries", strings.NewReader(`{"user_id":"student-1","entry_text":"too large","mood_level":5,"energy_level":5}`))
	server.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", recorder.Code)
	}
}

func TestCreateEntryRejectsTrailingJSON(t *testing.T) {
	store, err := NewFileStore(filepath.Join(t.TempDir(), "journals.json"), mustKey(t))
	if err != nil {
		t.Fatalf("NewFileStore() error = %v", err)
	}

	server := NewServer(store)
	recorder := httptest.NewRecorder()
	body := `{"user_id":"student-1","entry_text":"I studied calmly today.","mood_level":6,"energy_level":5}{}`
	request := httptest.NewRequest(http.MethodPost, "/v1/entries", strings.NewReader(body))
	server.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", recorder.Code)
	}
}

func TestServerHandlesOptionsPreflight(t *testing.T) {
	store, err := NewFileStore(filepath.Join(t.TempDir(), "journals.json"), mustKey(t))
	if err != nil {
		t.Fatalf("NewFileStore() error = %v", err)
	}

	server := NewServerWithConfig(store, ServerConfig{
		AllowedOrigins: []string{"https://wellness.example"},
	})
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodOptions, "/v1/entries", nil)
	request.Header.Set("Origin", "https://wellness.example")
	server.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", recorder.Code)
	}
	if got := recorder.Header().Get("Access-Control-Allow-Methods"); got != "GET, POST, OPTIONS" {
		t.Fatalf("unexpected CORS methods header %q", got)
	}
}

func TestHTTPRejectsUnsupportedMethods(t *testing.T) {
	store, err := NewFileStore(filepath.Join(t.TempDir(), "journals.json"), mustKey(t))
	if err != nil {
		t.Fatalf("NewFileStore() error = %v", err)
	}

	server := NewServer(store)
	tests := []struct {
		name   string
		method string
		path   string
	}{
		{"health post", http.MethodPost, "/healthz"},
		{"entries delete", http.MethodDelete, "/v1/entries"},
		{"trends post", http.MethodPost, "/v1/trends"},
		{"analysis post", http.MethodPost, "/v1/analysis"},
		{"coping post", http.MethodPost, "/v1/coping"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(tt.method, tt.path, nil)
			server.ServeHTTP(recorder, request)
			if recorder.Code != http.StatusMethodNotAllowed {
				t.Fatalf("expected 405, got %d: %s", recorder.Code, recorder.Body.String())
			}
		})
	}
}

func TestHTTPRejectsMissingQueryParameters(t *testing.T) {
	store, err := NewFileStore(filepath.Join(t.TempDir(), "journals.json"), mustKey(t))
	if err != nil {
		t.Fatalf("NewFileStore() error = %v", err)
	}

	server := NewServer(store)
	tests := []struct {
		name string
		path string
	}{
		{"list entries missing user", "/v1/entries"},
		{"analysis missing user", "/v1/analysis?entry_id=entry-1"},
		{"analysis missing entry", "/v1/analysis?user_id=student-1"},
		{"coping missing user", "/v1/coping"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodGet, tt.path, nil)
			server.ServeHTTP(recorder, request)
			if recorder.Code != http.StatusBadRequest {
				t.Fatalf("expected 400, got %d: %s", recorder.Code, recorder.Body.String())
			}
		})
	}
}

func mustKey(t *testing.T) []byte {
	t.Helper()

	key := bytes.Repeat([]byte{7}, AES256KeySize)
	encoded := base64.StdEncoding.EncodeToString(key)
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		t.Fatalf("decode key: %v", err)
	}

	return decoded
}
