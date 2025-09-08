package router

import (
	"backend/internal/handler"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
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

	fmt.Println("*************** SERVER LISTENING ON PORT 3301 ***************")
	log.Fatal(fiber_app.Listen(":3002"))
}
func setupRouter(fiber_app *fiber.App) {
	api := fiber_app.Group("/api")
	api.Get("/device_history", handler.GetAllDeviceHistory)
}
