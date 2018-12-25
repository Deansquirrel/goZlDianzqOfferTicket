package Object

import "time"

type TktCompeleteInfo struct {
	TktKeyInfo TktKeyInfo
	AppId      string
	AccId      string
	TkeKind    string
	EffDate    time.Time
	DeadLine   time.Time
	CrYwLsh    string
	CrBr       string
	ChTs       int
}
