package main

import (
	"backend/app"
	"backend/internal/mqtt"
	"backend/internal/socketio"
	"backend/router"
	"fmt"
	"log"
)

func main() {
	fmt.Println("*************** SERVER MODE ***************")
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Panic recovered: %v", r)
		}
	}()
	app.Setup()
	socketio.Setup() // Setup socketio trước để đăng ký broadcaster
	mqtt.Setup()     // ← Setup sau để sử dụng broadcaster
	router.Setup()
}
