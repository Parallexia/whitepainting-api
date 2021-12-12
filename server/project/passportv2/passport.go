package passportv2

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"html"
	"log"
	"net/http"
	"regexp"

	//"strconv"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"

	//"github.com/mojocn/base64Captcha"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

//初始化sql数据库
var db *sql.DB
var website = "localhost"

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

//连接redis数据库
func DialRedis() (redis.Conn, error) {
	redisconnect := "127.0.0.1:6379"
	//redispassword := "123456"
	redis.DialConnectTimeout(time.Duration(10))
	//option := new(redis.DialOption)

	c, err := redis.Dial("tcp", redisconnect /*, redis.DialPassword(redispassword), *option*/)
	if err != nil {
		log.Print("Connect to redis error", err)
		return nil, err
	}
	return c, nil
}

func Register(c *gin.Context) {
	//获取用户名
	username, exist := c.GetPostForm("username")
	if !exist || username == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": "300",
			"msg":  "没有输入用户名",
		})
		return
	}
	//检测用户名长度
	if len(username) < 6 || len(username) > 32 {
		c.JSON(http.StatusOK, gin.H{
			"code": "307",
			"msg":  "用户名过短或过长",
		})
		return
	}

	var nullSHA256 = "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855" //需要修改
	//获取密码
	password, exist := c.GetPostForm("password")
	if !exist || password == nullSHA256 {
		c.JSON(http.StatusOK, gin.H{
			"code": "301",
			"msg":  "没有输入密码",
		})
		return
	}

	//获取邮箱
	email, exist := c.GetPostForm("email")
	if !exist || email == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": "302",
			"msg":  "没有输入邮箱",
		})
		return
	}
	//测试邮箱格式
	isemail := VerifyEmailFormat(email)
	if !isemail {
		c.JSON(http.StatusOK, gin.H{
			"code": "305",
			"msg":  "邮箱格式错误",
		})
		return
	}
	username = html.EscapeString(username)
	DBemail := QueryEmailIsExist(email)

	if DBemail {
		c.JSON(http.StatusOK, gin.H{
			"code": "231",
			"msg":  "邮箱已被注册",
		})
		return
	}

	DBuser := QueryUsernameIsExist(username)
	if DBuser {
		c.JSON(http.StatusOK, gin.H{
			"code": "230",
			"msg":  "用户名已被注册",
		})
		return
	}

	//返回token

	token, err := NewUserSession(username)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": "500",
			"msg":  "未知错误",
		})

		return
	}

	err = AddUser(username, password, email)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": "500",
			"msg":  "未知错误",
		})
		log.Print(err)
		return
	}

	c.SetCookie(username, token, 24*60*60, "/", website, false, true)
	c.JSON(http.StatusOK, gin.H{
		"code": "200",
		"msg":  "注册成功",
	})
}

//邮箱正则匹配
func VerifyEmailFormat(email string) bool {
	pattern := `^[a-zA-Z0-9_-]+@[a-zA-Z0-9_-]+(\.[a-zA-Z0-9_-]+)+$` //匹配电子邮箱
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(email)
}

func VerifyHex(hex string) bool {
	pattern := `^[A-Fa-f0-9]+$`
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(hex)
}

//未完成
//接收邮件之后打开网址，路由应该设置到网址所指向的位置
func RegVerify(c *gin.Context) {
	//添加用户
	//发送token标记
}

func Login(c *gin.Context) {
	//读取用户名/邮箱
	username, usernameExist := c.GetPostForm("username")
	email, exist2 := c.GetPostForm("email")
	cookiRequset := c.Request.Cookies()

	for i := 0; i < len(cookiRequset); i++ {
		item := cookiRequset[i]
		//
		isLogged, err := VerifyUserSession(item.Name, item.Value)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": "500",
				"msg":  "未知错误",
			})
			return
		}
		if isLogged {
			c.JSON(http.StatusOK, gin.H{
				"code": "208",
				"msg":  "已登录",
			})
			return
		}
	}

	if !usernameExist && !exist2 || email == "" && username == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": "305",
			"msg":  "没有输入用户名或者邮箱",
		})
		return
	}

	if usernameExist && exist2 {
		c.JSON(http.StatusOK, gin.H{
			"code": "310",
			"msg":  "输入了多个用户标识",
		})
		return
	}
	//读取密码
	passwd, exist := c.GetPostForm("password")
	if !exist || passwd == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": "306",
			"msg":  "没有输入密码",
		})
		return
	}

	if !VerifyHex(passwd) {
		c.JSON(http.StatusOK, gin.H{
			"code": "317",
			"msg":  "密码错误",
		})
		return
	}

	//与数据库比较
	DBpasswd := ""
	salt := ""
	if !usernameExist {
		DBpasswd = hex.EncodeToString(QueryPasswd(QueryUsername(email)))
		salt = hex.EncodeToString(QuerySalt(QueryUsername(email)))
	} else {
		DBpasswd = hex.EncodeToString(QueryPasswd(username))
		salt = hex.EncodeToString(QuerySalt(username))
	}

	compare := Encrypt(passwd, salt)
	if compare != DBpasswd {
		c.JSON(http.StatusOK, gin.H{
			"code": "201",
			"msg":  "密码错误",
		})
		return
	}

	//返回token
	token, err := NewUserSession(username)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": "500",
			"msg":  "未知错误",
		})
		return
	}

	c.SetCookie(username, token, 24*60*60, "/", website, false, true)
	c.JSON(http.StatusOK, gin.H{
		"code": "200",
		"msg":  "登录成功",
	})
}

func CheckLogin(c *gin.Context) {
	cookiRequset := c.Request.Cookies()
	for i := 0; i < len(cookiRequset); i++ {
		item := cookiRequset[i]
		//
		isLogged, err := VerifyUserSession(item.Name, item.Value)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": "500",
				"msg":  "未知错误",
			})
			return
		}
		if isLogged {
			//ExpireSession(item.Value)
			c.JSON(http.StatusOK, gin.H{
				"code": "200",
				"msg":  "已登录",
			})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"code": "201",
		"msg":  "未登录",
	})
}

func ExitLogin(c *gin.Context) {
	cookiRequset := c.Request.Cookies()

	for i := 0; i < len(cookiRequset); i++ {
		item := cookiRequset[i]
		//
		isLogged, err := VerifyUserSession(item.Name, item.Value)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": "500",
				"msg":  "未知错误",
			})
			return
		}
		if isLogged {
			ExpireSession(item.Name)
			c.SetCookie(item.Name, "Shimin Li", 1, "/", website, false, true)
			c.JSON(http.StatusOK, gin.H{
				"code": "200",
				"msg":  "已退出登录",
			})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"code": "201",
		"msg":  "未登录",
	})
}

//未完成架构有问题
//通过邮箱验证改变密码
func ChangePasswd(c *gin.Context) {
	email, exist := c.GetPostForm("email")
	if !exist || email == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": "300",
			"msg":  "没有输入邮箱",
		})
		return
	}
	_ = email
	//验证码
	verify, exist := c.GetPostForm("verify")
	if !exist || verify == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": "301",
			"msg":  "没有输入验证码",
		})
		return
	}
	_ = verify
}

//添加用户数据
func AddUser(username string, password string, email string) error {
	salt := SaltDefault()
	enc := Encrypt(password, hex.EncodeToString(salt))
	encHex, err := hex.DecodeString(enc)
	if err != nil {
		log.Print(err)
		return err
	}

	sqlstr := `insert into users(username,email,salt,password) values(?,?,?,?)`
	_, err = db.Exec(sqlstr, username, email, salt, encHex)
	if err != nil {
		log.Print(err)
		return err
	}
	return nil
}

//生成盐
func SaltDefault() []byte {
	result, err := rand.Prime(rand.Reader, 32)
	if err != nil {
		log.Print(err)
	}

	return result.Bytes()
}

//生成验证码
func VerifyCodeImg() {

}

func VerifyCodeAdu() {

}

func Encrypt(password string, salt string) string {
	enc := password + salt

	o := sha256.New()

	o.Write([]byte(enc))

	encpasswd := hex.EncodeToString(o.Sum(nil))

	return encpasswd
}
