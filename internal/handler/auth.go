package handler

import (
	"backend/app"
	"backend/internal/auth"
	"backend/internal/consts"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Token    string `json:"token"`
	Username string `json:"username"`
	UserID   string `json:"user_id"`
}

func Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, consts.GetFailed)
	}
	adminUsername := app.Config("ADMIN_USERNAME")
	adminPassword := app.Config("ADMIN_PASSWORD")
	if adminUsername == "" {
		adminUsername = "admin"
	}
	if adminPassword == "" {
		adminPassword = "123456"
	}
	if req.Username != adminUsername || req.Password != adminPassword {
		return ResponseError(c, fiber.StatusUnauthorized, "Tài khoản hoặc mật khẩu không đúng", consts.GetFailed)
	}
	// Generate JWT token
	userID := uuid.New() // Tạo UUID cho admin user
	token, err := auth.GenerateToken(userID, req.Username)
	if err != nil {
		return ResponseError(c, fiber.StatusInternalServerError, "Lỗi tạo token", consts.GetFailed)
	}
	response := LoginResponse{
		Token:    token,
		Username: req.Username,
		UserID:   userID.String(),
	}
	return ResponseSuccess(c, fiber.StatusOK, "Đăng nhập thành công", response)
}
