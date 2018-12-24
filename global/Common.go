package global

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/Deansquirrel/go-tool"
)

var Config SysConfig
var Redis go_tool.MyRedis

func MyLog(s string) {
	fmt.Println(s)
	//err := go_tool.Log(s)
	//if err!= nil {
	//	fmt.Println(err," - ",s)
	//}
	//}
}

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
		fmt.Println(string(configJson))
	}
	Redis := go_tool.MyRedis{}
	Redis.Server = Config.RedisConfig.Server
	Redis.Auth = Config.RedisConfig.Password
	return nil
}
