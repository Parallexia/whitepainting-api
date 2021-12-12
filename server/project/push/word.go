package push

import (
	//"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func PushInspiritWords(c *gin.Context) {
	id := c.Query("id")
	sqlstr := `select content from words where data_id=?`
	word := ""

	err := db.QueryRow(sqlstr, id).Scan(&word)
	if err != nil {
		log.Print("push.go" + err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": "500",
			"msg":  "获取失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    "200",
		"msg":     "获取成功",
		"context": word,
	})
}
