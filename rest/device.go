package rest

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/gin-gonic/gin"
)

// untuk mengambil fingerprint device
func GetUserDeviceFingeprint(ctx *gin.Context) string {
	deviceId := ctx.GetHeader("X-Device-Id")
	userAgent := ctx.GetHeader("User-Agent")
	ip := ctx.ClientIP()

	// menggabungkan string untuk diubah menjadi fingerprint
	fingerprint := fmt.Sprintf("%s:%s:%s", deviceId, userAgent, ip)

	// melakukan hashing pada string
	hash := sha256.New()
	hash.Write([]byte(fingerprint))

	return hex.EncodeToString(hash.Sum(nil))
}
