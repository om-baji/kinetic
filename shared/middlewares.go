package shared

import (
	"log"
	"runtime/debug"
	"time"

	"github.com/gofiber/fiber/v3"
)

func Logger() fiber.Handler {
	return func(c fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		log.Printf("%s %s %d %s", c.Method(), c.OriginalURL(), c.Response().StatusCode(), time.Since(start))
		return err
	}
}

func Recovery() fiber.Handler {
	return func(c fiber.Ctx) (err error) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("panic recovered: %v\n%s", r, debug.Stack())
				err = c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
					Error: "internal server error",
				})
			}
		}()
		return c.Next()
	}
}
