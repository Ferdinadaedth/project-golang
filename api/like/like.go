package like

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golandprojects/dao"
	"golandprojects/dao/redis"
	"golandprojects/utils"
	"net/http"
)

func Like(c *gin.Context) {
	const notificationType = "给你点赞"

	value, exists := c.Get("username")
	if !exists {
		// 变量不存在，处理错误
		utils.RespFail(c, "username not found")

		return
	}
	username, ok := value.(string)
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "username is not a string"})
		return
	}
	if username == "" {
		utils.RespFail(c, "你还未登录，不能参与带点赞")
		return
	}
	questionID := c.PostForm("questionid")
	err := redis.Like(username, questionID)
	if err != nil {

		utils.RespFail(c, "internal error")
		return
	}

	dao.Notification(questionID, username, notificationType)
	fmt.Println("存储通知成功")

	utils.RespSuccess(c, "点赞成功")
}
func UnLike(c *gin.Context) {
	value, exists := c.Get("username")
	if !exists {
		// 变量不存在，处理错误
		utils.RespFail(c, "username not found")

		return
	}
	username, ok := value.(string)
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "username is not a string"})
		return
	}
	if username == "" {
		utils.RespFail(c, "你还未登录，不能参与带点赞")
		return
	}
	questionID := c.PostForm("questionid")
	err := redis.Unlike(username, questionID)
	if err != nil {
		utils.RespFail(c, "internal error")
		return
	}
	utils.RespSuccess(c, "取消点赞成功")
}
func NumberLike(c *gin.Context) {
	questionID := c.PostForm("questionid")
	likes, err := redis.GetLikes(questionID)
	if err != nil {
		utils.RespFail(c, "internal error")
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"likes": likes,
	})
}

func IsLike(c *gin.Context) {
	value, exists := c.Get("username")
	if !exists {
		// 变量不存在，处理错误
		utils.RespFail(c, "username not found")

		return
	}
	username, ok := value.(string)
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "username is not a string"})
		return
	}
	if username == "" {
		utils.RespFail(c, "你还没有登录")
		return
	}
	questionID := c.PostForm("questionid")

	isLike, err := redis.IsLike(username, questionID)
	if err != nil {
		utils.RespFail(c, "internal error")
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"isLike": isLike,
	})
}
