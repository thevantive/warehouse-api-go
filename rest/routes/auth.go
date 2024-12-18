package routes

import (
	"errors"
	"gudang-api-go-v1/core"
	"gudang-api-go-v1/models"
	"gudang-api-go-v1/rest/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Auth(router *gin.RouterGroup, dbs map[string]*gorm.DB) {

	// endpoint untuk kebutuhan validasi token
	router.GET("auth/validate", func(ctx *gin.Context) {
		var db = dbs["main"].Session(&gorm.Session{})
		userId := ctx.GetUint("user_id")

		// memastikan user_id terdaftar pada database
		var user models.User
		if err := db.First(&user, userId).Error; err != nil {
			response.Fatal(ctx, 500, err.Error())
			return
		}
		if user.Id == 0 {
			response.Unauthorized(ctx, 401, "user not found")
			return
		}

		// melakukan generate pada token
		token, err := core.GenerateToken(uint(userId), user.Fullname)
		if err != nil {
			response.Fatal(ctx, 500, err.Error())
			return
		}

		response.OkWithData(ctx, 200, "token valid", gin.H{
			"token": token,
		})
	})
}

func AuthPublic(router *gin.RouterGroup, dbs map[string]*gorm.DB) {
	router.POST("auth/login", func(ctx *gin.Context) {
		var db = dbs["main"].Session(&gorm.Session{})

		type RequestBody struct {
			Username string `json:"username" binding:"required"`
			Password string `json:"password" binding:"required"`
		}
		var requestBody RequestBody
		if err := ctx.ShouldBindBodyWithJSON(&requestBody); err != nil {
			response.BadRequest(ctx, 400, "please read documentation")
			return
		}

		// mengambil data pengguna
		var user models.User
		if err := db.First(&user, "username = ?", requestBody.Username).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				response.BadRequest(ctx, 401, "invalid username or password")
				return
			}
			response.Fatal(ctx, 500, err.Error())
			return
		}

		// melakukan validasi pada hashing password
		if !core.ValidatePassword(user.Password, requestBody.Password) {
			response.BadRequest(ctx, 402, "invalid username or password")
			return
		}

		// membuat token baru pada database
		token, err := core.GenerateToken(user.Id, user.Fullname)
		if err != nil {
			response.Fatal(ctx, 500, err.Error())
			return
		}

		response.OkWithData(ctx, 200, "login successfully", gin.H{
			"token":    token,
			"fullname": user.Fullname,
		})
	})

	router.POST("auth/logout", func(ctx *gin.Context) {
		response.Ok(ctx, 200, "logout succesfully")
	})
}
