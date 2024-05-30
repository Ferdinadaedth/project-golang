package primess

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"golandprojects/dao"
	"golandprojects/model"
	"golandprojects/utils"
	"log"
	"net/http"
	"strconv"
)

const (
	userName = "root"
	Password = "h74o+JIi5SpSY3MU"
	ip       = "47.108.208.111"
	port     = "3306"
	dbName   = "userdb"
)

func GetNotification(c *gin.Context) {
	// 打开数据库连接
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", userName, Password, ip, port, dbName))
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	var notifications []model.Notification
	value, exists := c.Get("username")
	if !exists {
		// 变量不存在，处理错误
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "username not found",
		})
		return
	}
	username, ok := value.(string)
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "username is not a string"})
		return
	}
	fmt.Println(username)

	str1 := "SELECT userid from user WHERE username = ?"
	row, err := db.Query(str1, username)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	defer row.Close()
	var rID string
	row.Next()
	err = row.Scan(&rID)
	fmt.Println(rID)
	if err != nil {
		panic(err)
	}

	str2 := "SELECT notificationID,recipientUserID,notificationType,notificationTime FROM notification WHERE recipientUserID=?"
	rows, err := db.Query(str2, rID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to execute query"})
		return
	}

	defer rows.Close()
	for rows.Next() {
		var notification model.Notification
		err = rows.Scan(&notification.NotificationID, &notification.RecipientUserID, &notification.NotificationType, &notification.NotificationTime)
		notification.SenderUserName = username
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan result"})
			return
		}
		notifications = append(notifications, notification)
	}

	// 检查是否发生错误
	if err = rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred during rows iteration"})
		return
	}

	// 返回JSON响应
	c.JSON(http.StatusOK, gin.H{"notifications": notifications})
}

func Message(c *gin.Context) {
	if err := c.ShouldBind(&model.Message{}); err != nil {
		utils.RespSuccess(c, "verification failed")
		return
	}
	message := c.PostForm("message")
	receiver := c.PostForm("receiver")
	value, exists := c.Get("username")
	if !exists {
		// 变量不存在，处理错误
		utils.RespFail(c, "receiver not found")

		return
	}
	sender, ok := value.(string)
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "receiver is not a string"})
		return
	}
	dao.Addmessage(receiver, sender, message)
	utils.RespSuccess(c, "成功发送私信")
}
func Getmessage(c *gin.Context) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", userName, Password, ip, port, dbName))
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	var allmessage []model.Getmessage
	value, exists := c.Get("username")
	if !exists {
		// 变量不存在，处理错误
		utils.RespFail(c, "receiver not found")

		return
	}
	username, ok := value.(string)
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "receiver is not a string"})
		return
	}
	rows, errq := db.Query("select messageid,sender,receiver,content from message where receiver = ?", username)
	if errq != nil {
		log.Fatal(errq.Error)
		return
	}
	//遍历结果
	for rows.Next() {
		var u model.Getmessage
		errn := rows.Scan(&u.Messageid, &u.Sender, &u.Receiver, &u.Content)
		if errn != nil {
			fmt.Printf("%v", errn)
		}

		allmessage = append(allmessage, u)
	}

	c.JSON(http.StatusOK, gin.H{"message": allmessage})

}
func Getsendingmessage(c *gin.Context) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", userName, Password, ip, port, dbName))
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	var allmessage []model.Getmessage
	value, exists := c.Get("username")
	if !exists {
		// 变量不存在，处理错误
		utils.RespFail(c, "receiver not found")

		return
	}
	username, ok := value.(string)
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "receiver is not a string"})
		return
	}
	rows, errq := db.Query("select messageid,sender,receiver,content from message where sender = ?", username)
	if errq != nil {
		log.Fatal(errq.Error)
		return
	}
	//遍历结果
	for rows.Next() {
		var u model.Getmessage
		errn := rows.Scan(&u.Messageid, &u.Sender, &u.Receiver, &u.Content)
		if errn != nil {
			fmt.Printf("%v", errn)
		}

		allmessage = append(allmessage, u)
	}

	c.JSON(http.StatusOK, gin.H{"message": allmessage})

}
func Updatem(c *gin.Context) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", userName, Password, ip, port, dbName))
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	content := c.PostForm("content")
	messageid := c.PostForm("messageid")
	stmt, err := db.Prepare("UPDATE `message` SET `content` =? WHERE `messageid` = ?")
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()
	// 执行更新操作
	_, err = stmt.Exec(content, messageid)
	if err != nil {
		utils.RespFail(c, "内部错误")
		return
	}
	utils.RespSuccess(c, "成功修改")
}
func Deletem(c *gin.Context) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", userName, Password, ip, port, dbName))
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	message := c.PostForm("messageid")
	messageid, abc := strconv.Atoi(message)
	if abc != nil {
		panic(abc)
	}
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM message WHERE messageid=?", messageid).Scan(&count)
	if err != nil {
		panic(err.Error())
	}
	if count > 0 {
		_, err = db.Exec("delete from message where messageid=?", messageid)
		if err != nil {
			panic(err.Error())
		}
		utils.RespSuccess(c, "删除成功")
	} else {
		utils.RespFail(c, "no such message")
	}
}
