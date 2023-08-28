package api

import (
	"github.com/gin-gonic/gin"
	"golandprojects/api/like"
	"golandprojects/api/middleware"
	"golandprojects/api/primess"
)

func InitRouter() {
	r := gin.Default()
	r.Use(middleware.CORS())
	r.POST("/verify", middleware.JWTAuthMiddleware())
	r.POST("/register", register) // 注册
	r.POST("/login", login)       // 登录
	r.POST("/findpassword", findpassword)
	r.POST("/getallquestions", GetALlQuestions)
	r.POST("/changepassword", middleware.JWTAuthMiddleware(), changepassword)
	r.POST("/question", middleware.JWTAuthMiddleware(), question)
	r.POST("/gptanswer", answer)
	r.POST("/answer", middleware.JWTAuthMiddleware(), answer)
	r.POST("/getquestion", middleware.JWTAuthMiddleware(), getquestion)
	r.POST("/getanswer", getanswer)
	r.POST("/userquestion", middleware.JWTAuthMiddleware(), getuserquestion)
	r.POST("/useranswer", middleware.JWTAuthMiddleware(), getuseranswer)
	r.POST("/modifyq", middleware.JWTAuthMiddleware(), modifyq)
	r.POST("/modifya", middleware.JWTAuthMiddleware(), modifya)
	r.POST("/deleteq", middleware.JWTAuthMiddleware(), deleteq)
	r.POST("/deletea", middleware.JWTAuthMiddleware(), deletea)
	r.POST("/like", middleware.JWTAuthMiddleware(), like.Like)
	r.POST("/unlike", middleware.JWTAuthMiddleware(), like.UnLike)
	r.POST("/likenumber", like.NumberLike)
	r.POST("/islike", middleware.JWTAuthMiddleware(), like.IsLike)
	r.POST("/follow", middleware.JWTAuthMiddleware(), like.Follow)
	r.POST("/unfollow", middleware.JWTAuthMiddleware(), like.Unfollow)
	r.POST("/isfollowed", middleware.JWTAuthMiddleware(), like.IsFollowed)
	r.POST("/message", middleware.JWTAuthMiddleware(), primess.Message)
	r.POST("/getmessage", middleware.JWTAuthMiddleware(), primess.Getmessage)
	r.POST("/getsendingmessage", middleware.JWTAuthMiddleware(), primess.Getsendingmessage)
	r.POST("/modifym", middleware.JWTAuthMiddleware(), primess.Updatem)
	r.POST("/deletem", middleware.JWTAuthMiddleware(), primess.Deletem)
	UserRouter := r.Group("/user")
	{
		UserRouter.Use(middleware.JWTAuthMiddleware())
		UserRouter.POST("/get", getUsernameFromToken)
	}

	r.Run(":8088") // 跑在 8088 端口上
}
