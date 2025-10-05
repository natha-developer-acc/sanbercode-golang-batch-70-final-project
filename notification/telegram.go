package notification

import (
    "fmt"
    "log"
    "os"
    "strconv"

    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// SendTelegram kirim pesan ke Telegram berdasarkan token dan chatID dari env
func SendTelegram(chatID, message string) {
    token := os.Getenv("TELEGRAM_TOKEN")
    if token == "" {
        log.Println("TELEGRAM_TOKEN belum diatur di .env")
        return
    }

    bot, err := tgbotapi.NewBotAPI(token)
    if err != nil {
        log.Println("Gagal konek Telegram:", err)
        return
    }

    id, err := strconv.ParseInt(chatID, 10, 64)
    if err != nil {
        log.Println("ChatID tidak valid:", chatID)
        return
    }

    msg := tgbotapi.NewMessage(id, message)
    _, err = bot.Send(msg)
    if err != nil {
        log.Println("Gagal kirim pesan Telegram:", err)
    } else {
        fmt.Println("Pesan terkirim ke Telegram:", chatID)
    }
}
