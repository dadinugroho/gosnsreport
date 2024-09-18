package main

import (
	"fmt"
	"time"

	"github.com/tealeg/xlsx"
)

func generateExcelReport(stockDataToday []Stock, salesDataToday []Sales, stockDataYesterday []Stock, salesDataYesterday []Sales) (string, error) {
	today := time.Now().Format("2-Jan-06")
	yesterday := time.Now().AddDate(0, 0, -1).Format("2-Jan-06")

	file := xlsx.NewFile()

	stockSheetToday, err := file.AddSheet("Stok " + today)
	if err != nil {
		return "", fmt.Errorf("failed to add stock sheet: %v", err)
	}
	addStockRows(stockSheetToday, stockDataToday)

	salesSheetToday, err := file.AddSheet("Penjualan " + today)
	if err != nil {
		return "", fmt.Errorf("failed to add sales sheet: %v", err)
	}
	addSalesRows(salesSheetToday, salesDataToday)

	stockSheetYesterday, err := file.AddSheet("Stok " + yesterday)
	if err != nil {
		return "", fmt.Errorf("failed to add stock sheet: %v", err)
	}
	addStockRows(stockSheetYesterday, stockDataYesterday)

	salesSheetYesterday, err := file.AddSheet("Penjualan " + yesterday)
	if err != nil {
		return "", fmt.Errorf("failed to add sales sheet: %v", err)
	}
	addSalesRows(salesSheetYesterday, salesDataYesterday)

	filePath := "report_" + time.Now().Format("2006-01-02") + ".xlsx"
	err = file.Save(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to save Excel file: %v", err)
	}

	return filePath, nil
}

func addStockRows(sheet *xlsx.Sheet, stockData []Stock) {
	header := []string{"NO", "BARANG", "SATUAN", "STOK"}
	headerRow := sheet.AddRow()
	headerStyle := xlsx.NewStyle()
	headerStyle.Font.Bold = true
	headerStyle.Border = *xlsx.NewBorder("thin", "thin", "thin", "thin")
	headerStyle.ApplyBorder = true
	headerStyle.Alignment.Horizontal = "center"

	for _, h := range header {
		cell := headerRow.AddCell()
		cell.Value = h
		cell.SetStyle(headerStyle)
	}

	rowStyle := xlsx.NewStyle()
	rowStyle.Border = *xlsx.NewBorder("thin", "thin", "thin", "thin")
	rowStyle.ApplyBorder = true

	rowCenterStyle := xlsx.NewStyle()
	rowCenterStyle.Border = *xlsx.NewBorder("thin", "thin", "thin", "thin")
	rowCenterStyle.ApplyBorder = true
	headerStyle.Alignment.Horizontal = "center"

	for i, stock := range stockData {
		row := sheet.AddRow()

		numberCell := row.AddCell()
		numberCell.SetValue(i + 1)
		numberCell.SetStyle(rowCenterStyle)

		itemCell := row.AddCell()
		itemCell.SetValue(stock.Item)
		itemCell.SetStyle(rowStyle)

		unitCell := row.AddCell()
		unitCell.SetValue(stock.Unit)
		unitCell.SetStyle(rowStyle)

		stockInHandCell := row.AddCell()
		stockInHandCell.SetFloat(stock.StockInHand)
		stockInHandCell.NumFmt = "#,##0"
		stockInHandCell.SetStyle(rowStyle)
	}

	maxLengths := make([]int, len(header))
	for _, row := range sheet.Rows {
		for colIdx, cell := range row.Cells {
			length := len(cell.String())
			if length > maxLengths[colIdx] {
				maxLengths[colIdx] = length
			}
		}
	}

	for i, length := range maxLengths {
		// Add extra space for padding
		sheet.Col(i).Width = float64(length + 2)
	}
}

func addSalesRows(sheet *xlsx.Sheet, salesData []Sales) {
	header := []string{"NO", "FJ", "PELANGGAN", "SALES", "BARANG", "SATUAN", "JUMLAH", "HARGA", "DISKON", "SUBTOTAL"}
	headerRow := sheet.AddRow()
	headerStyle := xlsx.NewStyle()
	headerStyle.Font.Bold = true
	headerStyle.Border = *xlsx.NewBorder("thin", "thin", "thin", "thin")
	headerStyle.ApplyBorder = true
	headerStyle.Alignment.Horizontal = "center"

	for _, h := range header {
		cell := headerRow.AddCell()
		cell.Value = h
		cell.SetStyle(headerStyle)
	}

	rowStyle := xlsx.NewStyle()
	rowStyle.Border = *xlsx.NewBorder("thin", "thin", "thin", "thin")
	rowStyle.ApplyBorder = true

	rowCenterStyle := xlsx.NewStyle()
	rowCenterStyle.Border = *xlsx.NewBorder("thin", "thin", "thin", "thin")
	rowCenterStyle.ApplyBorder = true
	headerStyle.Alignment.Horizontal = "center"

	var grandTotal float64
	for i, sale := range salesData {
		row := sheet.AddRow()

		numberCell := row.AddCell()
		numberCell.SetValue(i + 1)
		numberCell.SetStyle(rowCenterStyle)

		salesNoCell := row.AddCell()
		salesNoCell.SetValue(sale.SalesNo)
		salesNoCell.SetStyle(rowStyle)

		customerCell := row.AddCell()
		customerCell.SetValue(sale.Customer)
		customerCell.SetStyle(rowStyle)

		salesmanCell := row.AddCell()
		salesmanCell.SetValue(sale.Salesman)
		salesmanCell.SetStyle(rowStyle)

		itemCell := row.AddCell()
		itemCell.SetValue(sale.Item)
		itemCell.SetStyle(rowStyle)

		unitCell := row.AddCell()
		unitCell.SetValue(sale.Unit)
		unitCell.SetStyle(rowStyle)

		quantityCell := row.AddCell()
		quantityCell.SetFloat(sale.Quantity)
		quantityCell.NumFmt = "#,##0"
		quantityCell.SetStyle(rowStyle)

		unitPriceCell := row.AddCell()
		unitPriceCell.SetFloat(sale.UnitPrice)
		unitPriceCell.NumFmt = "#,##0"
		unitPriceCell.SetStyle(rowStyle)

		discountCell := row.AddCell()
		discountCell.SetFloat(sale.Discount)
		discountCell.NumFmt = "#,##0"
		discountCell.SetStyle(rowStyle)

		subtotal := sale.Quantity * (sale.UnitPrice - sale.Discount)
		subtotalCell := row.AddCell()
		subtotalCell.SetFloat(subtotal)
		subtotalCell.NumFmt = "#,##0"
		subtotalCell.SetStyle(rowStyle)

		grandTotal += subtotal
	}

	totalRow := sheet.AddRow()

	grandTotalLabelCell := totalRow.AddCell()
	grandTotalLabelCell.Merge(8, 0)
	grandTotalLabelCell.SetString("Grand Total")

	boldStyle := xlsx.NewStyle()
	boldStyle.Font.Bold = true
	boldStyle.Alignment.Horizontal = "center"
	boldStyle.Border = *xlsx.NewBorder("thin", "thin", "thin", "thin")
	boldStyle.ApplyBorder = true
	grandTotalLabelCell.SetStyle(boldStyle)

	grandTotalCell := totalRow.AddCell()
	grandTotalCell.SetFloat(grandTotal)
	grandTotalCell.NumFmt = "#,##0"

	numberStyle := xlsx.NewStyle()
	numberStyle.Font.Bold = true
	numberStyle.Alignment.Horizontal = "right"
	numberStyle.Border = *xlsx.NewBorder("thin", "thin", "thin", "thin")
	numberStyle.ApplyBorder = true
	grandTotalCell.SetStyle(numberStyle)

	maxLengths := make([]int, len(header))
	for _, row := range sheet.Rows {
		for colIdx, cell := range row.Cells {
			length := len(cell.String())
			if length > maxLengths[colIdx] {
				maxLengths[colIdx] = length
			}
		}
	}

	for i, length := range maxLengths {
		sheet.Col(i).Width = float64(length + 2)
	}
}
