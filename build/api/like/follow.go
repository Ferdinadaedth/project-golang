package like

import (
	"github.com/gin-gonic/gin"
	"golandprojects/dao/redis"
	"golandprojects/utils"
	"net/http"
)

func Follow(c *gin.Context) {
	followed := c.PostForm("followed")
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
	err := redis.Follow(username, followed)
	if err != nil {
		utils.RespFail(c, "internal error")
		return
	}
	utils.RespSuccess(c, "关注成功")
}

func Unfollow(c *gin.Context) {
	followed := c.PostForm("followed")
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
	err := redis.UnFollow(username, followed)
	if err != nil {
		utils.RespFail(c, "internal error")
		return
	}
	utils.RespFail(c, "取消关注成功")
}

func IsFollowed(c *gin.Context) {
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

	BeUsername := c.PostForm("BeUsername")

	followed, err := redis.IsFollowed(username, BeUsername)
	if err != nil {
		utils.RespFail(c, "internal error")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"IsFollowed": followed,
	})

}
