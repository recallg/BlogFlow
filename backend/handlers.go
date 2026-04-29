package main

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func getArticleListHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
		search := c.Query("search")

		if page < 1 {
			page = 1
		}
		if pageSize < 1 || pageSize > 50 {
			pageSize = 10
		}

		offset := (page - 1) * pageSize

		var total int
		var rows *sql.Rows
		var err error

		if search != "" {
			query := "SELECT COUNT(*) FROM articles WHERE title LIKE ? OR content LIKE ?"
			err = db.QueryRow(query, "%"+search+"%", "%"+search+"%").Scan(&total)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
				return
			}
			query = "SELECT id, title, content, summary, author, tags, cover_url, created_at, updated_at FROM articles WHERE title LIKE ? OR content LIKE ? ORDER BY created_at DESC LIMIT ? OFFSET ?"
			rows, err = db.Query(query, "%"+search+"%", "%"+search+"%", pageSize, offset)
		} else {
			err = db.QueryRow("SELECT COUNT(*) FROM articles").Scan(&total)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
				return
			}
			rows, err = db.Query("SELECT id, title, content, summary, author, tags, cover_url, created_at, updated_at FROM articles ORDER BY created_at DESC LIMIT ? OFFSET ?", pageSize, offset)
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
			return
		}
		defer rows.Close()

		articles := make([]Article, 0)
		for rows.Next() {
			var a Article
			err := rows.Scan(&a.ID, &a.Title, &a.Content, &a.Summary, &a.Author, &a.Tags, &a.CoverURL, &a.CreatedAt, &a.UpdatedAt)
			if err != nil {
				continue
			}
			articles = append(articles, a)
		}

		c.JSON(http.StatusOK, ArticleListResponse{
			Articles: articles,
			Total:    total,
			Page:     page,
			PageSize: pageSize,
		})
	}
}

func getArticleHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
			return
		}

		var a Article
		err = db.QueryRow("SELECT id, title, content, summary, author, tags, cover_url, created_at, updated_at FROM articles WHERE id = ?", id).
			Scan(&a.ID, &a.Title, &a.Content, &a.Summary, &a.Author, &a.Tags, &a.CoverURL, &a.CreatedAt, &a.UpdatedAt)
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
			return
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
			return
		}

		c.JSON(http.StatusOK, a)
	}
}

func createArticleHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var a Article
		if err := c.ShouldBindJSON(&a); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "参数无效: " + err.Error()})
			return
		}

		// If user is authenticated, use their username as author
		if userID, exists := c.Get("user_id"); exists {
			if username, ok := c.Get("username"); ok {
				a.Author = username.(string)
			}
			_ = userID
		} else if a.Author == "" {
			a.Author = "Anonymous"
		}

		now := time.Now()
		result, err := db.Exec(
			"INSERT INTO articles (title, content, summary, author, tags, cover_url, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
			a.Title, a.Content, a.Summary, a.Author, a.Tags, a.CoverURL, now, now,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "创建失败"})
			return
		}

		id, _ := result.LastInsertId()
		a.ID = int(id)
		a.CreatedAt = now
		a.UpdatedAt = now

		c.JSON(http.StatusCreated, a)
	}
}

func updateArticleHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
			return
		}

		var a Article
		if err := c.ShouldBindJSON(&a); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "参数无效: " + err.Error()})
			return
		}

		now := time.Now()
		result, err := db.Exec(
			"UPDATE articles SET title=?, content=?, summary=?, author=?, tags=?, cover_url=?, updated_at=? WHERE id=?",
			a.Title, a.Content, a.Summary, a.Author, a.Tags, a.CoverURL, now, id,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败"})
			return
		}

		affected, _ := result.RowsAffected()
		if affected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
			return
		}

		a.ID = id
		a.UpdatedAt = now
		c.JSON(http.StatusOK, a)
	}
}

func deleteArticleHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
			return
		}

		result, err := db.Exec("DELETE FROM articles WHERE id = ?", id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
			return
		}

		affected, _ := result.RowsAffected()
		if affected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
	}
}
