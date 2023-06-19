package api

import (
	"github.com/gin-gonic/gin"
	"golandprojects/api/middleware"
)

func InitRouter() {
	r := gin.Default()
	r.Use(middleware.CORS())
	r.POST("/verify", middleware.JWTAuthMiddleware())
	r.POST("/register", register) // 注册
	r.POST("/login", login)       // 登录
	r.POST("/login/changepassword", middleware.JWTAuthMiddleware(), changepassword)
	r.POST("/findpassword", findpassword)
	r.POST("/login/question", middleware.JWTAuthMiddleware(), question)
	r.POST("/login/answer", middleware.JWTAuthMiddleware(), answer)
	r.POST("/login/getquestion", middleware.JWTAuthMiddleware(), getquestion)
	r.POST("login/getanswer", middleware.JWTAuthMiddleware(), getanswer)
	r.POST("/login/qa", middleware.JWTAuthMiddleware(), getqa)
	r.POST("/login/modifyq", middleware.JWTAuthMiddleware(), modifyq)
	r.POST("/login/modifya", middleware.JWTAuthMiddleware(), modifya)
	r.POST("/login/deleteq", middleware.JWTAuthMiddleware(), deleteq)
	r.POST("/login/deletea", middleware.JWTAuthMiddleware(), deletea)
	UserRouter := r.Group("/user")
	{
		UserRouter.Use(middleware.JWTAuthMiddleware())
		UserRouter.GET("/get", getUsernameFromToken)
	}

	r.Run(":8088") // 跑在 8088 端口上
}
