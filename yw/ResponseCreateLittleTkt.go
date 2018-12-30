package yw

import (
	"database/sql"
	"encoding/json"
	"github.com/Deansquirrel/goZl"
	"github.com/Deansquirrel/goZlDianzqOfferTicket/Object"
	"github.com/Deansquirrel/goZlDianzqOfferTicket/common"
	"github.com/Deansquirrel/goZlDianzqOfferTicket/global"
	"github.com/Deansquirrel/goZlDianzqOfferTicket/repository"
	"github.com/kataras/iris"
	"github.com/kataras/iris/core/errors"
	"strconv"
	"strings"
)

type ResponseCreateLittleTkt struct {
	TmpCol    int                    `json:"tmpcol"`
	TktReturn []Object.TktReturnInfo `json:"tktReturn"`

	ErrorModel  Object.ErrorModel `json:"errormodel"`
	Description string            `json:"description"`

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
	tktNos, err := Object.GetTktNoMulti(len(request.Body.CrmCardInfo))
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

	_, err = global.Redis.Set(strconv.Itoa(global.RedisDbId1), request.AppId+request.Body.YwInfo.OprYwSno, string(jsonNoList))
	if err != nil {
		response = GetResponseCreateLittleTktError(request, err, ctx.GetStatusCode())
		return
	}
	defer func() {
		err = global.Redis.Del(strconv.Itoa(global.RedisDbId1), request.AppId+request.Body.YwInfo.OprYwSno)
		if err != nil {
			common.MyLog(err.Error())
		}
	}()

	common.MyLog(string(jsonNoList))

	//生成电子券系统流水号
	common.MyLog("准备生成电子券系统流水号")
	sno, err := Object.GetSno("CT")
	if err != nil {
		response = GetResponseCreateLittleTktError(request, err, ctx.GetStatusCode())
		return
	}
	common.MyLog(sno)

	pzR := repository.PeiZhRepository{}
	dbConnList, err := pzR.GetXtMappingDbConnInfo(request.AppId, "DB_TicketHx", "TicketHx")
	if err != nil {
		response = GetResponseCreateLittleTktError(request, err, ctx.GetStatusCode())
		return
	}

	if len(dbConnList) < 1 {
		err = errors.New("传入的APPID无效（APPID错误或配置库缺少配置）")
		response = GetResponseCreateLittleTktError(request, err, ctx.GetStatusCode())
		return
	}

	t := conftools.ConfTools{}
	tktInfos := make([]Object.TktInfo, 0)
	tktReturnInfos := make([]Object.TktReturnInfo, 0)
	tktModels := make([]Object.TktModel, 0)
	j := 0

	for _, val := range request.Body.CrmCardInfo {
		common.MyLog("CardNo : " + val.CardNo)
		accIdInput, err := t.DecryptFromBase64Format(val.CardNo, "accid")
		if err != nil {
			response = GetResponseCreateLittleTktError(request, err, ctx.GetStatusCode())
			return
		}
		common.MyLog(accIdInput)
		accIdInputLong, err := strconv.Atoi(accIdInput)
		if err != nil {
			response = GetResponseCreateLittleTktError(request, err, ctx.GetStatusCode())
			return
		}
		var tktInfo Object.TktInfo
		if request.Body.TktInfo == nil || len(request.Body.TktInfo) < 1 {
			break
		} else {
			tktInfo = request.Body.TktInfo[0]
		}
		tktItem := Object.TktInfo{
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

		tktRetItem := Object.TktReturnInfo{
			Sn:     val.Sn,
			TktNo:  tktNos[j],
			TktSno: sno,
		}
		tktReturnInfos = append(tktReturnInfos, tktRetItem)

		tktModel := Object.TktModel{
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

	crTktInfo := Object.TktCreateInfo{
		TktInfo:       make([]Object.TktInfo, 0),
		TktYwInfo:     Object.YwInfo{},
		TktReturnInfo: make([]Object.TktReturnInfo, 0),
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

	tkTs := Object.TktModels{
		TktModels: make([]Object.TktModel, 0),
	}
	for _, val := range tktModels {
		tkTs.TktModels = append(tkTs.TktModels, val)
	}

	//for _, db := range hxDbList {
	//	verInfo, err := hxR.GetVerInfo(db)
	//	if err != nil {
	//		response = GetResponseCreateLittleTktError(request, err, ctx.GetStatusCode())
	//		return
	//	}
	//	fmt.Println(verInfo)
	//}

	if err != nil {
		common.MyLog(err.Error())
		response = GetResponseCreateLittleTktError(request, err, ctx.GetStatusCode())
		return
	}
	response = createLittleTktCreate(request.ReturnTktNo, crTktInfo, tkTs, ctx, request)
	return
}

func createLittleTktCreate(returnTktNo int, crTktInfo Object.TktCreateInfo, tktModels Object.TktModels, ctx iris.Context, request *RequestCreateLittleTkt) (response ResponseCreateLittleTkt) {
	if crTktInfo.TktInfo == nil || len(crTktInfo.TktInfo) < 1 {
		response = GetResponseCreateLittleTktError(request, errors.New("TktInfo列表不能为空"), ctx.GetStatusCode())
		return
	}
	appId := crTktInfo.TktInfo[0].AppId
	for _, val := range crTktInfo.TktInfo {
		if val.AppId != appId {
			response = GetResponseCreateLittleTktError(request, errors.New("一次请求APPID必须相同"), ctx.GetStatusCode())
			return
		}
	}

	pzR := repository.PeiZhRepository{}
	dbConnList, err := pzR.GetXtMappingDbConnInfo(request.AppId, "DB_TicketHx", "TicketHx")
	if err != nil {
		response = GetResponseCreateLittleTktError(request, err, ctx.GetStatusCode())
		return
	}

	if len(dbConnList) < 1 {
		err = errors.New("传入的APPID无效（APPID错误或配置库缺少配置）")
		response = GetResponseCreateLittleTktError(request, err, ctx.GetStatusCode())
		return
	}

	hxR := repository.HeXRepository{}
	dbConn := make([]*sql.DB, 0)
	for _, val := range dbConnList {
		db, err := hxR.GetDbConnByString(val.MConnStr)
		if err != nil {
			response = GetResponseCreateLittleTktError(request, err, ctx.GetStatusCode())
			return
		}
		dbConn = append(dbConn, db)
	}
	if dbConn == nil || len(dbConn) < 1 {
		response = GetResponseCreateLittleTktError(request, errors.New("未获取到有效的hx库连接"), ctx.GetStatusCode())
		return
	}

	//执行Hx库存储过程
	hxDb := dbConn[0]
	err = hxR.CreateLittleTktCreate(hxDb, crTktInfo.TktInfo)
	if err != nil {
		response = GetResponseCreateLittleTktError(request, err, ctx.GetStatusCode())
		return
	}

	response = GetResponseCreateLittleTktError(request, errors.New("Test End"), ctx.GetStatusCode())
	return
}

func GetResponseCreateLittleTktError(request *RequestCreateLittleTkt, err error, httpCode int) (response ResponseCreateLittleTkt) {
	response.BaseFunc(request, &response)
	response.HttpCode = httpCode
	response.ErrorModel = Object.ErrorModel{
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
