package service

import (
	"fmt"
	"net/smtp"
	"os"
)

type EmailService struct {
	host string
	port string
	user string
	pass string
	from string
}

func NewEmailService() *EmailService {
	return &EmailService{
		host: os.Getenv("SMTP_HOST"),
		port: os.Getenv("SMTP_PORT"),
		user: os.Getenv("SMTP_USER"),
		pass: os.Getenv("SMTP_PASS"),
		from: os.Getenv("SMTP_FROM"),
	}
}

func (e *EmailService) Send(to, subject, body string) error {
	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n%s",
		e.from, to, subject, body)

	auth := smtp.PlainAuth("", e.user, e.pass, e.host)
	addr := fmt.Sprintf("%s:%s", e.host, e.port)

	return smtp.SendMail(addr, auth, e.from, []string{to}, []byte(msg))
}

func (e *EmailService) SendAlert(to, locationLabel, status, category string, hotspotCount int64) error {
	subject := fmt.Sprintf("[Cek Karhutla] Status %s di %s", status, locationLabel)
	body := fmt.Sprintf(
		"Halo,\n\nStatus wilayah %s saat ini: %s\n\nKualitas Udara: %s\nTitik Panas Terdeteksi: %d titik\n\nSilakan cek kondisi terkini dan ambil langkah pencegahan yang diperlukan.\n\n— Cek Karhutla",
		locationLabel, status, category, hotspotCount,
	)
	return e.Send(to, subject, body)
}

func (e *EmailService) SendConfirmation(to, locationLabel string) error {
	subject := "Konfirmasi Langganan - Cek Karhutla"
	body := fmt.Sprintf(
		"Halo,\n\nEmail ini terdaftar untuk menerima notifikasi status karhutla dan kualitas udara di wilayah: %s\n\nKamu akan menerima email jika status wilayah tersebut berubah menjadi Waspada atau Tidak Aman.\n\n— Cek Karhutla",
		locationLabel,
	)
	return e.Send(to, subject, body)
}
type DigestLocationSummary struct {
	LocationLabel string
	Status        string
	AQICategory   string
	AQIValue      int
	HotspotCount  int64
}

func (e *EmailService) SendWeeklyDigest(to string, summaries []DigestLocationSummary) error {
	subject := "Ringkasan Mingguan - Cek Karhutla"

	body := "Halo,\n\nBerikut ringkasan status wilayah yang kamu pantau minggu ini:\n\n"
	for _, s := range summaries {
		body += fmt.Sprintf(
			"📍 %s\n   Status: %s\n   Kualitas Udara: %s (ISPU %d)\n   Titik Panas: %d titik\n\n",
			s.LocationLabel, s.Status, s.AQICategory, s.AQIValue, s.HotspotCount,
		)
	}
	body += "Kamu menerima email ini karena berlangganan notifikasi untuk wilayah di atas.\n\n— Cek Karhutla"

	return e.Send(to, subject, body)
}