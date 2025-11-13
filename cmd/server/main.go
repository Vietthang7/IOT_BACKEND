package main

import (
	"backend/app"
	"backend/internal/mqtt"
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
	mqtt.Setup() // ← Setup sau để sử dụng broadcaster
	router.Setup()
}
