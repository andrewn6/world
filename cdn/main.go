package main

import (
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gofiber/fiber/v2"
	"nijmeh.cloud/routes"
)

func setupRoutes(app *fiber.App) {
	app.Get("/ping", routes.Ping)
	app.Post("/upload", routes.uploadFile)
	app.Delete("/delete", routes.deleteFile)
}

func main() {
	key := "DO008YGT3UV8M3GNA6Z2"
	secret := os.Getenv("SPACES_SECRET")

	s3Config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(key, secret, ""),
		Endpoint:         aws.String("https://nyc3.digitaloceanspaces.com"),
		S3ForcePathStyle: aws.Bool(false),
		Region:           aws.String("us-east-1"),
	}

	newSession := session.New(s3Config)
	s3Client := s3.New(newSession)

	app := fiber.New()

	setupRoutes(app)

	log.Fatal(app.Listen(":8080"))
}
