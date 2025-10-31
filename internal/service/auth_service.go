package service

import (
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"lims_auth_service/internal/dto"
	"lims_auth_service/internal/model"
	"lims_auth_service/internal/repository"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AuthService struct {
	Repo      *repository.UserRepository
	jwtSecret string
}

func NewAuthService(repo *repository.UserRepository, jwtSecret string) *AuthService {
	return &AuthService{Repo: repo, jwtSecret: jwtSecret}
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

func (s *AuthService) Login(email, password string) (string, error) {
	var errorMessage = "Incorrect email or password"

	// проверка существования пользователя в системе
	user, err := s.Repo.GetByEmail(email)
	if err != nil {
		return "", fmt.Errorf(errorMessage)
	}

	// проверка пароля
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New(errorMessage)
	}

	if user.IsActive == false {
		return "", errors.New("User is not active")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
