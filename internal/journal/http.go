package journal

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

type Server struct {
	service *Service
	mux     *http.ServeMux
}

func NewServer(repo Repository) *Server {
	service := NewService(repo)
	mux := http.NewServeMux()
	server := &Server{service: service, mux: mux}

	mux.HandleFunc("/healthz", server.handleHealth)
	mux.HandleFunc("/v1/entries", server.handleEntries)
	mux.HandleFunc("/v1/trends", server.handleTrends)
	mux.HandleFunc("/v1/analysis", server.handleAnalysis)
	mux.HandleFunc("/v1/coping", server.handleCoping)

	return server
}

func (s *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if request.Method == http.MethodOptions {
		writer.WriteHeader(http.StatusOK)
		return
	}
	s.mux.ServeHTTP(writer, request)
}

func (s *Server) handleHealth(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		writeJSONError(writer, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	writeJSON(writer, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) handleEntries(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodPost:
		s.handleCreateEntry(writer, request)
	case http.MethodGet:
		s.handleListEntries(writer, request)
	default:
		writeJSONError(writer, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (s *Server) handleCreateEntry(writer http.ResponseWriter, request *http.Request) {
	var payload CreateEntryRequest
	if err := decodeJSON(request, &payload); err != nil {
		writeJSONError(writer, http.StatusBadRequest, err.Error())
		return
	}

	entry, err := s.service.CreateEntry(request.Context(), payload)
	if err != nil {
		writeJSONError(writer, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(writer, http.StatusCreated, entry)
}

func (s *Server) handleListEntries(writer http.ResponseWriter, request *http.Request) {
	userID := strings.TrimSpace(request.URL.Query().Get("user_id"))
	if userID == "" {
		writeJSONError(writer, http.StatusBadRequest, "user_id query parameter is required")
		return
	}

	entries, err := s.service.ListEntries(request.Context(), userID)
	if err != nil {
		writeJSONError(writer, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(writer, http.StatusOK, map[string]any{"entries": entries})
}

func (s *Server) handleTrends(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		writeJSONError(writer, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	userID := strings.TrimSpace(request.URL.Query().Get("user_id"))
	if userID == "" {
		writeJSONError(writer, http.StatusBadRequest, "user_id query parameter is required")
		return
	}

	response, err := s.service.GetTrends(request.Context(), userID)
	if err != nil {
		writeJSONError(writer, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(writer, http.StatusOK, response)
}

func (s *Server) handleAnalysis(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		writeJSONError(writer, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	entryID := strings.TrimSpace(request.URL.Query().Get("entry_id"))
	if entryID == "" {
		writeJSONError(writer, http.StatusBadRequest, "entry_id query parameter is required")
		return
	}

	userID := strings.TrimSpace(request.URL.Query().Get("user_id"))
	if userID == "" {
		writeJSONError(writer, http.StatusBadRequest, "user_id query parameter is required")
		return
	}

	response, err := s.service.GetAnalysis(request.Context(), entryID, userID)
	if err != nil {
		writeJSONError(writer, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(writer, http.StatusOK, response)
}

func (s *Server) handleCoping(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		writeJSONError(writer, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	userID := strings.TrimSpace(request.URL.Query().Get("user_id"))
	if userID == "" {
		writeJSONError(writer, http.StatusBadRequest, "user_id query parameter is required")
		return
	}

	response, err := s.service.GetCopingGuidance(request.Context(), userID)
	if err != nil {
		writeJSONError(writer, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(writer, http.StatusOK, response)
}

func decodeJSON(request *http.Request, destination any) error {
	defer request.Body.Close()
	dec := json.NewDecoder(request.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(destination); err != nil {
		return errors.New("invalid JSON payload")
	}
	if dec.More() {
		return errors.New("invalid JSON payload")
	}

	return nil
}

func writeJSON(writer http.ResponseWriter, status int, payload any) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	_ = json.NewEncoder(writer).Encode(payload)
}

func writeJSONError(writer http.ResponseWriter, status int, message string) {
	writeJSON(writer, status, map[string]string{"error": message})
}
