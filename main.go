package main

import (
	"encoding/json"
	"fmt"
	"github.com/Deansquirrel/go-tool"
	"github.com/Deansquirrel/goZlDianzqOfferTicket/common"
	"github.com/Deansquirrel/goZlDianzqOfferTicket/global"
	"github.com/Deansquirrel/goZlDianzqOfferTicket/repository"
	"github.com/Deansquirrel/goZlDianzqOfferTicket/yw"
	"github.com/kataras/iris"
	"github.com/kataras/iris/core/errors"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {

	err := refreshConfig()
	if err != nil {
		errMsg := "获取配置时遇到问题:" + err.Error()
		fmt.Println(errMsg)
		common.MyLog(errMsg)
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
	addr := ":" + strconv.Itoa(global.Config.TotalConfig.Port)
	err = app.Run(iris.Addr(addr), iris.WithoutServerError(iris.ErrServerClosed), iris.WithOptimizations)
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

func getResponse(ctx iris.Context) (response yw.ResponseCreateLittleTkt) {
	request, err := yw.GetRequestCreateLittleTktByContext(ctx)
	if err != nil {
		return getErrorResponse(request, ctx, err)
	}
	err = request.CheckRequest()
	if err != nil {
		return getErrorResponse(request, ctx, err)
	}
	return yw.GetResponseCreateLittleTkt(ctx, &request)
}

func getResponseData(response yw.ResponseCreateLittleTkt) []byte {
	data, err := json.Marshal(response)
	if err != nil {
		common.MyLog(err.Error())
		return []byte(err.Error())
	} else {
		return data
	}
}

func getErrorResponse(request yw.RequestCreateLittleTkt, ctx iris.Context, err error) (response yw.ResponseCreateLittleTkt) {
	common.MyLog(err.Error())
	response = yw.GetResponseCreateLittleTktError(&request, err, ctx.GetStatusCode())
	return response
}

func refreshConfig() error {
	//获取toml配置
	err := global.GetConfig()
	if err != nil {
		fmt.Println(err)
		common.MyLog(err.Error())
		return err
	}

	pZhR := repository.PeiZhRepository{}

	//获取Redis连接信息
	redisConfigStr, err := pZhR.GetXtWxAppIdJoinInfo(global.Config.TotalConfig.JPeiZh, "SERedis", 0)
	if err != nil {
		return err
	}
	redisConfig := strings.Split(redisConfigStr, "|")
	if len(redisConfig) != 2 {
		return errors.New("redis配置参数异常.expected 2 , got " + strconv.Itoa(len(redisConfig)))
	}
	global.Redis.Server = redisConfig[0]
	global.Redis.Auth = redisConfig[1]

	redisDbId1Str, err := pZhR.GetXtWxAppIdJoinInfo(global.Config.TotalConfig.JPeiZh, "RedisDbId1", 0)
	if err != nil {
		return err
	}
	global.RedisDbId1, err = strconv.Atoi(redisDbId1Str)
	if err != nil {
		return err
	}

	redisDbId2Str, err := pZhR.GetXtWxAppIdJoinInfo(global.Config.TotalConfig.JPeiZh, "RedisDbId2", 0)
	if err != nil {
		return err
	}
	global.RedisDbId2, err = strconv.Atoi(redisDbId2Str)
	if err != nil {
		return err
	}

	//获取RabbitMQ连接信息
	rabbitMQConfigStr, err := pZhR.GetXtWxAppIdJoinInfo(global.Config.TotalConfig.JPeiZh, "RabbitConnection", 0)
	if err != nil {
		return err
	}
	rabbitMQConfig := strings.Split(rabbitMQConfigStr, "|")
	if len(rabbitMQConfig) != 5 {
		return errors.New("rabbitMQ配置参数异常.expected 5 , got " + strconv.Itoa(len(rabbitMQConfig)))
	}
	rabbitMQPort, err := strconv.Atoi(rabbitMQConfig[1])
	if err != nil {
		return err
	}
	global.RabbitMQ = go_tool.NewRabbitMQ(rabbitMQConfig[3], rabbitMQConfig[4], rabbitMQConfig[0], rabbitMQPort, rabbitMQConfig[2], time.Second*60, time.Millisecond*500, 3, time.Second*5)

	//获取SnoServer信息
	global.SnoServer, err = pZhR.GetXtWxAppIdJoinInfo(global.Config.TotalConfig.JPeiZh, "SnoServer", 0)
	if err != nil {
		return err
	}
	snoWorkIdStr, err := pZhR.GetXtWxAppIdJoinInfo(global.Config.TotalConfig.JPeiZh, "WorkerId", 0)
	if err != nil {
		return err
	}
	global.SnoWorkerId, err = strconv.Atoi(snoWorkIdStr)
	if err != nil {
		return err
	}
	global.SnoWorkerId, err = strconv.Atoi(snoWorkIdStr)
	if err != nil {
		return err
	}

	err = rabbitMqInit()
	if err != nil {
		return err
	}

	return nil
}

func rabbitMqInit() error {
	conn, err := global.RabbitMQ.GetConn()
	if err != nil {
		return err
	}
	defer func() {
		_ = conn.Close()
	}()

	err = global.RabbitMQ.QueueDeclareSimple(conn, "TktCreateYwdetail")
	if err != nil {
		return err
	}

	err = global.RabbitMQ.QueueBind(conn, "TktCreateYwdetail", "", "amq.fanout", true)
	if err != nil {
		return err
	}

	err = global.RabbitMQ.AddProducer("")
	if err != nil {
		return err
	}

	return nil
}
