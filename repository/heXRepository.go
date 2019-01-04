package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/Deansquirrel/goZlDianzqOfferTicket/Object"
	"github.com/Deansquirrel/goZlDianzqOfferTicket/common"
	"strconv"
	"strings"
	"time"
)

//var heXDbConn []*sql.DB

type HeXRepository struct {
}

type VersionInfo struct {
	Name string
	Ver  string
	Date time.Time
}

func (hx *HeXRepository) CreateLittleTktCreate(conn *sql.DB, tktInfo []Object.TktInfo) error {

	if conn == nil {
		return errors.New("数据库连接不能为空")
	}
	if tktInfo == nil || len(tktInfo) < 1 {
		return errors.New("传入列表不能为空")
	}

	ctx := context.TODO()
	defer ctx.Done()
	c, err := conn.Conn(ctx)
	if err != nil {
		return err
	}
	defer func() {
		errLs := c.Close()
		if errLs != nil {
			common.MyLog(errLs.Error())
		}
	}()

	tx,err := conn.Begin()
	if err != nil {
		return err
	}
	defer func(){
		if err != nil {
			_ = tx.Rollback()
		} else {
			_ = tx.Commit()
		}
	}()

	stmt,err := tx.Prepare(getCreateTempTableTktInfoSqlStr() + " " + getInsertTempTableTktInfoSqlStr() + " " + getExecProc() + " " + getDropTmepTableTktInfoSqlStr())
	if err != nil {
		return err
	}
	defer func(){
		_ = stmt.Close()
	}()

	for _,val := range tktInfo{
		_,err = stmt.Exec(val.AppId,val.AccId,val.TktNo,val.CashMy,val.AddMy,val.TktName,val.TktKind,val.PCno,val.EffDate,val.Deadline,val.CrYwLsh,val.CrBr)
		if err != nil {
			return err
		}
	}

	return nil
}

func (hx *HeXRepository) GetVerInfo(conn *sql.DB) (ver VersionInfo, err error) {
	stmt, err := conn.Prepare("" +
		"select svname,svver,svdate from xtselfver")
	if err != nil {
		common.MyLog(err.Error())
		return
	}
	defer func() {
		errLs := stmt.Close()
		if errLs != nil {
			common.MyLog(errLs.Error())
		}
	}()
	rows, err := stmt.Query()
	if err != nil {
		common.MyLog(err.Error())
		return
	}
	defer func() {
		errLs := rows.Close()
		if errLs != nil {
			common.MyLog(errLs.Error())
		}
	}()
	for rows.Next() {
		err = rows.Scan(&ver.Name, &ver.Ver, &ver.Date)
		if err != nil {
			return
		}
	}
	return
}

//解析配置字符串,并获取连接
func (hx *HeXRepository) GetDbConnByString(s string) (*sql.DB, error) {
	config := strings.Split(s, "|")
	if len(config) != 5 {
		common.MyLog("数据库配置串解析失败 - " + s)
		return nil, errors.New("数据库配置串解析失败")
	}
	port, err := strconv.Atoi(config[1])
	if err != nil {
		common.MyLog("数据库配置串端口解析失败 - " + s)
		return nil, errors.New("数据库配置串端口解析失败")
	}
	return GetDbConn(config[0], port, config[2], config[3], config[4])
}

func getCreateTempTableTktInfoSqlStr() string {
	sqlStr := "" +
		"CREATE TABLE #TktInfo"+
		"("+
		"    Appid	varchar(30),"+
		"    Accid	bigint,"+
		"    Tktno	varchar(30),"+
		"    Cashmy	decimal(18,2),"+
		"    Addmy	decimal(18,2),"+
		"    Tktname	nvarchar(30),"+
		"    TktKind	varchar(30),"+
		"    Pcno	varchar(30),"+
		"    EffDate	smalldatetime,"+
		"    Deadline	smalldatetime,"+
		"    CrYwlsh	varchar(12),"+
		"    CrBr	varchar(30)"+
		")"
	return sqlStr
}

func getInsertTempTableTktInfoSqlStr() string {
	sqlStr := "" +
		"insert into #TktInfo(Appid,Accid,Tktno,Cashmy,Addmy,Tktname,TktKind,Pcno,EffDate,Deadline,CrYwlsh,CrBr) " +
		"select ?,?,?,?,?,?,?,?,?,?,?,?"
	return sqlStr
}

func getExecProc() string {
	sqlStr := "" +
		"exec pr_CreateLittleTkt_Create"
	return sqlStr
}

func getDropTmepTableTktInfoSqlStr() string {
	sqlStr := "" +
		"Drop table #TktInfo"
	return sqlStr
}

