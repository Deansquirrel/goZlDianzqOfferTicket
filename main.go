package main

import (
	"encoding/json"
	"fmt"
	"github.com/Deansquirrel/goZlDianzqOfferTicket/Object"
	"github.com/Deansquirrel/goZlDianzqOfferTicket/global"
	"github.com/kataras/iris"
	"os"
)

func main() {
	global.MyLog("程序启动")

	defer func() {
		if err := recover(); err != nil {
			errMsg := "Error:未处理的异常 - " + fmt.Sprint(err)
			global.MyLog(errMsg)

			os.Exit(-1)
		}
		global.MyLog("程序退出")
	}()

	err := global.GetConfig()
	if err != nil {
		fmt.Println(err)
		global.MyLog(err.Error())
		return
	}

	app := iris.New()
	app.Post("/", Handler)
	err = app.Run(iris.Addr(":8080"), iris.WithoutServerError(iris.ErrServerClosed), iris.WithOptimizations)
	if err != nil {
		global.MyLog(err.Error())
	}
}

func Handler(ctx iris.Context) {
	//===========================================================================================
	//接受请求数据
	request, err := Object.GetRequestByContext(ctx)
	if err != nil {
		global.MyLog(err.Error())
		_, err = ctx.Write([]byte(err.Error()))
		if err != nil {
			global.MyLog(err.Error())
		}
		return
	}
	data, err := json.Marshal(request)
	if err != nil {
		global.MyLog(err.Error())
		_, err = ctx.Write([]byte(err.Error()))
		if err != nil {
			global.MyLog(err.Error())
		}
		return
	}
	//===========================================================================================
	_, err = ctx.Write(data)
	if err != nil {
		global.MyLog(err.Error())
	}

	return
}
