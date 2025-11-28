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

func (s *Server) Register(ctx context.Context, request *authorizationservicev1.RegisterRequest) (*authorizationservicev1.RegisterResponse, error) {
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "request is nil")
	}

	// Basic input validation on transport layer.
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
