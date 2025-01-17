package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// User 结构体定义用户模型
type User struct {
	ID    uint   `gorm:"primaryKey"`
	Name  string `gorm:"not null"`
	Email string `gorm:"uniqueIndex"`
	Age   int
}

// CreateUser 创建用户
func CreateUser(db *gorm.DB, user User) {
	result := db.Create(&user)
	if result.Error != nil {
		fmt.Printf("创建用户失败: %v\n", result.Error)
	}
}

// ReadUser 读取用户
func ReadUser(db *gorm.DB, id uint) (User, error) {
	var user User
	result := db.First(&user, id)
	if result.Error != nil {
		return user, result.Error
	}
	return user, nil
}

// UpdateUser 更新用户
func UpdateUser(db *gorm.DB, id uint, newUser User) {
	result := db.Model(&User{}).Where("id =?", id).Updates(newUser)
	if result.Error != nil {
		fmt.Printf("更新用户失败: %v\n", result.Error)
	}
}

// DeleteUser 删除用户
func DeleteUser(db *gorm.DB, id uint) {
	result := db.Delete(&User{}, id)
	if result.Error != nil {
		fmt.Printf("删除用户失败: %v\n", result.Error)
	}
}

// GetAllUsers 获取所有用户
func GetAllUsers(db *gorm.DB) ([]User, error) {
	var users []User
	result := db.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

func main() {
	// 数据库连接字符串
	dsn := "user:password@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
	// 连接数据库
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("连接数据库失败: " + err.Error())
	}

	// 自动迁移模式，创建表
	db.AutoMigrate(&User{})

	router := gin.Default()

	// 创建用户接口
	router.POST("/users", func(c *gin.Context) {
		var newUser User
		if err := c.ShouldBindJSON(&newUser); err != nil {
			c.JSON(400, gin.H{"error": "请求参数错误"})
			return
		}
		CreateUser(db, newUser)
		c.JSON(201, gin.H{"message": "用户创建成功", "user": newUser})
	})

	// 获取单个用户接口
	router.GET("/users/:id", func(c *gin.Context) {
		id := c.Param("id")
		var userID uint
		_, err := fmt.Sscanf(id, "%d", &userID)
		if err != nil {
			c.JSON(400, gin.H{"error": "ID参数格式错误"})
			return
		}
		user, err := ReadUser(db, userID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(404, gin.H{"error": "用户未找到"})
			} else {
				c.JSON(500, gin.H{"error": "获取用户失败"})
			}
			return
		}
		c.JSON(200, gin.H{"user": user})
	})

	// 更新用户接口
	router.PUT("/users/:id", func(c *gin.Context) {
		id := c.Param("id")
		var userID uint
		_, err := fmt.Sscanf(id, "%d", &userID)
		if err != nil {
			c.JSON(400, gin.H{"error": "ID参数格式错误"})
			return
		}
		var updatedUser User
		if err := c.ShouldBindJSON(&updatedUser); err != nil {
			c.JSON(400, gin.H{"error": "请求参数错误"})
			return
		}
		UpdateUser(db, userID, updatedUser)
		c.JSON(200, gin.H{"message": "用户更新成功"})
	})

	// 删除用户接口
	router.DELETE("/users/:id", func(c *gin.Context) {
		id := c.Param("id")
		var userID uint
		_, err := fmt.Sscanf(id, "%d", &userID)
		if err != nil {
			c.JSON(400, gin.H{"error": "ID参数格式错误"})
			return
		}
		DeleteUser(db, userID)
		c.JSON(200, gin.H{"message": "用户删除成功"})
	})

	// 获取所有用户接口
	router.GET("/users", func(c *gin.Context) {
		users, err := GetAllUsers(db)
		if err != nil {
			c.JSON(500, gin.H{"error": "获取用户列表失败"})
			return
		}
		c.JSON(200, gin.H{"users": users})
	})

	router.Run(":8080")
}
