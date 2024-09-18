package main

import (
    "log"
    "time"
)

func main() {
    today := time.Now().Format("2006-01-02") 
    yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02") 

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
    excelFile, err := generateExcelReport(stockDataToday, salesDataToday, stockDataYesterday, salesDataYesterday)
    if err != nil {
        log.Fatalf("Failed to generate Excel report: %v", err)
    }

    // Send the report via Gmail
    email := "adinugro@gmail.com"
    subject := "Daily Stock and Sales Report"
    body := "Please find the attached stock and sales report."
    err = sendEmailWithGmail(email, subject, body, excelFile)
    if err != nil {
        log.Fatalf("Failed to send email: %v", err)
    }

    log.Println("Report sent successfully")
}
