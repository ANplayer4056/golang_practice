package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/fatih/structs"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

/**
* gin :
* 	golang 所開發出來的 web framework
*	用來 建立 API
*
*
* gorm :
* 	是Golang語言中一款效能極好的ORM庫
*	用來與 db 溝通及處理資料表等使用
*
 */

// main ===> main.go 進入點
func main() {

	// defer 不管有無錯誤都會 run
	defer func() {
		if err := recover(); err != nil {
			// 補上將err傳至telegram
			fmt.Println("[❌ Fatal❌ ] HTTP:", err)
		}
	}()

	if err := DBcheckTable(); err != nil {
		return
	}

	// 使用 gin 來製作 CRUD API
	r := gin.Default()
	r.POST("/createUser", CreateUser)
	r.DELETE("/deleteUser", DeleteUser)
	r.PUT("/updateUser", UpdateUser)
	r.POST("/queryUser", QueryUser)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

/*
UserLists ===> Binding from JSON (keyword: bind json)
               定義 接收傳送過來的數據(params) / json:"API的欄位名稱" /  binding:"required" 必要欄位
*/
type UserLists struct {
	User     string `form:"user" json:"user" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

/*
DBcheckTable ===> deal with db AutoMigrate
                  處理 AutoMigrate
*/
func DBcheckTable() error {

	type UserList struct {
		ID        int    `gorm:"priamrykey"`
		Username  string `gorm:"column:username"`
		Password  string `gorm:"column:password"`
		CreatedAt time.Time
		UpdatedAt time.Time
		DeletedAt gorm.DeletedAt `gorm:"index"`
	}

	db, err := connectDB()
	if err != nil {
		fmt.Println("DB connect failed ===> ", err)
	}

	if err = db.AutoMigrate(&UserList{}); err != nil {
		fmt.Println("DB Migrate failed ===> ", err)
	}

	return err
}

// CreateUser ===>  新增使用者的 api
func CreateUser(c *gin.Context) {

	// get json data
	var json UserLists
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// new db connect
	db, err := connectDB()
	if err != nil {
		fmt.Println("DB connect failed ===> ", err)
	}

	type UserList struct {
		ID       int    `gorm:"priamrykey"`
		Username string `gorm:"column:username"`
		Password string `gorm:"column:password"`
	}

	if err = db.Model(&UserLists{}).Create(map[string]interface{}{
		"username": json.User, "password": json.Password,
	}).Error; err != nil {

		c.JSON(200, gin.H{
			"statusCode": 1001,
			"message":    "create faile",
		})
		return
	}

	c.JSON(200, gin.H{
		"statusCode": 200,
		"userName":   json.User,
		"Password":   json.Password,
	})

}

// DeleteUser ===>  刪除使用者的 api
func DeleteUser(c *gin.Context) {

	type ReqUser struct {
		ID int `form:"id" json:"id" binding:"required"`
	}

	// get json data
	req := ReqUser{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// new db connect
	db, err := connectDB()
	if err != nil {
		fmt.Println("DB connect failed ===> ", err)
	}

	type UserList struct {
		ID       int    `gorm:"priamrykey"`
		Username string `gorm:"column:username"`
		Password string `gorm:"column:password"`
	}

	dbUser := UserList{
		ID: req.ID,
	}

	if err = db.Delete(&dbUser).Error; err != nil {
		log.Printf("Error Message is %v ", err.Error())
		c.JSON(200, gin.H{
			"statusCode": 1001,
			"message":    "delete faile",
		})
		return

	}
	c.JSON(200, gin.H{
		"statusCode": 200,
		"message":    "delete success",
	})

}

// UpdateUser ===>  更新單一使用者的 api
func UpdateUser(c *gin.Context) {

	type ReqUser struct {
		Username string `form:"Username" json:"Username" binding:"required"`
		Password string `form:"Password" json:"Password"`
	}

	// get json data
	reqParmams := ReqUser{}
	if err := c.ShouldBindJSON(&reqParmams); err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}

	// new db connect
	db, err := connectDB()
	if err != nil {
		fmt.Println("DB connect failed ===> ", err)
	}

	type UserList struct {
		Password string `gorm:"column:password"`
	}

	strDB := &UserList{
		Password: reqParmams.Password,
	}
	fmt.Println("strDB ===>", strDB)

	mapDB := structs.Map(strDB)
	fmt.Println("===>", mapDB)

	if err = db.Model(&strDB).Where("Username = ?", reqParmams.Username).Updates(&mapDB).Error; err != nil {
		log.Printf("Error Message is %v ", err.Error())

		c.JSON(200, gin.H{
			"statusCode": 1001,
			"message":    "update Error",
		})
		return
	}

	c.JSON(200, gin.H{
		"statusCode": 200,
		"message":    "update success",
	})

}

// QueryUser ===>  查詢指定用戶的 api
func QueryUser(c *gin.Context) {

	type ReqUser struct {
		Username string `form:"Username" json:"Username" binding:"required"`
	}

	// get json data
	reqParmams := ReqUser{}
	if err := c.ShouldBindJSON(&reqParmams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// new db connect
	db, err := connectDB()
	if err != nil {
		fmt.Println("DB connect failed ===> ", err)
	}

	type UserList struct {
		ID       int    `gorm:"priamrykey"`
		Username string `gorm:"column:username"`
		Password string `gorm:"column:password"`
	}

	dbaccept := []UserList{}

	result := db.Where("Username = ?", reqParmams.Username).Find(&dbaccept)
	fmt.Println(dbaccept)

	if result.Error != nil {
		log.Printf("Error Message is %v ", result.Error)
		c.JSON(200, gin.H{
			"statusCode": 1001,
			"message":    "query faile",
		})
		return
	}

	c.JSON(200, gin.H{
		"statusCode": 200,
		"message":    dbaccept,
	})

}

// connectDB ===> 建立 db 的連線
func connectDB() (*gorm.DB, error) {
	dsn := "root:example@tcp(127.0.0.1:3306)/backend_user?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	tx := db.Debug()

	return tx, err
}

// db.create(&xxx) delet(&xxx) update(&xxx) ... &xxx ===> into func to do somethings
// db.find(&xxx)  &xxx ===> func result set value into &xxx
