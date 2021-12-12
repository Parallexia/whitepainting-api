package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"main.go/passportv2"
	"main.go/push"
	"main.go/userdata"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	err := passportv2.InitDB()
	if err != nil {
		print("\n" + err.Error() + "\n")
	}

	err = push.InitDB()
	if err != nil {
		print("\n" + err.Error() + "\n")
	}

	err = userdata.InitDB()
	if err != nil {
		print("\n" + err.Error() + "\n")
	}

	main := gin.Default()
	main.MaxMultipartMemory = 8 << 20 //设置文件最大为8mb
	defer main.Run(":9000")

	main.Use(Cors())

	main.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"code": "404",
			"msg":  "请求的资源未找到",
		})
	})
	//登录页面
	passportPage := main.Group("/passport")

	passportPage.POST("/login", passportv2.Login)
	passportPage.DELETE("/login", passportv2.ExitLogin)
	passportPage.GET("/login", passportv2.CheckLogin)
	passportPage.POST("/register", passportv2.Register)

	//用户信息页面api
	usertables := main.Group("/userinfo")
	usertables.Use(VerifyToken())
	//修改用户头像
	//usertables.POST("/profileimage", userdata.UploadProfilePicture)
	//文字推送
	usertables.GET("/push", push.PushInspiritWords)
	//专注数据后台交互
	//开始专注
	usertables.POST("/focuson", userdata.StartConsentraterRecord)
	//专注数据
	usertables.PUT("/focuson", userdata.KeepConsentraterRecord)
	//结束专注
	usertables.DELETE("/focuson", userdata.StopConsentraterRecord)

	//社区页面api
	community := main.Group("/community")
	community.Use(VerifyToken())
	//获取社区消息
	community.GET("/message", push.GetMessage)
	//新建社区消息
	community.POST("/message", push.PostMessage)
}

//跨域请求
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin") //请求头部
		if origin != "" {
			//接收客户端发送的origin （重要！）
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			//服务器支持的所有跨域请求的方法
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE")
			//允许跨域设置可以返回其他子段，可以自定义字段
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session,X_Requested_With,Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language,DNT, X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, Content-Type, Pragma")
			// 允许浏览器（客户端）可以解析的头部 （重要）
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers")
			//设置缓存时间
			c.Header("Access-Control-Max-Age", "172800")
			//允许客户端传递校验信息比如 cookie (重要)
			c.Header("Access-Control-Allow-Credentials", "true")
		}

		//允许类型校验
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "ok!")
		}

		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic info is: %v", err)
			}
		}()
		c.Next()
	}
}

//需要登录的操作的验证中间件
func VerifyToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookies := c.Request.Cookies()

		access := false

		for i := 0; i < len(cookies); i++ {
			item := cookies[i]
			token, err := passportv2.VerifyUserSession(item.Name, item.Value)
			if err != nil {
				log.Print("main.go" + err.Error())
			}

			if token {
				access = true
				c.Set("username", item.Name)
				c.Set("token", item.Value)
				c.Next()
			}
		}

		if access {
			return
		}

		c.JSON(http.StatusNotFound, gin.H{
			"code": "404",
			"msg":  "请求的资源不可用",
		})

		c.Abort()
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic info is: %v", err)
			}
		}()
	}
}
