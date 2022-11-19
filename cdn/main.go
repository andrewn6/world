package main

import (
  "log"
  "os"
  "fmt"
  "strings"

  "github.com/gofiber/fiber/v2"
  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/aws/credentials"
  "github.com/aws/aws-sdk-go/aws/session"
  "github.com/aws/aws-sdk-go/service/s3"
  "nijmeh.cloud/routes"
)

func setupRoutes(app *fiber.App) {
  app.Get("/ping", routes.Ping)
}

func main() {
  key := "DO008YGT3UV8M3GNA6Z2"
  secret := os.Getenv("SPACES_SECRET")
  
  s3Config := &aws.Config {
    Credentials: credentials.NewStaticCredentials(key, secret, ""),
    Endpoint: aws.String("https://nyc3.digitaloceanspaces.com"),
    S3ForcePathStyle: aws.Bool(false),
    Region: aws.String("us-east-1"),
  }
  
  newSession := session.New(s3Config)
  s3Client := s3.New(newSession)
  
  object := s3.PutObjectInput {
    Bucket: aws.String("nijmeh"),
    Key: aws.String("folder/flag-orpheus-top.svg"),
    Body: strings.NewReader("Hello, World!"),
    ACL: aws.String("private"),
    Metadata: map[string]*string{
                            "x-amz-meta-my-key": aws.String("your-value"),
                    },
  }

  _, err := s3Client.PutObject(&object)
  if err != nil {
    fmt.Println(err.Error())
  }
  

  app := fiber.New()

  setupRoutes(app)

  log.Fatal(app.Listen(":8080"))
}
