package main

import (
	"context"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/vctrl/currency-service/gateway/internal/config"
	"github.com/vctrl/currency-service/gateway/internal/handler"
	"github.com/vctrl/currency-service/gateway/internal/middleware"
	"github.com/vctrl/currency-service/gateway/internal/repository"
	"github.com/vctrl/currency-service/gateway/internal/service"
	"github.com/vctrl/currency-service/pkg/currency"
	"github.com/vctrl/currency-service/pkg/grpc_client"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthClientAdapter struct {
	grpcClient currency.AuthServiceClient
}

func (a *AuthClientAdapter) GenerateToken(ctx context.Context, login, password string) (string, error) {
	resp, err := a.grpcClient.Login(ctx, &currency.LoginRequest{
		Login:    login,
		Password: password,
	})
	if err != nil {
		return "", err
	}
	return resp.Token, nil
}

func (a *AuthClientAdapter) ValidateToken(ctx context.Context, token string) error {
	resp, err := a.grpcClient.ValidateToken(ctx, &currency.ValidateTokenRequest{Token: token})
	if err != nil {
		return err
	}

	if !resp.Valid {
		return fmt.Errorf("invalid token")
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err.Error())
	}
}

func run() error {
	configPath := flag.String("config", "./config", "path to the config file")

	flag.Parse()

	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	logger, err := zap.NewProduction()
	if err != nil {
		return fmt.Errorf("init logger: %w", err)
	}

	logger.Info("Server initializing with config", zap.Any("config", cfg))

	router := gin.New()
	router.Use(ginzap.Ginzap(logger, time.RFC3339, true))
	router.Use(ginzap.RecoveryWithZap(logger, true))

	// auth middleware
	authConn, err := grpc.NewClient(
		"localhost:8082",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return fmt.Errorf("cannot connect to AuthService: %w", err)
	}
	defer authConn.Close()

	authClient := currency.NewAuthServiceClient(authConn)

	authAdapter := &AuthClientAdapter{grpcClient: authClient}
	authMiddleware := middleware.NewAuthorization(authAdapter, shouldSkipAuthMiddleware, logger)
	router.Use(authMiddleware.Authorize())

	currencyClient, conn, err := grpc_client.NewCurrencyServiceClient(cfg.GRPC.CurrencyServiceURL)
	if err != nil {
		return fmt.Errorf("grpc_client.NewCurrencyServiceClient: %w", err)
	}

	defer func() {
		if err := conn.Close(); err != nil {
			logger.Warn("Cannot close GRPC Client for auth service", zap.Error(err))
		}
	}()

	/*
		userRepo := user.NewRepository()
		authService := auuth.NewService(authClient, userRepo)
		currencyService := currency.NewService(currencyClient)
	*/

	db, err := sql.Open("postgres", cfg.DatabaseConfig.ToDSN())
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	userRepo := repository.NewUser(db)
	authService := service.NewAuth(authAdapter, userRepo)
	currencyService := service.NewCurrency(currencyClient)

	srv := &http.Server{
		Addr:    cfg.Server.Port,
		Handler: router,
	}

	handler.RegisterRoutes(authService, currencyService, router, logger)

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("srv.ListenAndServe: %s\n", err)
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()
	log.Println("shutting down gracefully, press Ctrl+C again to force")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("shutdown: %w", err)
	}

	log.Println("Server exiting")
	return nil
}

func shouldSkipAuthMiddleware(c *gin.Context) bool {
	if strings.HasSuffix(c.Request.URL.Path, "/login") ||
		strings.HasSuffix(c.Request.URL.Path, "/register") {
		return true
	}

	return false
}
