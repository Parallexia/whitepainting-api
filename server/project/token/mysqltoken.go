package token

import (
	"database/sql"
	"log"

	tsgutils "github.com/typa01/go-utils"
)

var db *sql.DB

//两个值的token，一个是value，一个是GUID
func NewSqlToken2Values(form1 string, value1 string, form2 string, value2 string) error {
	sqlstr := `INSERT INTO token(?,?,?) VALUES(?,?,?); `
	guid := tsgutils.GUID()
	_, err := db.Exec(sqlstr, form1, form2, "token", value1, value2, guid)
	if err != nil {
		log.Print("tokenMysql" + err.Error())
		return err
	}
	return nil
}

//验证两个值的Token
func VerifySqlToken2Values(form1 string, value1 string, form2 string, value2 string, token string) (bool, error) {
	result := ""
	sqlstr := `select token from token where ?=? and where ?=?`
	err := db.QueryRow(sqlstr, form1, value1, form2, value2).Scan(&result)
	if err != nil {
		log.Print("tokenMysql" + err.Error())
		return false, err
	}

	if result == token {
		return true, nil
	}
	return false, nil
}
