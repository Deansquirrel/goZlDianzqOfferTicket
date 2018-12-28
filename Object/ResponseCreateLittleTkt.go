package Object

import (
	"encoding/json"
	"fmt"
	"github.com/Deansquirrel/goZl"
	"github.com/Deansquirrel/goZlDianzqOfferTicket/common"
	"github.com/Deansquirrel/goZlDianzqOfferTicket/global"
	"github.com/Deansquirrel/goZlDianzqOfferTicket/repository"
	"github.com/kataras/iris"
	"github.com/kataras/iris/core/errors"
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
	response.refresh()

	//生成券号码
	tktNos, err := GetTktNoMulti(len(request.Body.CrmCardInfo))
	if err != nil {
		response = GetResponseCreateLittleTktError(request, err, ctx.GetStatusCode())
		return
	}
	common.MyLog("生成的券号码")
	for index := range tktNos {
		common.MyLog(tktNos[index])
	}

	//记录redis,防止重复提交
	common.MyLog("准备记录Reids,防止重复提交")
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
			common.MyLog(err.Error())
		}
	}()

	common.MyLog(string(jsonNoList))

	//生成电子券系统流水号
	common.MyLog("准备生成电子券系统流水号")
	sno, err := GetSno("CT")
	if err != nil {
		response = GetResponseCreateLittleTktError(request, err, ctx.GetStatusCode())
		return
	}
	common.MyLog(sno)

	hxR := repository.HeXRepository{}
	hxDbList, err := repository.GetHxDbConn(global.Config.TotalConfig.AppId)
	if err != nil {
		response = GetResponseCreateLittleTktError(request, err, ctx.GetStatusCode())
		return
	}
	if len(hxDbList) < 1 {
		err = errors.New("传入的APPID无效（APPID错误或配置库缺少配置）")
		response = GetResponseCreateLittleTktError(request, err, ctx.GetStatusCode())
		return
	}

	t := conftools.ConfTools{}
	tktInfos := make([]TktInfo, 0)
	tktReturnInfos := make([]TktReturnInfo, 0)
	tktModels := make([]TktModel, 0)
	j := 0
	for _, val := range request.Body.CrmCardInfo {
		accIdInput, err := t.DecryptFromBase64Format(val.CardNo, "accid")
		accIdInputLong, err := strconv.Atoi(accIdInput)
		if err != nil {
			response = GetResponseCreateLittleTktError(request, err, ctx.GetStatusCode())
			return
		}
		var tktInfo TktInfo
		if request.Body.TktInfo == nil || len(request.Body.TktInfo) < 1 {
			break
		} else {
			tktInfo = request.Body.TktInfo[0]
		}
		tktItem := TktInfo{
			AppId:    request.AppId,
			AccId:    accIdInputLong,
			TktKind:  tktInfo.TktKind,
			EffDate:  tktInfo.EffDate,
			Deadline: tktInfo.Deadline,
			CrYwLsh:  request.Body.YwInfo.OprYwSno,
			CrBr:     request.Body.YwInfo.OprBrid,
			CashMy:   tktInfo.CashMy,
			AddMy:    tktInfo.AddMy,
			TktName:  tktInfo.TktName,
			PCno:     tktInfo.PCno,
			TktNo:    tktNos[j],
		}
		tktInfos = append(tktInfos, tktItem)

		tktRetItem := TktReturnInfo{
			Sn:     val.Sn,
			TktNo:  tktNos[j],
			TktSno: sno,
		}
		tktReturnInfos = append(tktReturnInfos, tktRetItem)

		tktModel := TktModel{
			AccId:    accIdInputLong,
			AddMy:    tktInfo.AddMy,
			AppId:    request.AppId,
			CashMy:   tktInfo.CashMy,
			DeadLine: tktInfo.Deadline,
			PcNo:     tktInfo.PCno,
			EffTime:  tktInfo.EffDate,
			TktKind:  tktInfo.TktKind,
			TktName:  tktInfo.TktName,
			TktNo:    tktNos[j],
		}
		tktModels = append(tktModels, tktModel)

		j++
	}

	crTktInfo := TktCreateInfo{
		TktInfo:       make([]TktInfo, 0),
		TktYwInfo:     YwInfo{},
		TktReturnInfo: make([]TktReturnInfo, 0),
		CzLx:          2,
		CzLxSm:        "",
	}

	for _, val := range tktInfos {
		crTktInfo.TktInfo = append(crTktInfo.TktInfo, val)
	}
	crTktInfo.TktYwInfo = request.Body.YwInfo
	for _, val := range tktReturnInfos {
		crTktInfo.TktReturnInfo = append(crTktInfo.TktReturnInfo, val)
	}

	tkTs := TktModels{
		TktModels: make([]TktModel, 0),
	}
	for _, val := range tktModels {
		tkTs.TktModels = append(tkTs.TktModels, val)
	}

	for _, db := range hxDbList {
		verInfo, err := hxR.GetVerInfo(db)
		if err != nil {
			response = GetResponseCreateLittleTktError(request, err, ctx.GetStatusCode())
			return
		}
		fmt.Println(verInfo)
	}

	if err != nil {
		common.MyLog(err.Error())
		response = GetResponseCreateLittleTktError(request, err, ctx.GetStatusCode())
		return
	}

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
