package api

import (
	"database/sql"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golandprojects/api/middleware"
	"golandprojects/dao"
	"golandprojects/model"
	"golandprojects/utils"
	"log"
	"net/http"
	"strconv"
	"time"
)

const (
	userName = "root"
	Password = "yx041110"
	ip       = "127.0.0.1"
	port     = "3306"
	dbName   = "userdb"
)

func question(c *gin.Context) {
	if err := c.ShouldBind(&model.Question{}); err != nil {
		utils.RespSuccess(c, "verification failed")
		return
	}
	question := c.PostForm("question")
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
	dao.Addquestion(username, question)
	utils.RespSuccess(c, "successfully post the question")
}
func getquestion(c *gin.Context) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", userName, Password, ip, port, dbName))
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	var allquestion []model.Getquestion
	rows, errq := db.Query("select id,username,question from question")
	if errq != nil {
		log.Fatal(errq.Error)
		return
	}
	//遍历结果
	for rows.Next() {
		var u model.Getquestion
		errn := rows.Scan(&u.Id, &u.Username, &u.Question)
		if errn != nil {
			fmt.Printf("%v", errn)
		}

		allquestion = append(allquestion, u)
	}

	c.JSON(http.StatusOK, gin.H{"res": allquestion})

}
func getanswer(c *gin.Context) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", userName, Password, ip, port, dbName))
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	var allanswer []model.Getanswer
	rows, errq := db.Query("select id,questionid,username,answer from answer")
	if errq != nil {
		log.Fatal(errq.Error)
		return
	}
	//遍历结果
	for rows.Next() {
		var u model.Getanswer
		errn := rows.Scan(&u.Id, &u.Questionid, &u.Username, &u.Answer)
		if errn != nil {
			fmt.Printf("%v", errn)
		}

		allanswer = append(allanswer, u)
	}

	c.JSON(http.StatusOK, gin.H{"res": allanswer})
}
func answer(c *gin.Context) {
	if err := c.ShouldBind(&model.Answer{}); err != nil {
		utils.RespSuccess(c, "verification failed")
		return
	}
	answer := c.PostForm("answer")
	question := c.PostForm("questionid")
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
	questionid, abc := strconv.Atoi(question)
	if abc != nil {
		panic(abc)
	}
	flag := dao.Selectquestion(questionid)
	if !flag {
		utils.RespFail(c, "question doesn't exists")
	}
	err := dao.Addanswer(questionid, username, answer)
	if err != nil {
		utils.RespFail(c, "unable to answer")
		return
	}
	utils.RespSuccess(c, "successfully answered the question")

}
func modifya(c *gin.Context) {
	answer := c.PostForm("answer")
	questionstr := c.PostForm("questionid")
	answerstr := c.PostForm("id")
	value, exists := c.Get("username")
	if !exists {
		// 变量不存在，处理错误
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "username not found",
		})
		return
	}
	questionid, err := strconv.Atoi(questionstr)
	if err != nil {
		// 处理错误
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid questionid",
		})
		return
	}
	answerid, err := strconv.Atoi(answerstr)
	if err != nil {
		// 处理错误
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid questionid",
		})
		return
	}
	username, ok := value.(string)
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "username is not a string"})
		return
	}
	dao.Updateanswer(username, answerid, questionid, answer)
	utils.RespSuccess(c, "successfully modify the answer")
}
func modifyq(c *gin.Context) {
	questionstr := c.PostForm("id")
	question := c.PostForm("question")
	value, exists := c.Get("username")
	if !exists {
		// 变量不存在，处理错误
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "username not found",
		})
		return
	}
	questionid, err := strconv.Atoi(questionstr)
	if err != nil {
		// 处理错误
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid questionid",
		})
		return
	}
	username, ok := value.(string)
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "username is not a string"})
		return
	}
	dao.UpdateQuestin(username, questionid, question)
	utils.RespSuccess(c, "successfully modify the question")
}
func getqa(c *gin.Context) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", userName, Password, ip, port, dbName))
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	var allanswer []model.Getanswer
	var allquestion []model.Getquestion
	value, exists := c.Get("username")
	if !exists {
		// 变量不存在，处理错误
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "username not found",
		})
		return
	}
	rowsa, errqa := db.Query("select id,questionid,username,answer from answer where username = ?", value)
	if errqa != nil {
		log.Fatal(errqa.Error)
		return
	}
	rowsq, errqq := db.Query("select id,username,question from question where username = ?", value)
	if errqq != nil {
		log.Fatal(errqq.Error)
		return
	}
	for rowsq.Next() {
		var u model.Getquestion
		errn := rowsq.Scan(&u.Id, &u.Username, &u.Question)
		if errn != nil {
			fmt.Printf("%v", errn)
		}

		allquestion = append(allquestion, u)
	}

	c.JSON(http.StatusOK, gin.H{"res": allquestion})
	c.JSON(http.StatusOK, gin.H{"res": allanswer})
	for rowsa.Next() {
		var u model.Getanswer
		errn := rowsa.Scan(&u.Id, &u.Questionid, &u.Username, &u.Answer)
		if errn != nil {
			fmt.Printf("%v", errn)
		}

		allanswer = append(allanswer, u)
	}

	c.JSON(http.StatusOK, gin.H{"res": allanswer})
}
func register(c *gin.Context) {
	if err := c.ShouldBind(&model.User{}); err != nil {
		utils.RespSuccess(c, "verification failed")
		return
	}
	// 传入用户名和密码
	username := c.PostForm("username")
	password := c.PostForm("password")

	// 验证用户名是否重复
	flag := dao.SelectUser(username)
	fmt.Println(flag)
	if flag {
		// 以 JSON 格式返回信息
		utils.RespFail(c, "user already exists")
		return
	}

	dao.AddUser(username, password)
	// 以 JSON 格式返回信息
	utils.RespSuccess(c, "add user successful")
}
func findpassword(c *gin.Context) {
	// 传入用户名和密码
	username := c.PostForm("username")
	selectPassword := dao.SelectPasswordFromUsername(username)
	flag := dao.SelectUser(username)
	// 不存在则退出
	if !flag {
		// 以 JSON 格式返回信息
		utils.RespFail(c, "user doesn't exists")
		return
	}
	c.SetCookie("gin_demo_cookie", "test", 3600, "/", "localhost", false, true)
	utils.RespSuccess(c, fmt.Sprintf("find successful,%s", selectPassword))
}

// 仅有登录部分有改动
func login(c *gin.Context) {
	if err := c.ShouldBind(&model.User{}); err != nil {
		utils.RespFail(c, "verification failed")
		return
	}
	// 传入用户名和密码
	username := c.PostForm("username")
	password := c.PostForm("password")

	// 验证用户名是否存在
	flag := dao.SelectUser(username)
	// 不存在则退出
	if !flag {
		// 以 JSON 格式返回信息
		utils.RespFail(c, "user doesn't exists")
		return
	}

	// 查找正确的密码
	selectPassword := dao.SelectPasswordFromUsername(username)
	// 若不正确则传出错误
	if selectPassword != password {
		// 以 JSON 格式返回信息
		utils.RespFail(c, "wrong password")
		return
	}
	c.SetCookie("gin_demo_cookie", "test", 3600, "/", "localhost", false, true)
	utils.RespSuccess(c, "login successful")
	// 正确则登录成功
	// 创建一个我们自己的声明
	claim := model.MyClaims{
		Username: username, // 自定义字段
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 2).Unix(), // 过期时间
			Issuer:    "ferdinand",                          // 签发人
		},
	}
	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	// 使用指定的secret签名并获得完整的编码后的字符串token
	tokenString, _ := token.SignedString(middleware.Secret)
	utils.RespSuccess(c, tokenString)
}
func changepassword(c *gin.Context) {
	username := c.PostForm("username")
	oldpassword := c.PostForm("oldpassword")
	newpassword := c.PostForm("newpassword")
	// 验证用户是否存在
	flag := dao.SelectUser(username)
	// 不存在则退出
	if !flag {
		// 以 JSON 格式返回信息
		utils.RespFail(c, "user doesn't exists")
		return
	}
	// 验证旧密码是否正确
	selectPassword := dao.SelectPasswordFromUsername(username)
	if selectPassword != oldpassword {
		utils.RespFail(c, "old password is incorrect")
		return
	}
	// 更新密码
	err := dao.UpdatePassword(username, newpassword, oldpassword)
	if err != nil {
		utils.RespFail(c, "unable to change password")
		return
	}
	// 成功更新密码，返回成功响应
	utils.RespSuccess(c, "password changed")
}

/*func AddMessage(c *gin.Context) {
	Id := c.PostForm("id")
	messageText := c.PostForm("text")
	message := model.Message{
		Id:        Id,
		Text:      messageText,
		CreatedAt: time.Now(),
	}
	dao.AddMessage(message)
	utils.RespSuccess(c, "add message successful")
}*/

// 新增以下代码
func getUsernameFromToken(c *gin.Context) {
	username, _ := c.Get("username")
	utils.RespSuccess(c, username.(string))
}
