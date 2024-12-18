package core

import (
	"fmt"
	"log"
	"gudang-api-go-v1/models"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var jwtSecret = "main"

// membuat user default
func CreateDefaultUser(db *gorm.DB) {
	password, err := HashPassword("admin")
	if err != nil {
		log.Fatalf("unable to hash password: %s", err.Error())
		return
	}

	// menyiapkan data user default
	var user = models.User{
		Fullname: "Superadmin",
		Username: "admin",
		Password: password,
	}

	// membuat transaksi database
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// memeriksa pada database
	var count int64
	if err := db.Model(&models.User{}).Where("status = ?", "active").Count(&count).Error; err != nil {
		log.Fatalf("unable to check user count: %s", err.Error())
		tx.Rollback()
		return
	}
	if count > 0 {
		tx.Rollback()
		return
	}

	// menyimpan user ke database
	if err := tx.Save(&user).Error; err != nil {
		log.Fatalf("unable to save user to database: %s", err.Error())
		tx.Rollback()
		return
	}

	// melakukan commit transaksi
	fmt.Println("user default created")
	tx.Commit()

}

// untuk membuat jwt token
func GenerateToken(userId uint, fullname string) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  userId,
		"fullname": fullname,
		"exp":      time.Now().Add(time.Hour * 1).Unix(),
	}

	// membuat claim token baru
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// untuk melakukan hashing pada password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// untuk melakukan validasi pada hashing password
func ValidatePassword(hashedPassword, plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	return err == nil
}
