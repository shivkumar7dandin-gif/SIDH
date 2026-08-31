package service

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/shivkumar7dandin-gif/students-api/internal/college/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	collegeRepo *repository.CollegeRepository
	jwtSecret   []byte
}

type Claims struct {
	CollegeID string `json:"college_id"`
	Username  string `json:"username"`
	Role      string `json:"role"`

	jwt.RegisteredClaims
}

func NewAuthService(
	collegeRepo *repository.CollegeRepository,
	jwtSecret string,
) *AuthService {

	return &AuthService{
		collegeRepo: collegeRepo,
		jwtSecret:   []byte(jwtSecret),
	}
}

func (s *AuthService) Login(
	ctx context.Context,
	username string,
	password string,
) (string, string, error) {

	college, err := s.collegeRepo.GetByUsername(
		ctx,
		username,
	)
	if err != nil {
		return "", "", fmt.Errorf(
			"invalid username or password",
		)
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(college.PasswordHash),
		[]byte(password),
	)
	if err != nil {
		return "", "", fmt.Errorf(
			"invalid username or password",
		)
	}

	claims := Claims{
		CollegeID: college.ID.Hex(),
		Username:  college.Username,
		Role:      college.Role,

		RegisteredClaims: jwt.RegisteredClaims{
			Subject: college.ID.Hex(),

			IssuedAt: jwt.NewNumericDate(
				time.Now(),
			),

			ExpiresAt: jwt.NewNumericDate(
				time.Now().Add(24 * time.Hour),
			),
		},
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	tokenString, err := token.SignedString(
		s.jwtSecret,
	)
	if err != nil {
		return "", "", err
	}

	return tokenString, college.Role, nil
}
