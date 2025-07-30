package handler

import (
	"context"
	"database/sql"

	"github.com/vctrl/currency-service/pkg/currency"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/vctrl/currency-service/auth/internal/service"
)

type AuthServer struct {
	currency.UnimplementedAuthServiceServer
	userService *service.UserService
}

func NewAuthServer(secret string, db *sql.DB) *AuthServer {
	return &AuthServer{
		userService: service.NewUserService(secret, db),
	}
}

func (s *AuthServer) Login(ctx context.Context, req *currency.LoginRequest) (*currency.LoginResponse, error) {
	if !s.userService.Authenticate(req.Login, req.Password) {
		return nil, status.Error(codes.Unauthenticated, "invalid credentials")
	}
	token, err := s.userService.GenerateJWT(req.Login)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to generate token")
	}
	return &currency.LoginResponse{Token: token}, nil
}

func (s *AuthServer) ValidateToken(ctx context.Context, req *currency.ValidateTokenRequest) (*currency.ValidateTokenResponse, error) {
	_, err := s.userService.ValidateJWT(req.Token)
	if err != nil {
		return &currency.ValidateTokenResponse{Valid: false}, nil
	}
	return &currency.ValidateTokenResponse{Valid: true}, nil
}
