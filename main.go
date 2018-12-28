package main

import (
	"encoding/json"
	"fmt"
	"github.com/Deansquirrel/goZlDianzqOfferTicket/Object"
	"github.com/Deansquirrel/goZlDianzqOfferTicket/common"
	"github.com/Deansquirrel/goZlDianzqOfferTicket/global"
	"github.com/kataras/iris"
	"os"
)

func main() {

	err := global.GetConfig()
	if err != nil {
		fmt.Println(err)
		common.MyLog(err.Error())
		return
	}

	common.MyLog("程序启动")

	defer func() {
		if err := recover(); err != nil {
			errMsg := "Error:未处理的异常 - " + fmt.Sprint(err)
			common.MyLog(errMsg)

			os.Exit(-1)
		}
		common.MyLog("程序退出")
	}()

	app := iris.New()
	app.Post("/", Handler)
	err = app.Run(iris.Addr(":8080"), iris.WithoutServerError(iris.ErrServerClosed), iris.WithOptimizations)
	if err != nil {
		common.MyLog(err.Error())
	}
}

func Handler(ctx iris.Context) {
	response := getResponse(ctx)
	_, err := ctx.Write(getResponseData(response))
	if err != nil {
		common.MyLog(err.Error())
	}
	return

	////===========================================================================================
	////接受请求数据
	//request, err := Object.GetRequestCreateLittleTktByContext(ctx)
	//if err != nil {
	//	common.MyLog(err.Error())
	//	_, err = ctx.Write([]byte(err.Error()))
	//	if err != nil {
	//		common.MyLog(err.Error())
	//	}
	//	return
	//}
	//data, err := json.Marshal(request)
	//if err != nil {
	//	common.MyLog(err.Error())
	//	_, err = ctx.Write([]byte(err.Error()))
	//	if err != nil {
	//		common.MyLog(err.Error())
	//	}
	//	return
	//}
	////===========================================================================================
	//_, err = ctx.Write(data)
	//_, err = ctx.Write([]byte(strconv.Itoa(ctx.GetStatusCode())))
	//if err != nil {
	//	common.MyLog(err.Error())
	//}
	//return
}

func getResponse(ctx iris.Context) (response Object.ResponseCreateLittleTkt) {
	request, err := Object.GetRequestCreateLittleTktByContext(ctx)
	if err != nil {
		return getErrorResponse(request, ctx, err)
	}
	err = request.CheckRequest()
	if err != nil {
		return getErrorResponse(request, ctx, err)
	}
	return Object.GetResponseCreateLittleTkt(ctx, &request)
}

func getResponseData(response Object.ResponseCreateLittleTkt) []byte {
	data, err := json.Marshal(response)
	if err != nil {
		common.MyLog(err.Error())
		return []byte(err.Error())
	} else {
		return data
	}
}

func getErrorResponse(request Object.RequestCreateLittleTkt, ctx iris.Context, err error) (response Object.ResponseCreateLittleTkt) {
	common.MyLog(err.Error())
	response = Object.GetResponseCreateLittleTktError(&request, err, ctx.GetStatusCode())
	return response
}
