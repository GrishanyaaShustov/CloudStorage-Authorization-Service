package authentication

import (
	"context"
	"log/slog"

	authorizationservicev1 "github.com/GrishanyaaShustov/CloudStorage-Protos-Service/gen/go/authorization-service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Service describes authentication business logic.
// gRPC handlers should be thin and delegate all work to this interface.
type Service interface {
	Register(ctx context.Context, request *authorizationservicev1.RegisterRequest) (*authorizationservicev1.RegisterResponse, error)
	VerifyEmail(ctx context.Context, request *authorizationservicev1.VerifyEmailRequest) (*authorizationservicev1.VerifyEmailResponse, error)
	Login(ctx context.Context, request *authorizationservicev1.LoginRequest) (*authorizationservicev1.LoginResponse, error)
	RefreshToken(ctx context.Context, request *authorizationservicev1.RefreshTokenRequest) (*authorizationservicev1.RefreshTokenResponse, error)
	Logout(ctx context.Context, request *authorizationservicev1.LogoutRequest) (*authorizationservicev1.LogoutResponse, error)
}

// Server is a gRPC transport for AuthenticationService.
// It delegates all business logic to the Service interface.
type Server struct {
	authorizationservicev1.UnimplementedAuthenticationServiceServer
	log     *slog.Logger
	service Service
}

// NewServer constructs a new Authentication gRPC server.
// Logger and Service are injected from the caller (main or wiring layer).
func NewServer(log *slog.Logger, service Service) *Server {
	return &Server{
		log:     log,
		service: service,
	}
}

// Register handles user registration and delegates all business logic
// to the Service implementation.
func (s *Server) Register(ctx context.Context, request *authorizationservicev1.RegisterRequest) (*authorizationservicev1.RegisterResponse, error) {
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "request is nil")
	}

	// Basic input validation on transport layer.

	if request.GetLogin() == "" {
		return nil, status.Error(codes.InvalidArgument, "login is required")
	}

	if request.GetEmail() == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	if request.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	s.log.InfoContext(ctx, "Register called",
		"email", request.GetEmail(),
		"login", request.GetLogin(),
	)

	resp, err := s.service.Register(ctx, request)
	if err != nil {
		// The service is expected to return a gRPC-aware error (status.Error),
		// but just in case we still log it here.
		s.log.ErrorContext(ctx, "Register failed",
			"email", request.GetEmail(),
			"login", request.GetLogin(),
		)
		return nil, err
	}
	return resp, nil
}

// Login authenticates user by email and password and delegates
// the actual logic to the Service implementation.
func (s *Server) Login(ctx context.Context, request *authorizationservicev1.LoginRequest) (*authorizationservicev1.LoginResponse, error) {
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "request is nil")
	}

	if request.GetEmail() == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	if request.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	s.log.InfoContext(ctx, "Login called",
		"email", request.GetEmail(),
		"client_id", request.GetClientId(),
	)
	resp, err := s.service.Login(ctx, request)

	if err != nil {
		// The service is expected to return a gRPC-aware error (status.Error),
		// but just in case we still log it here.
		s.log.ErrorContext(ctx, "Login failed failed",
			"email", request.GetEmail(),
			"client_id", request.GetClientId(),
		)

		return nil, err
	}

	return resp, err
}

// RefreshToken exchanges a refresh token for a new set of tokens.
// The handler only does basic validation and logging and delegates
// the actual refresh flow to the Service implementation.
func (s *Server) RefreshToken(ctx context.Context, request *authorizationservicev1.RefreshTokenRequest) (*authorizationservicev1.RefreshTokenResponse, error) {
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "request is nil")
	}

	if request.GetRefreshToken() == "" {
		return nil, status.Error(codes.InvalidArgument, "refresh token is required")
	}

	s.log.InfoContext(ctx, "RefreshToken called",
		"client_id", request.GetClientId(),
	)

	resp, err := s.service.RefreshToken(ctx, request)

	if err != nil {
		s.log.ErrorContext(ctx, "RefreshToken failed",
			"client_id", request.GetClientId(),
		)

		return nil, err
	}
	return resp, nil
}

// Logout performs explicit logout by revoking the refresh token / session.
// The handler validates input, logs the call and delegates the actual
// revocation logic to the Service implementation.
func (s *Server) Logout(ctx context.Context, request *authorizationservicev1.LogoutRequest) (*authorizationservicev1.LogoutResponse, error) {
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "request is nil")
	}

	if request.GetRefreshToken() == "" {
		return nil, status.Error(codes.InvalidArgument, "refresh_token is required")
	}

	s.log.InfoContext(ctx, "Logout called",
		"client_id", request.GetClientId(),
	)

	resp, err := s.service.Logout(ctx, request)

	if err != nil {
		s.log.ErrorContext(ctx, "Logout failed",
			"client_id", request.GetClientId(),
		)
		return nil, err
	}

	return resp, nil
}

// VerifyEmail confirms user email using a verification code or flow identifier.
// The handler only checks basic input and delegates the rest to the Service.
func (s *Server) VerifyEmail(ctx context.Context, request *authorizationservicev1.VerifyEmailRequest) (*authorizationservicev1.VerifyEmailResponse, error) {
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "request is nil")
	}

	if request.GetVerificationCode() == "" {
		return nil, status.Error(codes.InvalidArgument, "verification_code is required")
	}

	if request.GetFlowId() == "" {
		return nil, status.Error(codes.InvalidArgument, "flow_id is required")
	}

	s.log.InfoContext(ctx, "VerifyEmail called")

	resp, err := s.service.VerifyEmail(ctx, request)
	if err != nil {
		s.log.ErrorContext(ctx, "VerifyEmail failed")
		return nil, err
	}

	return resp, nil
}
