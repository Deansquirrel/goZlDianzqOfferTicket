package repository

import (
	"bytes"
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

	stmt, err := conn.Prepare(getQueryString(len(tktInfo)))
	if err != nil {
		return err
	}
	defer func() {
		errLs := stmt.Close()
		if errLs != nil {
			common.MyLog(errLs.Error())
		}
	}()

	var c = make([]interface{}, 0)
	var val Object.TktInfo
	for i := 0; i < len(tktInfo); i++ {
		val = tktInfo[i]
		c = append(c, val.AppId)
		c = append(c, val.AccId)
		c = append(c, val.TktNo)
		c = append(c, val.CashMy)
		c = append(c, val.AddMy)
		c = append(c, val.TktName)
		c = append(c, val.TktKind)
		c = append(c, val.PCno)
		c = append(c, val.EffDate)
		c = append(c, val.Deadline)
		c = append(c, val.CrYwLsh)
		c = append(c, val.CrBr)
	}

	_, err = stmt.Exec(c...)
	if err != nil {
		return err
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

func getQueryString(n int) string {
	if n > 0 {
		var buffer bytes.Buffer
		buffer.WriteString(getCreateTempTableTktInfoSqlStr())
		for i := 0; i < n; i++ {
			buffer.WriteString(getInsertTempTableTktInfoSqlStr())
		}
		buffer.WriteString(getExecProc())
		buffer.WriteString(getDropTempTableTktInfoSqlStr())
		return buffer.String()
	} else {
		return ""
	}
}

func getCreateTempTableTktInfoSqlStr() string {

	var buffer bytes.Buffer
	buffer.WriteString("CREATE TABLE #TktInfo")
	buffer.WriteString("(")
	buffer.WriteString("    Appid varchar(30),")
	buffer.WriteString("    Accid bigint,")
	buffer.WriteString("    Tktno varchar(30),")
	buffer.WriteString("    Cashmy decimal(18,2),")
	buffer.WriteString("    Addmy decimal(18,2),")
	buffer.WriteString("    Tktname nvarchar(30),")
	buffer.WriteString("    TktKind	varchar(30),")
	buffer.WriteString("    Pcno varchar(30),")
	buffer.WriteString("    EffDate smalldatetime,")
	buffer.WriteString("    Deadline smalldatetime,")
	buffer.WriteString("    CrYwlsh varchar(12),")
	buffer.WriteString("    CrBr varchar(30)")
	buffer.WriteString(") ")
	return buffer.String()
}

func getInsertTempTableTktInfoSqlStr() string {
	var buffer bytes.Buffer
	buffer.WriteString("insert into #TktInfo(Appid,Accid,Tktno,Cashmy,Addmy,Tktname,TktKind,Pcno,EffDate,Deadline,CrYwlsh,CrBr) ")
	buffer.WriteString("select ?,?,?,?,?,?,?,?,?,?,?,? ")
	return buffer.String()
}

func getExecProc() string {
	sqlStr := "exec pr_CreateLittleTkt_Create "
	return sqlStr
}

func getDropTempTableTktInfoSqlStr() string {
	sqlStr := "Drop table #TktInfo "
	return sqlStr
}
