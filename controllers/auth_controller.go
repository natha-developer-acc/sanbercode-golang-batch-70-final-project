package controllers

import (
	"net/http"
	"os"
	"time"

	"sanbercode-golang-batch-70-final-project/config"
	"sanbercode-golang-batch-70-final-project/models"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// ==================== REGISTER ====================

// Register godoc
// @Summary Register user baru
// @Description Membuat akun user baru dengan role default "user"
// @Tags Auth
// @Accept json
// @Produce json
// @Param register body models.RegisterInput true "Data user baru" example({"name":"Admin Satu","email":"admin@mail.com","password":"admin123"})
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Router /users/register [post]
func Register(c *gin.Context) {
	var input models.RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal hashing password"})
		return
	}

	user := models.User{
		RoleID:   3, // default role user
		Name:     input.Name,
		Email:    input.Email,
		Password: string(hashedPassword),
	}

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email sudah digunakan"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Registrasi berhasil",
		"user": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"role":  user.RoleID,
		},
	})
}

// ==================== LOGIN ====================

// LoginInput digunakan untuk dokumentasi swagger (supaya ada contoh payload)
type LoginInput struct {
	Email    string `json:"email" example:"admin@mail.com"`
	Password string `json:"password" example:"admin123"`
}

// Login godoc
// @Summary Login user
// @Description Login user menggunakan email dan password, menghasilkan JWT token
// @Tags Auth
// @Accept json
// @Produce json
// @Param login body LoginInput true "Data login user"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /users/login [post]
func Login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := config.DB.Preload("Role").Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// cek hash password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// generate JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":   user.ID,
		"role": user.Role.Name,
		"exp":  time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
		"user": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"role":  user.Role.Name,
		},
	})
}
