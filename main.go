package main

import (
	"fmt"
	"github.com/Deansquirrel/goZlDianzqOfferTicket/global"
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

}
