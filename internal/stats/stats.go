package stats

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/sidx1/sure/internal/config"
)

type CommandEvent struct {
	Timestamp  time.Time `json:"timestamp"`
	Command    string    `json:"command"`
	Warned     bool      `json:"warned"`
	Proceeded  bool      `json:"proceeded"`
	Category   string    `json:"category,omitempty"` // "destructive", "typo", "suspicious", "invalid_syntax", "risky"
	Reason     string    `json:"reason,omitempty"`
	Confidence float64   `json:"confidence,omitempty"`
}

type Summary struct {
	TotalAnalyzed    int            `json:"total_analyzed"`
	TotalWarned      int            `json:"total_warned"`
	MistakesCaught   int            `json:"mistakes_caught"` // Warned & Cancelled (User chose No)
	WarningsProceeded int           `json:"warnings_proceeded"` // Warned & User chose Yes anyway
	SafeAllowed      int            `json:"safe_allowed"`
	CategoryBreakdown map[string]int `json:"category_breakdown"`
	MistakeRate      float64        `json:"mistake_rate"`       // Mistakes caught / Total analyzed (%)
	AccuracyScore    float64        `json:"accuracy_score"`     // 100 - MistakeRate (%)
	ImprovementTrend string         `json:"improvement_trend"` // "improving", "stable", "attention_needed"
	RecentMistakes   []CommandEvent `json:"recent_mistakes"`
}

type Tracker interface {
	Record(event CommandEvent) error
	GetSummary() (*Summary, error)
	GetEvents(limit int) ([]CommandEvent, error)
	Reset() error
}

type FileTracker struct {
	mu       sync.Mutex
	filePath string
}

func NewFileTracker() *FileTracker {
	dir := config.DataDir()
	_ = os.MkdirAll(dir, 0755)
	return &FileTracker{
		filePath: filepath.Join(dir, "events.jsonl"),
	}
}

func (t *FileTracker) Record(event CommandEvent) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	dir := filepath.Dir(t.filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now()
	}

	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	f, err := os.OpenFile(t.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(append(data, '\n'))
	return err
}

func (t *FileTracker) GetEvents(limit int) ([]CommandEvent, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	return t.readEvents(limit)
}

func (t *FileTracker) readEvents(limit int) ([]CommandEvent, error) {
	data, err := os.ReadFile(t.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return []CommandEvent{}, nil
		}
		return nil, err
	}

	var events []CommandEvent
	lines := splitLines(data)
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		var ev CommandEvent
		if err := json.Unmarshal(line, &ev); err == nil {
			events = append(events, ev)
		}
	}

	if limit > 0 && len(events) > limit {
		return events[len(events)-limit:], nil
	}
	return events, nil
}

func (t *FileTracker) GetSummary() (*Summary, error) {
	events, err := t.GetEvents(0)
	if err != nil {
		return nil, err
	}

	summary := &Summary{
		CategoryBreakdown: make(map[string]int),
		RecentMistakes:   make([]CommandEvent, 0),
	}

	if len(events) == 0 {
		summary.AccuracyScore = 100.0
		summary.ImprovementTrend = "perfect"
		return summary, nil
	}

	summary.TotalAnalyzed = len(events)

	for _, ev := range events {
		if ev.Warned {
			summary.TotalWarned++
			cat := ev.Category
			if cat == "" {
				cat = "risky"
			}
			summary.CategoryBreakdown[cat]++

			if !ev.Proceeded {
				summary.MistakesCaught++
				summary.RecentMistakes = append(summary.RecentMistakes, ev)
			} else {
				summary.WarningsProceeded++
			}
		} else {
			summary.SafeAllowed++
		}
	}

	if summary.TotalAnalyzed > 0 {
		summary.MistakeRate = (float64(summary.MistakesCaught) / float64(summary.TotalAnalyzed)) * 100.0
		summary.AccuracyScore = 100.0 - summary.MistakeRate
	}

	// Calculate trend by comparing first half of history vs second half
	if len(events) >= 10 {
		mid := len(events) / 2
		firstHalfMistakes := 0
		secondHalfMistakes := 0

		for i := 0; i < mid; i++ {
			if events[i].Warned && !events[i].Proceeded {
				firstHalfMistakes++
			}
		}
		for i := mid; i < len(events); i++ {
			if events[i].Warned && !events[i].Proceeded {
				secondHalfMistakes++
			}
		}

		firstRate := float64(firstHalfMistakes) / float64(mid)
		secondRate := float64(secondHalfMistakes) / float64(len(events)-mid)

		if secondRate < firstRate {
			summary.ImprovementTrend = "improving"
		} else if secondRate == firstRate {
			summary.ImprovementTrend = "stable"
		} else {
			summary.ImprovementTrend = "attention_needed"
		}
	} else {
		summary.ImprovementTrend = "calibrating"
	}

	// Keep only the 10 most recent mistakes
	if len(summary.RecentMistakes) > 10 {
		summary.RecentMistakes = summary.RecentMistakes[len(summary.RecentMistakes)-10:]
	}

	return summary, nil
}

func (t *FileTracker) Reset() error {
	t.mu.Lock()
	defer t.mu.Unlock()
	return os.Remove(t.filePath)
}

func splitLines(data []byte) [][]byte {
	var lines [][]byte
	start := 0
	for i, b := range data {
		if b == '\n' {
			lines = append(lines, data[start:i])
			start = i + 1
		}
	}
	if start < len(data) {
		lines = append(lines, data[start:])
	}
	return lines
}

var (
	defaultTracker Tracker
	trackerOnce    sync.Once
)

func DefaultTracker() Tracker {
	trackerOnce.Do(func() {
		defaultTracker = NewFileTracker()
	})
	return defaultTracker
}

func Record(ev CommandEvent) error {
	return DefaultTracker().Record(ev)
}

func GetSummary() (*Summary, error) {
	return DefaultTracker().GetSummary()
}

func GetEvents(limit int) ([]CommandEvent, error) {
	return DefaultTracker().GetEvents(limit)
}
