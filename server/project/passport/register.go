package passport

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"html"
	"log"
	"strconv"

	"net/http"

	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	//没有输入邮箱
	email, exist := c.GetPostForm("email")
	if !exist {
		c.JSON(http.StatusOK, gin.H{
			"msg":  "NoInput",
			"code": "300",
		})

		return
	}

	username, exist := c.GetPostForm("username")
	if !exist {
		c.JSON(http.StatusOK, gin.H{
			"msg":  "NoInput",
			"code": "300",
		})

		return

	} else {
		//转义HTML和JS字符，防止XSS
		username = html.UnescapeString(username)
	}
	password, exist := c.GetPostForm("password")
	if !exist {
		c.JSON(http.StatusOK, gin.H{
			"msg":  "NoInput",
			"code": "300",
		})

		return

	} else {
		salt, _ := SaltDefault()
		err := StoreUserAccount(username, password, email, salt, db)

		if errors.Is(err, errTheSameUserNameExpection) {
			c.JSON(http.StatusOK, gin.H{
				"msg":  "UserAlreadyReg",
				"code": "201",
			})

			return
		}

		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": "301",
				"msg":  "Exception",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"msg":  "RegSuccess",
			"code": "200",
		})

		return
	}

}

//加盐函数
func SaltDefault() (string, error) {
	result, err := rand.Prime(rand.Reader, 32)
	if err != nil {
		log.Print(err)
		return "", err
	}

	return strconv.FormatInt(result.Int64(), 16), nil
}

func StoreUserAccount(username string, password string, email string, salt string, db *sql.DB) error {

	//使用全局变量，确保查询
	var isexistid string
	//查询数据库是否有相同用户名的用户
	sqlstr := `select userid from userlogin where username=?`
	//从连接池中查询单条记录并释放连接
	db.QueryRow(sqlstr, username).Scan(&isexistid)
	//对于用户名已经注册的处理
	if isexistid != "" {
		return errTheSameUserNameExpection
	}

	//使用全局变量，确保查询
	var isexistemail string
	//查询数据库是否有相同邮箱的用户
	sqlstr = `select useremail from userlogin where username=?`
	//从连接池中查询单条记录并释放连接
	db.QueryRow(sqlstr, username).Scan(&isexistemail)
	//对于用户名已经注册的处理
	if isexistid != "" {
		return errTheSameUserEmailExpection
	}

	//将加盐后的密码转换成字节
	enc, err := hex.DecodeString(Encrypt(password, salt))
	if err != nil {
		log.Print(err)
	}
	//将盐值转换成字节
	bytesalt, err := hex.DecodeString(salt)
	if err != nil {
		log.Print(err)
	}

	//插入数据库
	sqlstr = `insert into userlogin(username,password,useremail,salt) values(?,?,?,?)`
	//执行数据库语句
	ret, err := db.Exec(sqlstr, username, enc, email, bytesalt)
	//如果数据库插入数据失败，则执行这条语句
	if err != nil {
		log.Print(err)
		return err
	}

	id, _ := ret.LastInsertId()

	//fmt.Printf("用户%d创建", id)
	log.Print(strconv.FormatInt(id, 10) + "被创建")
	return nil
}
