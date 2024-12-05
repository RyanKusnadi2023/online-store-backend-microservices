package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"synapsis-online-store/services/user-service/internal/middleware"
	"synapsis-online-store/services/user-service/internal/proto"
	"synapsis-online-store/services/user-service/internal/repository"
	"synapsis-online-store/services/user-service/internal/service"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type userServer struct {
	proto.UnimplementedUserServiceServer
	service service.UserService
}

// Register handles user registration
func (s *userServer) Register(ctx context.Context, req *proto.RegisterRequest) (*proto.RegisterResponse, error) {
	err := s.service.Register(ctx, req.Username, req.Email, req.Password)
	if err != nil {
		return nil, err
	}

	return &proto.RegisterResponse{
		Message: "User registered successfully!",
	}, nil
}

// Login handles user authentication
func (s *userServer) Login(ctx context.Context, req *proto.LoginRequest) (*proto.LoginResponse, error) {
	token, err := s.service.Login(ctx, req.Email, req.Password)
	if err != nil {
		return nil, err
	}

	return &proto.LoginResponse{
		Token:  token,
		UserId: req.Email,
	}, nil
}

func main() {
	// Load environment variables from .env file (optional)
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("No .env file found, using system environment variables.")
	}

	// Ensure JWT environment variables are set
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatalf("JWT_SECRET must be set in the environment")
	}

	// Connect to the database
	db, err := getDBConnection()
	if err != nil {
		log.Fatalf("failed to connect to the database: %v", err)
	}
	defer db.Close()

	// Initialize repository and service
	repo := repository.NewUserRepository(db)
	userService := service.NewUserService(repo)

	// Set up a gRPC server WITHOUT middleware for public routes (Register & Login)
	publicServer := grpc.NewServer()
	reflection.Register(publicServer)
	proto.RegisterUserServiceServer(publicServer, &userServer{service: userService})

	// Set up a gRPC server WITH middleware for protected routes
	protectedServer := grpc.NewServer(
		grpc.UnaryInterceptor(middleware.UnaryJWTInterceptor(jwtSecret)),
	)
	reflection.Register(protectedServer)
	proto.RegisterUserServiceServer(protectedServer, &userServer{service: userService})

	// Listen on two different ports
	go func() {
		lis, err := net.Listen("tcp", ":50051")
		if err != nil {
			log.Fatalf("failed to listen on port 50051: %v", err)
		}

		log.Println("Public User Service is running on port 50051...")
		if err := publicServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve public server: %v", err)
		}
	}()

	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("failed to listen on port 50052: %v", err)
	}

	log.Println("Protected User Service is running on port 50052...")
	if err := protectedServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve protected server: %v", err)
	}
}

// getDBConnection initializes a PostgreSQL connection pool
func getDBConnection() (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		getEnv("DB_USER"),
		getEnv("DB_PASSWORD"),
		getEnv("DB_HOST"),
		getEnv("DB_PORT"),
		getEnv("DB_NAME"),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}
	return db, nil
}

// getEnv retrieves the value of the environment variable, or throws an error if not set
func getEnv(key string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	log.Fatalf("Environment variable %s not set", key)
	return "" // This line will never be reached due to log.Fatalf
}
