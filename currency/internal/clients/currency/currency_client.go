package currency

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/vctrl/currency-service/currency/internal/config"

	"go.uber.org/zap"
)

type Currency struct {
	baseURL    *url.URL
	httpClient *http.Client
	logger     *zap.Logger
}

func New(cfg config.APIConfig, logger *zap.Logger) (Currency, error) {
	baseURL, err := url.Parse(cfg.BaseURL)
	if err != nil {
		return Currency{}, fmt.Errorf("invalid base URL: %w", err)
	}

	return Currency{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: time.Duration(cfg.TimeoutSeconds) * time.Second,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // todo ⚠️ НЕ ДЛЯ ПРОДА
			},
		},
		logger: logger,
	}, nil
}

type RatesResponse struct {
	Date string             `json:"date"`
	Rub  map[string]float64 `json:"rub"`
}

func (c *Currency) FetchRatesByDate(ctx context.Context, date, currency string) { // todo
	// todo set metrics for date with value

}

func (c *Currency) FetchCurrentRates(ctx context.Context, currency string) (RatesResponse, error) {
	relativeCurrencyPath, _ := url.Parse(fmt.Sprintf("/v1/currencies/%s.json", strings.ToLower(currency)))
	fullURL := *c.baseURL.ResolveReference(relativeCurrencyPath)

	fullURLStr := fullURL.String()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fullURLStr, nil)
	if err != nil {
		return RatesResponse{}, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)

	if err != nil {
		return RatesResponse{}, fmt.Errorf("failed to make request to currency API: %w", err)
	}

	defer func() {
		err := resp.Body.Close()
		if err != nil {
			c.logger.Error("failed to close response body", zap.Error(err))
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return RatesResponse{}, fmt.Errorf("received non-200 response code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return RatesResponse{}, fmt.Errorf("failed to read response body: %w", err)
	}

	var rateResponse RatesResponse
	if err := json.Unmarshal(body, &rateResponse); err != nil {
		return RatesResponse{}, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// todo set metrics for date with value time.Now()
	return rateResponse, nil
}
