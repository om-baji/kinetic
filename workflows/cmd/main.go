package main

import (
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/om-baji/kinetic/shared"
	"github.com/om-baji/kinetic/workflows/internal"
)

func main() {
	app := fiber.New()

	app.Use(shared.Recovery())
	app.Use(shared.Logger())

	svc := internal.NewWorkflowService()
	ctrl := internal.NewWorkflowController(svc)

	internal.RegisterRoutes(app, ctrl)

	log.Fatal(app.Listen(":3000"))
}
