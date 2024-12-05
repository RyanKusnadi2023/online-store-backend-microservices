package service

import (
	"context"
	"errors"
	"os"
	"strconv"
	"time"

	"synapsis-online-store/services/user-service/internal/repository"
	"synapsis-online-store/services/user-service/internal/utils"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Register(ctx context.Context, username, email, password string) error
	Login(ctx context.Context, email, password string) (string, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

// Register handles user registration logic
func (s *userService) Register(ctx context.Context, username, email, password string) error {
	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Check if email already exists
	existingUser, _ := s.repo.GetUserByEmail(ctx, email)
	if existingUser != nil {
		return errors.New("email already in use")
	}

	// Save the user
	user := &repository.User{
		Username: username,
		Email:    email,
		Password: string(hashedPassword),
	}
	return s.repo.CreateUser(ctx, user)
}

// Login handles user login logic
func (s *userService) Login(ctx context.Context, email, password string) (string, error) {
	// Retrieve the user
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil || user == nil {
		return "", errors.New("invalid email or password")
	}

	// Compare the password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	// Generate JWT
	secret := os.Getenv("JWT_SECRET")
	expirationStr := os.Getenv("JWT_EXPIRATION")
	if secret == "" || expirationStr == "" {
		return "", errors.New("JWT_SECRET or JWT_EXPIRATION not set")
	}

	expiration, err := time.ParseDuration(expirationStr)
	if err != nil {
		return "", errors.New("invalid JWT_EXPIRATION format")
	}

	token, err := utils.GenerateJWT(strconv.Itoa(user.ID), secret, expiration)
	if err != nil {
		return "", err
	}

	return token, nil
}
