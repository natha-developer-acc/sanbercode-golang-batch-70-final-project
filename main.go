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
	// ✅ 1. Load file .env (abaikan jika tidak ada, misal di Railway)
	if err := godotenv.Load(); err != nil {
		log.Println("[INFO] .env file tidak ditemukan — menggunakan environment variables dari Railway")
	}

	// ✅ 2. Koneksi ke database
	config.ConnectDB()

	// ✅ 3. Inisialisasi WhatsApp client (trigger login / QR code)
	go func() {
		fmt.Println("🚀 Inisialisasi WhatsApp client...")
		notification.SendWhatsApp("init", "init") // hanya trigger, bukan kirim pesan sungguhan
	}()

	// ✅ 4. Setup router
	r := routes.SetupRouter()

	// ✅ 5. Tambahkan endpoint Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// ✅ 6. Baca PORT dari environment (Railway pakai PORT, bukan APP_PORT)
	port := os.Getenv("PORT")
	if port == "" {
		port = os.Getenv("APP_PORT") // fallback kalau di lokal
		if port == "" {
			port = "8080"
		}
	}

	// ✅ 7. Jalankan server
	fmt.Printf("✅ Server berjalan di http://localhost:%s\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("❌ Gagal menjalankan server: %v", err)
	}
}
