package dao

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"golandprojects/DFA"
	"golandprojects/model"
	"gopkg.in/gomail.v2"
	"log"
	"net/http"
	"time"
)

// 数据库连接信息
const (
	userName = "root"
	Password = "h74o+JIi5SpSY3MU"
	ip       = "47.108.208.111"
	port     = "3306"
	dbName   = "userdb"
)

func Notification(questionID, currentUsername, notificationType string) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", userName, Password, ip, port, dbName))
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	str1 := "SELECT username from question where questionid=?"
	row := db.QueryRow(str1, questionID)
	var username string
	if err = row.Scan(&username); err != nil {
		panic(err.Error())
	}

	str2 := "SELECT userid from user where username=?"
	row = db.QueryRow(str2, username)
	var recipientUserID string
	if err = row.Scan(&recipientUserID); err != nil {
		panic(err.Error())
	}

	str3 := "SELECT userid from user where username=?"
	row = db.QueryRow(str3, currentUsername)
	var senderUserID string
	if err = row.Scan(&senderUserID); err != nil {
		panic(err.Error())
	}
	now := time.Now()
	notificationTime := now.Format("2006-01-02 15:04:05")

	str4 := "INSERT INTO notification (recipientUserID,senderUserID,notificationType,notificationTime) VALUES (?,?,?,?)"
	_, err = db.Exec(str4, recipientUserID, senderUserID, notificationType, notificationTime)
}

func GetAllQuestions() ([]model.Getquestion, error) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", userName, Password, ip, port, dbName))
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	var allquestion []model.Getquestion
	rows, errq := db.Query("select questionid,username,question,label from question")
	if errq != nil {
		log.Fatal(errq.Error)
		return allquestion, err
	}
	//遍历结果
	for rows.Next() {
		var u model.Getquestion
		errn := rows.Scan(&u.Questionid, &u.Username, &u.Question, &u.Label)
		if errn != nil {
			fmt.Printf("%v", errn)
		}

		allquestion = append(allquestion, u)
	}

	if err = rows.Err(); err != nil {
		return allquestion, err
	}

	return allquestion, nil
}

func Addquestion(username, description string) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", userName, Password, ip, port, dbName))
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	var ques string
	ques = DFA.Filter(description)
	fmt.Println(ques)
	_, err = db.Exec("INSERT INTO question (question,username) VALUES (?,?)", ques, username)
	if err != nil {
		panic(err.Error())
	}
	data := map[string]interface{}{
		"model":    "gpt-3.5-turbo",
		"messages": []map[string]string{{"role": "user", "content": description}, {"role": "user", "content": "，请给前面这句话的内容总结为一个中文标签，只返回这个词语不要有其他内容,不超过20个词"}},
	}
	payload, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}

	// 发送 POST 请求
	url := "https://api.aiproxy.io/v1/chat/completions"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer sk-7YXV6hBPcbWhXS0z3vGTZk3t5PvRpQ1e7XipJHqWRO9ZwUbc") // 替换为您的 API KEY

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	// 处理响应
	var result map[string]interface{}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&result)
	if err != nil {
		fmt.Println("Error decoding response:", err)
		return
	}

	fmt.Println(result) // 打印响应结果
	jsonData, err := json.Marshal(result)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}

	// 解码 result 中的 JSON 数据
	var decodedResult map[string]interface{}
	err = json.Unmarshal(jsonData, &decodedResult)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}

	// 提取并设置content到变量中
	content := decodedResult["choices"].([]interface{})[0].(map[string]interface{})["message"].(map[string]interface{})["content"].(string)
	fmt.Println(content)
	_, err = db.Exec("UPDATE question SET label = ? WHERE question = ? AND username = ?", content, ques, username)
	if err != nil {
		panic(err.Error())
	}

}
func Selectquestion(questionid int) bool {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", userName, Password, ip, port, dbName))
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM question WHERE questionid=?", questionid).Scan(&count)
	if err != nil {
		panic(err.Error())
	}
	if count > 0 {
		return true
	} else {
		return false
	}
}
func Selectanswer(pid int) bool {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", userName, Password, ip, port, dbName))
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM answer WHERE answerid=?", pid).Scan(&count)
	if err != nil {
		panic(err.Error())
	}
	if count > 0 {
		return true
	} else {
		return false
	}
}
func Addanswer(questionid int, username, description string, pid int) error {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", userName, Password, ip, port, dbName))
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	var ques string
	ques = DFA.Filter(description)
	_, err = db.Exec("INSERT INTO answer (questionid,username,answer,pid) VALUES (?, ?,?,?)", questionid, username, ques, pid)
	if err != nil {
		panic(err.Error())
	}
	return nil
}
func Findpassword(email string, password string) {
	m := gomail.NewMessage()

	//发送人
	m.SetHeader("From", "2794954964@qq.com")
	//接收人
	m.SetHeader("To", email)
	//主题
	m.SetHeader("Subject", "找回密码")
	//内容
	m.SetBody("text", password)
	//拿到token，并进行连接,第4个参数是填授权码
	d := gomail.NewDialer("smtp.qq.com", 587, "2794954964@qq.com", "zkfsgcjtbapiddha")

	// 发送邮件
	if err := d.DialAndSend(m); err != nil {
		fmt.Printf("DialAndSend err %v:", err)
		panic(err)
	}
	fmt.Printf("send mail success\n")
}
func Selectemail(username string) string {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", userName, Password, ip, port, dbName))
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	// 查询用户名是否存在
	var email string
	err = db.QueryRow("SELECT email FROM user WHERE username=?", username).Scan(&email)
	if err != nil {
		panic(err.Error())
	}
	return email

}

// SelectUser 根据用户名查询用户是否存在
func SelectUser(username string) bool {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", userName, Password, ip, port, dbName))
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	// 查询用户名是否存在
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM user WHERE username=?", username).Scan(&count)
	if err != nil {
		panic(err.Error())
	}
	if count > 0 {
		return true
	} else {
		return false
	}
}

// AddUser 添加用户
func AddUser(username, password, email string) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", userName, Password, ip, port, dbName))
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	var ques string
	ques = DFA.Filter(username)
	// 插入用户记录
	_, err = db.Exec("INSERT INTO user (username, password,email) VALUES (?, ?,?)", ques, password, email)
	if err != nil {
		panic(err.Error())
	}
}
func Addmessage(reveiver, sender, message string) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", userName, Password, ip, port, dbName))
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	// 插入用户记录
	_, err = db.Exec("INSERT INTO message (sender, receiver,content) VALUES (?, ?,?)", sender, reveiver, message)
	if err != nil {
		panic(err.Error())
	}
}

// SelectPasswordFromUsername 根据用户名查询密码
func SelectPasswordFromUsername(username string) string {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", userName, Password, ip, port, dbName))
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	// 查询密码
	var password string
	err = db.QueryRow("SELECT password FROM user WHERE username=?", username).Scan(&password)
	if err != nil {
		panic(err.Error())
	}
	return password
}
func UpdatePassword(username, newPassword, oldpassword string) error {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", userName, Password, ip, port, dbName))
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	stmt, err := db.Prepare("UPDATE `user` SET `password` = ? WHERE `username` = ?")
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()
	// 执行更新操作
	_, err = stmt.Exec(newPassword, username)
	if err != nil {
		return err
	}

	return nil
}
func UpdateQuestin(username string, questionid int, question string) error {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", userName, Password, ip, port, dbName))
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	stmt, err := db.Prepare("UPDATE `question` SET `question` =? WHERE `username` = ? AND `questionid` = ?")
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()
	// 执行更新操作
	_, err = stmt.Exec(question, username, questionid)
	if err != nil {
		return err
	}

	return nil
}
func Updateanswer(username string, id int, answer string) error {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", userName, Password, ip, port, dbName))
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	stmt, err := db.Prepare("UPDATE `answer` SET `answer` =? WHERE `username` = ? AND `answerid` = ? ")
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()
	// 执行更新操作
	_, err = stmt.Exec(answer, username, id)
	if err != nil {
		return err
	}

	return nil
}
