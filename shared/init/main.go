package main

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("[KINETIC] Running Initial Checks and Migrations!")
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	fmt.Println("[KINETIC] Checks Complete")
}
