// internal/handler/device_control.go
package handler

import (
	"backend/internal/consts"
	"backend/internal/mqtt"
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func ControlDevice(c *fiber.Ctx) error {
	var req struct {
		DeviceName string `json:"device_name"`
		Action     string `json:"action"`
	}
	// if c.Locals("MQTT_status") != consts.MQTT_CONNECTED {
	// 	return ResponseError(c, fiber.StatusInternalServerError,
	// 		"Lỗi điều khiển thiết bị", consts.GetFailed)
	// }
	if err := c.BodyParser(&req); err != nil {
		return ResponseError(c, fiber.StatusBadRequest,
			fmt.Sprintf("%s: %s", consts.InvalidInput, err.Error()), consts.GetFailed)
	}

	if req.DeviceName == "" || (req.Action != consts.ACTION_ON && req.Action != consts.ACTION_OFF) {
		return ResponseError(c, fiber.StatusBadRequest, consts.InvalidInput, consts.GetFailed)
	}
	ctx, cancel := context.WithTimeout(context.Background(), consts.TIMEOUT)
	defer cancel()
	// Gửi lệnh và chờ xác nhận
	if err := mqtt.PublishCommandAndWait(ctx, req.DeviceName, req.Action); err != nil {
		logrus.Error(err)
		if err = context.DeadlineExceeded; err != nil {
			return ResponseError(c, fiber.StatusRequestTimeout, " Lỗi gửi lệnh: Hết thời gian chờ phản hồi từ thiết bị", consts.GetFailed)
		}
		return ResponseError(c, fiber.StatusInternalServerError,
			fmt.Sprintf("Lỗi điều khiển thiết bị: %s", err.Error()), consts.GetFailed)
	}
	// CHỈ KHI ESP32 XÁC NHẬN ĐÃ THỰC THI THÀNH CÔNG
	return ResponseSuccess(c, fiber.StatusOK,
		fmt.Sprintf("Thiết bị %s đã %s thành công", req.DeviceName, req.Action), nil)
}
