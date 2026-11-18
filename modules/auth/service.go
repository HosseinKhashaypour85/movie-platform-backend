package auth

import (
	"errors"

	"movie-platform/modules/users"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var ErrEmailAlreadyExists = errors.New("email already exists")

type Service struct {
	Users *users.Repository
}

func NewService(repo *users.Repository) *Service {
	return &Service{Users: repo}
}

func (s *Service) Register(name, email, password string) (*users.User, error) {
	if s == nil || s.Users == nil {
		return nil, errors.New("user repository is not configured")
	}

	_, err := s.Users.FindByEmail(email)
	switch {
	case err == nil:
		return nil, ErrEmailAlreadyExists
	case !errors.Is(err, gorm.ErrRecordNotFound):
		return nil, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &users.User{
		Name:     name,
		Email:    email,
		Password: string(hashedPassword),
	}
	if err := s.Users.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}
