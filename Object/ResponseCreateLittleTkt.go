package Object

type ResponseCreateLittleTkt struct {
	TmpCol    int             `json:"tmpcol"`
	TktReturn []TktReturnInfo `json:"TktReturn"`

	dbCommitted bool

	ErrorModel  ErrorModel `json:"errormodel"`
	Description string     `json:"description"`

	HttpCode    int  `josn:"httpcode"`
	DBCommitted bool `json:"dbcommitted"`

	IsSuccess bool `json:"issuccess"`
	FindAccid bool `json:"findaccid"`
	HasCommit bool `json:"hascommit"`
}

func NewErrorResponse(model ErrorModel, httpCode int) *ResponseCreateLittleTkt {
	r := new(ResponseCreateLittleTkt)
	r.ErrorModel = model

	//if (httpCode == 503 || httpCode == 403) && r.Description
}
