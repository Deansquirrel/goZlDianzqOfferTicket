package Object

type TktCreateRequest struct {
	Appid string

	tmpcol      int
	TktInfo     []TktCompeleteInfo
	CrmFqYwInfo CrmCrTktYwInfo
	CrmCardInfo []CrmCrTktCardInfo
	MdFqYwInfo  []MdCrTktYwInfo
	YwInfo      TktYwCommon
}
