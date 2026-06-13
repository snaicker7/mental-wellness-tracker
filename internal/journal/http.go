package journal

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"strings"
)

const defaultMaxJSONBodyBytes int64 = 1 << 20

var defaultAllowedOrigins = []string{
	"http://localhost:5173",
	"http://127.0.0.1:5173",
}

type ServerConfig struct {
	AllowedOrigins   []string
	MaxJSONBodyBytes int64
}

type Server struct {
	service          *Service
	mux              *http.ServeMux
	allowedOrigins   map[string]struct{}
	allowAnyOrigin   bool
	maxJSONBodyBytes int64
}

func NewServer(repo Repository) *Server {
	return NewServerWithConfig(repo, ServerConfig{
		AllowedOrigins: splitCSV(os.Getenv("ALLOWED_ORIGINS")),
	})
}

func NewServerWithConfig(repo Repository, config ServerConfig) *Server {
	allowedOrigins := config.AllowedOrigins
	if len(allowedOrigins) == 0 {
		allowedOrigins = defaultAllowedOrigins
	}

	service := NewService(repo)
	mux := http.NewServeMux()
	server := &Server{
		service:          service,
		mux:              mux,
		allowedOrigins:   allowedOriginSet(allowedOrigins),
		maxJSONBodyBytes: config.MaxJSONBodyBytes,
	}
	if server.maxJSONBodyBytes <= 0 {
		server.maxJSONBodyBytes = defaultMaxJSONBodyBytes
	}
	if _, ok := server.allowedOrigins["*"]; ok {
		server.allowAnyOrigin = true
	}

	mux.HandleFunc("/healthz", server.handleHealth)
	mux.HandleFunc("/v1/entries", server.handleEntries)
	mux.HandleFunc("/v1/trends", server.handleTrends)
	mux.HandleFunc("/v1/analysis", server.handleAnalysis)
	mux.HandleFunc("/v1/coping", server.handleCoping)

	return server
}

func (s *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	s.writeSecurityHeaders(writer)
	s.writeCORSHeaders(writer, request)
	if request.Method == http.MethodOptions {
		writer.WriteHeader(http.StatusOK)
		return
	}
	s.mux.ServeHTTP(writer, request)
}

func (s *Server) writeSecurityHeaders(writer http.ResponseWriter) {
	header := writer.Header()
	header.Set("X-Content-Type-Options", "nosniff")
	header.Set("X-Frame-Options", "DENY")
	header.Set("Referrer-Policy", "no-referrer")
	header.Set("Cache-Control", "no-store")
}

func (s *Server) writeCORSHeaders(writer http.ResponseWriter, request *http.Request) {
	origin := request.Header.Get("Origin")
	if origin == "" {
		return
	}

	if s.allowAnyOrigin {
		writer.Header().Set("Access-Control-Allow-Origin", "*")
	} else if _, ok := s.allowedOrigins[origin]; ok {
		writer.Header().Set("Access-Control-Allow-Origin", origin)
		writer.Header().Add("Vary", "Origin")
	} else {
		return
	}

	writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	writer.Header().Set("Access-Control-Max-Age", "600")
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
	if err := decodeJSON(request, s.maxJSONBodyBytes, &payload); err != nil {
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

func decodeJSON(request *http.Request, maxBytes int64, destination any) error {
	defer request.Body.Close()
	if maxBytes <= 0 {
		maxBytes = defaultMaxJSONBodyBytes
	}

	request.Body = http.MaxBytesReader(nil, request.Body, maxBytes)
	dec := json.NewDecoder(request.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(destination); err != nil {
		if errors.Is(err, io.EOF) {
			return errors.New("request body is required")
		}
		return errors.New("invalid JSON payload")
	}
	var extra any
	if err := dec.Decode(&extra); !errors.Is(err, io.EOF) {
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

func splitCSV(value string) []string {
	parts := strings.Split(value, ",")
	values := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			values = append(values, part)
		}
	}
	return values
}

func allowedOriginSet(origins []string) map[string]struct{} {
	allowed := make(map[string]struct{}, len(origins))
	for _, origin := range origins {
		origin = strings.TrimSpace(origin)
		if origin != "" {
			allowed[origin] = struct{}{}
		}
	}
	return allowed
}
