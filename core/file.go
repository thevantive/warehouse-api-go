package core

import (
	"errors"
	"fmt"
	"mime/multipart"
	"gudang-api-go-v1/models"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SaveUploadedFile(ctx *gin.Context, db *gorm.DB, row *models.File, name string, group string, file *multipart.FileHeader, ext string) error {

	// menyimpan file ke local storage
	filename := fmt.Sprintf("%s-%s-%s%s", name, group, time.Now().Format("20060102150405"), ext)
	if err := ctx.SaveUploadedFile(file, "uploads/"+filename); err != nil {
		return fmt.Errorf("failed to save file: %s", err.Error())
	}

	// menyimpan informasi file ke database ke database
	*row = models.File{
		Name:      name,
		Filename:  filename,
		Filetype:  file.Header.Get("Content-Type"),
		Size:      0,
		Status:    "active",
		CreatedBy: uint(1),
	}
	if err := db.Create(row).Error; err != nil {
		return errors.New("failed to create attachment record")
	}

	return nil
}
