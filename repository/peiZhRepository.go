package repository

import (
	"database/sql"
	"github.com/Deansquirrel/goZlDianzqOfferTicket/global"
	"time"
)

var peiZhDbConn *sql.DB

type peiZhRepository struct {
}

//获取配置库连接对象
func GetPeiZhDbConn() (*sql.DB, error) {
	if peiZhDbConn != nil {
		err := peiZhDbConn.Ping()
		if err != nil {
			return peiZhDbConn, nil
		}
	}
	conn, err := GetDbConn(global.Config.PeiZhDbConfig.Server, global.Config.PeiZhDbConfig.Port, global.Config.PeiZhDbConfig.DbName, global.Config.PeiZhDbConfig.User, global.Config.PeiZhDbConfig.PassWord)
	if err != nil {
		return nil, err
	}

	err = conn.Ping()
	if err != nil {
		return nil, err
	}

	conn.SetMaxIdleConns(30)
	conn.SetMaxOpenConns(30)
	conn.SetConnMaxLifetime(time.Second * 60 * 10)
	peiZhDbConn = conn

	return conn, nil
}
