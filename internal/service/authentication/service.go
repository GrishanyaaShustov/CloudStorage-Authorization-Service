package authentication

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	authorizationservicev1 "github.com/GrishanyaaShustov/CloudStorage-Protos-Service/gen/go/authorization-service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"authorization-service/internal/domain"
	grpcauth "authorization-service/internal/grpc/authentication"
	userrepo "authorization-service/internal/repository/user"
)

// AuthService is a concrete implementation of the authentication Service.
type AuthService struct {
	log   *slog.Logger
	users userrepo.Repository
}

func NewAuthService(log *slog.Logger, users userrepo.Repository) *AuthService {
	return &AuthService{
		log:   log,
		users: users,
	}
}

// Make sure AuthService implements the grpcauth.Service interface.
var _ grpcauth.Service = (*AuthService)(nil)

// --- Stub implementations ---

func (s *AuthService) Register(
	ctx context.Context,
	request *authorizationservicev1.RegisterRequest,
) (*authorizationservicev1.RegisterResponse, error) {

	// 1. Check, user with that email does not exist
	_, err := s.users.GetByEmail(ctx, request.GetEmail())
	if err == nil {
		// пользователь найден → ошибка
		return nil, status.Error(codes.AlreadyExists, "email is already registered")
	}
	if !errors.Is(err, userrepo.ErrNotFound) {
		return nil, status.Error(codes.Internal, "failed to check email")
	}

	// 2. Password hashing
	hash, err := hashPassword(request.GetPassword())
	if err != nil {
		s.log.ErrorContext(ctx, "failed to hash password", slog.Any("err", err))
		return nil, status.Error(codes.Internal, "failed to hash password")
	}

	// 3. Create domain user model
	user := domain.User{
		Email:         request.GetEmail(),
		Login:         request.GetLogin(),
		PasswordHash:  hash,
		EmailVerified: false,
	}

	// 4. Write user in DB
	created, err := s.users.Create(ctx, user)
	if err != nil {
		s.log.ErrorContext(ctx, "failed to create user", slog.Any("err", err))
		return nil, status.Error(codes.Internal, "failed to create user")
	}

	// 5. Turn domain.User to protobuf User
	resp := &authorizationservicev1.RegisterResponse{
		User: &authorizationservicev1.User{
			UserId:        fmt.Sprintf("%d", created.ID),
			Email:         created.Email,
			Login:         created.Login,
			EmailVerified: created.EmailVerified,
			CreatedAt:     timestamppb.New(created.CreatedAt),
			UpdatedAt:     timestamppb.New(created.UpdatedAt),
		},
	}

	s.log.InfoContext(ctx, "Register completed",
		slog.Int64("user_id", created.ID),
	)

	return resp, nil
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
