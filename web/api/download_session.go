package api

import (
	"fmt"
	"math/rand/v2"
	"sync"
	"time"
)

const (
	downloadStatusPending = "pending"
	downloadStatusRunning = "running"
	downloadStatusSuccess = "success"
	downloadStatusError   = "error"
)

type downloadEvent struct {
	Type        string `json:"type"`
	SessionID   string `json:"sessionId"`
	Target      string `json:"target,omitempty"`
	Title       string `json:"title,omitempty"`
	Status      string `json:"status,omitempty"`
	Progress    int    `json:"progress,omitempty"`
	Current     int    `json:"current,omitempty"`
	Total       int    `json:"total,omitempty"`
	CurrentName string `json:"currentName,omitempty"`
	Level       string `json:"level,omitempty"`
	Message     string `json:"message,omitempty"`
	OutputDir   string `json:"outputDir,omitempty"`
	Timestamp   int64  `json:"timestamp"`
}

type downloadSession struct {
	id          string
	target      string
	title       string
	outputDir   string
	status      string
	progress    int
	current     int
	total       int
	currentName string
	createdAt   time.Time
	updatedAt   time.Time
	history     []downloadEvent
	listeners   map[chan downloadEvent]struct{}
	mu          sync.RWMutex
}

type downloadSessionStore struct {
	sessions map[string]*downloadSession
	mu       sync.RWMutex
}

func newDownloadSessionStore() *downloadSessionStore {
	return &downloadSessionStore{
		sessions: map[string]*downloadSession{},
	}
}

func (s *downloadSessionStore) Create(target, title, outputDir string) *downloadSession {
	now := time.Now()
	session := &downloadSession{
		id:        fmt.Sprintf("%d-%d", now.UnixNano(), rand.Uint64()),
		target:    target,
		title:     title,
		outputDir: outputDir,
		status:    downloadStatusPending,
		createdAt: now,
		updatedAt: now,
		listeners: map[chan downloadEvent]struct{}{},
	}

	s.mu.Lock()
	s.sessions[session.id] = session
	s.mu.Unlock()

	session.publish(downloadEvent{
		Type:      "start",
		Status:    downloadStatusPending,
		Message:   "下载任务已创建",
		OutputDir: outputDir,
	})

	return session
}

func (s *downloadSessionStore) Get(sessionID string) (*downloadSession, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	session, ok := s.sessions[sessionID]
	return session, ok
}

func (s *downloadSessionStore) Subscribe(sessionID string) (*downloadSession, chan downloadEvent, []downloadEvent, bool) {
	session, ok := s.Get(sessionID)
	if !ok {
		return nil, nil, nil, false
	}

	ch := make(chan downloadEvent, 64)

	session.mu.Lock()
	session.listeners[ch] = struct{}{}
	history := append([]downloadEvent(nil), session.history...)
	session.mu.Unlock()

	return session, ch, history, true
}

func (s *downloadSessionStore) Unsubscribe(sessionID string, ch chan downloadEvent) {
	session, ok := s.Get(sessionID)
	if !ok {
		return
	}

	session.mu.Lock()
	delete(session.listeners, ch)
	session.mu.Unlock()
	close(ch)
}

func (s *downloadSessionStore) Cleanup(expireAfter time.Duration) {
	if expireAfter <= 0 {
		return
	}

	cutoff := time.Now().Add(-expireAfter)
	s.mu.Lock()
	defer s.mu.Unlock()

	for id, session := range s.sessions {
		session.mu.RLock()
		shouldDelete := (session.status == downloadStatusSuccess || session.status == downloadStatusError) && session.updatedAt.Before(cutoff)
		session.mu.RUnlock()
		if shouldDelete {
			delete(s.sessions, id)
		}
	}
}

func (s *downloadSessionStore) StartCleanupLoop(interval, expireAfter time.Duration) {
	if interval <= 0 || expireAfter <= 0 {
		return
	}
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for range ticker.C {
			s.Cleanup(expireAfter)
		}
	}()
}

func (s *downloadSession) Snapshot() downloadEvent {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return downloadEvent{
		Type:        "snapshot",
		SessionID:   s.id,
		Target:      s.target,
		Title:       s.title,
		Status:      s.status,
		Progress:    s.progress,
		Current:     s.current,
		Total:       s.total,
		CurrentName: s.currentName,
		OutputDir:   s.outputDir,
		Timestamp:   s.updatedAt.UnixMilli(),
	}
}

func (s *downloadSession) publish(event downloadEvent) {
	s.mu.Lock()
	defer s.mu.Unlock()

	event.SessionID = s.id
	event.Target = s.target
	event.Title = s.title
	if event.OutputDir == "" {
		event.OutputDir = s.outputDir
	}
	if event.Timestamp == 0 {
		event.Timestamp = time.Now().UnixMilli()
	}

	switch event.Type {
	case "start":
		s.status = downloadStatusPending
	case "progress":
		s.status = downloadStatusRunning
		s.progress = event.Progress
		s.current = event.Current
		s.total = event.Total
		s.currentName = event.CurrentName
	case "done":
		s.status = downloadStatusSuccess
		s.progress = 100
	case "error":
		s.status = downloadStatusError
	case "log":
		if s.status == downloadStatusPending {
			s.status = downloadStatusRunning
		}
	}

	event.Status = s.status
	s.updatedAt = time.Now()
	if len(s.history) >= 300 {
		s.history = append(s.history[:0], s.history[1:]...)
	}
	s.history = append(s.history, event)

	for listener := range s.listeners {
		select {
		case listener <- event:
		default:
		}
	}
}
