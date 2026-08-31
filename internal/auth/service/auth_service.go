package service

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	userRepository "github.com/shivkumar7dandin-gif/students-api/internal/user/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo  *userRepository.UserRepository
	jwtSecret []byte
}

type Claims struct {
	UserID      string `json:"user_id"`
	ReferenceID string `json:"reference_id"`
	Username    string `json:"username"`
	Role        string `json:"role"`

	jwt.RegisteredClaims
}

func NewAuthService(
	userRepo *userRepository.UserRepository,
	jwtSecret string,
) *AuthService {

	return &AuthService{
		userRepo:  userRepo,
		jwtSecret: []byte(jwtSecret),
	}
}

func (s *AuthService) Login(
	ctx context.Context,
	username string,
	password string,
) (string, string, error) {

	// Find user from common users collection
	user, err := s.userRepo.GetByUsername(
		ctx,
		username,
	)
	if err != nil {
		return "", "", fmt.Errorf(
			"invalid username or password",
		)
	}

	// Compare entered password with hashed password
	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(password),
	)
	if err != nil {
		return "", "", fmt.Errorf(
			"invalid username or password",
		)
	}

	// Create JWT claims
	claims := Claims{
		UserID:      user.ID.Hex(),
		ReferenceID: user.ReferenceID.Hex(),
		Username:    user.Username,
		Role:        user.Role,

		RegisteredClaims: jwt.RegisteredClaims{
			Subject: user.ID.Hex(),

			IssuedAt: jwt.NewNumericDate(
				time.Now(),
			),

			ExpiresAt: jwt.NewNumericDate(
				time.Now().Add(24 * time.Hour),
			),
		},
	}

	// Create JWT token
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	// Sign JWT
	tokenString, err := token.SignedString(
		s.jwtSecret,
	)
	if err != nil {
		return "", "", err
	}

	// Return token and user role
	return tokenString, user.Role, nil
}
