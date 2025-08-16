package handler

import (
	"net/http"

	"github.com/vctrl/currency-service/gateway/internal/dto"

	"context"

	"github.com/gin-gonic/gin"
	"github.com/vctrl/currency-service/pkg/currency"
)

type registerRequest struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}

type AuthService struct {
	authClient currency.AuthServiceClient
}

func (s *AuthService) Login(ctx context.Context, login, password string) (string, error) {
	resp, err := s.authClient.Login(ctx, &currency.LoginRequest{
		Login:    login,
		Password: password,
	})
	if err != nil {
		return "", err
	}
	return resp.Token, nil
}

func (s *controller) Register(c *gin.Context) {
	var req registerRequest
	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = s.authService.Register(dto.RegisterRequest(req))
	if err != nil {
		s.handleError(c, err)
		return
	}

	c.Status(http.StatusCreated)
}

type loginRequest struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}

func (s *controller) Login(c *gin.Context) {
	var req loginRequest
	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	token, err := s.authService.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		s.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (s *controller) Logout(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Authorization token is required"})
		return
	}

	err := s.authService.Logout(token)
	if err != nil {
		s.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "logout successful"})
}
