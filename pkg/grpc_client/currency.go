package grpc_client

import (
	"fmt"

	"github.com/vctrl/currency-service/pkg/currency"

	"google.golang.org/grpc/credentials/insecure"

	"google.golang.org/grpc"
)

func NewCurrencyServiceClient(address string) (currency.CurrencyServiceClient, *grpc.ClientConn, error) {
	dialOptions := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	conn, err := grpc.NewClient(address, dialOptions...)
	if err != nil {
		return nil, nil, fmt.Errorf("grpc.NewClient: %w", err)
	}

	client := currency.NewCurrencyServiceClient(conn)
	return client, conn, nil
}
