package authentication

import (
	"context"
	"log/slog"

	authorizationservicev1 "github.com/GrishanyaaShustov/CloudStorage-Protos-Service/gen/go/authorization-service"

	grpcauth "authorization-service/internal/grpc/authentication"
)

// AuthService is a concrete implementation of the authentication Service.
type AuthService struct {
	log *slog.Logger
}

func NewAuthService(log *slog.Logger) *AuthService {
	return &AuthService{
		log: log,
	}
}

// Make sure AuthService implements the grpcauth.Service interface.
var _ grpcauth.Service = (*AuthService)(nil)

// --- Stub implementations ---

func (s *AuthService) Register(
	ctx context.Context,
	request *authorizationservicev1.RegisterRequest,
) (*authorizationservicev1.RegisterResponse, error) {
	// TODO: implement registration using Kratos.
	s.log.InfoContext(ctx, "Register not implemented yet")
	return nil, statusUnimplemented("Register")
}

func (s *AuthService) VerifyEmail(
	ctx context.Context,
	request *authorizationservicev1.VerifyEmailRequest,
) (*authorizationservicev1.VerifyEmailResponse, error) {
	// TODO: implement email verification using Kratos.
	s.log.InfoContext(ctx, "VerifyEmail not implemented yet")
	return nil, statusUnimplemented("VerifyEmail")
}

func (s *AuthService) Login(
	ctx context.Context,
	request *authorizationservicev1.LoginRequest,
) (*authorizationservicev1.LoginResponse, error) {
	// TODO: implement login via Kratos + token issuance via Hydra.
	s.log.InfoContext(ctx, "Login not implemented yet")
	return nil, statusUnimplemented("Login")
}

func (s *AuthService) RefreshToken(
	ctx context.Context,
	request *authorizationservicev1.RefreshTokenRequest,
) (*authorizationservicev1.RefreshTokenResponse, error) {
	// TODO: implement refresh via Hydra.
	s.log.InfoContext(ctx, "RefreshToken not implemented yet")
	return nil, statusUnimplemented("RefreshToken")
}

func (s *AuthService) Logout(
	ctx context.Context,
	request *authorizationservicev1.LogoutRequest,
) (*authorizationservicev1.LogoutResponse, error) {
	// TODO: implement logout (refresh token/session revocation) via Hydra.
	s.log.InfoContext(ctx, "Logout not implemented yet")
	return nil, statusUnimplemented("Logout")
}
