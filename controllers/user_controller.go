package controllers

import (
	"net/http"

	"sanbercode-golang-batch-70-final-project/config"
	"sanbercode-golang-batch-70-final-project/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// ===== Struct tambahan untuk dokumentasi Swagger =====

type UserInput struct {
	RoleID   uint   `json:"role_id" example:"3"`
	Name     string `json:"name" example:"User Baru"`
	Email    string `json:"email" example:"baru@mail.com"`
	Password string `json:"password" example:"password123"`
}

type UserUpdateInput struct {
	RoleID   uint   `json:"role_id" example:"3"`
	Name     string `json:"name" example:"User Baru"`
	Email    string `json:"email" example:"baru@mail.com"`
	Password string `json:"password,omitempty" example:"newpassword123"`
}

// =====================================================

// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user (only admin can access)
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body UserInput true "User input"
// @Success 201 {object} models.User
// @Failure 400 {object} map[string]string
// @Router /users/ [post]
func CreateUser(c *gin.Context) {
	var input UserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal hashing password"})
		return
	}

	user := models.User{
		RoleID:   input.RoleID,
		Name:     input.Name,
		Email:    input.Email,
		Password: string(hashedPassword),
	}

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email sudah digunakan"})
		return
	}

	config.DB.Preload("Role").First(&user, user.ID)
	c.JSON(http.StatusCreated, user)
}

// GetUsers godoc
// @Summary Get all users
// @Description Get all users (only admin can access)
// @Tags Users
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.User
// @Router /users/ [get]
func GetUsers(c *gin.Context) {
	var users []models.User
	config.DB.Preload("Role").Find(&users)
	c.JSON(http.StatusOK, users)
}

// GetUserByID godoc
// @Summary Get user by ID
// @Description Get user detail by ID (only admin can access)
// @Tags Users
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Success 200 {object} models.User
// @Failure 404 {object} map[string]string
// @Router /users/{id} [get]
func GetUserByID(c *gin.Context) {
	var user models.User
	if err := config.DB.Preload("Role").First(&user, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

// UpdateUser godoc
// @Summary Update user by ID
// @Description Update user data (only admin can access)
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Param request body UserUpdateInput true "User update input"
// @Success 200 {object} models.User
// @Failure 404 {object} map[string]string
// @Router /users/{id} [put]
func UpdateUser(c *gin.Context) {
	var user models.User
	if err := config.DB.First(&user, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var input UserUpdateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user.RoleID = input.RoleID
	user.Name = input.Name
	user.Email = input.Email

	if input.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal hashing password"})
			return
		}
		user.Password = string(hashedPassword)
	}

	if err := config.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal update user"})
		return
	}

	config.DB.Preload("Role").First(&user, user.ID)
	c.JSON(http.StatusOK, user)
}

// DeleteUser godoc
// @Summary Delete user by ID
// @Description Delete user (only admin can access)
// @Tags Users
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Success 200 {object} map[string]string
// @Router /users/{id} [delete]
func DeleteUser(c *gin.Context) {
	var user models.User
	if err := config.DB.Delete(&user, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
}
