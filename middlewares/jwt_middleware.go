package middlewares

import (
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(roleRequired string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			c.Abort()
			return
		}

		// Pastikan format "Bearer <token>"
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		claims := jwt.MapClaims{}
		_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Ambil role dan user_id dari claims
		role, ok := claims["role"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token payload"})
			c.Abort()
			return
		}

		// --- Bagian penting: konversi ID dari float64 ke uint ---
		var userID uint
		switch idValue := claims["id"].(type) {
		case float64:
			userID = uint(idValue)
		case int:
			userID = uint(idValue)
		case uint:
			userID = idValue
		default:
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID type in token"})
			c.Abort()
			return
		}

		// Cek role kalau endpoint terbatas
		if roleRequired != "" && role != roleRequired {
			c.JSON(http.StatusForbidden, gin.H{"error": "Endpoint hanya bisa diakses Admin!"})
			c.Abort()
			return
		}

		// Simpan ke context supaya bisa diakses di controller
		c.Set("user_id", userID)
		c.Set("role", role)

		c.Next()
	}
}
