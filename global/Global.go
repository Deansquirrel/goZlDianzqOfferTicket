package global

import (
	"encoding/json"
	"github.com/BurntSushi/toml"
	"github.com/Deansquirrel/go-tool"
	"github.com/Deansquirrel/goZlDianzqOfferTicket/common"
	"github.com/kataras/iris/core/errors"
)

var Config SysConfig
var Redis go_tool.MyRedis
var RabbitMQ *go_tool.MyRabbitMQ
var RedisDbId1 int
var RedisDbId2 int
var SnoServer string
var SnoWorkerId int

func GetConfig() (err error) {
	_, err = toml.DecodeFile("config.toml", &Config)
	if err != nil {
		err = errors.New("配置文件获取异常:" + err.Error())
		return
	}
	configJson, err := json.Marshal(Config)
	if err != nil {
		err = errors.New("配置文件转JSON时遇到异常:" + err.Error())
	} else {
		common.MyLog(string(configJson))
	}
	common.IsDebug = Config.TotalConfig.IsDebug
	return
}
