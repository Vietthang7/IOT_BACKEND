package router

import (
	"backend/internal/handler"
	"backend/internal/socketio"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
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
	api := fiber_app.Group("/api")
	api.Get("/device_history", handler.GetDeviceHistory)
	api.Get("/data_sensor", handler.GetDataSensor)
	api.Get("/list_devices", handler.ListDevices)
	api.Post("/control_device", handler.ControlDevice)
	// Cấu hình Socket.IO với CORS headers
	fiber_app.Use("/socket.io/*", func(c *fiber.Ctx) error {
		// Thêm CORS headers cho Socket.IO
		c.Set("Access-Control-Allow-Origin", "http://localhost:3000")
		c.Set("Access-Control-Allow-Credentials", "true")
		c.Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		c.Set("Access-Control-Allow-Headers", "Origin,Content-Type,Accept,Authorization")

		if c.Method() == "OPTIONS" {
			return c.SendStatus(200)
		}

		return adaptor.HTTPHandler(socketio.ServeHTTP())(c)
	})
	fiber_app.Use("/socket.io/", adaptor.HTTPHandler(socketio.ServeHTTP()))
}
