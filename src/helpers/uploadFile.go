package helpers

import (
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func UploadFile(ctx *gin.Context, name string) (string, error) {
	file, err := ctx.FormFile(name)
	if err != nil {
		return "No file is received", err
	}

	filename := filepath.Ext(file.Filename)
	filename = strconv.Itoa(time.Now().Nanosecond()) + filename
	path := "./public/uploads/" + filename
	if err := ctx.SaveUploadedFile(file, path); err != nil {
		return "Unable to save the file", err
	}

	return filename, nil
}
