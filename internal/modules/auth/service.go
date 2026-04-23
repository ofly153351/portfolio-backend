package auth

import (
	"crypto/rand"
	"encoding/hex"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/bcrypt"
	"portfolio-backend/internal/config"
)

type Service struct {
	user         User
	passwordHash []byte
	tokenTTL     time.Duration
	mu           sync.RWMutex
	sessions     map[string]sessionState
}

type sessionState struct {
	user      User
	expiresAt time.Time
}

func NewService(cfg config.Config) (*Service, error) {
	hash := strings.TrimSpace(cfg.AdminPasswordHash)
	if hash == "" {
		encoded, err := bcrypt.GenerateFromPassword([]byte(cfg.AdminPassword), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		hash = string(encoded)
	}

	return &Service{
		user: User{
			ID:       "u_001",
			Username: cfg.AdminUsername,
			Role:     "admin",
		},
		passwordHash: []byte(hash),
		tokenTTL:     time.Duration(cfg.AuthTokenTTLMinutes) * time.Minute,
		sessions:     make(map[string]sessionState),
	}, nil
}

func (s *Service) Login(username, password string) (string, User, error) {
	if s == nil {
		return "", User{}, ErrUnauthorized
	}
	if username != s.user.Username {
		return "", User{}, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword(s.passwordHash, []byte(password)); err != nil {
		return "", User{}, ErrInvalidCredentials
	}

	token, err := generateToken()
	if err != nil {
		return "", User{}, err
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.sessions[token] = sessionState{
		user:      s.user,
		expiresAt: time.Now().Add(s.tokenTTL),
	}

	return token, s.user, nil
}

func (s *Service) Validate(token string) (User, error) {
	if s == nil || strings.TrimSpace(token) == "" {
		return User{}, ErrUnauthorized
	}

	s.mu.RLock()
	session, ok := s.sessions[token]
	s.mu.RUnlock()
	if !ok {
		return User{}, ErrUnauthorized
	}
	if time.Now().After(session.expiresAt) {
		s.mu.Lock()
		delete(s.sessions, token)
		s.mu.Unlock()
		return User{}, ErrUnauthorized
	}

	return session.user, nil
}

func (s *Service) Logout(token string) {
	if s == nil || strings.TrimSpace(token) == "" {
		return
	}
	s.mu.Lock()
	delete(s.sessions, token)
	s.mu.Unlock()
}

func generateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
