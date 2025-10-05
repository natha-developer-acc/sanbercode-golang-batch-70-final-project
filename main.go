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
    // ✅ Load .env (abaikan kalau di Railway)
    _ = godotenv.Load()

    // ✅ Koneksi database
    config.ConnectDB()

    // ✅ Inisialisasi WhatsApp client (background)
    go func() {
        fmt.Println("🚀 Inisialisasi WhatsApp client...")
        notification.SendWhatsApp("init", "init")
    }()

    // ✅ Setup router (Swagger sudah ditangani di routes)
    r := routes.SetupRouter()

    // ✅ Ambil port dari env (Railway wajib pakai PORT)
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    fmt.Printf("✅ Server berjalan di port %s\n", port)
    if err := r.Run(":" + port); err != nil {
        log.Fatalf("❌ Gagal menjalankan server: %v", err)
    }
}
