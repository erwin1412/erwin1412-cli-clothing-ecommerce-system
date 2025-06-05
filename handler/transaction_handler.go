package handler

import (
	"database/sql"
	"fmt"
	"pairingproject/config"
	"pairingproject/entity"

	"github.com/jung-kurt/gofpdf"
	"github.com/xuri/excelize/v2"
)

func CreateTransaction() error {
	db := config.InitDB()

	var userId int
	fmt.Println("=== Create Transaction ===")

	fmt.Println("Creating transaction...")
	fmt.Println("Enter User ID:")
	fmt.Scanln(&userId)

	// Insert transaction with initial total = 0
	insertQuery := "INSERT INTO transactions (userId, total) VALUES (?, ?)"
	result, err := db.Exec(insertQuery, userId, 0)
	if err != nil {
		fmt.Println("Error creating transaction:", err)
		return err
	}

	transactionId, err := result.LastInsertId()
	if err != nil {
		fmt.Println("Error getting last insert ID:", err)
		return err
	}

	var total int
	for {
		var clothId, price, qty, subtotal int

		fmt.Println("Enter Cloth ID:")
		fmt.Scanln(&clothId)

		fmt.Println("Enter Quantity:")
		fmt.Scanln(&qty)

		// Get price from clothes
		queryPrice := "SELECT price FROM cloths WHERE id = ?"
		row := db.QueryRow(queryPrice, clothId)
		err := row.Scan(&price)
		if err != nil {
			if err == sql.ErrNoRows {
				fmt.Println("Cloth not found.")
			} else {
				fmt.Println("Error fetching price:", err)
			}
			continue
		}

		var stock int
		stockQuery := "SELECT stock FROM cloths WHERE id = ?"
		err = db.QueryRow(stockQuery, clothId).Scan(&stock)
		if err != nil {
			fmt.Println("Error checking stock:", err)
			continue
		}

		if stock < qty {
			fmt.Printf("Not enough stock. Available: %d, Requested: %d\n", stock, qty)
			continue
		}

		subtotal = price * qty
		total += subtotal
		detailQuery := "INSERT INTO transaction_details (transactionId, clothId, price, qty, total) VALUES (?, ?, ?, ?, ?)"
		_, err = db.Exec(detailQuery, transactionId, clothId, price, qty, subtotal)
		if err != nil {
			fmt.Println("Error creating transaction detail:", err)
			return err
		}

		// Kurangi stok cloth
		reduceQtyQuery := "UPDATE cloths SET stock = stock - ? WHERE id = ?"
		_, err = db.Exec(reduceQtyQuery, qty, clothId)
		if err != nil {
			fmt.Println("Error updating cloth quantity:", err)
			return err
		}
		var choice string
		fmt.Println("Do you want to add another detail? (yes/no)")
		fmt.Scanln(&choice)
		if choice != "yes" {
			break
		}
	}

	// Update transaction total
	updateQuery := "UPDATE transactions SET total = ? WHERE id = ?"
	_, err = db.Exec(updateQuery, total, transactionId)
	if err != nil {
		fmt.Println("Error updating transaction total:", err)
		return err
	}

	fmt.Println("Transaction successfully created.")
	return nil
}

func ReportTransactionDetail() {
	db := config.InitDB()

	rows, err := db.Query(`
		SELECT t.id, up.name, t.total 
		FROM transactions t
		JOIN users u ON t.userId = u.id
		JOIN user_profiles up ON up.userId = u.id
	`)
	if err != nil {
		fmt.Println("Error querying transactions:", err)
		return
	}
	defer rows.Close()

	var transactions []entity.Transaction

	for rows.Next() {
		var trx entity.Transaction
		err := rows.Scan(&trx.ID, &trx.User, &trx.Total)
		if err != nil {
			fmt.Println("Error scanning transaction:", err)
			return
		}

		// Ambil detail
		detailRows, err := db.Query(`
			SELECT td.id, td.transactionId, c.name, td.price, td.qty, td.total
			FROM transaction_details td
			JOIN cloths c ON td.clothId = c.id
			WHERE td.transactionId = ?
		`, trx.ID)
		if err != nil {
			fmt.Println("Error querying transaction details:", err)
			return
		}

		var details []entity.TransactionDetail
		for detailRows.Next() {
			var d entity.TransactionDetail
			err := detailRows.Scan(&d.ID, &d.TransactionId, &d.ClothName, &d.Price, &d.Qty, &d.Total)
			if err != nil {
				fmt.Println("Error scanning detail:", err)
				return
			}
			details = append(details, d)
		}
		detailRows.Close()

		trx.Details = details
		transactions = append(transactions, trx)
	}

	ExportToPDF(transactions)
	ExportToExcel(transactions)

	fmt.Println("Report exported to PDF and Excel.")
}

func ExportToPDF(transactions []entity.Transaction) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(40, 10, "Transaction Report")
	pdf.Ln(12)

	pdf.SetFont("Arial", "", 10)

	for _, trx := range transactions {
		pdf.Cell(60, 10, fmt.Sprintf("Transaction ID: %d | User: %s | Total: %d", trx.ID, trx.User, trx.Total))
		pdf.Ln(8)

		pdf.SetFont("Arial", "B", 10)
		pdf.Cell(40, 8, "Cloth Name")
		pdf.Cell(20, 8, "Price")
		pdf.Cell(15, 8, "Qty")
		pdf.Cell(20, 8, "Total")
		pdf.Ln(7)

		pdf.SetFont("Arial", "", 10)
		for _, d := range trx.Details {
			pdf.Cell(40, 8, d.ClothName)
			pdf.Cell(20, 8, fmt.Sprintf("%d", d.Price))
			pdf.Cell(15, 8, fmt.Sprintf("%d", d.Qty))
			pdf.Cell(20, 8, fmt.Sprintf("%d", d.Total))
			pdf.Ln(6)
		}

		pdf.Ln(5)
	}

	err := pdf.OutputFileAndClose("transaction_report.pdf")
	if err != nil {
		fmt.Println("Error generating PDF:", err)
	}
}

func ExportToExcel(transactions []entity.Transaction) {
	file := excelize.NewFile()
	sheet := "Sheet1"

	headers := []string{"TransactionID", "User", "Total", "ClothName", "Price", "Qty", "DetailTotal"}

	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		file.SetCellValue(sheet, cell, header)
	}

	row := 2
	for _, trx := range transactions {
		for _, d := range trx.Details {
			file.SetCellValue(sheet, fmt.Sprintf("A%d", row), trx.ID)
			file.SetCellValue(sheet, fmt.Sprintf("B%d", row), trx.User)
			file.SetCellValue(sheet, fmt.Sprintf("C%d", row), trx.Total)
			file.SetCellValue(sheet, fmt.Sprintf("D%d", row), d.ClothName)
			file.SetCellValue(sheet, fmt.Sprintf("E%d", row), d.Price)
			file.SetCellValue(sheet, fmt.Sprintf("F%d", row), d.Qty)
			file.SetCellValue(sheet, fmt.Sprintf("G%d", row), d.Total)
			row++
		}
	}

	if err := file.SaveAs("transaction_report.xlsx"); err != nil {
		fmt.Println("Error saving Excel file:", err)
	}
}

type TransactionDetailInput struct {
	ClothId int
	Qty     int
}

func CreateTransactionTest(db *sql.DB, trx entity.Transaction) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	result, err := tx.Exec("INSERT INTO transactions (userId, total) VALUES (?, ?)", trx.UserID, trx.Total)
	if err != nil {
		tx.Rollback()
		return err
	}

	transactionID, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, detail := range trx.Details {
		_, err := tx.Exec("INSERT INTO transaction_details (transactionId, clothId, price, qty, total) VALUES (?, ?, ?, ?, ?)",
			transactionID, detail.ClothId, detail.Price, detail.Qty, detail.Total,
		)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}
	fmt.Println("Transaction created successfully with ID:", transactionID)
	return nil
}
