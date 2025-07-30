package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/vctrl/currency-service/pkg/currency"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	ErrUnexpectedStatusCode  = errors.New("unexpected status code")
	ErrInvalidCredentials    = errors.New("invalid credentials")
	ErrTokenGeneration       = errors.New("failed to generate token")
	ErrTokenNotFound         = errors.New("token not found")
	ErrInvalidOrExpiredToken = errors.New("token is invalid or expired")
)

func Example() {
	conn, err := grpc.NewClient(
		"localhost:8082",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	authClient := currency.NewAuthServiceClient(conn)
	resp, err := authClient.Login(context.Background(), &currency.LoginRequest{
		Login:    "user",
		Password: "pass",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("Token:", resp.Token)

	// Используйте resp.Token, resp.Success, resp.Error
}
