package main

import (
	"database/sql"
	"fmt"

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
	result, err := db.Exec("SELECT * from user where id=?", 1)

	fmt.Println("db ===> ", result)

	fmt.Println("&result ===> ", &result)
	fmt.Println("*result ===> ", *&result)

	// res := &result
	// fmt.Printf("ptr = %p\n", &result)
	var p = &result
	fmt.Printf("5. main  -- p %T: &p=%p p=&i=%p  *p=i=%v\n", p, &p, p, *p)

	checkErr(err)
}
