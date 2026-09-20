package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"

	"sosu-backend/internal/service"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("no .env found")
	}

	fmt.Println("SMTP_HOST:", os.Getenv("SMTP_HOST"))
	fmt.Println("SMTP_USER:", os.Getenv("SMTP_USER"))
	fmt.Println("SMTP_FROM:", os.Getenv("SMTP_FROM"))

	emailService := service.NewEmailService()
	err := emailService.Send("nickymanullang2121@gmail.com", "Test SMTP", "Ini email test dari sosu backend.")
	if err != nil {
		fmt.Println("GAGAL KIRIM:", err)
	} else {
		fmt.Println("BERHASIL KIRIM")
	}
}