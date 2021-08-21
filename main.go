package main

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// return sum of a and b
	r.GET("/getSum", func(c *gin.Context) {
		StringA := c.Query("NumA")
		StringB := c.Query("NumB")

		IntA, err := strconv.Atoi(StringA)
		IntB, err := strconv.Atoi(StringB)

		fmt.Println("====>", IntA, IntB)

		if err != nil {
			fmt.Println("====>", err)
			c.JSON(200, gin.H{
				"message": "NumA error",
			})
			return
		}

		if err != nil {
			fmt.Println("====>", err)
			c.JSON(200, gin.H{
				"message": "NumB error",
			})
			return
		}

		c.JSON(200, gin.H{
			"sum": IntA + IntB,
		})
	})

	// return typeOf value
	r.GET("/getType", func(c *gin.Context) {
		ValueA := c.Query("value")

		_, err := strconv.ParseBool(ValueA)
		if err != nil {

			_, err := strconv.Atoi(ValueA)
			if err != nil {
				c.JSON(200, gin.H{
					"message": "type is string",
				})
				return
			}
			c.JSON(200, gin.H{
				"message": "type is int",
			})
			return
		}

		c.JSON(200, gin.H{
			"message": "type is boolean",
		})

	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

// 1. gin framework >> http serve
// 2. go mod init >> created go.mod
