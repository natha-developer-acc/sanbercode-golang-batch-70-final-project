package main

import (
	"fmt"
	"log"
	"os"

	"sanbercode-golang-batch-70-final-project/config"
	_ "sanbercode-golang-batch-70-final-project/docs"
	"sanbercode-golang-batch-70-final-project/notification"
	"sanbercode-golang-batch-70-final-project/routes"

	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Surat Notifikasi API
// @version 1.0
// @description API untuk notifikasi pengajuan surat via Telegram & WhatsApp
// @host localhost:8080
// @BasePath /api
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// ✅ Load file .env (abaikan jika tidak ditemukan, misal di Railway)
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found — using environment variables from Railway")
	}

	// ✅ Koneksi ke database
	config.ConnectDB()

	// ✅ Inisialisasi WhatsApp client (tampilkan QR jika belum login)
	go func() {
		fmt.Println("Inisialisasi WhatsApp client...")
		notification.SendWhatsApp("init", "init") // aman, hanya trigger
	}()

	// ✅ Setup router
	r := routes.SetupRouter()

	// ✅ Swagger endpoint
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// ✅ Cek port (Railway pakai env PORT)
	port := os.Getenv("PORT")
	if port == "" {
		port = os.Getenv("APP_PORT")
		if port == "" {
			port = "8080"
		}
	}

	fmt.Printf("Server berjalan di http://localhost:%s\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Gagal menjalankan server: %v", err)
	}
}
