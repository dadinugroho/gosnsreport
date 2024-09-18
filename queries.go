package main

import (
    "fmt"
)

type Stock struct {
    Item        string
    Unit        string
    StockInHand float64
}

type Sales struct {
    SalesNo     string
    Customer    string
    Salesman    string
    Item        string
    Unit        string
    Quantity    float64
    UnitPrice   float64
    Discount    float64
}

func fetchStockReport(journalDate string) ([]Stock, error) {
    db, err := connectToDB()
    if err != nil {
        return nil, err
    }
    defer db.Close()

    query := `
        SELECT i.name, i.baseUnit, SUM(j.type * j.quantity) 
        FROM journal j 
        LEFT JOIN item i ON i.id = j.itemFk 
        WHERE j.shopFk = 1 AND j.accountFk = 2 AND j.locationFk = 11 
        AND j.journalDate <= ? 
        GROUP BY j.itemFk
    `

    rows, err := db.Query(query, journalDate)
    if err != nil {
        return nil, fmt.Errorf("failed to execute stock query: %v", err)
    }
    defer rows.Close()

    var stocks []Stock
    for rows.Next() {
        var stock Stock
        if err := rows.Scan(&stock.Item, &stock.Unit, &stock.StockInHand); err != nil {
            return nil, err
        }
        stocks = append(stocks, stock)
    }

    return stocks, nil
}

func fetchSalesReport(formDate string) ([]Sales, error) {
    db, err := connectToDB()
    if err != nil {
        return nil, err
    }
    defer db.Close()

    query := `
        SELECT h.headerNo, p.name, s.name, i.name, d.unit, d.quantity, d.unitPrice, d.discount FROM form_detail d
        LEFT JOIN form_header h ON h.id = d.headerFk
        LEFT JOIN partner p ON p.id = h.partnerFk
        LEFT JOIN partner s ON s.id = h.salesmanFk
        LEFT JOIN item i ON i.id = d.itemFk
        WHERE 
            h.formType = 5 AND 
            h.postingStatus = 1 AND
            h.formDate = ? AND
            h.shopFk = 1
    `
    rows, err := db.Query(query, formDate)
    if err != nil {
        return nil, fmt.Errorf("failed to execute sales query: %v", err)
    }
    defer rows.Close()

    var sales []Sales
    for rows.Next() {
        var sale Sales
        if err := rows.Scan(&sale.SalesNo, &sale.Customer, &sale.Salesman,  &sale.Item,  &sale.Unit,  &sale.Quantity,  &sale.UnitPrice,  &sale.Discount); err != nil {
            return nil, err
        }
        sales = append(sales, sale)
    }

    return sales, nil
}
