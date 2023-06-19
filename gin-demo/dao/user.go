package dao

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

// 数据库连接信息
const (
	userName = "root"
	Password = "123456"
	ip       = "ten.ferdinandaedth.top"
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
