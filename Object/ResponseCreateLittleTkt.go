package Object

import (
	"encoding/json"
	"github.com/Deansquirrel/goZlDianzqOfferTicket/global"
	"github.com/kataras/iris"
	"strconv"
	"strings"
)

type ResponseCreateLittleTkt struct {
	TmpCol    int             `json:"tmpcol"`
	TktReturn []TktReturnInfo `json:"tktReturn"`

	ErrorModel  ErrorModel `json:"errormodel"`
	Description string     `json:"description"`

	HttpCode    int  `josn:"httpcode"`
	DBCommitted bool `json:"dbcommitted"`

	IsSuccess bool   `json:"issuccess"`
	FindAccid bool   `json:"findaccid"`
	HasCommit bool   `json:"hascommit"`
	Guid      string `json:"requestGuid"`
}

func GetResponseCreateLittleTkt(ctx iris.Context, request *RequestCreateLittleTkt) (response ResponseCreateLittleTkt) {
	response.BaseFunc(request, &response)

	//生成券号码
	tktNos, err := GetTktNoMulti(len(request.Body.CrmCardInfo))
	if err != nil {
		response = GetResponseCreateLittleTktError(request, err, ctx.GetStatusCode())
		return
	}
	global.MyLog("生成的券号码")
	for index := range tktNos {
		global.MyLog(tktNos[index])
	}

	//记录redis,防止重复提交
	global.MyLog("准备记录Reids,防止重复提交")
	jsonNoList, err := json.Marshal(tktNos)
	if err != nil {
		response = GetResponseCreateLittleTktError(request, err, ctx.GetStatusCode())
		return
	}
	_, err = global.Redis.Set(strconv.Itoa(global.Config.RedisConfig.DbId1), request.AppId+request.Body.YwInfo.OprYwSno, string(jsonNoList))
	if err != nil {
		response = GetResponseCreateLittleTktError(request, err, ctx.GetStatusCode())
		return
	}
	defer func() {
		err = global.Redis.Del(strconv.Itoa(global.Config.RedisConfig.DbId1), request.AppId+request.Body.YwInfo.OprYwSno)
		if err != nil {
			global.MyLog(err.Error())
		}
	}()

	global.MyLog(string(jsonNoList))

	//生成电子券系统流水号
	global.MyLog("准备生成电子券系统流水号")
	sno, err := GetSno("CT")
	if err != nil {
		response = GetResponseCreateLittleTktError(request, err, ctx.GetStatusCode())
		return
	}
	global.MyLog(sno)

	//通过配置库获取Hx库连接信息
	pz := PeiZh{
		AppId:    global.Config.PeiZhDbConfig.AppId,
		Server:   global.Config.PeiZhDbConfig.Server,
		Port:     global.Config.PeiZhDbConfig.Port,
		DbName:   global.Config.PeiZhDbConfig.DbName,
		User:     global.Config.PeiZhDbConfig.User,
		PassWord: global.Config.PeiZhDbConfig.PassWord,
	}
	d, err := pz.GetHxDbInfo()
	if err != nil {
		global.MyLog(err.Error())
		response = GetResponseCreateLittleTktError(request, err, ctx.GetStatusCode())
		return
	}
	global.MyLog(strconv.Itoa(len(d)))

	return
}

func GetResponseCreateLittleTktError(request *RequestCreateLittleTkt, err error, httpCode int) (response ResponseCreateLittleTkt) {
	response.BaseFunc(request, &response)
	response.HttpCode = httpCode
	response.ErrorModel = ErrorModel{
		ErrType: "1",
		Desc:    err.Error(),
	}
	response.refresh()
	return
}

func (response *ResponseCreateLittleTkt) BaseFunc(req *RequestCreateLittleTkt, resp *ResponseCreateLittleTkt) {
	resp.Guid = req.Guid
}

func (response *ResponseCreateLittleTkt) refresh() {
	//if (response.HttpCode == 503 || response.HttpCode == 403) && response.Description != "" {
	//	if strings.Contains(response.Description,"账号无法找到") {
	//		panic(response.Description)
	//	}
	//}

	if response.ErrorModel.ErrType == "" && response.ErrorModel.Desc == "" {
		response.DBCommitted = true
	} else {
		response.DBCommitted = false
	}

	if response.ErrorModel.ErrType != "" || response.ErrorModel.Desc != "" {
		response.IsSuccess = false
	} else if response.HttpCode != 200 {
		response.IsSuccess = false
	} else {
		response.IsSuccess = true
	}

	if response.HttpCode == 503 && strings.Contains(response.Description, "账号无法找到") {
		response.FindAccid = false
	} else {
		response.FindAccid = true
	}
}
