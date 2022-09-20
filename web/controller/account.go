package controller

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"sgserver/constant"
	"sgserver/server/common"
	"sgserver/web/logic"
	"sgserver/web/model"
)

var DefaultAccountController = &AccountController{}

type AccountController struct {
}

func (a *AccountController) Register(ctx *gin.Context) {
	rq := &model.RegisterReq{}
	err := ctx.ShouldBind(rq)
	if err != nil {
		log.Println("参数格式不合法", err)
		ctx.JSON(http.StatusOK, common.Error(constant.InvalidParam, "参数不合法"))
		return
	}
	//一般web服务 错误格式 自定义
	err = logic.DefaultAccountLogic.Register(rq)
	if err != nil {
		log.Println("注册业务出错", err)
		ctx.JSON(http.StatusOK, common.Error(err.(*common.MyError).Code(), err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, common.Success(constant.OK, nil))
}
