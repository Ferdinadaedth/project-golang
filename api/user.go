package api

import (
	"database/sql"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golandprojects/api/middleware"
	"golandprojects/cache"
	"golandprojects/dao"
	"golandprojects/model"
	"golandprojects/utils"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"strconv"
	"time"
)

const (
	userName = "root"
	Password = "h74o+JIi5SpSY3MU"
	ip       = "47.108.208.111"
	port     = "3306"
	dbName   = "userdb"
)

func GetALlQuestions(c *gin.Context) {
	var questions []model.Getquestion
	err := cache.GetCache("allQuestions", &questions)
	if err == nil {
		c.JSON(http.StatusOK, gin.H{
			"status":    200,
			"questions": questions,
		})
		return
	}

	questions, err = dao.GetAllQuestions()
	if err != nil {
		utils.RespFail(c, "get all questions error")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":    200,
		"questions": questions,
	})
	err = cache.SetCache("allQuestions", &questions)
	if err != nil {
		return
	}
}

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
		//c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "username is not a string"})
		return
	}
	dao.Addquestion(username, question)
	err := cache.DeleteCache("allQuestions")
	if err != nil {
		log.Fatal(err.Error)
	}

	utils.RespSuccess(c, "successfully post the question")

}
func getquestion(c *gin.Context) {
	questionid := c.PostForm("questionid")
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", userName, Password, ip, port, dbName))
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	var allquestion []model.Getquestion
	rows, errq := db.Query("select questionid,username,question from question where questionid = ?", questionid)
	if errq != nil {
		log.Fatal(errq.Error)
		return
	}
	//遍历结果
	for rows.Next() {
		var u model.Getquestion
		errn := rows.Scan(&u.Questionid, &u.Username, &u.Question)
		if errn != nil {
			fmt.Printf("%v", errn)
		}

		allquestion = append(allquestion, u)
	}

	c.JSON(http.StatusOK, gin.H{"question": allquestion})

}
func deleteq(c *gin.Context) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", userName, Password, ip, port, dbName))
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	question := c.PostForm("questionid")
	questionid, abc := strconv.Atoi(question)
	if abc != nil {
		panic(abc)
	}
	flag := dao.Selectquestion(questionid)
	if !flag {
		utils.RespFail(c, "question doesn't exists")
	}
	_, err = db.Exec("delete from question where questionid=?", questionid)
	if err != nil {
		panic(err.Error())
	}
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM answer WHERE questionid=?", questionid).Scan(&count)
	if err != nil {
		panic(err.Error())
	}
	if count > 0 {
		_, err = db.Exec("delete from answer where questionid=?", questionid)
		if err != nil {
			panic(err.Error())
		}
	}
	err = cache.DeleteCache("allQuestions")
	if err != nil {
		log.Fatal(err.Error())
	}
	utils.RespSuccess(c, "successfully delete the question")
}
func deletea(c *gin.Context) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", userName, Password, ip, port, dbName))
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	answer := c.PostForm("answerid")
	answerid, abc := strconv.Atoi(answer)
	if abc != nil {
		panic(abc)
	}
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM answer WHERE answerid=?", answerid).Scan(&count)
	if err != nil {
		panic(err.Error())
	}
	if count > 0 {
		_, err = db.Exec("delete from answer where answerid=?", answerid)
		if err != nil {
			panic(err.Error())
		}
		utils.RespSuccess(c, "successfully delete the answer")
	} else {
		utils.RespFail(c, "no such answer")
	}
}
func getanswer(c *gin.Context) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", userName, Password, ip, port, dbName))
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	var allanswer []model.Getanswer
	questionidstr := c.PostForm("questionid")
	questionid, abc := strconv.Atoi(questionidstr)
	if abc != nil {
		panic(abc)
	}
	rows, errq := db.Query("select answerid,questionid,username,answer,pid from answer WHERE questionid = ?", questionid)
	if errq != nil {
		log.Fatal(errq.Error)
		return
	}
	//遍历结果
	for rows.Next() {
		var u model.Getanswer
		errn := rows.Scan(&u.Id, &u.Questionid, &u.Username, &u.Answer, &u.Pid)
		if errn != nil {
			fmt.Printf("%v", errn)
		}

		allanswer = append(allanswer, u)
	}

	c.JSON(http.StatusOK, gin.H{"res": allanswer})
}
func answer(c *gin.Context) {
	const notificationType = "给你评论了"
	if err := c.ShouldBind(&model.Getanswer{}); err != nil {
		utils.RespSuccess(c, "verification failed")
		return
	}
	answer := c.PostForm("answer")
	question := c.PostForm("questionid")
	pidstr := c.PostForm("pid")

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
	pid, abc := strconv.Atoi(pidstr)
	if abc != nil {
		panic(abc)
	}
	flag := dao.Selectquestion(questionid)
	if !flag {
		utils.RespFail(c, "question doesn't exists")
		return
	}
	if pid != 0 {
		flag1 := dao.Selectanswer(pid)
		if !flag1 {
			utils.RespFail(c, "answer doesn't exists")
		}
	}

	err := dao.Addanswer(questionid, username, answer, pid)
	if err != nil {
		utils.RespFail(c, "unable to answer")
		return
	}
	dao.Notification(question, username, notificationType)
	fmt.Println("存储通知成功")

	utils.RespSuccess(c, "successfully answered")

}
func modifya(c *gin.Context) {
	answer := c.PostForm("answer")
	answerstr := c.PostForm("answerid")
	value, exists := c.Get("username")
	if !exists {
		// 变量不存在，处理错误
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "username not found",
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
	dao.Updateanswer(username, answerid, answer)
	utils.RespSuccess(c, "successfully modify the answer")
}
func modifyq(c *gin.Context) {
	questionstr := c.PostForm("questionid")
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
	err = cache.DeleteCache("allQuestions")
	if err != nil {
		log.Fatal(err.Error())
	}

	utils.RespSuccess(c, "successfully modify the question")
}
func getuserquestion(c *gin.Context) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", userName, Password, ip, port, dbName))
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	var allquestion []model.Getquestion
	value, exists := c.Get("username")
	if !exists {
		// 变量不存在，处理错误
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "username not found",
		})
		return
	}
	rowsq, errqq := db.Query("select questionid,username,question from question where username = ?", value)
	if errqq != nil {
		log.Fatal(errqq.Error)
		return
	}
	for rowsq.Next() {
		var u model.Getquestion
		errn := rowsq.Scan(&u.Questionid, &u.Username, &u.Question)
		if errn != nil {
			fmt.Printf("%v", errn)
		}

		allquestion = append(allquestion, u)
	}
	c.JSON(http.StatusOK, gin.H{"res": allquestion})
}
func getuseranswer(c *gin.Context) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", userName, Password, ip, port, dbName))
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	var allanswer []model.Getanswer
	value, exists := c.Get("username")
	if !exists {
		// 变量不存在，处理错误
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "username not found",
		})
		return
	}
	rowsa, errqa := db.Query("select answerid,questionid,username,answer,pid from answer where username = ?", value)
	if errqa != nil {
		log.Fatal(errqa.Error)
		return
	}
	for rowsa.Next() {
		var u model.Getanswer
		errn := rowsa.Scan(&u.Id, &u.Questionid, &u.Username, &u.Answer, &u.Pid)
		if errn != nil {
			fmt.Printf("%v", errn)
		}

		allanswer = append(allanswer, u)
	}

	c.JSON(http.StatusOK, gin.H{"res": allanswer})
}
func register(c *gin.Context) {
	if err := c.ShouldBind(&model.Userre{}); err != nil {
		utils.RespFail(c, "verification failed")
		return
	}
	// 传入用户名和密码
	username := c.PostForm("username")
	password := c.PostForm("password")
	email := c.PostForm("email")
	// 验证用户名是否重复
	flag := dao.SelectUser(username)
	fmt.Println(flag)
	if flag {
		// 以 JSON 格式返回信息
		utils.RespFail(c, "user already exists")
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost) //加密处理
	if err != nil {
		fmt.Println(err)
	}
	encodePWD := string(hash)
	dao.AddUser(username, encodePWD, email)
	// 以 JSON 格式返回信息
	utils.RespSuccess(c, "add user successful")
}
func findpassword(c *gin.Context) {
	// 传入用户名和密码
	username := c.PostForm("username")
	flag := dao.SelectUser(username)
	// 不存在则退出
	if !flag {
		// 以 JSON 格式返回信息
		utils.RespFail(c, "user doesn't exists")
		return
	}
	email := dao.Selectemail(username)
	selectPassword := dao.SelectPasswordFromUsername(username)
	dao.Findpassword(email, selectPassword)
	utils.RespSuccess(c, fmt.Sprintf("successfully post the password to the email"))
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
	err := bcrypt.CompareHashAndPassword([]byte(selectPassword), []byte(password)) //验证（对比）
	if err != nil {
		utils.RespFail(c, "wrong password")
		return
	}
	c.SetCookie("gin_demo_cookie", "test", 3600, "/", "localhost", false, true)
	//utils.RespSuccess(c, "login successful")
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
	//utils.RespSuccess(c, tokenString)
	c.JSON(http.StatusOK, gin.H{
		"status":  "200",
		"message": "登录成功",
		"token":   tokenString,
	})
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
	err := bcrypt.CompareHashAndPassword([]byte(selectPassword), []byte(oldpassword)) //验证（对比）
	if err != nil {
		utils.RespFail(c, "wrong oldpassword")
		return
	}
	// 更新密码
	hash, err2 := bcrypt.GenerateFromPassword([]byte(newpassword), bcrypt.DefaultCost) //加密处理
	if err2 != nil {
		fmt.Println(err)
	}
	encodePWD := string(hash)
	err1 := dao.UpdatePassword(username, encodePWD, oldpassword)
	if err1 != nil {
		utils.RespFail(c, "unable to change password")
		return
	}
	// 成功更新密码，返回成功响应
	utils.RespSuccess(c, "password changed")
}

// 新增以下代码
func getUsernameFromToken(c *gin.Context) {
	username, _ := c.Get("username")
	c.JSON(http.StatusOK, gin.H{
		"status":   "200",
		"username": username,
	})
}
