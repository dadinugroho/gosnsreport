package main

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	sendEmailTo := os.Getenv("SEND_EMAIL_TO")

	today := time.Now().Format("2006-01-02")
	todayString := time.Now().Format("2-Jan-06")
	weekday := time.Now().Weekday()

	var yesterday string
	var yesterdayString string
	if weekday == time.Monday {
		// If today is Monday, set yesterday to Saturday
		yesterday = time.Now().AddDate(0, 0, -2).Format("2006-01-02")
		yesterdayString = time.Now().AddDate(0, 0, -2).Format("2-Jan-06")
	} else {
		yesterday = time.Now().AddDate(0, 0, -1).Format("2006-01-02")
		yesterdayString = time.Now().AddDate(0, 0, -1).Format("2-Jan-06")
	}

	stockDataYesterday, err := fetchStockReport(yesterday)
	if err != nil {
		log.Fatalf("Failed to fetch stock report: %v", err)
	}

	salesDataYesterday, err := fetchSalesReport(yesterday)
	if err != nil {
		log.Fatalf("Failed to fetch sales report: %v", err)
	}

	stockDataToday, err := fetchStockReport(today)
	if err != nil {
		log.Fatalf("Failed to fetch stock report: %v", err)
	}

	salesDataToday, err := fetchSalesReport(today)
	if err != nil {
		log.Fatalf("Failed to fetch sales report: %v", err)
	}

	// Generate Excel report
	excelFile, err := generateExcelReport(stockDataToday, salesDataToday, stockDataYesterday, salesDataYesterday, todayString, yesterdayString)
	if err != nil {
		log.Fatalf("Failed to generate Excel report: %v", err)
	}

	// Send the report via Gmail
	emails := strings.Split(sendEmailTo, ",")
	subject := "Laporan Stok dan Penjualan per " + time.Now().Format("2 Jan 2006")
	body := "Terlampir " + subject
	for _, email := range emails {
		err := sendEmailWithGmail(strings.TrimSpace(email), subject, body, excelFile)
		if err != nil {
			log.Printf("Failed to send email to %s: %v", email, err)
		} else {
			log.Printf("Email sent to %s", email)
		}
	}

	log.Println("Report sent successfully")
}
