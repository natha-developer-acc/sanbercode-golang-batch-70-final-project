package controllers

import (
	"fmt"
	"net/http"

	"sanbercode-golang-batch-70-final-project/config"
	"sanbercode-golang-batch-70-final-project/models"
	"sanbercode-golang-batch-70-final-project/notification"

	"github.com/gin-gonic/gin"
)

// ===============================
// Struct Input
// ===============================

// LetterCreateInput digunakan untuk membuat surat baru
type LetterCreateInput struct {
	UserID uint `json:"user_id,omitempty" example:"5"`
	TypeID uint `json:"type_id" example:"1" binding:"required"`
}

// LetterUpdateInput digunakan untuk update surat
type LetterUpdateInput struct {
	UserID       uint   `json:"user_id,omitempty" example:"4"`
	TypeID       uint   `json:"type_id,omitempty" example:"3"`
	Status       string `json:"status,omitempty" example:"accepted"`
	RejectReason string `json:"reject_reason,omitempty" example:"Ditolak untuk testing"`
}

// ===============================
// Create Letter
// ===============================

// CreateLetter godoc
// @Summary Create a new letter
// @Description Buat pengajuan surat baru (user & admin bisa)
// @Tags Letters
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body LetterCreateInput true "Letter create payload"
// @Success 201 {object} models.Letter
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /letters/ [post]
func CreateLetter(c *gin.Context) {
	roleVal, _ := c.Get("role")
	role := fmt.Sprintf("%v", roleVal)

	// Reviewer tidak boleh bikin surat
	if role == "reviewer" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Reviewer hanya bisa menerima atau menolak pengajuan surat!"})
		return
	}

	var input LetterCreateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var userID uint
	if role == "user" {
		// Ambil user_id dari token (bukan dari body)
		uid, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "User ID tidak ditemukan di token"})
			return
		}
		userID = uid.(uint)
	} else if role == "admin" {
		// Admin wajib isi user_id
		if input.UserID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Admin harus menentukan user_id untuk surat ini"})
			return
		}
		userID = input.UserID
	} else {
		c.JSON(http.StatusForbidden, gin.H{"error": "Role tidak dikenali"})
		return
	}

	// Validasi user
	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User tidak ditemukan"})
		return
	}

	// Validasi tipe surat
	var letterType models.LetterType
	if err := config.DB.First(&letterType, input.TypeID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Letter type tidak ditemukan"})
		return
	}

	// Buat surat baru
	letter := models.Letter{
		UserID:       userID,
		TypeID:       input.TypeID,
		Status:       "pending",
		RejectReason: "",
	}

	if err := config.DB.Create(&letter).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat surat"})
		return
	}

	config.DB.Preload("User.Role").Preload("LetterType").First(&letter, letter.ID)

	// Kirim notifikasi ke semua reviewer yang aktif
	var settings []models.Setting
	config.DB.Preload("User.Role").Where("allow_telegram = 'yes' OR allow_wa = 'yes'").Find(&settings)

	message := fmt.Sprintf("📩 Pengajuan surat baru dari *%s* untuk jenis surat *%s* (status: pending).",
		user.Name, letterType.Name)

	for _, s := range settings {
		if s.User.Role.Name == "reviewer" {
			if s.AllowTelegram == "yes" && s.TelegramChatID != "" {
				go notification.SendTelegram(s.TelegramChatID, message)
			}
			if s.AllowWA == "yes" && s.WANumber != "" {
				go notification.SendWhatsApp(s.WANumber, message)
			}
		}
	}

	c.JSON(http.StatusCreated, letter)
}

// ===============================
// Get All Letters
// ===============================

// GetLetters godoc
// @Summary Get all letters
// @Description Get all letters (admin & reviewer only)
// @Tags Letters
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.Letter
// @Router /letters/ [get]
func GetLetters(c *gin.Context) {
	role, _ := c.Get("role")

	if role == "user" {
		uid, _ := c.Get("user_id")
		var letters []models.Letter
		config.DB.Preload("User.Role").Preload("LetterType").
			Where("user_id = ?", uid).Find(&letters)
		c.JSON(http.StatusOK, letters)
		return
	}

	var letters []models.Letter
	config.DB.Preload("User.Role").Preload("LetterType").Find(&letters)
	c.JSON(http.StatusOK, letters)
}

// ===============================
// Get Letter By ID
// ===============================

func GetLetterByID(c *gin.Context) {
	var letter models.Letter
	if err := config.DB.Preload("User.Role").Preload("LetterType").
		First(&letter, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Letter not found"})
		return
	}
	c.JSON(http.StatusOK, letter)
}

// ===============================
// Update Letter
// ===============================

func UpdateLetter(c *gin.Context) {
	role, _ := c.Get("role")

	var letter models.Letter
	if err := config.DB.Preload("User").Preload("LetterType").
		First(&letter, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Letter not found"})
		return
	}

	var input LetterUpdateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	switch role {
	case "admin":
		if input.UserID != 0 {
			letter.UserID = input.UserID
		}
		if input.TypeID != 0 {
			letter.TypeID = input.TypeID
		}
		if input.Status != "" {
			letter.Status = input.Status
			if input.Status == "accepted" {
				letter.RejectReason = ""
			} else if input.Status == "rejected" {
				letter.RejectReason = input.RejectReason
			}
		}
	case "reviewer":
		if input.Status == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Reviewer hanya bisa mengubah status surat"})
			return
		}
		if input.UserID != 0 || input.TypeID != 0 {
			c.JSON(http.StatusForbidden, gin.H{"error": "Reviewer tidak boleh ubah ID pengguna/jenis surat"})
			return
		}
		if input.Status == "accepted" {
			letter.Status = "accepted"
			letter.RejectReason = ""
		} else if input.Status == "rejected" {
			letter.Status = "rejected"
			letter.RejectReason = input.RejectReason
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Status tidak valid untuk reviewer"})
			return
		}
	case "user":
		c.JSON(http.StatusForbidden, gin.H{"error": "User hanya bisa mengajukan surat!"})
		return
	default:
		c.JSON(http.StatusForbidden, gin.H{"error": "Role tidak dikenali"})
		return
	}

	if err := config.DB.Save(&letter).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal update surat"})
		return
	}

	config.DB.Preload("User.Role").Preload("LetterType").First(&letter, letter.ID)

	// Kirim notifikasi ke user
	var setting models.Setting
	if err := config.DB.Where("user_id = ?", letter.UserID).First(&setting).Error; err == nil {
		message := fmt.Sprintf("📢 Status surat kamu (%s) kini: *%s*.",
			letter.LetterType.Name, letter.Status)

		if letter.Status == "rejected" && letter.RejectReason != "" {
			message += fmt.Sprintf("\nAlasan: %s", letter.RejectReason)
		}

		if setting.AllowTelegram == "yes" && setting.TelegramChatID != "" {
			go notification.SendTelegram(setting.TelegramChatID, message)
		}
		if setting.AllowWA == "yes" && setting.WANumber != "" {
			go notification.SendWhatsApp(setting.WANumber, message)
		}
	}

	c.JSON(http.StatusOK, letter)
}

// ===============================
// Delete Letter
// ===============================

func DeleteLetter(c *gin.Context) {
	role, _ := c.Get("role")

	if role == "reviewer" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Reviewer hanya bisa menerima atau menolak pengajuan surat!"})
		return
	}
	if role == "user" {
		c.JSON(http.StatusForbidden, gin.H{"error": "User hanya bisa mengajukan surat!"})
		return
	}

	if err := config.DB.Delete(&models.Letter{}, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Letter not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Letter deleted"})
}
