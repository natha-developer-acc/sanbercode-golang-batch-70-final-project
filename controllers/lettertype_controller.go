package controllers

import (
	"net/http"

	"sanbercode-golang-batch-70-final-project/config"
	"sanbercode-golang-batch-70-final-project/models"

	"github.com/gin-gonic/gin"
)

// LetterTypeInput digunakan untuk create & update letter type
type LetterTypeInput struct {
	Name        string `json:"name" example:"Surat Test"`
	Description string `json:"description" example:"Deskripsi surat Test"`
}

// CreateLetterType godoc
// @Summary Create a new letter type
// @Description Create a new letter type (admin only)
// @Tags Letter Types
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body LetterTypeInput true "Letter Type input payload"
// @Success 201 {object} models.LetterType
// @Failure 400 {object} map[string]string
// @Router /letter_types/ [post]
func CreateLetterType(c *gin.Context) {
	var input LetterTypeInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	lt := models.LetterType{
		Name:        input.Name,
		Description: input.Description,
	}
	config.DB.Create(&lt)
	c.JSON(http.StatusCreated, lt)
}

// GetLetterTypes godoc
// @Summary Get all letter types
// @Description Get all letter types (admin only)
// @Tags Letter Types
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.LetterType
// @Router /letter_types/ [get]
func GetLetterTypes(c *gin.Context) {
	var lts []models.LetterType
	config.DB.Find(&lts)
	c.JSON(http.StatusOK, lts)
}

// GetLetterTypeByID godoc
// @Summary Get letter type by ID
// @Description Get letter type by ID (admin only)
// @Tags Letter Types
// @Produce json
// @Security BearerAuth
// @Param id path int true "Letter Type ID"
// @Success 200 {object} models.LetterType
// @Failure 404 {object} map[string]string
// @Router /letter_types/{id} [get]
func GetLetterTypeByID(c *gin.Context) {
	var lt models.LetterType
	if err := config.DB.First(&lt, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Letter Type not found"})
		return
	}
	c.JSON(http.StatusOK, lt)
}

// UpdateLetterType godoc
// @Summary Update letter type
// @Description Update letter type (admin only)
// @Tags Letter Types
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Letter Type ID"
// @Param request body LetterTypeInput true "Letter Type update payload"
// @Success 200 {object} models.LetterType
// @Failure 404 {object} map[string]string
// @Router /letter_types/{id} [put]
func UpdateLetterType(c *gin.Context) {
	var lt models.LetterType
	if err := config.DB.First(&lt, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Letter Type not found"})
		return
	}

	var input LetterTypeInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	lt.Name = input.Name
	lt.Description = input.Description
	config.DB.Save(&lt)
	c.JSON(http.StatusOK, lt)
}

// DeleteLetterType godoc
// @Summary Delete letter type
// @Description Delete letter type (admin only)
// @Tags Letter Types
// @Produce json
// @Security BearerAuth
// @Param id path int true "Letter Type ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /letter_types/{id} [delete]
func DeleteLetterType(c *gin.Context) {
	var lt models.LetterType
	if err := config.DB.Delete(&lt, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Letter Type not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Letter Type deleted"})
}
