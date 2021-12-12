package passport

import (
	"encoding/hex"
	"html"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	username, exist := c.GetPostForm("username")
	if !exist {
		c.JSON(http.StatusOK, gin.H{
			"code": "300",
			"msg":  "NoInput",
		})

		return

	}

	password, exist := c.GetPostForm("password")

	/*
		xss过滤模块
	*/

	//没有输入用户名或密码
	if !exist {
		c.JSON(http.StatusOK, gin.H{
			"code": "300",
			"msg":  "NoInput",
		})
	} else {

		err := CheckLogin(username, password)
		if err == errWrongUserNameOrPasswd {
			c.JSON(http.StatusOK, gin.H{
				"code": "201",
				"msg":  "WrongUsernameorPassword",
			})
			return
		}

		if err != nil {
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code": "200",
			"msg":  "LoginSuccess",
		})
	}
}

func CheckLogin(username string, password string) error {
	username = html.EscapeString(username)
	//使用全局变量，确保连接被释放
	//盐值
	var loginQuery []byte
	//密码
	var checkPasswd []byte

	sqlStr := `select salt from userlogin where username=?`

	err := db.QueryRow(sqlStr, username).Scan(&loginQuery)

	if loginQuery == nil {
		return errWrongUserNameOrPasswd
	}

	if err != nil {
		log.Print(err)
		return err
	}

	enc := Encrypt(password, hex.EncodeToString(loginQuery))
	if err != nil {
		log.Print(err)
		return err
	}

	sqlStr = `select password from userlogin where username=?`

	err = db.QueryRow(sqlStr, username).Scan(&checkPasswd)
	if err != nil {
		log.Print(err)
		return err
	}

	cmp := hex.EncodeToString(checkPasswd)

	if enc != cmp {
		return errWrongUserNameOrPasswd
	}

	return nil
}
