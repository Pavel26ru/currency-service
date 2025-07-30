package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"
	"github.com/vctrl/currency-service/auth/internal/handler"
	"github.com/vctrl/currency-service/pkg/config"
	"github.com/vctrl/currency-service/pkg/currency"
	"google.golang.org/grpc"
)

type appConfig struct {
	Database config.DatabaseConfig `mapstructure:"database"`
}

func main() {
	configPath := flag.String("config", "./config", "path to the config file")
	flag.Parse()

	cfg, err := loadConfig(*configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	db, err := sql.Open("postgres", cfg.Database.ToDSN())
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()

	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	currency.RegisterAuthServiceServer(grpcServer, handler.NewAuthServer("jwt_secret_key_SoftTradeIT", db))
	log.Println("Auth gRPC server listening on :8080")

	// Запускаем HTTP сервер для метрик в горутине
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		log.Println("Prometheus metrics server running on :8081")
		if err := http.ListenAndServe(":8081", nil); err != nil {
			log.Fatalf("Error starting Prometheus metrics server: %s", err)
		}
	}()

	// Запускаем gRPC сервер
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func loadConfig(path string) (appConfig, error) {
	var config appConfig

	viper.SetConfigFile(path)

	if err := viper.ReadInConfig(); err != nil {
		return config, fmt.Errorf("error reading config file: %w", err)
	}

	if err := viper.Unmarshal(&config); err != nil {
		return config, fmt.Errorf("unable to unmarshal config: %w", err)
	}

	return config, nil
}
