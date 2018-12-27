package Object

import (
	"database/sql"
	"fmt"
	"github.com/Deansquirrel/goZlDianzqOfferTicket/global"
	_ "github.com/denisenkom/go-mssqldb"
)

type PeiZh struct {
	AppId    string
	Server   string
	Port     int
	DbName   string
	User     string
	PassWord string

	hxDbInfo []hxDbInfo
}

type hxDbInfo struct {
	Id       int
	Server   string
	Port     int
	DbName   string
	User     string
	PassWord string
}

func (pz *PeiZh) GetHxDbInfo() ([]hxDbInfo, error) {
	if pz.hxDbInfo == nil {
		err := pz.getHxDbInfo()
		if err != nil {
			return nil, err
		} else {
			return pz.hxDbInfo, nil
		}
	}
	return pz.hxDbInfo, nil
}

func (pz *PeiZh) getHxDbInfo() error {
	connString := fmt.Sprintf("server=%s;database=%s;user id=%s;password=%s;port=%d;encrypt=disable", pz.Server, pz.DbName, pz.User, pz.PassWord, pz.Port)
	global.MyLog(connString)
	conn, err := sql.Open("mssql", connString)
	if err != nil {
		global.MyLog("Open Error")
		return err
	}
	defer func() {
		err := conn.Close()
		if err != nil {
			global.MyLog(err.Error())
		}
	}()

	stmt, err := conn.Prepare(`select * from zlaccount40`)
	if err != nil {
		global.MyLog(err.Error())
		return err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	defer rows.Close()

	//err = conn.Ping()
	//
	//if err != nil {
	//	global.MyLog("Ping Error")
	//	return err
	//}
	//global.MyLog("TestSuccess")

	//
	//rows,err := stmt.Query()
	//if err != nil {
	//	return err
	//}
	//fmt.Println(rows)

	//rows,err := conn.Query("select miid,mconnstr from xtmappingdbconn where appid = ? and miidtype= ?",pz.AppId,"TicketHx")
	//if err != nil {
	//	return err
	//}
	//fmt.Println(rows)
	return nil
}
