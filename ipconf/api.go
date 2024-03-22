package ipconf

import (
	"blue/ipconf/domain"
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

type Response struct {
	Message string      `json:"message"`
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
}

func GetIpInfoList(c context.Context, ctx *app.RequestContext) {
	//只是单次请求有问题,但不应该让整个进程都退出
	defer func() {
		if err := recover(); err != nil {
			ctx.JSON(consts.StatusBadRequest, utils.H{"err": err})
		}
	}()
	//Step0: 封装context请求对象
	ipConfCtx := domain.BuildIpConfContext(&c, ctx)
	//Step1: 根据ip信息进行调度
	endPoints := domain.Dispatch(ipConfCtx)
	//Step2: 获取前五个并返回
	ipConfCtx.AppCtx.JSON(consts.StatusOK, packResp(top5EndPort(endPoints)))
}
