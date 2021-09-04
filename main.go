package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "root:example@tcp(127.0.0.1:3306)/backend")

	if err := db.Ping(); err != nil {
		fmt.Println("資料庫連線錯誤，原因為：", err.Error())
		return
	}
	checkErr(err)
	queryUser(db)

	db.Close()
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func createUser(db *sql.DB) {
	// create 資料
	_, err := db.Exec("insert INTO user(username,password) values(?,?)", "test", "00000000")
	checkErr(err)
}

func updateUser(db *sql.DB) {
	// update 資料
	_, err := db.Exec("UPDATE user set password=? where username=?", "teatabc123", "test")
	checkErr(err)
}

func deleteUser(db *sql.DB) {
	// delete 資料
	_, err := db.Exec("delete from user where username=?", "test")
	checkErr(err)
}

func queryUser(db *sql.DB) {
	// select 資料
	result, err := db.Query("SELECT * from user")

	for result.Next() {

		var (
			id       int
			username string
			password string
		)

		if err := result.Scan(&id, &username, &password); err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%v,%v,%v\n", id, username, password)
	}

	checkErr(err)
}
