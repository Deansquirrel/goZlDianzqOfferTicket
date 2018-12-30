package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
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

	_, err = c.ExecContext(ctx, "CREATE TABLE #TktInfo"+
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
		")")
	if err != nil {
		return err
	}

	//stmt,err := c.PrepareContext(ctx,"" +
	//	"insert into #TktInfo(Appid,Accid,Tktno,Cashmy,Addmy,Tktname,TktKind,Pcno,EffDate,Deadline,CrYwlsh,CrBr)" +
	//	"select ?,?,?,?,?,?,?,?,?,?,?,?")
	//if err != nil {
	//	return nil
	//}
	//defer func(){
	//	errLs := stmt.Close()
	//	if errLs != nil {
	//		common.MyLog(errLs.Error())
	//	}
	//}()

	//for _,val := range tktInfo {
	//	_,err = c.ExecContext(ctx,"insert into #TktInfo(Appid,Accid,Tktno,Cashmy,Addmy,Tktname,TktKind,Pcno,EffDate,Deadline,CrYwlsh,CrBr)" +
	//		"select ?,?,?,?,?,?,?,?,?,?,?,?",val.AppId,val.AccId,val.TktNo,val.CashMy,val.AddMy,val.TktName,val.TktKind,val.PCno,val.EffDate,val.Deadline,val.CrYwLsh,val.CrBr)
	//	if err != nil {
	//		return err
	//	}
	//}

	val := tktInfo[0]
	//_,err = c.ExecContext(ctx,"insert into #TktInfo(Appid,Accid,Tktno,Cashmy,Addmy,Tktname,TktKind,Pcno,EffDate,Deadline,CrYwlsh,CrBr)" +
	//	"select ?,?,?,?,?,?,?,?,?,?,?,?",val.AppId,val.AccId,val.TktNo,val.CashMy,val.AddMy,val.TktName,val.TktKind,val.PCno,val.EffDate,val.Deadline,val.CrYwLsh,val.CrBr)
	//if err != nil {
	//	return err
	//}

	s, err := json.Marshal(val)
	fmt.Println(string(s))

	_, err = c.ExecContext(ctx, ""+
		"exec pr_CreateLittleTkt_Create")
	if err != nil {
		return nil
	}

	_, err = c.ExecContext(ctx, ""+
		"drop table #TktInfo")
	if err != nil {
		return nil
	}

	//
	//
	//_,err = c.Exec("CREATE TABLE #TktInfo" +
	//	"(" +
	//	"    Appid	varchar(30)," +
	//	"    Accid	bigint," +
	//	"    Tktno	varchar(30)," +
	//	"    Cashmy	decimal(18,2)," +
	//	"    Addmy	decimal(18,2)," +
	//	"    Tktname	nvarchar(30)," +
	//	"    TktKind	varchar(30)," +
	//	"    Pcno	varchar(30)," +
	//	"    EffDate	smalldatetime," +
	//	"    Deadline	smalldatetime," +
	//	"    CrYwlsh	varchar(12)," +
	//	"    CrBr	varchar(30)" +
	//	")")
	//if err != nil {
	//	return err
	//}
	//_,err = tx.Exec("exec pr_CreateLittleTkt_Create")
	//if err != nil {
	//	return err
	//}
	//_,errLs := tx.Exec("drop table #TktInfo")
	//if errLs != nil {
	//	common.MyLog(errLs.Error())
	//}
	//
	//tx,err := conn.Begin()
	//if err != nil {
	//	return err
	//}
	//defer func(){
	//	var errLs error
	//	switch{
	//	case err != nil:
	//		errLs = tx.Rollback()
	//	default:
	//		errLs = tx.Commit()
	//	}
	//	if errLs != nil {
	//		common.MyLog(errLs.Error())
	//	}
	//}()
	//
	//_,err = tx.Exec("CREATE TABLE #TktInfo" +
	//	"(" +
	//	"    Appid	varchar(30)," +
	//	"    Accid	bigint," +
	//	"    Tktno	varchar(30)," +
	//	"    Cashmy	decimal(18,2)," +
	//	"    Addmy	decimal(18,2)," +
	//	"    Tktname	nvarchar(30)," +
	//	"    TktKind	varchar(30)," +
	//	"    Pcno	varchar(30)," +
	//	"    EffDate	smalldatetime," +
	//	"    Deadline	smalldatetime," +
	//	"    CrYwlsh	varchar(12)," +
	//	"    CrBr	varchar(30)" +
	//	")")
	//if err != nil {
	//	return err
	//}
	//_,err = tx.Exec("exec pr_CreateLittleTkt_Create")
	//if err != nil {
	//	return err
	//}
	//_,errLs := tx.Exec("drop table #TktInfo")
	//if errLs != nil {
	//	common.MyLog(errLs.Error())
	//}
	//

	//stepOne,err := conn.Prepare("" +
	//	"CREATE TABLE #TktInfo" +
	//	"(" +
	//	"    Appid	varchar(30)," +
	//	"    Accid	bigint," +
	//	"    Tktno	varchar(30)," +
	//	"    Cashmy	decimal(18,2)," +
	//	"    Addmy	decimal(18,2)," +
	//	"    Tktname	nvarchar(30)," +
	//	"    TktKind	varchar(30)," +
	//	"    Pcno	varchar(30)," +
	//	"    EffDate	smalldatetime," +
	//	"    Deadline	smalldatetime," +
	//	"    CrYwlsh	varchar(12)," +
	//	"    CrBr	varchar(30)" +
	//	")")
	//if err != nil {
	//	return err
	//}
	//
	//stepTwo,err := conn.Prepare("" +
	//	"exec pr_CreateLittleTkt_Create")
	//if err != nil {
	//	return err
	//}
	//
	//stepThree,err := conn.Prepare("" +
	//	"drop table #TktInfo")
	//if err != nil {
	//	return err
	//}

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

////获取核销库连接对象
//func GetHxDbConn(appId string) ([]*sql.DB, error) {
//	//if heXDbConn != nil {
//	//	flag := false
//	//	for _, conn := range heXDbConn {
//	//		err := conn.Ping()
//	//		if err != nil {
//	//			flag = true
//	//			break
//	//		}
//	//	}
//	//	if !flag {
//	//		return heXDbConn, nil
//	//	}
//	//}
//
//	pZhR := PeiZhRepository{}
//	dbConnInfoList,err := pZhR.GetXtMappingDbConnInfo(appId,"DB_TicketHx","TicketHx")
//	if err != nil {
//		return nil,err
//	}
//	if dbConnInfoList == nil || len(dbConnInfoList)<1 {
//		return nil,errors.New("未获取到核销库连接信息")
//	}
//
//	heXDbConn := make([]*sql.DB, 0)
//
//	if len(dbConnInfoList) > 0 {
//		for _, val := range dbConnInfoList {
//			conn, err := GetDbConnByString(val.MConnStr)
//			if err != nil {
//				return nil, err
//			}
//			heXDbConn = append(heXDbConn, conn)
//		}
//	}
//	return heXDbConn, nil
//}

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
