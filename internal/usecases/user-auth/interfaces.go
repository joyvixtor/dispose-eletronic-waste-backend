package userauth

import "context"

type AuthUseCase interface {
	RegisterUser(ctx context.Context, req *RegisterRequest) (*AuthResponse, error)
	LoginUser(ctx context.Context, req *LoginRequest) (*AuthResponse, error)
}
