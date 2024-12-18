package routes

import (
	"gudang-api-go-v1/models"
	"gudang-api-go-v1/rest/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Product(router *gin.RouterGroup, dbs map[string]*gorm.DB) {

	// endpoint untuk mengambil daftar produk
	router.GET("/products", func(ctx *gin.Context) {
		db := dbs["main"].Session(&gorm.Session{})

		var products []models.Product
		if err := db.Find(&products).Error; err != nil {
			response.Fatal(ctx, 500, err.Error())
			return
		}

		response.OkWithTabledata(ctx, 200, "products loaded", products, gin.H{})
	})

}
