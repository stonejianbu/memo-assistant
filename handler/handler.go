package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/stonejianbu/memo-assistant/pkg/utils"
	"github.com/stonejianbu/memo-assistant/service"
)

type TrainTextReq struct {
	Input []string `json:"input"`
}

// TrainText 训练文本
func TrainText(ctx *gin.Context) {
	req := TrainTextReq{}
	if err := ctx.ShouldBind(&req); err != nil {
		utils.BadRequest(ctx, err)
		return
	}
	if err := service.Srv.TextManger.Train(ctx, req.Input); err != nil {
		utils.InternalError(ctx, err)
		return
	}
	utils.OK(ctx)
}

type RetrieveTextReq struct {
	Prompt string `json:"prompt"`
}

// RetrieveText 检索文本
func RetrieveText(ctx *gin.Context) {
	req := RetrieveTextReq{}
	if err := ctx.ShouldBind(&req); err != nil {
		utils.BadRequest(ctx, err)
		return
	}
	answer, err := service.Srv.TextManger.Query(ctx, req.Prompt)
	if err != nil {
		utils.InternalError(ctx, err)
		return
	}
	utils.Resp(ctx, map[string]interface{}{"answer": answer})
}
