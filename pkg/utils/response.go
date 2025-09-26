package utils

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

func OK(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, map[string]interface{}{"code": 200, "msg": "ok"})
}

func Resp(ctx *gin.Context, data any) {
	dataByte, err := json.Marshal(data)
	if err != nil {
		logrus.Warnf("json.Marshal failed, err: %v", err)
		InternalError(ctx, fmt.Errorf("resp data json.Marshal failed"))
		return
	}
	ctx.JSON(http.StatusOK, map[string]interface{}{"code": 200, "msg": "ok", "data": string(dataByte)})
}

func BadRequest(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusBadRequest, map[string]interface{}{"code": 400, "msg": err.Error()})
}

func InternalError(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusInternalServerError, map[string]interface{}{"code": 500, "msg": err.Error()})
}
