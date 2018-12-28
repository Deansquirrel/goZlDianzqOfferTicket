package repository

import (
	"database/sql"
	"errors"
	"github.com/Deansquirrel/goZlDianzqOfferTicket/common"
	"strconv"
	"strings"
	"time"
)

var heXDbConn []*sql.DB

type HeXRepository struct {
}

type VersionInfo struct {
	Name string
	Ver  string
	Date time.Time
}

func (hx *HeXRepository) GetVerInfo(conn *sql.DB) (ver VersionInfo, err error) {
	stmt, err := conn.Prepare("select svname,svver,svdate from xtselfver")
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

//获取核销库连接对象
func GetHxDbConn(appId string) ([]*sql.DB, error) {
	if heXDbConn != nil {
		flag := false
		for _, conn := range heXDbConn {
			err := conn.Ping()
			if err != nil {
				flag = true
				break
			}
		}
		if !flag {
			return heXDbConn, nil
		}
	}

	heXDbConn = make([]*sql.DB, 0)

	conn, err := GetPeiZhDbConn()
	if err != nil {
		return nil, err
	}
	defer func() {
		errLs := conn.Close()
		if errLs != nil {
			common.MyLog(errLs.Error())
		}
	}()

	stmt, err := conn.Prepare("select mconnstr from xtmappingdbconn where appid = ? and miidtype = 'TicketHx'")
	if err != nil {
		return nil, err
	}
	defer func() {
		errLs := stmt.Close()
		if errLs != nil {
			common.MyLog(errLs.Error())
		}
	}()

	rows, err := stmt.Query(appId)
	if err != nil {
		return nil, err
	}
	defer func() {
		errLs := rows.Close()
		if errLs != nil {
			common.MyLog(errLs.Error())
		}
	}()

	var infoList []string
	for rows.Next() {
		var info string
		err := rows.Scan(&info)
		if err != nil {
			return nil, err
		}
		infoList = append(infoList, info)
	}
	if len(infoList) > 0 {
		for _, s := range infoList {
			conn, err := getDbConnByString(s)
			if err != nil {
				return nil, err
			}
			heXDbConn = append(heXDbConn, conn)
		}
	}
	return heXDbConn, nil
}

//解析配置字符串,并获取连接
func getDbConnByString(s string) (*sql.DB, error) {
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
