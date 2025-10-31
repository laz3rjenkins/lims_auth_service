package service

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"lims_auth_service/internal/dto"
	"lims_auth_service/internal/model"
	"lims_auth_service/internal/repository"
)

type AuthService struct {
	Repo *repository.UserRepository
}

func NewAuthService(repo *repository.UserRepository) *AuthService {
	return &AuthService{Repo: repo}
}

func (s *AuthService) Register(body dto.RegisterRequest) error {
	// Проверка существующего пользователя
	if _, err := s.Repo.GetByEmail(body.Email); err == nil {
		return fmt.Errorf("user already exists")
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	user := &model.User{Email: body.Email, Password: string(hash), IsActive: true}

	return s.Repo.CreateUser(user)
}
