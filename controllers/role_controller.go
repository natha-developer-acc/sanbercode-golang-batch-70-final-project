package controllers

import (
    "net/http"

    "sanbercode-golang-batch-70-final-project/config"
    "sanbercode-golang-batch-70-final-project/models"

    "github.com/gin-gonic/gin"
)

// ==== Tambahan struct untuk Swagger ====
type RoleInput struct {
    Name string `json:"name" example:"test"`
}

// CreateRole godoc
// @Summary Create new role
// @Description Create a new role (only admin can access)
// @Tags Roles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body RoleInput true "Role input"
// @Success 200 {object} models.Role
// @Failure 400 {object} map[string]string
// @Router /roles/ [post]
func CreateRole(c *gin.Context) {
    var role models.Role
    if err := c.ShouldBindJSON(&role); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    config.DB.Create(&role)
    c.JSON(http.StatusOK, role)
}

// GetRoles godoc
// @Summary Get all roles
// @Description Get list of all roles
// @Tags Roles
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.Role
// @Router /roles/ [get]
func GetRoles(c *gin.Context) {
    var roles []models.Role
    config.DB.Find(&roles)
    c.JSON(http.StatusOK, roles)
}

// GetRoleByID godoc
// @Summary Get role by ID
// @Description Get detail of a specific role
// @Tags Roles
// @Produce json
// @Security BearerAuth
// @Param id path int true "Role ID"
// @Success 200 {object} models.Role
// @Failure 404 {object} map[string]string
// @Router /roles/{id} [get]
func GetRoleByID(c *gin.Context) {
    var role models.Role
    if err := config.DB.First(&role, c.Param("id")).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Role not found"})
        return
    }
    c.JSON(http.StatusOK, role)
}

// UpdateRole godoc
// @Summary Update role by ID
// @Description Update role name (only admin can access)
// @Tags Roles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Role ID"
// @Param request body RoleInput true "Role update input"
// @Success 200 {object} models.Role
// @Failure 404 {object} map[string]string
// @Router /roles/{id} [put]
func UpdateRole(c *gin.Context) {
    var role models.Role
    if err := config.DB.First(&role, c.Param("id")).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Role not found"})
        return
    }

    c.BindJSON(&role)
    config.DB.Save(&role)
    c.JSON(http.StatusOK, role)
}

// DeleteRole godoc
// @Summary Delete role by ID
// @Description Delete role (only admin can access)
// @Tags Roles
// @Produce json
// @Security BearerAuth
// @Param id path int true "Role ID"
// @Success 200 {object} map[string]string
// @Router /roles/{id} [delete]
func DeleteRole(c *gin.Context) {
    var role models.Role
    if err := config.DB.Delete(&role, c.Param("id")).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Role not found"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "Role deleted"})
}
