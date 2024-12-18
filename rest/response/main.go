package response

import "github.com/gin-gonic/gin"

type Meta struct {
	Status  bool   `json:"status"`
	Code    uint   `json:"code"`
	Message string `json:"message"`
	Detail  string `json:"detail,omitempty"`
}

type Response struct {
	Meta Meta        `json:"meta"`
	Data interface{} `json:"data,omitempty"`
	Info interface{} `json:"info,omitempty"`
}

func Ok(ctx *gin.Context, code uint, message string) {
	ctx.JSON(200, Response{
		Meta: Meta{
			Status:  true,
			Code:    code,
			Message: message,
		},
	})
}

func OkWithData(ctx *gin.Context, code uint, message string, data interface{}) {
	ctx.JSON(200, Response{
		Meta: Meta{
			Status:  true,
			Code:    code,
			Message: message,
		},
		Data: data,
	})
}

func OkWithTabledata(ctx *gin.Context, code uint, message string, data interface{}, info interface{}) {
	ctx.JSON(200, Response{
		Meta: Meta{
			Status:  true,
			Code:    code,
			Message: message,
		},
		Data: data,
		Info: info,
	})
}

func Fatal(ctx *gin.Context, code uint, message string) {
	ctx.JSON(500, Response{
		Meta: Meta{
			Status:  false,
			Code:    code,
			Message: message,
		},
	})
}

func BadRequest(ctx *gin.Context, code uint, message string) {
	ctx.JSON(400, Response{
		Meta: Meta{
			Status:  false,
			Code:    code,
			Message: message,
		},
	})
}

func Unauthorized(ctx *gin.Context, code uint, message string) {
	ctx.JSON(401, Response{
		Meta: Meta{
			Status:  false,
			Code:    code,
			Message: message,
		},
	})
}

func Notfound(ctx *gin.Context, code uint, message string) {
	ctx.JSON(404, Response{
		Meta: Meta{
			Status:  false,
			Code:    code,
			Message: message,
		},
	})
}
