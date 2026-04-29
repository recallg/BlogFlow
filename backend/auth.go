package main

import (
	"database/sql"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("blogflow-secret-key-2024")

type Claims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type LoginInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterInput struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type ChangePasswordInput struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

// generateToken creates a JWT token valid for 72 hours
func generateToken(userID int, username string) (string, error) {
	claims := &Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(72 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

// AuthMiddleware validates the JWT token and sets user info in context
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未提供认证令牌"})
			c.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "认证令牌无效或已过期"})
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Next()
	}
}

// getCurrentUserID retrieves the user ID from the Gin context
func getCurrentUserID(c *gin.Context) (int, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}
	id, ok := userID.(int)
	return id, ok
}

// registerHandler handles user registration
func registerHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input RegisterInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "参数无效: " + err.Error()})
			return
		}

		// Check if username or email already exists
		var existing int
		err := db.QueryRow("SELECT COUNT(*) FROM users WHERE username = ? OR email = ?", input.Username, input.Email).Scan(&existing)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "注册失败"})
			return
		}
		if existing > 0 {
			c.JSON(http.StatusConflict, gin.H{"error": "用户名或邮箱已被注册"})
			return
		}

		// Hash password and insert user
		hashedPassword := hashPassword(input.Password)
		now := time.Now().Format("2006-01-02 15:04:05")
		result, err := db.Exec(
			"INSERT INTO users (username, email, password, created_at, updated_at) VALUES (?, ?, ?, ?, ?)",
			input.Username, input.Email, hashedPassword, now, now,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "注册失败: " + err.Error()})
			return
		}

		userID, _ := result.LastInsertId()

		// Generate token
		token, err := generateToken(int(userID), input.Username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "生成令牌失败"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"token": token,
			"user": UserResponse{
				ID:        int(userID),
				Username:  input.Username,
				Email:     input.Email,
				Avatar:    "",
				Bio:       "",
				CreatedAt: time.Now(),
			},
		})
	}
}

// loginHandler handles user login
func loginHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input LoginInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "参数无效: " + err.Error()})
			return
		}

		// Find user by email
		var user User
		err := db.QueryRow(
			"SELECT id, username, email, password, avatar, bio, created_at, updated_at FROM users WHERE email = ?",
			input.Email,
		).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Avatar, &user.Bio, &user.CreatedAt, &user.UpdatedAt)
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "邮箱或密码错误"})
			return
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "登录失败"})
			return
		}

		// Verify password
		if hashPassword(input.Password) != user.Password {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "邮箱或密码错误"})
			return
		}

		// Generate token
		token, err := generateToken(user.ID, user.Username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "生成令牌失败"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"token": token,
			"user": UserResponse{
				ID:        user.ID,
				Username:  user.Username,
				Email:     user.Email,
				Avatar:    user.Avatar,
				Bio:       user.Bio,
				CreatedAt: user.CreatedAt,
			},
		})
	}
}

// getMeHandler returns the current authenticated user's info
func getMeHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetInt("user_id")

		var user User
		err := db.QueryRow(
			"SELECT id, username, email, password, avatar, bio, created_at, updated_at FROM users WHERE id = ?",
			userID,
		).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Avatar, &user.Bio, &user.CreatedAt, &user.UpdatedAt)
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
			return
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
			return
		}

		c.JSON(http.StatusOK, UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			Avatar:    user.Avatar,
			Bio:       user.Bio,
			CreatedAt: user.CreatedAt,
		})
	}
}

// changePasswordHandler handles password change for authenticated users
func changePasswordHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetInt("user_id")

		var input ChangePasswordInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "参数无效: " + err.Error()})
			return
		}

		// Get current password from DB
		var currentPassword string
		err := db.QueryRow("SELECT password FROM users WHERE id = ?", userID).Scan(&currentPassword)
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
			return
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "查询用户失败"})
			return
		}

		// Verify old password
		if hashPassword(input.OldPassword) != currentPassword {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "原密码错误"})
			return
		}

		// Update password
		newHashed := hashPassword(input.NewPassword)
		now := time.Now().Format("2006-01-02 15:04:05")
		_, err = db.Exec("UPDATE users SET password = ?, updated_at = ? WHERE id = ?", newHashed, now, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "密码修改失败"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "密码修改成功"})
	}
}
