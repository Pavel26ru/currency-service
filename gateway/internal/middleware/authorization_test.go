package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http/httptest"
	"testing"
)

func TestAuth(t *testing.T) {
	client := NewMOckCllient()
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request.Header.Set("Authorization", "Bearer TOKEN")
	auth := NewAuthorization()

	// todo add auth header to
	auth.Authorize()(ctx)

}
