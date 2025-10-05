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
	// load env
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// init database
	config.ConnectDB()

	// ✅ Inisialisasi WhatsApp client di awal (tampilkan QR jika belum login)
	go func() {
		fmt.Println("🚀 Inisialisasi WhatsApp client...")
		notification.SendWhatsApp("6281217741759", "init") // cukup panggil tanpa `_ =`
	}()

	// setup router
	r := routes.SetupRouter()

	// swagger endpoint
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// run server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("🌐 Server berjalan di http://localhost:%s\n", port)
	r.Run(":" + port)
}
