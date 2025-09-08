package router

import (
	"backend/internal/handler"
	"backend/internal/socketio"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
)

func Setup() {
	fiber_app := fiber.New(fiber.Config{
		Prefork:       false,
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "Fiber",
		AppName:       "IOT Backend v1.0.0",
	})

	setupRouter(fiber_app)

	log.Fatal(fiber_app.Listen(":3002"))
}
func setupRouter(fiber_app *fiber.App) {
	api := fiber_app.Group("/api")
	api.Get("/device_history", handler.GetDeviceHistory)
	api.Get("/data_sensor", handler.GetDataSensor)
	api.Get("/list_devices", handler.ListDevices)
	api.Post("/control_device", handler.ControlDevice)
	// Cấu hình Socket.IO
	fiber_app.Use("/socket.io/", adaptor.HTTPHandler(socketio.ServeHTTP()))
}
