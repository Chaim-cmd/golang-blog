package handler

import (
	"errors"
	"net/http"

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
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "参数不合法:" + err.Error()})
		return
	}
	u, err := h.svc.Register(req.Username, req.Email, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrUserExists) {
			ctx.JSON(http.StatusConflict, gin.H{
				"error": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "注册失败",
		})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{
		"id":       u.ID,
		"username": u.Username,
		"email":    u.Email,
	})

}
func (h *UserHandler) Login(ctx *gin.Context) {
	var req LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "参数不合法:" + err.Error()})
		return
	}
	u, err := h.svc.Login(req.Username, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrBadCredential) {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "登录失败",
		})
		return
	}
	token, err := h.jwtSvc.GenerateToken(u.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "签发 token失败",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
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
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "用户不存在",
		})
	}
	ctx.JSON(http.StatusOK, gin.H{
		"id":       u.ID,
		"username": u.Username,
		"email":    u.Email,
	})
}
