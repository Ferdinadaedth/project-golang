package dao

//
//import (
//	"bufio"
//	"encoding/csv"
//	"fmt"
//	"io"
//	"os"
//	"strings"
//)
//
//// 文件路径
//const path = "D:\\Edc\\golandprojects\\go1.20.3\\gin-demo\\dao\\data\\users.csv"
//
////从文件中加载所有用户数据
//
//func loadUsers() (map[string]string, error) {
//	file, err := os.Open(path)
//	if err != nil {
//		return nil, err
//	}
//	defer file.Close()
//
//	csvReader := csv.NewReader(file)
//	csvReader.FieldsPerRecord = 2
//	records, err := csvReader.ReadAll()
//	if err != nil {
//		return nil, err
//	}
//
//	users := make(map[string]string)
//	for _, record := range records {
//		users[record[0]] = record[1]
//	}
//
//	return users, nil
//}
//
//func Loadpassword(file string) (map[string]string, error) {
//	// 定义用于存储密码信息的映射
//	password := make(map[string]string)
//
//	// 打开密码文件
//	f, err := os.Open(path)
//	if err != nil {
//		return password, err
//	}
//	defer f.Close()
//
//	// 创建一个 bufio.Scanner 对象来逐行读取文件内容
//	scanner := bufio.NewScanner(f)
//	scanner.Split(bufio.ScanLines)
//
//	// 逐行读取文件内容，将每个用户名和密码添加到映射中
//	for scanner.Scan() {
//		line := scanner.Text()
//		fields := strings.Split(line, ",")
//		if len(fields) < 2 {
//			continue // 遇到不合法的行跳过
//		}
//		password[fields[0]] = fields[1]
//	}
//
//	if err := scanner.Err(); err != nil {
//		return password, err
//	}
//
//	return password, nil
//}
//
//// 保存所有用户数据到文件
//func saveUsers(users map[string]string) error {
//	// 创建目录
//	err := os.MkdirAll("./data", 0755)
//	if err != nil {
//		return err
//	}
//
//	// 创建或打开文件
//	file, err := os.Create(path)
//	if err != nil {
//		return err
//	}
//	defer file.Close()
//
//	csvWriter := csv.NewWriter(file)
//	for username, password := range users {
//		err := csvWriter.Write([]string{username, password})
//		if err != nil {
//			return err
//		}
//	}
//	csvWriter.Flush()
//
//	return nil
//}
//func UpdatePassword(username string, newPassword string) error {
//	// 打开用户数据文件，以读写方式打开
//	f, err := os.OpenFile(path, os.O_RDWR, 0666)
//	if err != nil {
//		return err
//	}
//	defer f.Close()
//
//	// 构造一个 *bufio.Reader 对象，以便读取文件中的数据
//	reader := bufio.NewReader(f)
//
//	// 读取文件中的每一行数据
//	for {
//		// 在文件中查找指定的用户名
//		line, err := reader.ReadString('\n')
//		if err == io.EOF {
//			break // 已到文件末尾
//		} else if err != nil {
//			return err
//		}
//		fields := strings.Split(line, ",") // 将用户名和密码分离
//		if strings.TrimSpace(fields[0]) == username {
//			// 找到指定的用户，更新密码字段并写回文件
//			fields[1] = newPassword
//			_, err = f.Seek(-int64(len(line)), io.SeekCurrent)
//			if err != nil {
//				return err
//			}
//			_, err = f.WriteString(strings.Join(fields, ","))
//			if err != nil {
//				return err
//			}
//			return nil // 密码更新成功，返回 nil
//		}
//	}
//
//	return fmt.Errorf("user not found: %s", username)
//}
