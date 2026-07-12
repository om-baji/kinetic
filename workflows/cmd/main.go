package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/om-baji/kinetic/shared"
	"github.com/om-baji/kinetic/workflows/internal"
)

func main() {
	app := fiber.New()

	app.Use(shared.Recovery())
	app.Use(shared.Logger())

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/kinetic?sslmode=disable"
	}

	db := shared.ConnectDB(dbURL)
	shared.MigrateDB(db)

	svc := internal.NewWorkflowService(db)
	ctrl := internal.NewWorkflowController(svc)

	internal.RegisterRoutes(app, ctrl)

	log.Fatal(app.Listen(":3000"))
}
