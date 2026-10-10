package handler

import (
	"errors"
	"net/http"

	"github.com/Chaim-cmd/golang-blog/internal/response"
	"github.com/Chaim-cmd/golang-blog/internal/service"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	svc    *service.UserService
	jwtSvc *service.JWTService
}

func NewUserHandler(svc *service.UserService, jwtSvc *service.JWTService) *UserHandler {
	return &UserHandler{
		svc:    svc,
		jwtSvc: jwtSvc,
	}
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *UserHandler) Register(ctx *gin.Context) {
	var req RegisterRequest
	//ShouldBindJSON 解析请求体 JSON 并执行binding 校验
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, response.CodeInvalidParams, "参数错误: ...")
		return
	}
	u, err := h.svc.Register(req.Username, req.Email, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrUserExists) {
			response.Error(ctx, http.StatusConflict, response.CodeUserExists, "用户名或邮箱已注册")
			return
		}
		response.Error(ctx, http.StatusInternalServerError, response.CodeServerError, "服务器内部错误")
		return
	}
	response.Success(ctx, gin.H{
		"id":       u.ID,
		"username": u.Username,
		"email":    u.Email,
	})

}
func (h *UserHandler) Login(ctx *gin.Context) {
	var req LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, response.CodeInvalidParams, "参数错误: ...")
		return
	}
	u, err := h.svc.Login(req.Username, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrBadCredential) {
			response.Error(ctx, http.StatusBadRequest, response.CodeInvalidCredentials, "用户名或密码错误")
			return
		}
		response.Error(ctx, http.StatusInternalServerError, response.CodeServerError, "服务器错误")
		return
	}
	token, err := h.jwtSvc.GenerateToken(u.ID)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, response.CodeServerError, "token 提取失败")
		return
	}
	response.Success(ctx, gin.H{
		"message":  "登录成功",
		"token":    token,
		"user_id":  u.ID,
		"username": u.Username,
	})

}

func (h *UserHandler) Me(ctx *gin.Context) {
	//user_id 是 authMiddleware 验签通过后塞进上下文的
	userID := ctx.GetUint("user_id")

	u, err := h.svc.GetByID(userID)
	if err != nil {
		response.Error(ctx, http.StatusNotFound, response.CodeNotFound, "用户不存在")
		return
	}

	response.Success(ctx, gin.H{
		"id":       u.ID,
		"username": u.Username,
		"email":    u.Email,
	})
}
