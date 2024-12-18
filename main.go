package main

import (
	"fmt"
	"gudang-api-go-v1/core"
	"gudang-api-go-v1/middlewares"
	"gudang-api-go-v1/models"
	"gudang-api-go-v1/rest/routes"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// memuat file environment
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("unable to load env file: %s", err.Error())
	}

	// mengambil semua koneksi database
	dbs := core.GetAllDatabaseConnections()

	// melakukan auto migrasi database
	if err := dbs["main"].AutoMigrate(
		models.User{},
		models.Product{},
		models.Stock{},
	); err != nil {
		log.Fatalf("unable to migrate table to database: %s", err.Error())
	}

	// membuat user apabila database kosong
	core.CreateDefaultUser(dbs["main"])

	// membuat gin dengan setelan cors dan setting
	// kebutuhan device id pada setiap request
	router := gin.Default()
	router.Use(middlewares.Cors())

	// menambahkan endpotin yang dapat diakses public
	public := router.Group("")
	routes.AuthPublic(public, dbs)
	routes.Product(public, dbs)
	routes.Stock(public, dbs)

	// membuat router dengan kebutuhan token
	protected := router.Group("")
	protected.Use(middlewares.RequireDeviceId())
	protected.Use(middlewares.RequireToken())
	routes.Auth(protected, dbs)

	server := "localhost"
	host := "8021"

	// menjalankan server gin sesuai env
	router.Use(gin.Logger())
	router.Run(fmt.Sprintf("%s:%s", server, host))
}
