package handler

import (
	"backend/app"
	"backend/internal/model"

	"github.com/gofiber/fiber/v2"
)

func GetAllSensorData(c *fiber.Ctx) error {
	var sensors []model.DataSensor

	// Get latest 50 records ordered by time desc
	result := app.Database.DB.Order("time DESC").Limit(50).Find(&sensors)
	if result.Error != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to get sensor data",
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   sensors,
		"count":  len(sensors),
	})
}
