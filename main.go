package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {

	r := gin.Default()
	r.POST("/createUser", CreateUser)
	r.POST("/deleteUser", DeleteUser)
	r.POST("/updateUser", UpdateUser)
	r.POST("/queryUser", QueryUser)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

// UserLists ===> Binding from JSON (keyword: bind json)
type UserLists struct {
	User     string `form:"user" json:"user" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

// DBcheckTable ===>  deal with db AutoMigrate
func DBcheckTable() {

}

// CreateUser ===>  create a new user api
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

	if err = db.AutoMigrate(&UserList{}); err != nil {
		fmt.Println("DB Migrate failed ===> ", err)
	}

	// result := db.Create()
	if err = db.Model(&UserLists{}).Create(map[string]interface{}{
		"username": json.User, "password": json.Password,
	}).Error; err != nil {

		c.JSON(200, gin.H{
			"statusCode": 1001,
			"message":    "create faile",
		})

	} else {

		c.JSON(200, gin.H{
			"statusCode": 200,
			"userName":   json.User,
			"Password":   json.Password,
		})
	}

}

// DeleteUser ===>  delete selected user api
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

	if err = db.AutoMigrate(&UserList{}); err != nil {
		fmt.Println("DB Migrate failed ===> ", err)
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

	} else {
		c.JSON(200, gin.H{
			"statusCode": 200,
			"message":    "delete success",
		})
	}

}

// UpdateUser ===>  update single user password api
func UpdateUser(c *gin.Context) {

	type ReqUser struct {
		Username string `form:"Username" json:"Username" binding:"required"`
		Password string `form:"Password" json:"Password" binding:"required"`
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

	if err = db.AutoMigrate(&UserList{}); err != nil {
		fmt.Println("DB Migrate failed ===> ", err)
	}

	dbUpdate := UserList{
		Password: reqParmams.Password,
	}

	if err = db.Model(&dbUpdate).Where("Username = ?", reqParmams.Username).Updates(&dbUpdate).Error; err != nil {
		log.Printf("Error Message is %v ", err.Error())

		c.JSON(200, gin.H{
			"statusCode": 1001,
			"message":    "update Error",
		})

	} else {
		c.JSON(200, gin.H{
			"statusCode": 200,
			"message":    "update success",
		})
	}

}

// QueryUser ===>  select single user informations api
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

	if err = db.AutoMigrate(&UserList{}); err != nil {
		fmt.Println("DB Migrate failed ===> ", err)
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
	} else {
		c.JSON(200, gin.H{
			"statusCode": 200,
			"message":    dbaccept,
		})
	}
}

func connectDB() (*gorm.DB, error) {
	dsn := "root:example@tcp(127.0.0.1:3306)/backend_user?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	tx := db.Debug()

	return tx, err
}

// db.create(&xxx) delet(&xxx) update(&xxx) ... &xxx ===> into func to do somethings
// db.find(&xxx)  &xxx ===> func result set value into &xxx
