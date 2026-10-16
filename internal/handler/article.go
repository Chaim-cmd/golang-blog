package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Chaim-cmd/golang-blog/internal/response"
	"github.com/Chaim-cmd/golang-blog/internal/service"
	"github.com/gin-gonic/gin"
)

type ArticleHandler struct {
	svc *service.ArticlesService
}

var Req struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

func NewArticleHandler(svc *service.ArticlesService) *ArticleHandler {
	return &ArticleHandler{
		svc: svc,
	}
}

func (h *ArticleHandler) Create(ctx *gin.Context) {

	if err := ctx.ShouldBindJSON(&Req); err != nil {
		response.Error(ctx, http.StatusBadRequest, http.StatusBadRequest, err.Error())
		return
	}
	userID := ctx.GetUint("user_id")
	art, err := h.svc.Create(userID, Req.Title, Req.Content)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, http.StatusBadRequest, "创建文章失败")
		return
	}
	response.Success(ctx, gin.H{
		"user_id": userID,
		"msg":     "创建成功",
		"data":    art,
	})

}

func (h *ArticleHandler) Update(ctx *gin.Context) {
	idStr := ctx.Param("id")
	//转成uint
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, http.StatusBadRequest, err.Error())
		return
	}
	if err := ctx.ShouldBindJSON(&Req); err != nil {
		response.Error(ctx, http.StatusBadRequest, http.StatusBadRequest, err.Error())
		return
	}
	userID := ctx.GetUint("user_id")
	if err := h.svc.Update(uint(id), userID, Req.Title, Req.Content); err != nil {
		if errors.Is(err, service.ErrArticleNoPerm) {
			response.Error(ctx, http.StatusForbidden, response.CodeForbidden, "没有权限修改文章")
			return
		}
		response.Error(ctx, http.StatusNotFound, response.CodeNotFound, "更新失败")
		return
	}
	response.Success(ctx, gin.H{
		"user_id": userID,
		"msg":     "文章更新成功",
	})
}

func (h *ArticleHandler) Delete(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, http.StatusBadRequest, err.Error())
		return
	}
	userID := ctx.GetUint("user_id")
	if err := h.svc.Delete(uint(id), userID); err != nil {
		if errors.Is(err, service.ErrArticleNotFound) {
			response.Error(ctx, http.StatusNotFound, response.CodeNotFound, "文章不存在")
			return
		} else if errors.Is(err, service.ErrArticleNoPerm) {
			response.Error(ctx, http.StatusForbidden, response.CodeInvalidParams, "无权删除文章")
			return
		} else {
			response.Error(ctx, http.StatusInternalServerError, response.CodeServerError, err.Error())
			return
		}

	}

}

func (h *ArticleHandler) List(ctx *gin.Context) {

}

func (h *ArticleHandler) Detail(ctx *gin.Context) {}
