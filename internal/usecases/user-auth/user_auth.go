package userauth

import (
	"context"
	"fmt"
	"time"

	"github.com/joyvixtor/dispose-eletronic-waste-backend/internal/adapters/auth"
	"github.com/joyvixtor/dispose-eletronic-waste-backend/internal/adapters/repository"
	"github.com/joyvixtor/dispose-eletronic-waste-backend/internal/domain/models"
)

type authUseCase struct {
	userRepo      *repository.UserRepository
	jwtSecret     string
	jwtExpiration time.Duration
}

func NewAuthService(userRepo *repository.UserRepository, jwtSecret string, jwtExpiration time.Duration) AuthUseCase {
	return &authUseCase{
		userRepo:      userRepo,
		jwtSecret:     jwtSecret,
		jwtExpiration: jwtExpiration,
	}
}

func (s *authUseCase) RegisterUser(ctx context.Context, req *RegisterRequest) (*AuthResponse, error) {
	user := &models.User{
		FullName:  req.FullName,
		UFPEEmail: req.UFPEEmail,
		Password:  req.Password,
		Workplace: req.Workplace,
		Role:      req.Role,
	}

	if err := s.userRepo.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	token, err := auth.GenerateToken(user.Id, user.UFPEEmail, user.Role, s.jwtSecret, s.jwtExpiration)
	if err != nil {
		return nil, err
	}

	user.Password = ""

	return &AuthResponse{
		Token:     token,
		FullName:  user.FullName,
		UFPEEmail: user.UFPEEmail,
		Workplace: user.Workplace,
		Role:      user.Role,
	}, nil
}

func (s *authUseCase) LoginUser(ctx context.Context, req *LoginRequest) (*AuthResponse, error) {
	user, err := s.userRepo.GetUserByEmail(ctx, req.UFPEEmail)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	if err := s.userRepo.VerifyPassword(user.Password, req.Password); err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	token, err := auth.GenerateToken(user.Id, user.UFPEEmail, user.Role, s.jwtSecret, s.jwtExpiration)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	user.Password = ""

	return &AuthResponse{
		Token:     token,
		FullName:  user.FullName,
		UFPEEmail: user.UFPEEmail,
		Workplace: user.Workplace,
		Role:      user.Role,
	}, nil
}
