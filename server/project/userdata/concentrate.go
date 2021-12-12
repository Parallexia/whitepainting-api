package userdata

import (
	"html"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	tsgutils "github.com/typa01/go-utils"
	"main.go/push"
	"main.go/token"
	//"main.go/passportv2"
	//"main.go/token"
)

/*专注计划交互函数*/

//新建专注计划
func CreateConsentratePlan(c *gin.Context) {
	name, exist := c.GetPostForm("name")
	if !exist || name == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": "221",
			"msg":  "没有输入计划名",
		})
		return
	}

	name = html.EscapeString(name)

	startTime, exist := c.GetPostForm("start")
	if !exist || startTime == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": "222",
			"msg":  "没有输入起始时间",
		})
		return
	}

	stopTime, exist := c.GetPostForm("stop")
	if !exist || stopTime == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": "223",
			"msg":  "没有输入结束时间",
		})
		return
	}

	description, exist := c.GetPostForm("description")
	if !exist || description == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": "223",
			"msg":  "没有输入结束时间",
		})
		return
	}
}

/*专注数据交互函数*/

//获得开始专注的消息,POST请求
func StartConsentraterRecord(c *gin.Context) {
	//usernameInterface, _ := c.Get("username")
	//username := push.Strval(usernameInterface)
	guid := tsgutils.GUID()
	startTime := time.Now().UTC().Format("2006-01-02 15:04:05")
	err := token.NewTokenFromTokenValue(guid, startTime)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": "500",
			"msg":  "未知错误",
		})
		return
	}

	token.ResetExpireTime(guid, "120")
	c.JSON(http.StatusOK, gin.H{
		"code":    "200",
		"msg":     "开始记录专注数据",
		"session": guid,
	})
}

//结束专注的消息,DELETE请求
func StopConsentraterRecord(c *gin.Context) {
	usernameInterface, _ := c.Get("username")
	username := push.Strval(usernameInterface)
	session := c.Query("session")

	res, err := token.GetValueFromTokenKey(session)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": "500",
			"msg":  "未知错误",
		})

		return
	}

	if res == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": "300",
			"msg":  "没有找到会话",
		})

		return
	}
	startTime, err := token.GetValueFromTokenKey(session)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": "500",
			"msg":  "未知错误",
		})

		return
	}

	submitTime := time.Now().UTC()
	//插入专注记录数据库
	err = RecodeConCentRecode(session, startTime, submitTime, username)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": "500",
			"msg":  "未知错误",
		})

		return
	}

	_ = token.DeleteToken(session)

	c.JSON(http.StatusOK, gin.H{
		"code": "200",
		"msg":  "数据交互成功",
	})
}

//保持专注消息的函数,PUT请求
func KeepConsentraterRecord(c *gin.Context) {
	session := c.PostForm("session")

	usernameInterface, _ := c.Get("username")
	username := push.Strval(usernameInterface)

	res, err := token.GetValueFromTokenKey(session)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": "500",
			"msg":  "未知错误",
		})

		return
	}

	if res == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": "300",
			"msg":  "会话错误",
		})

		return
	}

	token.ResetExpireTime(session, "120")
	startTime, err := token.GetValueFromTokenKey(session)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": "500",
			"msg":  "未知错误",
		})

		return
	}
	stopTime := time.Now().UTC()

	RecodeConCentRecode(session, startTime, stopTime, username)
	c.JSON(http.StatusOK, gin.H{
		"code": "200",
		"msg":  "打卡成功",
	})
}
