package routes

import (
  "math/rand"
  "time"
)
  
type JSONResponse struct {
  Success bool `json:"success"`
  Code int `json:"code"`
  Message string `json:"message,omitempty"`
  Data interface{} `json:"data,omitempty"`
}

func NewResponse(code int, message string) *JSONResponse {
  return &JSONResponse {
    Message: message,
    Code: code,
  }
}

func NewResponseByError(code int, err error) *JSONResponse {
  return &JSONResponse {
    Message: err.Error(),
    Code: code,
  }
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
// @acollier17
func randSeq(n int) string {
  rand.Seed(time.Now().UnixNano())
  b := make([]rune, n)

  for i := range b {
    b[i] = letters[rand.Intn(len(letters))]
  }

  return string(b)
}

