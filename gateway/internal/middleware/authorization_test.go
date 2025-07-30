package middleware

import (
	"net/http/httptest"
	"testing"

	"context"
	"errors"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type MockAuthClient struct{}

func NewMockClient() *MockAuthClient {
	return &MockAuthClient{}
}

func (m *MockAuthClient) ValidateToken(ctx context.Context, token string) error {
	if token == "TOKEN" {
		return nil
	}
	return errors.New("invalid token")
}

func TestAuth(t *testing.T) {
	client := NewMockClient()
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request.Header.Set("Authorization", "Bearer TOKEN")
	// auth := NewAuthorization()
	auth := NewAuthorization(
		client,
		func(c *gin.Context) bool { return false }, // не пропускаем без авторизации
		zap.NewNop(),
	)

	auth.Authorize()(ctx)

}
