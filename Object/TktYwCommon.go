package Object

import "time"

type TktYwCommon struct {
	OprBrid   string
	OprYwSno  string
	OprPpid   string
	OprTime   time.Time
	OprHsDate time.Time
	OprId     string
	OprAccId  int64
	XxTsr     time.Time
	TsFs      int
	TsNr      string
}
