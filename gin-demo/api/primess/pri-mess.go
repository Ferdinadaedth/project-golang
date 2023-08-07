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
	Password = "yx041110"
	ip       = "127.0.0.1"
	port     = "3306"
	dbName   = "userdb"
)

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
