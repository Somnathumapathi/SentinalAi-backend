package github

import (
	"sync"
	"time"
)

// TokenCache represents a cached GitHub App installation token
type TokenCache struct {
	Token     string
	ExpiresAt time.Time
}

// TokenCacheService manages GitHub App installation tokens
type TokenCacheService struct {
	cache map[int64]*TokenCache
	mu    sync.RWMutex
}

var (
	tokenService *TokenCacheService
	once         sync.Once
)

// GetTokenCacheService returns a singleton instance of TokenCacheService
func GetTokenCacheService() *TokenCacheService {
	once.Do(func() {
		tokenService = &TokenCacheService{
			cache: make(map[int64]*TokenCache),
		}
	})
	return tokenService
}

// GetToken returns a valid token for the given installation ID
func (s *TokenCacheService) GetToken(installationID int64) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if cache, exists := s.cache[installationID]; exists {
		// Check if token is still valid (with 5-minute buffer)
		if time.Until(cache.ExpiresAt) > 5*time.Minute {
			return cache.Token, true
		}
	}
	return "", false
}

// SetToken caches a new token for the given installation ID
func (s *TokenCacheService) SetToken(installationID int64, token string, expiresAt time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.cache[installationID] = &TokenCache{
		Token:     token,
		ExpiresAt: expiresAt,
	}
}

// ClearToken removes a token from the cache
func (s *TokenCacheService) ClearToken(installationID int64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.cache, installationID)
}
