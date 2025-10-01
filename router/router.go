package router

import (
	"backend/internal/handler"
	"backend/internal/middleware"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func Setup() {
	fiber_app := fiber.New(fiber.Config{
		Prefork:       false,
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "Fiber",
		AppName:       "IOT Backend v1.0.0",
	})
	fiber_app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization,X-Requested-With",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))
	setupRouter(fiber_app)

	log.Fatal(fiber_app.Listen(":3002"))
}
func setupRouter(fiber_app *fiber.App) {
	auth := fiber_app.Group("/auth")
	auth.Post("/login", handler.Login)
	api := fiber_app.Group("/api")
	api.Use(middleware.JWTMiddleware())
	api.Get("/device_history", handler.GetDeviceHistory)
	api.Get("/data_sensor", handler.GetDataSensor)
	api.Get("/list_devices", handler.ListDevices)
	api.Post("/control_device", handler.ControlDevice)
}
