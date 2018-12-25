package Object

import (
	"fmt"
	"github.com/Deansquirrel/go-tool"
	"github.com/kataras/iris"
	"github.com/kataras/iris/core/errors"
)

type RequestCreateLittleTkt struct {
	Header requestHeader
	Body   requestBody
}

type requestHeader struct {
	AppId string
}

type requestBody struct {
	TmpCol      int
	TktInfo     []TktInfo
	CrmFqYwInfo CrmFqYwInfo
	CrmCardInfo []CrmCardInfo
	MdFqYwInfo  []MdFqYwInfo
	YwInfo      YwInfo
}

func GetRequestByContext(ctx iris.Context) (request RequestCreateLittleTkt, err error) {
	request.Header.AppId = ctx.GetHeader("appid")
	if request.Header.AppId == "" {
		err = errors.New("appid不允许为空")
	}
	err = ctx.ReadJSON(&request.Body)
	return
}

func (request *RequestCreateLittleTkt) Print() {
	fmt.Println("Header")
	fmt.Println("appid - " + request.Header.AppId)
	fmt.Println("Body")
	fmt.Println(request.Body.TmpCol)
	fmt.Println(go_tool.GetDateTimeStr(request.Body.CrmFqYwInfo.XtTsr))
}
