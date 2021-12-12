package userdata

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

//上传图像文件，POST请求
func UploadProfilePicture(c *gin.Context) {
	addr := c.PostForm("url")
	if addr == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": "400",
			"msg":  "没有输入网址",
		})
		return
	}

	//将网址放入数据库
}

//使用默认头像,GET请求
func UseDefaultProfilePicture(c *gin.Context) {

}

//读取用户头像
func ReadUserProfilePicture(c *gin.Context) {

}
