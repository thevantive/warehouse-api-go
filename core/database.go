package core

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func GetDatabaseConnection(name string) (*gorm.DB, error) {
	name = strings.ToUpper(fmt.Sprintf("db_%s", name))

	// mengambil konfigurasi database dari .env
	username := os.Getenv(fmt.Sprintf("%s_USERNAME", name))
	password := os.Getenv(fmt.Sprintf("%s_PASSWORD", name))
	host := os.Getenv(fmt.Sprintf("%s_HOST", name))
	dbname := os.Getenv(fmt.Sprintf("%s_NAME", name))

	// memastikan apabila konfigurasi tersedia
	if username == "" {
		return nil, errors.New("missing .env data")
	}

	// menyiapkan dsn untuk database
	var dsn string
	dsn = fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, host, dbname)
	dsn = fmt.Sprintf("%s?charset=utf8mb4&parseTime=True&loc=Local", dsn)
	fmt.Println(dsn)

	// membuat koneksi database
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, nil
}

// mengambil semua database yang tersedia
func GetAllDatabaseConnections() map[string]*gorm.DB {
	var dbs = make(map[string]*gorm.DB)

	// mengambil database utama transnusa
	var db, err = GetDatabaseConnection("main")
	if err != nil {
		log.Fatalf("unable to connect to database: %s", err.Error())
	}
	dbs["main"] = db

	return dbs
}
