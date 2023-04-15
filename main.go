package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func main() {
	// Connect to the database
	db, err := ConnectToDB()
	if err != nil {
		panic("Failed to connect to the database")
	}
	defer db.Close()
	// Set up Gin router
	r := gin.Default()

	// JWT Secret Key
	jwtSecretKey := []byte("secret key")

	// Login endpoint
	r.POST("/login", func(c *gin.Context) {
		var loginVals struct {
			Username string `json:"username" binding:"required"`
			Password string `json:"password" binding:"required"`
		}

		if err := c.ShouldBindJSON(&loginVals); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var user User
		if err := db.Where("username = ? AND password = ?", loginVals.Username, loginVals.Password).First(&user).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"id":  user.ID,
			"exp": time.Now().Add(time.Hour * 1).Unix(),
		})

		tokenString, err := token.SignedString(jwtSecretKey)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"token": tokenString})
	})

	// JWT Authorization Middleware
	jwtAuthMiddleware := func(c *gin.Context) {
		tokenStr := c.Request.Header.Get("Authorization")
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return jwtSecretKey, nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		c.Next()
	}

	// Secured routes
	secured := r.Group("/")
	secured.Use(jwtAuthMiddleware)
	{
		secured.GET("/jobs", func(c *gin.Context) {
			description := c.Query("description")
			location := c.Query("location")
			fullTime := c.Query("full_time") == "true"
			pageStr := c.DefaultQuery("page", "1")
			page, err := strconv.Atoi(pageStr)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page value"})
				return
			}

			jobs, err := GetJobs(description, location, fullTime, page)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to fetch job listings: %v", err)})
				return
			}

			c.JSON(http.StatusOK, jobs)
		})
		secured.GET("/jobs/:id", func(c *gin.Context) {
			id := c.Param("id")

			job, err := GetJobDetail(id)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch job details"})
				return
			}

			c.JSON(http.StatusOK, job)
		})
	}

	// Run the server
	r.Run(":8080")
}
