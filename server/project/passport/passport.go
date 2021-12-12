package passport

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"

	//"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

/*
func passportFmw() {
	//连接数据库
	err := InitDB()
	if err != nil {
		return
	}

	r := gin.Default()
	//登录页面
	passport := r.Group("/passport")

	passport.POST("/login", Login)
	passport.POST("/reg", Register)

	//运行数据库
	r.Run(":9000")
}
*/
//
func InitDB() (err error) {

	//sqlname : sqlpassword
	dsn := "root" + ":" + "MysqlTest" + "@tcp(127.0.0.1:3306)/test"

	db, err = sql.Open("mysql", dsn)
	if err != nil {
		fmt.Print(err.Error())
		log.Fatal(err)
		return err
	} else {
		//检测是否连接成功
		err := db.Ping()
		if err != nil {
			fmt.Print(err)
			log.Fatal(err)
			return err
		}
		db.SetMaxOpenConns(10)
		fmt.Print("连接成功")
	}
	return nil
}

//对密码进行加密储存的函数
func Encrypt(password string, salt string) string {
	enc := password + salt
	encpasswd := SHA256(enc)
	return encpasswd
}

//SHA1摘要算法
func SHA256(password string) string {

	o := sha256.New()

	o.Write([]byte(password))

	return hex.EncodeToString(o.Sum(nil))

}
