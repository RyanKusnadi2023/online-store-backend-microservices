package middleware

import (
	"context"
	"errors"
	"strings"

	"synapsis-online-store/services/user-service/internal/utils"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// UnaryJWTInterceptor validates JWT for incoming gRPC requests
func UnaryJWTInterceptor(secret string) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// Extract the token from metadata
		token, err := extractTokenFromMetadata(ctx)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, err.Error())
		}

		// Validate the JWT
		_, err = utils.ParseJWT(token, secret)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "invalid or expired token")
		}

		// Proceed to the actual handler
		return handler(ctx, req)
	}
}

// extractTokenFromMetadata extracts the JWT token from gRPC metadata
func extractTokenFromMetadata(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", errors.New("missing metadata")
	}

	authHeader := md["authorization"]
	if len(authHeader) == 0 {
		return "", errors.New("authorization token not provided")
	}

	// Extract the token (e.g., "Bearer <token>")
	parts := strings.Split(authHeader[0], " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", errors.New("invalid authorization header format")
	}

	return parts[1], nil
}
