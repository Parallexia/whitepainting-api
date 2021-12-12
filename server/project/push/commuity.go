package push

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"main.go/passportv2"
	"main.go/token"
)

type CommunityInfo struct {
	SubmitTime time.Time //递交时间

	User_id string //递交人id,使用id防止邮箱与用户名的更换

	Username string //递交人用户

	Content string //内容

	PicUrl string //包含的图片网址
}

//获取服务器的消息
func GetMessage(c *gin.Context) {
	//查询该会话的值，如果这个值没有对应的页码，就添加页码
	//usernameInterface, _ := c.Get("username")
	usertokenInterface, _ := c.Get("token")

	//username := Strval(usernameInterface)
	usertoken := Strval(usertokenInterface)
	//查询值
	pagesString, err := token.GetValueFromTokenKey(usertoken)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": "500",
			"msg":  "未知错误",
		})
		return
	}

	//如果没有数据则把页码设置为0
	if pagesString == "" {
		pagesString = "0"
	}

	pages, err := strconv.Atoi(pagesString)
	if err == strconv.ErrSyntax {
		c.JSON(http.StatusOK, gin.H{
			"code": "400",
			"msg":  "输入数据错误",
		})
		return
	}

	token.NewTokenFromTokenValue(usertoken, strconv.Itoa(pages))

	exist := c.Query("reset")
	//如果输入了请求，那么过期时间就重新设置
	if exist == "1" {
		pages = 0
	} else {
		pages += 1
	}

	token.ResetExpireTime(usertoken, "3600")
	//开始信息的查询，一次返回10行
	data, err := QueryCommunity(pages)
	if err != nil && data == nil {
		c.JSON(http.StatusOK, gin.H{
			"code": "500",
			"msg":  "未知错误",
		})
		return
	}

	for i := 0; i < len(data); i++ {
		data[i].Username = passportv2.QueryUsernameById(data[i].User_id)
		data[i].User_id = "0"
	}

	response, err := json.Marshal(data)
	if err != nil && response == nil {
		c.JSON(http.StatusOK, gin.H{
			"code": "500",
			"msg":  "未知错误",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    "200",
		"content": response,
	})
}

//上传数据，使用POST请求
func PostMessage(c *gin.Context) {
	usernameInterface, _ := c.Get("username")
	username := Strval(usernameInterface)

	content := c.PostForm("content")
	pic_url := c.PostForm("pic_url")

	if content == "" || username == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": "400",
			"msg":  "输入数据错误",
		})
		return
	}

	//为了确保数据与账号联系在一起，使用用户的id而不是用户名
	userid := passportv2.QueryIdByUsername(username)

	data := new(CommunityInfo)
	data.SubmitTime = time.Now()
	data.User_id = userid
	data.Content = content
	data.PicUrl = pic_url

	err := AddCommunityInfo(*data)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": "500",
			"msg":  "未知错误",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": "200",
		"msg":  "添加成功",
	})
}

//删除某条消息

func Strval(value interface{}) string {
	var key string
	if value == nil {
		return key
	}
	switch value.(type) {
	case float64:
		ft := value.(float64)
		key = strconv.FormatFloat(ft, 'f', -1, 64)
	case float32:
		ft := value.(float32)
		key = strconv.FormatFloat(float64(ft), 'f', -1, 64)
	case int:
		it := value.(int)
		key = strconv.Itoa(it)
	case uint:
		it := value.(uint)
		key = strconv.Itoa(int(it))
	case int8:
		it := value.(int8)
		key = strconv.Itoa(int(it))
	case uint8:
		it := value.(uint8)
		key = strconv.Itoa(int(it))
	case int16:
		it := value.(int16)
		key = strconv.Itoa(int(it))
	case uint16:
		it := value.(uint16)
		key = strconv.Itoa(int(it))
	case int32:
		it := value.(int32)
		key = strconv.Itoa(int(it))
	case uint32:
		it := value.(uint32)
		key = strconv.Itoa(int(it))
	case int64:
		it := value.(int64)
		key = strconv.FormatInt(it, 10)
	case uint64:
		it := value.(uint64)
		key = strconv.FormatUint(it, 10)
	case string:
		key = value.(string)
	case []byte:
		key = string(value.([]byte))
	default:
		newValue, _ := json.Marshal(value)
		key = string(newValue)
	}
	return key
}
