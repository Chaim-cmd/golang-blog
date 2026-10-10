// 统一响应
package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 业务错误码 0成功,1客户端业务错误,2 认证问题,5  服务端问题
const (
	CodeSuccess            = 0    // 成功
	CodeInvalidParams      = 1001 //参数校验失败
	CodeUserExists         = 1002 //用户名或邮箱已注册
	CodeNotFound           = 1003 //资源不存在
	CodeInvalidCredentials = 1004 // 用户或密码错误
	CodeUnauthorized       = 2001 //未登录 / token 无效或过期
	CodeServerError        = 5000 //服务端内部错误
)

// 全项目统一响应框架
// data 用omitempty : 失败响应 data 字段不直接出现, JSON 干净
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// 成功响应200
func Success(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, Response{Code: CodeSuccess, Message: "ok", Data: data})
}

// 创建资源201
func Created(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusCreated, Response{Code: CodeSuccess, Message: "ok", Data: data})

}

// Error 失败响应
func Error(ctx *gin.Context, httpStatus, code int, msg string) {
	ctx.JSON(httpStatus, Response{Code: code, Message: msg})
}
