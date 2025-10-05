package controllers

import (
	"net/http"

	"sanbercode-golang-batch-70-final-project/config"
	"sanbercode-golang-batch-70-final-project/models"

	"github.com/gin-gonic/gin"
)

// ==============================
// STRUCTS UNTUK INPUT PAYLOAD
// ==============================

// SettingCreateInput digunakan untuk membuat setting baru
type SettingCreateInput struct {
	UserID        uint   `json:"user_id" example:"5"`
	TelegramChat  string `json:"telegram_chatid" example:"123456789"`
	WANumber      string `json:"wa_number" example:"6281234567890"`
	AllowTelegram string `json:"allow_telegram" example:"yes"`
	AllowWA       string `json:"allow_wa" example:"yes"`
}

// SettingUpdateInput digunakan untuk mengupdate setting
type SettingUpdateInput struct {
	TelegramChat  string `json:"telegram_chatid" example:"1234567890"`
	WANumber      string `json:"wa_number" example:"62812345678900"`
	AllowTelegram string `json:"allow_telegram" example:"no"`
	AllowWA       string `json:"allow_wa" example:"no"`
}

// ==============================
// CREATE SETTING
// ==============================

// CreateSetting godoc
// @Summary Create a new setting
// @Description Buat pengaturan notifikasi baru (admin only)
// @Tags Settings
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body SettingCreateInput true "Setting create payload"
// @Success 201 {object} models.Setting
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /settings/ [post]
func CreateSetting(c *gin.Context) {
	role, _ := c.Get("role")
	if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Hanya admin yang bisa membuat setting!"})
		return
	}

	var input SettingCreateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// cek apakah user exist
	var user models.User
	if err := config.DB.Preload("Role").First(&user, input.UserID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User tidak ditemukan"})
		return
	}

	// cek apakah sudah ada setting untuk user ini
	var existing models.Setting
	if err := config.DB.Where("user_id = ?", input.UserID).First(&existing).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":      "Setting untuk User dengan ID ini sudah ada!",
			"user_id":    input.UserID,
			"setting_id": existing.ID,
		})
		return
	}

	setting := models.Setting{
		UserID:         input.UserID,
		TelegramChatID: input.TelegramChat,
		WANumber:       input.WANumber,
		AllowTelegram:  input.AllowTelegram,
		AllowWA:        input.AllowWA,
	}

	if err := config.DB.Create(&setting).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat setting"})
		return
	}

	config.DB.Preload("User.Role").First(&setting, setting.ID)
	c.JSON(http.StatusCreated, setting)
}

// ==============================
// GET ALL SETTINGS
// ==============================

// GetSettings godoc
// @Summary Get all settings
// @Description Ambil semua data setting (admin only)
// @Tags Settings
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.Setting
// @Failure 403 {object} map[string]string
// @Router /settings/ [get]
func GetSettings(c *gin.Context) {
	role, _ := c.Get("role")
	if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Hanya admin yang bisa melihat semua setting!"})
		return
	}

	var settings []models.Setting
	config.DB.Preload("User.Role").Find(&settings)
	c.JSON(http.StatusOK, settings)
}

// ==============================
// GET SETTING BY ID
// ==============================

// GetSettingByID godoc
// @Summary Get setting by ID
// @Description Ambil detail setting berdasarkan ID
// @Tags Settings
// @Produce json
// @Security BearerAuth
// @Param id path int true "Setting ID"
// @Success 200 {object} models.Setting
// @Failure 404 {object} map[string]string
// @Router /settings/{id} [get]
func GetSettingByID(c *gin.Context) {
	var setting models.Setting
	if err := config.DB.Preload("User.Role").First(&setting, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Setting not found"})
		return
	}
	c.JSON(http.StatusOK, setting)
}

// ==============================
// UPDATE SETTING
// ==============================

// UpdateSetting godoc
// @Summary Update setting
// @Description Ubah data setting (admin only)
// @Tags Settings
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Setting ID"
// @Param request body SettingUpdateInput true "Setting update payload"
// @Success 200 {object} models.Setting
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /settings/{id} [put]
func UpdateSetting(c *gin.Context) {
	role, _ := c.Get("role")
	if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Hanya admin yang bisa mengupdate setting!"})
		return
	}

	var setting models.Setting
	if err := config.DB.First(&setting, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Setting not found"})
		return
	}

	var input SettingUpdateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	setting.TelegramChatID = input.TelegramChat
	setting.WANumber = input.WANumber
	setting.AllowTelegram = input.AllowTelegram
	setting.AllowWA = input.AllowWA

	if err := config.DB.Save(&setting).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal update setting"})
		return
	}

	config.DB.Preload("User.Role").First(&setting, setting.ID)
	c.JSON(http.StatusOK, setting)
}

// ==============================
// DELETE SETTING
// ==============================

// DeleteSetting godoc
// @Summary Delete setting
// @Description Hapus setting (admin only)
// @Tags Settings
// @Produce json
// @Security BearerAuth
// @Param id path int true "Setting ID"
// @Success 200 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /settings/{id} [delete]
func DeleteSetting(c *gin.Context) {
	role, _ := c.Get("role")
	if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Hanya admin yang bisa menghapus setting!"})
		return
	}

	var setting models.Setting
	if err := config.DB.First(&setting, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Setting not found"})
		return
	}

	if err := config.DB.Delete(&setting).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus setting"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Setting deleted"})
}
