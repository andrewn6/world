package cdn

import ( 

  "github.com/gofiber/fiber/v2"
  "encoding/json"
  "github.com/cloudflare/cloudflare-go"

  "fmt"
  "bytes"
  "io"
  "os"
  "path/filepath"
  "net/http"
)


const uploadDirectory = "./uploads"

const (
	uploadURL = "https://api.cloudflare.com/client/v4/accounts/{account_id}/storage/kv/namespaces/{namespace_id}/values/{key}"
)

func main () {
  app := fiber.New()
  
  app.Get("/upload", func(c *fiber.Ctx) {
    c.Send("<form action='/upload' method='post' enctype='multipart/form-data'>" +
			"<input type='file' name='file'>" +
			"<input type='submit' value='Upload'>" +
			"</form>")
  })

  app.Post("/upload", func(c *fiber.Ctx) {
    file, err := c.FormFile("file")
    
    if err != nil {
      c.Send(fmt.Sprintf("Error getting file: %s", err))
      return
    }

    src, err := file.Open()

    if err != nil {
      c.Send(fmt.Sprintf("error opening file: %s", err))
      return
    }

    defer src.Close()

    size, err := file.Size()
    if err != nil {
      c.Send(fmt.Sprintf("error getting file: %s", err))
      return
    }

    buf := make([]byte, size)
    _, err = src.Read(buf)

    if err != nil {
      c.Send(fmt.Sprintf("Error reading file: %s", err))
      return
    }

    req.Header.Set("Content-Type", file.Header.Get("Content-Type"))
    req.Header.Set("Authorization", "Bearer " +os.Getenv("CLOUDFLARE_API_KEY"))

    client := &http.Client{}
    
    resp, err := client.Do()
  }
  

  app.Listen(8080)
}
