package main

import (
  "log"

  "github.com/gofiber/fiber/v2"
  "nijmeh.cloud/routes"
)

func setupRoutes(app *fiber.App) {
  app.Get("/ping", routes.Ping)
}

func main() {
  app := fiber.New()

  setupRoutes(app)

  log.Fatal(app.Listen(":8080"))
}
