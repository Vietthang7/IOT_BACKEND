// internal/handler/device_control.go
package handler

import (
	"backend/internal/consts"
	"backend/internal/mqtt"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func ControlDevice(c *fiber.Ctx) error {
	var req struct {
		DeviceName string `json:"device_name"`
		Action     string `json:"action"`
	}

	if err := c.BodyParser(&req); err != nil {
		return ResponseError(c, fiber.StatusBadRequest,
			fmt.Sprintf("%s: %s", consts.InvalidInput, err.Error()), consts.GetFailed)
	}

	// Kiểm tra dữ liệu
	if req.DeviceName == "" || (req.Action != consts.ACTION_ON && req.Action != consts.ACTION_OFF) {
		return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, consts.GetFailed)
	}

	if err := mqtt.PublishCommand(req.DeviceName, req.Action); err != nil {
		logrus.Error(err)
		return ResponseError(c, fiber.StatusInternalServerError,
			fmt.Sprintf("Lỗi gửi lệnh: %s", err.Error()), consts.GetFailed)
	}

	// Trả về thành công khi đã gửi lệnh
	return ResponseSuccess(c, fiber.StatusOK, "Đã gửi lệnh điều khiển thiết bị", nil)
}
