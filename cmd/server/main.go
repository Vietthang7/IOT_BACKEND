package main

import (
	"backend/app"
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

	router.Setup()
}
