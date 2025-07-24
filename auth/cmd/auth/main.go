package main

import (
	"log"
	"net"

	"github.com/vctrl/currency-service/auth/internal/handler"
	"github.com/vctrl/currency-service/pkg/currency"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	currency.RegisterAuthServiceServer(grpcServer, handler.NewAuthServer("your_jwt_secret"))
	log.Println("Auth gRPC server listening on :8080")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
