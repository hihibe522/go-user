package user

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"go-app/db"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
)

var mysql = db.GetDB()

// 用戶
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Alias    string `json:"alias"`
}

// 生成 session key
func generateSessionKey(userID int, username string) (string, error) {
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}

	hasher := sha256.New()
	hasher.Write([]byte(fmt.Sprintf("%d:%s:%s", userID, username, base64.URLEncoding.EncodeToString(salt))))

	sessionKey := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	return sessionKey, nil
}

// 驗證用戶輸入的資料
func validate(user *User) (bool, string) {
	// 內容必須至少一個英文和至少一個數字，長度不超過25
	pattern := "^[a-zA-Z0-9]{1,25}$"

	match, err := regexp.MatchString(pattern, user.Username)
	if err != nil || !match {
		return false, "Invalid username format. It should be alphanumeric and less than 25 characters"
	}

	match, err = regexp.MatchString(pattern, user.Password)
	if err != nil || !match {
		return false, "Invalid password format. It should be alphanumeric and less than 25 characters"
	}
	return true, ""
}

// 建立用戶
func CreateUser(c *gin.Context) {
	var user User
	err := c.BindJSON(&user)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// 驗證輸入的資料
	verityResult, errString := validate(&user)
	if !verityResult {
		c.JSON(http.StatusBadRequest, gin.H{"error": errString})
		return
	}

	// 檢查username是否重複
	var existingUser User
	if err := mysql.Where("username = ?", user.Username).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		return
	}

	result := mysql.Create(&user)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// 登入
func Login(c *gin.Context) {
	var user User
	err := c.BindJSON(&user)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	result := mysql.Where("username = ? AND password = ?", user.Username, user.Password).First(&user)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to login"})
		return
	}
	// 用取出的資料生成 session key
	sessionKey, _ := generateSessionKey(user.ID, user.Username)

	// 將 session key 存入 redis
	db.SetRedis(user.Username, sessionKey)

	// 將取到的資料，回傳給前端
	c.JSON(http.StatusOK, user)
	c.SetCookie("session", sessionKey, 36000, "/", "localhost", false, true)
}

// 更新會員資料
func UpdateUser(c *gin.Context) {
	userID := c.Param("id")
	var updatedUser User

	err := c.BindJSON(&updatedUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// 查找用户是否存在
	var user User
	if err := mysql.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// 使用 Updates 方法来自动更新用户信息
	if err := mysql.Model(&user).Updates(&updatedUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

// 取得所有會員資料
func GetUsers(c *gin.Context) {
	var users []User
	if err := mysql.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get users"})
		return
	}

	c.JSON(http.StatusOK, users)
}

// 刪除會員資料
func DeleteUser(c *gin.Context) {
	userID := c.Param("id")
	var targetUser User

	// 查找用户是否存在
	if err := mysql.Where("id = ?", userID).Take(&targetUser).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	mysql.Delete(&targetUser)
}

// 確認會員是否在線
func CheckUser(c *gin.Context) {
	username := c.Param("username")
	result := db.GetRedis(username)

	if result == "" {
		c.JSON(http.StatusOK, gin.H{"message": "User is not online"})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "User is online"})
	}
}
