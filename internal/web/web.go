package web

import (
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/sidx1/sure/internal/config"
	"github.com/sidx1/sure/internal/stats"
)

//go:embed templates/*
var templateFS embed.FS

type Server struct {
	port int
}

func NewServer(port int) *Server {
	if port <= 0 {
		port = 7873
	}
	return &Server{port: port}
}

func (s *Server) Start() error {
	mux := http.NewServeMux()

	// Static HTML Dashboard
	mux.HandleFunc("/", s.handleDashboard)

	// JSON API endpoints
	mux.HandleFunc("/api/stats", s.handleAPIStats)
	mux.HandleFunc("/api/events", s.handleAPIEvents)
	mux.HandleFunc("/api/badge", s.handleBadge)

	addr := fmt.Sprintf("0.0.0.0:%d", s.port)
	fmt.Printf("\n🚀 Sure Web Dashboard live at: http://localhost:%d\n", s.port)
	fmt.Printf("   Showcase Badge: http://localhost:%d/api/badge\n", s.port)
	fmt.Printf("   Press Ctrl+C to exit.\n\n")

	server := &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return server.ListenAndServe()
}

func (s *Server) handleDashboard(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	summary, err := stats.GetSummary()
	if err != nil {
		summary = &stats.Summary{AccuracyScore: 100.0, ImprovementTrend: "calibrating"}
	}

	cfg, _ := config.Load()

	tmpl, err := template.ParseFS(templateFS, "templates/dashboard.html")
	if err != nil {
		http.Error(w, fmt.Sprintf("Template error: %v", err), http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Summary": summary,
		"Config":  cfg,
		"Year":    time.Now().Year(),
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_ = tmpl.Execute(w, data)
}

func (s *Server) handleAPIStats(w http.ResponseWriter, r *http.Request) {
	summary, err := stats.GetSummary()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(summary)
}

func (s *Server) handleAPIEvents(w http.ResponseWriter, r *http.Request) {
	events, err := stats.GetEvents(50)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(events)
}

func (s *Server) handleBadge(w http.ResponseWriter, r *http.Request) {
	summary, err := stats.GetSummary()
	if err != nil {
		summary = &stats.Summary{MistakesCaught: 0, AccuracyScore: 100.0}
	}

	// Returns an SVG badge you can embed on GitHub README / portfolio to show off your command typing stats
	w.Header().Set("Content-Type", "image/svg+xml; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")

	caughtText := fmt.Sprintf("%d caught | %.1f%% safe", summary.MistakesCaught, summary.AccuracyScore)
	svg := fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" width="220" height="20" role="img" aria-label="sure: %s">
  <linearGradient id="b" x2="0" y2="100%%">
    <stop offset="0" stop-color="#bbb" stop-opacity=".1"/>
    <stop offset="1" stop-opacity=".1"/>
  </linearGradient>
  <clipPath id="a">
    <rect width="220" height="20" rx="3" fill="#fff"/>
  </clipPath>
  <g clip-path="url(#a)">
    <rect width="55" height="20" fill="#24292e"/>
    <rect x="55" width="165" height="20" fill="#007acc"/>
    <rect width="220" height="20" fill="url(#b)"/>
  </g>
  <g fill="#fff" text-anchor="middle" font-family="Verdana,Geneva,DejaVu Sans,sans-serif" text-rendering="geometricPrecision" font-size="110">
    <text x="285" y="140" transform="scale(.1)" fill="#fff" textLength="450">✋ sure</text>
    <text x="1360" y="140" transform="scale(.1)" fill="#fff" textLength="1500">%s</text>
  </g>
</svg>`, caughtText, caughtText)

	_, _ = w.Write([]byte(svg))
}
