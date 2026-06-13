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
