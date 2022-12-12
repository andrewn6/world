package routes

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gofiber/fiber/v2"
	"net/http"
  "path/filepath"
	"strings"
)

type File struct {
	ID      string `json:"id"`
	Ext     string `json:"ext"`
	Creator string `json:"owner"`
	Size    int64  `json:"size"`
}

func uploadFile(s *session.Session, ctx *fiber.Ctx) (string *JSONResponse) {
	Header, err := ctx.FormFile("file")
	if err != nil {
		return "", utils.NewResponse(fiber.StatusInternalServerError, "Failed to get uploaded file")
	}

	size := Header.Size
	buffer := make([]byte, size)

	uploadedFile, err := Header.Open()

	if err != nil {
		return "", NewResponse(fiber.StatusInternalServerError, "Failed to read!")
	}

	uploadedFile.Close()

	fileName := randSeq(10) + filepath.Ext(Header.Filename)

	object := s3.PutObjectInput{
		Bucket:               aws.String("nijmeh"),
		Key:                  aws.String(fileName),
		ACL:                  aws.String("public-read"),
    Body:                 strings.NewReader(string(buffer)),
		ContentLength:        aws.Int64(size),
		ContentType:          aws.String(http.DetectContentType(buffer)),
		ServerSideEncryption: aws.String("AES256"),
	}

	_, err = s3.New(s).PutObject(&object)

	if err != nil {
		return "", NewResponseByError(fiber.StatusInternalServerError, err)
	}
  return nil
}
