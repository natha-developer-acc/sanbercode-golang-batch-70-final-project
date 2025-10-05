package notification

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/mdp/qrterminal/v3"
	qrcode "github.com/skip2/go-qrcode"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"google.golang.org/protobuf/proto"

	_ "github.com/mattn/go-sqlite3"
)

var (
	client *whatsmeow.Client
	once   sync.Once
)

// initClient menginisialisasi client WhatsMeow sekali saja, menampilkan QR dan
// memonitor QR channel sehingga QR bisa ter-refresh otomatis saat timeout.
func initClient() *whatsmeow.Client {
	once.Do(func() {
		dbLog := waLog.Stdout("Database", "INFO", true)
		container, err := sqlstore.New(
			context.Background(),
			"sqlite3",
			"file:whatsapp.db?_foreign_keys=on",
			dbLog,
		)
		if err != nil {
			log.Fatalf("Gagal setup sqlstore: %v", err)
		}

		deviceStore, err := container.GetFirstDevice(context.Background())
		if err != nil {
			log.Fatalf("Gagal ambil device store: %v", err)
		}

		clientLog := waLog.Stdout("Client", "INFO", true)
		client = whatsmeow.NewClient(deviceStore, clientLog)

		// Event handler umum (tidak bergantung ke field yang mungkin berubah)
		client.AddEventHandler(func(evt interface{}) {
			switch v := evt.(type) {
			case *events.QR:
				fmt.Println("Event QR diterima:")
				if len(v.Codes) > 0 {
					code := v.Codes[0]
					qrterminal.GenerateWithConfig(code, qrterminal.Config{
						Level:     qrterminal.L,
						Writer:    os.Stdout,
						BlackChar: "█",
						WhiteChar: " ",
						QuietZone: 1,
					})
					_ = qrcode.WriteFile(code, qrcode.Medium, 256, "wa_login_qr.png")
					fmt.Println("QR juga disimpan ke 'wa_login_qr.png'")
				}
			case *events.PairSuccess:
				fmt.Println("Pair success (login WhatsApp berhasil)")
			case *events.StreamError:
				// jangan akses v.Error langsung — print structnya aman
				fmt.Printf("Stream error event: %v\n", v)
			case *events.Disconnected:
				fmt.Println("WhatsApp disconnected event:", v)
			case *events.ConnectFailure:
				fmt.Println("Connect failure event:", v)
			default:
				// untuk debug, jika perlu bisa di-uncomment
				// fmt.Printf("Event lain: %#v\n", v)
			}
		})

		// Gunakan GetQRChannel untuk dapat event code/timeout/success sehingga QR bisa auto-refresh
		qrChan, _ := client.GetQRChannel(context.Background())
		go func() {
			for evt := range qrChan {
				switch evt.Event {
				case "code":
					fmt.Println("QR baru diterbitkan (GetQRChannel):")
					qrterminal.GenerateWithConfig(evt.Code, qrterminal.Config{
						Level:     qrterminal.L,
						Writer:    os.Stdout,
						BlackChar: "█",
						WhiteChar: " ",
						QuietZone: 1,
					})
					_ = qrcode.WriteFile(evt.Code, qrcode.Medium, 256, "wa_login_qr.png")
				case "timeout":
					fmt.Println("⏱QR timeout, server akan menerbitkan QR baru otomatis")
				case "success":
					fmt.Println("QR pairing sukses (GetQRChannel)")
				default:
					// nothing
				}
			}
		}()

		// koneksi
		err = client.Connect()
		if err != nil {
			log.Fatalf("Gagal konek ke WhatsApp: %v", err)
		}
	})
	return client
}

// SendWhatsApp kirim pesan teks ke nomor WA tujuan (nomor tanpa + dan dengan kode negara)
func SendWhatsApp(phone, message string) {
	cli := initClient()

	jid := types.NewJID(phone, "s.whatsapp.net")
	msg := &waProto.Message{
		Conversation: proto.String(message),
	}

	_, err := cli.SendMessage(context.Background(), jid, msg)
	if err != nil {
		log.Println("Gagal kirim WA:", err)
	} else {
		fmt.Println("Pesan terkirim ke WA:", phone)
	}
}
