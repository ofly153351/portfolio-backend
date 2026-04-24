package publicauth

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"strings"
	"time"
)

type tokenClaims struct {
	Exp int64 `json:"exp"`
}

type Service struct {
	secret []byte
	ttl    time.Duration
}

func NewService(secret string, ttl time.Duration) *Service {
	normalizedSecret := strings.TrimSpace(secret)
	if normalizedSecret == "" {
		normalizedSecret = "dev-public-token-secret-change-me"
	}
	if ttl <= 0 {
		ttl = 5 * time.Minute
	}
	return &Service{
		secret: []byte(normalizedSecret),
		ttl:    ttl,
	}
}

func (s *Service) Issue() (string, time.Time, error) {
	expiresAt := time.Now().UTC().Add(s.ttl)
	claims := tokenClaims{Exp: expiresAt.Unix()}
	payload, err := json.Marshal(claims)
	if err != nil {
		return "", time.Time{}, err
	}

	encodedPayload := base64.RawURLEncoding.EncodeToString(payload)
	signature := s.sign(encodedPayload)
	encodedSignature := base64.RawURLEncoding.EncodeToString(signature)

	return encodedPayload + "." + encodedSignature, expiresAt, nil
}

func (s *Service) Validate(token string) error {
	parts := strings.Split(strings.TrimSpace(token), ".")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return ErrTokenInvalid
	}

	expected := s.sign(parts[0])
	received, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return ErrTokenInvalid
	}
	if subtle.ConstantTimeCompare(received, expected) != 1 {
		return ErrTokenInvalid
	}

	claimsBytes, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return ErrTokenInvalid
	}
	var claims tokenClaims
	if err := json.Unmarshal(claimsBytes, &claims); err != nil {
		return ErrTokenInvalid
	}
	if claims.Exp <= time.Now().UTC().Unix() {
		return ErrTokenExpired
	}
	return nil
}

func (s *Service) sign(payload string) []byte {
	mac := hmac.New(sha256.New, s.secret)
	mac.Write([]byte(payload))
	return mac.Sum(nil)
}

