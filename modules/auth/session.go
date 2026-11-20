package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"time"

	"gorm.io/gorm"
)

var ErrInvalidRefreshToken = errors.New("invalid refresh token")

const refreshTokenTTL = 30 * 24 * time.Hour

type Session struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	UserID    uint       `gorm:"not null;index" json:"user_id"`
	TokenHash string     `gorm:"size:64;not null;uniqueIndex" json:"-"`
	ExpiresAt time.Time  `gorm:"not null;index" json:"expires_at"`
	RevokedAt *time.Time `gorm:"index" json:"-"`
	UserAgent string     `gorm:"size:512" json:"user_agent,omitempty"`
	IPAddress string     `gorm:"size:45" json:"ip_address,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type SessionRepository struct{ DB *gorm.DB }

func NewSessionRepository(db *gorm.DB) *SessionRepository { return &SessionRepository{DB: db} }

func newRefreshToken() (string, string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", "", err
	}
	raw := base64.RawURLEncoding.EncodeToString(b)
	sum := sha256.Sum256([]byte(raw))
	return raw, base64.RawURLEncoding.EncodeToString(sum[:]), nil
}

func (r *SessionRepository) Create(userID uint, userAgent, ip string) (string, error) {
	raw, hash, err := newRefreshToken()
	if err != nil {
		return "", err
	}
	now := time.Now().UTC()
	s := Session{UserID: userID, TokenHash: hash, ExpiresAt: now.Add(refreshTokenTTL), UserAgent: truncate(userAgent, 512), IPAddress: truncate(ip, 45)}
	if err := r.DB.Create(&s).Error; err != nil {
		return "", err
	}
	return raw, nil
}

// Rotate performs a compare-and-swap update: only one concurrent refresh can consume a token.
func (r *SessionRepository) Rotate(raw string) (uint, string, error) {
	sum := sha256.Sum256([]byte(raw))
	oldHash := base64.RawURLEncoding.EncodeToString(sum[:])
	newRaw, newHash, err := newRefreshToken()
	if err != nil {
		return 0, "", err
	}
	now := time.Now().UTC()
	var session Session
	err = r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("token_hash = ? AND revoked_at IS NULL AND expires_at > ?", oldHash, now).First(&session).Error; err != nil {
			return ErrInvalidRefreshToken
		}
		result := tx.Model(&Session{}).Where("id = ? AND token_hash = ? AND revoked_at IS NULL AND expires_at > ?", session.ID, oldHash, now).
			Updates(map[string]any{"token_hash": newHash, "expires_at": now.Add(refreshTokenTTL), "updated_at": now})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected != 1 {
			return ErrInvalidRefreshToken
		}
		return nil
	})
	if err != nil {
		return 0, "", err
	}
	return session.UserID, newRaw, nil
}

func (r *SessionRepository) Revoke(raw string, userID uint) error {
	sum := sha256.Sum256([]byte(raw))
	hash := base64.RawURLEncoding.EncodeToString(sum[:])
	now := time.Now().UTC()
	return r.DB.Model(&Session{}).Where("token_hash = ? AND user_id = ? AND revoked_at IS NULL", hash, userID).Update("revoked_at", now).Error
}

func (r *SessionRepository) RevokeAll(userID uint) error {
	now := time.Now().UTC()
	return r.DB.Model(&Session{}).Where("user_id = ? AND revoked_at IS NULL", userID).Update("revoked_at", now).Error
}

func truncate(s string, max int) string {
	if len(s) > max {
		return s[:max]
	}
	return s
}
