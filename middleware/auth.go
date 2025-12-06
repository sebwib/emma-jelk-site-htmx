package middleware

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type contextKey string

const (
	SessionCookieName = "session_token"
	UserContextKey    = contextKey("user")
)

type Session struct {
	Token     string
	Username  string
	ExpiresAt time.Time
}

type SessionStore struct {
	sessions map[string]*Session
	mu       sync.RWMutex
}

func NewSessionStore() *SessionStore {
	store := &SessionStore{
		sessions: make(map[string]*Session),
	}
	// Cleanup expired sessions every hour
	go store.cleanupExpired()
	return store
}

func (s *SessionStore) CreateSession(username string) (string, error) {
	token, err := generateToken()
	if err != nil {
		return "", err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	session := &Session{
		Token:     token,
		Username:  username,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	s.sessions[token] = session

	return token, nil
}

func (s *SessionStore) GetSession(token string) (*Session, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	session, exists := s.sessions[token]
	if !exists || time.Now().After(session.ExpiresAt) {
		return nil, false
	}
	return session, true
}

func (s *SessionStore) DeleteSession(token string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.sessions, token)
}

func (s *SessionStore) cleanupExpired() {
	ticker := time.NewTicker(1 * time.Hour)
	for range ticker.C {
		s.mu.Lock()
		now := time.Now()
		for token, session := range s.sessions {
			if now.After(session.ExpiresAt) {
				delete(s.sessions, token)
			}
		}
		s.mu.Unlock()
	}
}

func generateToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// RequireAuth middleware checks if user is authenticated
func RequireAuth(store *SessionStore) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie(SessionCookieName)
			if err != nil {
				// No session cookie, redirect to login
				http.Redirect(w, r, "/login?redirect_to="+url.QueryEscape(r.URL.RequestURI()), http.StatusSeeOther)
				return
			}

			session, valid := store.GetSession(cookie.Value)
			if !valid {
				// Invalid or expired session
				http.Redirect(w, r, "/login?redirect_to="+url.QueryEscape(r.URL.RequestURI()), http.StatusSeeOther)
				return
			}

			// Add username to context
			ctx := context.WithValue(r.Context(), UserContextKey, session.Username)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
