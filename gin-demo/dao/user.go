//package dao

/*import (
	"database/sql"
	"fmt"
)*/

/*const path1 = "D:\\Edc\\golandprojects\\go1.20.3\\gin-demo\\dao\\data\\message.csv"

var (
	mutex    sync.Mutex
	database map[string]string
)

/*func init() {
	users, err := loadUsers()
	if err != nil {
		panic(err)
	}
	database = users
}*/
/*
func AddUser(username, password string) error {
	mutex.Lock()
	defer mutex.Unlock()
	database[username] = password
	err := saveUsers(database)
	if err != nil {
		return err
	}

	return nil
}

// 若没有这个用户返回 false，反之返回 true
func SelectUser(username string) bool {
	mutex.Lock()
	defer mutex.Unlock()

	if database[username] == "" {
		return false
	}
	return true
}

func SelectPasswordFromUsername(username string) string {
	mutex.Lock()
	defer mutex.Unlock()
	database, _ = Loadpassword(username)
	return database[username]
}
func AddMessage(message model.Message) {
	file, err := os.OpenFile(path1, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if err := writer.Write([]string{message.Id, message.Text, message.CreatedAt.Format("2006-01-02 15:04:05")}); err != nil {
		log.Fatal(err)
	}
}
*/
package dao

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

// 数据库连接信息
const (
	userName = "root"
	Password = "yx041110"
	ip       = "127.0.0.1"
	port     = "3306"
	dbName   = "userdb"
)

func Addquestion(username, description string) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", userName, Password, ip, port, dbName))
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	_, err = db.Exec("INSERT INTO question (question,username) VALUES (?,?)", description, username)
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
	err = db.QueryRow("SELECT COUNT(*) FROM question WHERE id=?", questionid).Scan(&count)
	if err != nil {
		panic(err.Error())
	}
	if count > 0 {
		return true
	} else {
		return false
	}
}
func Addanswer(questionid int, username, description string) error {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", userName, Password, ip, port, dbName))
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	_, err = db.Exec("INSERT INTO answer (questionid,username,answer) VALUES (?, ?,?)", questionid, username, description)
	if err != nil {
		panic(err.Error())
	}
	return nil
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
func AddUser(username, password string) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", userName, Password, ip, port, dbName))
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	// 插入用户记录
	_, err = db.Exec("INSERT INTO user (username, password) VALUES (?, ?)", username, password)
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
	stmt, err := db.Prepare("UPDATE `question` SET `question` =? WHERE `username` = ? AND `id` = ?")
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
func Updateanswer(username string, questionid int, id int, answer string) error {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", userName, Password, ip, port, dbName))
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	stmt, err := db.Prepare("UPDATE `answer` SET `answer` =? WHERE `username` = ? AND `id` = ? AND `questionid` = ?")
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()
	// 执行更新操作
	_, err = stmt.Exec(answer, username, id, questionid)
	if err != nil {
		return err
	}

	return nil
}
