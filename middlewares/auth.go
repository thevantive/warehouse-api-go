package middlewares

import (
	"gudang-api-go-v1/rest"
	"gudang-api-go-v1/rest/response"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

var jwtSecret = "main"

// untuk kebutuhan middleware autentikasi
// membutuhkan device id setiap pengguna
func RequireDeviceId() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		deviceId := ctx.GetHeader("X-Device-Id")

		// memastikan device id tidak kosong
		if deviceId == "" {
			response.BadRequest(ctx, 400, "invalid headers missing device id")
			ctx.Abort()
			return
		}

		// mengambil fingerprint user
		ctx.Set("fingerprint", rest.GetUserDeviceFingeprint(ctx))

		// apabila tersedia maka dilanjutkan
		ctx.Next()
	}
}

// membutuhkan token untuk setiap pengguna
func RequireToken() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")

		// memastikan token tidak kosong
		if len(authHeader) <= 7 {
			response.BadRequest(ctx, 401, "unauthorized: missing token")
			ctx.Abort()
			return
		}

		// memeriksa format token
		if !strings.HasPrefix(authHeader, "Bearer ") {
			response.BadRequest(ctx, 401, "unauthorized: invalid token format")
			ctx.Abort()
			return
		}

		// melakukan trimming pada token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// melakukan parsing pada token
		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})

		// token tidak valid
		if err != nil || !token.Valid {
			response.BadRequest(ctx, 401, "unauthorized: invalid token")
			ctx.Abort()
			return
		}

		// menambahkan user_id kepada context request
		userId := uint(claims["user_id"].(float64))
		ctx.Set("user_id", userId)
		ctx.Next()
	}
}
