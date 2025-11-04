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
	Repo             *repository.UserRepository
	jwtAccessSecret  string
	jwtRefreshSecret string
}

func NewAuthService(repo *repository.UserRepository, jwtAccessSecret, jwtRefreshSecret string) *AuthService {
	return &AuthService{Repo: repo, jwtAccessSecret: jwtAccessSecret, jwtRefreshSecret: jwtRefreshSecret}
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

func (s *AuthService) Login(email, password string) (string, string, error) {
	var errorMessage = "Incorrect email or password"

	// проверка существования пользователя в системе
	user, err := s.Repo.GetByEmail(email)
	if err != nil {
		return "", "", fmt.Errorf(errorMessage)
	}

	// проверка пароля
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", "", errors.New(errorMessage)
	}

	if user.IsActive == false {
		return "", "", errors.New("User is not active")
	}

	accessToken, err := generateToken(user.ID, user.Email, s.jwtAccessSecret, time.Now().Add(time.Hour*72).Unix())
	if err != nil {
		return "", "", err
	}

	refreshToken, err := generateToken(user.ID, user.Email, s.jwtRefreshSecret, time.Now().Add(7*24*time.Hour).Unix())
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func generateToken(userId uint, email, secret string, exp int64) (string, error) {
	refreshTokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userId,
		"email":   email,
		"exp":     exp,
	})

	token, err := refreshTokenClaims.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	
	return token, nil
}
