package api

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
	"time"
)

type qrLoginSession struct {
	ID           string
	Token        string
	QRCode       string
	QRCodeString string
	ExpiresAt    time.Time
}

type qrSessionStore struct {
	mu       sync.RWMutex
	sessions map[string]*qrLoginSession
}

func newQRSessionStore() *qrSessionStore {
	return &qrSessionStore{sessions: make(map[string]*qrLoginSession)}
}

func (s *qrSessionStore) Create(token, qrCode, qrCodeString string) *qrLoginSession {
	session := &qrLoginSession{
		ID:           randomID(),
		Token:        token,
		QRCode:       qrCode,
		QRCodeString: qrCodeString,
		ExpiresAt:    time.Now().Add(10 * time.Minute),
	}

	s.mu.Lock()
	s.sessions[session.ID] = session
	s.mu.Unlock()

	return session
}

func (s *qrSessionStore) Get(id string) (*qrLoginSession, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	session, ok := s.sessions[id]
	return session, ok
}

func (s *qrSessionStore) Delete(id string) {
	s.mu.Lock()
	delete(s.sessions, id)
	s.mu.Unlock()
}

func randomID() string {
	buf := make([]byte, 16)
	_, _ = rand.Read(buf)
	return hex.EncodeToString(buf)
}
