package main

import (
	"log"
	"net"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
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
	currency.RegisterAuthServiceServer(grpcServer, handler.NewAuthServer("jwt_secret_key_SoftTradeIT"))
	log.Println("Auth gRPC server listening on :8080")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		log.Println("Prometheus metrics server running on :8081")
		if err := http.ListenAndServe(":8081", nil); err != nil {
			log.Fatalf("Error starting Prometheus metrics server: %s", err)
		}
	}()
}
