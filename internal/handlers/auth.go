package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"go-backend-react-frontend/internal/auth"
	"go-backend-react-frontend/internal/db"
	"go-backend-react-frontend/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
)

type RegisterRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
	FirstName string `json:"firstName" binding:"required"`
	LastName  string `json:"lastName" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existing models.User
	if err := db.DB.Where("email = ?", req.Email).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already registered"})
		return
	}

	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	name := req.FirstName + " " + req.LastName
	user := models.User{
		ID:           uuid.New().String(),
		Email:        req.Email,
		Name:         &name,
		FirstName:    &req.FirstName,
		LastName:     &req.LastName,
		PasswordHash: &hash,
		Role:         models.RoleLearner,
		IsActive:     true,
	}

	if err := db.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	token, err := auth.GenerateToken(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.SetCookie("token", token, 3*24*3600, "/", "", false, true)
	c.JSON(http.StatusCreated, gin.H{
		"token": token,
		"user":  newAuthUserResponse(&user),
	})
}

func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := db.DB.Where("email = ? AND is_active = true", req.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	if user.PasswordHash == nil || !auth.CheckPassword(req.Password, *user.PasswordHash) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	now := time.Now()
	db.DB.Model(&user).Update("last_login_at", now)

	token, err := auth.GenerateToken(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.SetCookie("token", token, 3*24*3600, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user":  newAuthUserResponse(&user),
	})
}

func GetMe(c *gin.Context) {
	userID, _ := c.Get("userId")
	var user models.User
	if err := db.DB.First(&user, "id = ?", userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, newAuthUserResponse(&user))
}

func Logout(c *gin.Context) {
	c.SetCookie("token", "", -1, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Logged out"})
}

// Google OAuth

var googleOAuthState = "random-state-string" // In production, use a random state per request

func GoogleLogin(c *gin.Context) {
	config := auth.GoogleOAuthConfig()
	if config.ClientID == "" {
		c.JSON(http.StatusNotImplemented, gin.H{"error": "Google OAuth not configured"})
		return
	}
	url := config.AuthCodeURL(googleOAuthState, oauth2.AccessTypeOffline)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func GoogleCallback(c *gin.Context) {
	config := auth.GoogleOAuthConfig()

	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing code parameter"})
		return
	}

	token, err := config.Exchange(c.Request.Context(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange token"})
		return
	}

	client := config.Client(c.Request.Context(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var googleUser struct {
		ID      string `json:"id"`
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture string `json:"picture"`
	}
	json.Unmarshal(body, &googleUser)

	var user models.User
	err = db.DB.Where("email = ?", googleUser.Email).First(&user).Error
	if err != nil {
		user = models.User{
			ID:       uuid.New().String(),
			Email:    googleUser.Email,
			Name:     &googleUser.Name,
			Image:    &googleUser.Picture,
			Role:     models.RoleLearner,
			IsActive: true,
		}
		db.DB.Create(&user)
	}

	now := time.Now()
	db.DB.Model(&user).Update("last_login_at", now)

	jwtToken, err := auth.GenerateToken(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.SetCookie("token", jwtToken, 3*24*3600, "/", "", false, true)
	c.Redirect(http.StatusTemporaryRedirect, "/dashboard?token="+jwtToken)
}
