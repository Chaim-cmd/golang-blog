package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Chaim-cmd/golang-blog/internal/response"
	"github.com/Chaim-cmd/golang-blog/internal/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ArticleHandler struct {
	svc        *service.Articles
	commentSvc *service.Comments
}

func NewArticleHandler(svc *service.Articles, commentSvc *service.Comments) *ArticleHandler {
	return &ArticleHandler{
		svc:        svc,
		commentSvc: commentSvc,
	}
}

func (h *ArticleHandler) Create(ctx *gin.Context) {
	var Req struct {
		Title   string `json:"title" binding:"required"`
		Content string `json:"content" binding:"required"`
	}

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

	var Req struct {
		Title   string `json:"title" binding:"required"`
		Content string `json:"content" binding:"required"`
	}
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
			response.Error(ctx, http.StatusForbidden, response.CodeForbidden, "无权删除文章")
			return
		} else {
			response.Error(ctx, http.StatusInternalServerError, response.CodeServerError, err.Error())
			return
		}

	}

}

func (h *ArticleHandler) List(ctx *gin.Context) {
	pageStr := ctx.Query("page")
	pageSizeStr := ctx.Query("page_size")
	page, err := strconv.Atoi(pageStr)

	if err != nil {
		response.Error(ctx, http.StatusBadRequest, http.StatusBadRequest, err.Error())
		return
	}
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, http.StatusBadRequest, err.Error())
		return
	}
	articles, total, err := h.svc.List(page, pageSize)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, response.CodeServerError, "获取文章失败")
		return
	}
	response.Success(ctx, gin.H{
		"data":  articles,
		"total": total,
	})

}

func (h *ArticleHandler) Detail(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, response.CodeInvalidParams, err.Error())
		return
	}

	article, err := h.svc.GetByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Error(ctx, http.StatusNotFound, response.CodeNotFound, "文章不存在")
			return
		}
		response.Error(ctx, http.StatusInternalServerError, response.CodeServerError, "查询文章失败")
		return
	}

	response.Success(ctx, article)
}

func (h *ArticleHandler) CreateComments(ctx *gin.Context) {
	artIDStr := ctx.Param("id")
	artID, err := strconv.ParseUint(artIDStr, 10, 64)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, response.CodeInvalidParams, err.Error())
		return
	}
	userID, ok := ctx.Get("user_id")
	if !ok {
		response.Error(ctx, http.StatusBadRequest, response.CodeInvalidParams, err.Error())
		return
	}
	UID := userID.(uint)
	var Req struct {
		Content string `json:"content" gorm:"binding"`
	}
	if err := ctx.ShouldBindJSON(&Req); err != nil {
		response.Error(ctx, http.StatusBadRequest, response.CodeInvalidParams, err.Error())
		return

	}
	comment, err := h.commentSvc.Create(uint(artID), UID, Req.Content)
	if err != nil {
		response.Error(ctx, http.StatusNotFound, response.CodeNotFound, "文章不存在")
	}
	response.Success(ctx, gin.H{
		"data": comment,
	})

}

func (h *ArticleHandler) ListByArticle(ctx *gin.Context) {
	pageStr := ctx.Query("page")
	pageSizeStr := ctx.Query("pageSize")
	var title int64
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, response.CodeInvalidParams, err.Error())
	}
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, response.CodeInvalidParams, err.Error())
	}
	Comment, title, err := h.commentSvc.ListByArticle(page, pageSize)
	if err != nil {
		response.Error(ctx, http.StatusNotFound, response.CodeNotFound, "获取评论失败")
		return
	}
	if page < 1 || pageSize < 1 {
		pageSize = 10
		page = 1
	}
	if pageSize > 50 {
		pageSize = 50
	}
	response.Success(ctx, gin.H{
		"data":  Comment,
		"title": title,
	})

}
