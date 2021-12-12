package passport

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func changePassword(c *gin.Context) {
	email, exist := c.GetPostForm("email")
	if !exist {
		c.JSON(http.StatusOK, gin.H{
			"code": "301",
			"msg":  "EmailNotInput",
		})

		return
	}
	queryEmail(email)
}

func queryEmail(email string) {
	_ = email
}
