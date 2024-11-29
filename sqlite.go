package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3" // 导入SQLite驱动
)

var db *sql.DB

func Init() {
	DBFile := getDBFile()
	var err error
	db, err = sql.Open("sqlite3", DBFile)
	if err != nil {
		log.Fatal(err)
	}
	// 检查表是否存在，不存在就建表
	if !checkTableExist("account") {
		// 密码直接存明文了，你少管我那么多
		createTableSQL := `
		CREATE TABLE account (
			id TEXT NOT NULL PRIMARY KEY,
			password TEXT NOT NULL,
			home TEXT NOT NULL,
			writable INTEGER NOT NULL
		);
		`
		_, err := db.Exec(createTableSQL)
		if err != nil {
			log.Fatalf("创建数据表account失败: %v", err)
		}
	}
	if !checkTableExist("config") {
		createTableSQL := `
		CREATE TABLE config (
			id INTEGER NOT NULL PRIMARY KEY,
			sftp INTEGER NOT NULL,
			ftp INTEGER NOT NULL,
			webdav INTEGER NOT NULL,
		);
		`
		_, err := db.Exec(createTableSQL)
		if err != nil {
			log.Fatalf("创建数据表config失败: %v", err)
			return
		}
		insertSQL := "INSERT INTO config (id, sftp, ftp, webdav) VALUES (?, ?, ?, ?);"
		_, err = db.Exec(insertSQL, 1, 2121, 2222, 8080)
		if err != nil {
			log.Fatalf("初始化配置数据失败: %v", err)
		}
	}
}

func getDBFile() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("获取不到用户主目录: %v", err)
	}
	DBFilePath := filepath.Join(homeDir, "wavingleaf.db")

	if _, err := os.Stat(DBFilePath); os.IsNotExist(err) {
		// 文件不存在，创建文件
		err := os.WriteFile(DBFilePath, []byte{}, 0644)
		if err != nil {
			log.Fatalf("创建用户数据库失败: %v", err)
		}
	} else if err != nil {
		// 其他错误
		log.Fatalf("Error checking config file: %v", err)
	}
	return DBFilePath
}

func checkTableExist(tableName string) bool {
	query := fmt.Sprintf(`SELECT name FROM sqlite_master WHERE type='table' AND name='%s';`, tableName)
	rows, err := db.Query(query)
	if err != nil {
		return false
	}
	defer rows.Close()
	exists := rows.Next()
	return exists
}

func CloseDB() {
	db.Close()
}

func AddUser(id string, password string, home string, writable int) bool {
	insertSQL := "INSERT INTO account (id, password, home, writable) VALUES (?, ?, ?, ?);"
	_, err := db.Exec(insertSQL, id, password, home, writable)
	if err != nil {
		log.Fatalf("初始化配置数据失败: %v", err)
		return false
	}
	return true
}

func DeleteUser(id string) bool {
	deleteSQL := "DELETE FROM account WHERE id = ?;"
	_, err := db.Exec(deleteSQL, 12345)
	if err != nil {
		log.Fatalf("删除用户失败: %v", err)
		return false
	}
	return true
}

// id作为主键，不可更新，此处仅用于查找。只能更新用户的密码，主目录和权限
func UpdateUser(id string, password string, home string, writable int) bool {
	updateSQL := "UPDATE account SET password = ?, home = ?, writable = ? WHERE id = ?"
	_, err := db.Exec(updateSQL, password, home, writable, id)
	if err != nil {
		log.Fatalf("更新用户信息失败: %v", err)
		return false
	}
	return true
}

// 只返回用户ID，主目录和权限，密码不返回
func ListUser() {

}

func GetUser(id string) {

}

// 验证登录的
func CheckUser(id string, password string) {

}

// 用枚举实现一个查对应服务的端口
func GetPort() int {

}

// 根据枚举，更新对应服务的端口
func UpadtePort(enum, int) bool {

}
