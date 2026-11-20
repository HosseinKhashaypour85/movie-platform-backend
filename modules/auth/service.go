package auth

import (
	"errors"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"movie-platform/internal/config"
	"movie-platform/modules/users"
)

var (
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidInput       = errors.New("invalid input")
)

type Service struct {
	Users *users.Repository
	Sessions *SessionRepository
	Config config.Config
}

func NewService(u *users.Repository, sessions *SessionRepository, c config.Config) *Service {
	return &Service{Users: u, Sessions: sessions, Config: c}
}
func validate(name, email, password string) error {
	name = strings.TrimSpace(name)
	email = strings.ToLower(strings.TrimSpace(email))
	if !utf8.ValidString(name) || len(name) < 2 || len([]byte(name)) > 100 {
		return ErrInvalidInput
	}
	if len(email) > 254 || !strings.Contains(email, "@") || strings.HasPrefix(email, "@") || strings.HasSuffix(email, "@") {
		return ErrInvalidInput
	}
	if len(password) < 12 || len([]byte(password)) > 72 {
		return ErrInvalidInput
	}
	return nil
}
func (s *Service) Register(name, email, password string) (*users.User, error) {
	if e := validate(name, email, password); e != nil {
		return nil, e
	}
	email = strings.ToLower(strings.TrimSpace(email))
	hash, e := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if e != nil {
		return nil, e
	}
	u := &users.User{Name: strings.TrimSpace(name), Email: email, Password: string(hash), Role: "user", IsActive: true}
	if e = s.Users.Create(u); e != nil {
		if strings.Contains(strings.ToLower(e.Error()), "duplicate") {
			return nil, ErrEmailAlreadyExists
		}
		return nil, e
	}
	return u, nil
}
func (s *Service) Login(email, password string) (string, *users.User, error) {
	access, _, user, err := s.LoginWithSession(email, password, "", "")
	return access, user, err
}

func (s *Service) LoginWithSession(email, password, userAgent, ip string) (string, string, *users.User, error) {
	u, e := s.Users.FindByEmail(strings.ToLower(strings.TrimSpace(email)))
	if e != nil {
		return "", "", nil, ErrInvalidCredentials
	}
	if !u.IsActive || bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)) != nil {
		return "", "", nil, ErrInvalidCredentials
	}
	now := time.Now()
	c := jwt.MapClaims{"user_id": u.ID, "role": u.Role, "iss": s.Config.JWTIssuer, "sub": fmt.Sprint(u.ID), "iat": now.Unix(), "exp": now.Add(s.Config.JWTTTL).Unix()}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	signed, e := t.SignedString(s.Config.JWTSecret)
	if e != nil {
		return "", "", nil, e
	}
	refresh, err := s.Sessions.Create(u.ID, userAgent, ip)
	if err != nil { return "", "", nil, err }
	return signed, refresh, u, nil
}

func (s *Service) Refresh(raw string) (string, string, error) {
	userID, rotated, err := s.Sessions.Rotate(raw)
	if err != nil { return "", "", ErrInvalidRefreshToken }
	user, err := s.Users.FindByID(userID)
	if err != nil || !user.IsActive {
		_ = s.Sessions.Revoke(rotated, userID)
		return "", "", ErrInvalidRefreshToken
	}
	now := time.Now()
	claims := jwt.MapClaims{"user_id":user.ID, "role":user.Role, "iss":s.Config.JWTIssuer, "sub":fmt.Sprint(user.ID), "iat":now.Unix(), "exp":now.Add(s.Config.JWTTTL).Unix()}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	access, err := token.SignedString(s.Config.JWTSecret)
	if err != nil { return "", "", err }
	return access, rotated, nil
}
