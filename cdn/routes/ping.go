package routes

import (
    "github.com/gofiber/fiber/v2"
)

func Ping(app fiber.Router) {
    app.Get("/ping") {
      return c.SendString("pong")
    }
}
