package userdata

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"main.go/passportv2"
)

var db *sql.DB

func InitDB() error {
	var err error
	dsn := "root" + ":" + "MysqlTest" + "@tcp(127.0.0.1:3306)/test"

	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Print(err)
		return err
	} else {
		//检测是否连接成功
		err := db.Ping()
		if err != nil {
			log.Print(err)
			return err
		}
		db.SetMaxOpenConns(10)
		fmt.Print("连接成功\n")
	}
	return nil
}

func RecodeConCentRecode(session string, startTime string, submitTime time.Time, username string) error {
	userid := passportv2.QueryIdByUsername(username)
	sqlstr := `SELECT data_id FROM concentrate Where session=?`
	var dataId string
	_ = db.QueryRow(sqlstr, userid).Scan(&dataId)
	// if err != nil && err != sql.ErrNoRows {
	// 	log.Print(err)
	// 	return err
	// }

	//如果不存在记录
	if dataId == "" {
		sqlstr = `INSERT concentrate(timestart,timesub,session_id,user_id) VALUE(?,?,?,?)`
		submitTimestr := submitTime.Format("2006-01-02 15:04:05")
		_, err := db.Exec(sqlstr, startTime, submitTimestr, session, userid)
		if err != nil {
			log.Print("concerterr" + err.Error())
			return err
		}
	} else {
		sqlstr = `UPDATE concentrate SET timestart=?,timesub=? WHERE data_id=?`
		_, err := db.Exec(sqlstr, startTime, submitTime, dataId)
		if err != nil {
			log.Print("concerterr" + err.Error())
			return err
		}
	}

	return nil
}
