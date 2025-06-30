package handler

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vctrl/currency-service/currency/internal/handler/mocks"
	"github.com/vctrl/currency-service/pkg/currency"
	"go.uber.org/zap/zaptest"
	"testing"
)

func TestGetRate(t *testing.T) {
	ctrl := gomock.NewController(t)
	service := mocks.NewMockCurrencyService(ctrl)

	service.EXPECT().GetCurrencyRatesInInterval(gomock.Any(), &currency.GetRateRequest{
		Currency: "",
		DateFrom: nil,
		DateTo:   nil,
		// сюда добавишь новое поле когда
	}).
		Return(&currency.GetRateResponse{}, nil)
	//service := &CurrencyServiceMock{}
	//service.ExpectCall("GetCurrencyRatesInInterval")

	requestCount := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "test_request_count",
			Help: "Test count",
		},
		[]string{"method"},
	)
	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "test_request_duration",
			Help:    "Test duration",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)
	appUptime := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "test_app_uptime",
		Help: "Test app uptime",
	})

	logger := zaptest.NewLogger(t)
	server := NewCurrencyServer(service,
		logger,
		requestCount,
		requestDuration,
		appUptime,
	)

	expected := &currency.GetRateResponse{}

	ctx := context.Background()
	req := &currency.GetRateRequest{
		Currency: "",
		DateFrom: nil,
		DateTo:   nil,
	}

	fact, err := server.GetRate(ctx, req)

	require.NoError(t, err)

	assert.Equal(t, expected, fact)
}
