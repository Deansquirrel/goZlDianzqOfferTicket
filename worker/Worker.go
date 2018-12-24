package worker

import (
	"encoding/json"
	"github.com/Deansquirrel/goZlDianzqOfferTicket/Object"
	"github.com/Deansquirrel/goZlDianzqOfferTicket/global"
	"github.com/kataras/iris/core/errors"
	"strconv"
)

type Worker struct {
}

func (worker Worker) CreateLittleTktCreate(returnTktNo int, tktCreateRequest Object.TktCreateRequest) Object.TktCreateResponse {
	//检查请求
	err := checkRequest(tktCreateRequest)
	if err != nil {
		return getErrResponse("1", err.Error())
	}
	//生成券号码
	tktNum := len(tktCreateRequest.CrmCardInfo)
	noList, err := GetTktNoMulti(tktNum)
	if err != nil {
		return getErrResponse("1", err.Error())
	}
	if len(noList) < tktNum {
		return getErrResponse("1", "生成券号码出错")
	}

	//记录redis,防止重复提交
	jsonNoList, err := json.Marshal(noList)
	if err != nil {
		return getErrResponse("1", err.Error())
	}
	_, err = global.Redis.Set(strconv.Itoa(global.Config.RedisConfig.DbId1), tktCreateRequest.Appid+tktCreateRequest.YwInfo.OprYwSno, string(jsonNoList))
	if err != nil {
		return getErrResponse("1", err.Error())
	}

	//生成电子券系统流水号
	sno, err := GetSno("CT")
	if err != nil {
		return getErrResponse("1", err.Error())
	}
	global.MyLog("电子券系统流水号" + sno)

	return getErrResponse("1", err.Error())
}

//构造返回错误
func getErrResponse(errType string, desc string) Object.TktCreateResponse {
	resp := Object.ResponseResultInfo{
		ErrorModel: Object.ErrorModel{
			ErrType: errType,
			Desc:    desc,
		},
	}
	return Object.TktCreateResponse{
		ResponseResultInfo: resp,
	}
}

//检查请求
func checkRequest(req Object.TktCreateRequest) error {
	if req.TktInfo == nil || req.CrmCardInfo == nil {
		return errors.New("传入的记录集为空")
	}
	for index := range req.CrmCardInfo {
		if req.CrmCardInfo[index].CardType != 5 {
			return errors.New("会员券目前仅支持按账户id加密值操作")
		}
	}
	tktNum := len(req.CrmCardInfo)
	if tktNum > 100 {
		return errors.New("券发放（立即生效）禁止超过100张。")
	}

	val, err := global.Redis.Get(strconv.Itoa(global.Config.RedisConfig.DbId1), req.Appid+req.YwInfo.OprYwSno)
	if err != nil {
		return err
	}
	if val != "" {
		return errors.New("券发放请求重复提交")
	}
	return nil
}
