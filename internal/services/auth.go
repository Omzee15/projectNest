package services

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"lucid-lists-backend/internal/models"
	"lucid-lists-backend/internal/repositories"
	"lucid-lists-backend/pkg/logger"
)

type AuthService struct {
	userRepo  repositories.UserRepository
	jwtSecret string
}

func NewAuthService(userRepo repositories.UserRepository, jwtSecret string) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

type CustomClaims struct {
	UserID  int       `json:"user_id"`
	UserUID uuid.UUID `json:"user_uid"`
	Email   string    `json:"email"`
	Name    string    `json:"name"`
	jwt.RegisteredClaims
}

func (s *AuthService) Register(ctx context.Context, request models.RegisterRequest) (*models.AuthResponse, error) {
	logger.WithComponent("auth-service").
		WithFields(map[string]interface{}{
			"email": request.Email,
			"name":  request.Name,
		}).
		Info("Registering new user")

	// Check if user already exists
	existingUser, err := s.userRepo.GetByEmail(ctx, request.Email)
	if err == nil && existingUser != nil {
		logger.WithComponent("auth-service").
			WithFields(map[string]interface{}{
				"email": request.Email,
			}).
			Warn("User registration failed - email already exists")
		return nil, fmt.Errorf("user with this email already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.WithComponent("auth-service").
			WithFields(map[string]interface{}{
				"email": request.Email,
				"error": err.Error(),
			}).
			Error("Failed to hash password")
		return nil, fmt.Errorf("failed to hash password")
	}

	// Create user
	user := &models.User{
		Email:    request.Email,
		Password: string(hashedPassword),
		Name:     request.Name,
	}

	err = s.userRepo.Create(ctx, user)
	if err != nil {
		logger.WithComponent("auth-service").
			WithFields(map[string]interface{}{
				"email": request.Email,
				"error": err.Error(),
			}).
			Error("Failed to create user")
		return nil, fmt.Errorf("failed to create user")
	}

	// Generate token
	token, err := s.generateToken(user)
	if err != nil {
		logger.WithComponent("auth-service").
			WithFields(map[string]interface{}{
				"user_uid": user.UserUID.String(),
				"error":    err.Error(),
			}).
			Error("Failed to generate token")
		return nil, fmt.Errorf("failed to generate token")
	}

	response := &models.AuthResponse{
		User: models.UserResponse{
			UserUID:   user.UserUID,
			Name:      user.Name,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
		Token: token,
	}

	logger.WithComponent("auth-service").
		WithFields(map[string]interface{}{
			"user_uid": user.UserUID.String(),
			"email":    user.Email,
		}).
		Info("Successfully registered user")

	return response, nil
}

func (s *AuthService) Login(ctx context.Context, request models.LoginRequest) (*models.AuthResponse, error) {
	logger.WithComponent("auth-service").
		WithFields(map[string]interface{}{
			"email": request.Email,
		}).
		Info("User login attempt")

	// Get user by email
	user, err := s.userRepo.GetByEmail(ctx, request.Email)
	if err != nil {
		logger.WithComponent("auth-service").
			WithFields(map[string]interface{}{
				"email": request.Email,
				"error": err.Error(),
			}).
			Warn("Login failed - user not found")
		return nil, fmt.Errorf("invalid email or password")
	}

	// Check password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
	if err != nil {
		logger.WithComponent("auth-service").
			WithFields(map[string]interface{}{
				"email":    request.Email,
				"user_uid": user.UserUID.String(),
			}).
			Warn("Login failed - invalid password")
		return nil, fmt.Errorf("invalid email or password")
	}

	// Generate token
	token, err := s.generateToken(user)
	if err != nil {
		logger.WithComponent("auth-service").
			WithFields(map[string]interface{}{
				"user_uid": user.UserUID.String(),
				"error":    err.Error(),
			}).
			Error("Failed to generate token")
		return nil, fmt.Errorf("failed to generate token")
	}

	response := &models.AuthResponse{
		User: models.UserResponse{
			UserUID:   user.UserUID,
			Name:      user.Name,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
		Token: token,
	}

	logger.WithComponent("auth-service").
		WithFields(map[string]interface{}{
			"user_uid": user.UserUID.String(),
			"email":    user.Email,
		}).
		Info("Successfully logged in user")

	return response, nil
}

func (s *AuthService) generateToken(user *models.User) (string, error) {
	claims := CustomClaims{
		UserID:  user.ID,
		UserUID: user.UserUID,
		Email:   user.Email,
		Name:    user.Name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   user.UserUID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *AuthService) ValidateToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
