package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {

	r := gin.Default()
	r.POST("/createUser", CreateUser)
	r.POST("/deleteUser", DeleteUser)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

// UserLists ===> Binding from JSON (keyword: bind json)
type UserLists struct {
	User     string `form:"user" json:"user" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

// CreateUser ===>  create user api
func CreateUser(c *gin.Context) {

	// new db connect
	db, err := connectDB()
	if err != nil {
		fmt.Println("DB connect failed ===> ", err)
	}

	type UserList struct {
		gorm.Model
		ID       int    `gorm:"priamrykey"`
		Username string `gorm:"column:username"`
		Password string `gorm:"column:password"`
	}

	if err = db.AutoMigrate(&UserList{}); err != nil {
		fmt.Println("DB Migrate failed ===> ", err)
	}

	// get json data
	var json UserLists
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	memberName := json.User
	memberPassword := json.Password

	// result := db.Create()
	db.Model(&UserLists{}).Create(map[string]interface{}{
		"username": memberName, "password": memberPassword,
	})

	c.JSON(200, gin.H{
		"userName": memberName,
		"Password": memberPassword,
	})
}

// DeleteUser ===>  delete user api
func DeleteUser(c *gin.Context) {
	// new db connect
	db, err := connectDB()
	if err != nil {
		fmt.Println("DB connect failed ===> ", err)
	}

	type UserList struct {
		gorm.Model
		ID       int    `gorm:"priamrykey"`
		Username string `gorm:"column:username"`
		Password string `gorm:"column:password"`
	}

	if err = db.AutoMigrate(&UserList{}); err != nil {
		fmt.Println("DB Migrate failed ===> ", err)
	}

	// get json data
	var json UserList
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db.Where("ID = ?", json.ID).Find(&json)
	db.Delete(&json)

	c.JSON(200, gin.H{
		"message": "delete success",
	})
}

func connectDB() (*gorm.DB, error) {
	dsn := "root:example@tcp(127.0.0.1:3306)/backend_user?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	return db, err
}
