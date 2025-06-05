package handler

import (
	"bufio"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"pairingproject/config"
	"pairingproject/entity"
	"strconv"
	"strings"

	"github.com/jung-kurt/gofpdf"
	"github.com/xuri/excelize/v2"
)

func CreateCloth() {
	input := bufio.NewReader(os.Stdin)

	fmt.Println("Enter Cloth Name:")
	name, _ := input.ReadString('\n')
	name = strings.TrimSpace(name)
	fmt.Println("Enter Cloth Price:")
	price, _ := input.ReadString('\n')
	price = strings.TrimSpace(price)

	fmt.Println("Enter Cloth Size:")
	size, _ := input.ReadString('\n')
	size = strings.TrimSpace(size)
	fmt.Println("Enter Cloth Color:")
	color, _ := input.ReadString('\n')
	color = strings.TrimSpace(color)
	fmt.Println("Enter Cloth Stock:")
	stock, _ := input.ReadString('\n')
	stock = strings.TrimSpace(stock)
	fmt.Println("Enter Cloth Description:")
	description, _ := input.ReadString('\n')
	description = strings.TrimSpace(description)

	query := "INSERT INTO cloths (name, price, size, color, stock, description) VALUES (?, ?, ?, ?, ?, ?)"
	_, err := config.InitDB().Exec(query, name, price, size, color, stock, description)
	if err != nil {
		fmt.Println("Error creating cloth:", err)
		return
	}
	fmt.Println("Cloth created successfully!")

}

func ListCloth() {
	query := "SELECT id, name, price, size, color, stock, description FROM cloths"
	rows, err := config.InitDB().Query(query)
	if err != nil {
		fmt.Println("Error retrieving cloths:", err)
		return
	}
	defer rows.Close()

	fmt.Println("Cloth List:")
	for rows.Next() {
		var id int
		var name, size, color, stock, description string
		var price float64
		if err := rows.Scan(&id, &name, &price, &size, &color, &stock, &description); err != nil {
			fmt.Println("Error scanning row:", err)
			return
		}
		fmt.Printf("ID: %d, Name: %s, Price: %.2f, Size: %s, Color: %s, Stock: %s, Description: %s\n", id, name, price, size, color, stock, description)
	}
}

func DeleteCloth() {

	var id int
	fmt.Println("Enter Cloth ID to delete:")
	fmt.Scan(&id)

	query := "DELETE FROM cloths WHERE id = ?"
	_, err := config.InitDB().Exec(query, id)
	if err != nil {
		fmt.Println("Error deleting cloth:", err)
		return
	}
	fmt.Println("Cloth deleted successfully!")
}

func EditCloth() {
	input := bufio.NewReader(os.Stdin)
	var id int
	fmt.Println("Enter Cloth ID to edit:")
	fmt.Scan(&id)

	// Flush newline
	input.ReadString('\n')

	fmt.Println("Enter New Cloth Name:")
	name, _ := input.ReadString('\n')
	name = strings.TrimSpace(name)

	fmt.Println("Enter New Cloth Price:")
	priceStr, _ := input.ReadString('\n')
	priceStr = strings.TrimSpace(priceStr)
	price, err := strconv.Atoi(priceStr)
	if err != nil {
		fmt.Println("Invalid price:", err)
		return
	}

	fmt.Println("Enter New Cloth Size:")
	size, _ := input.ReadString('\n')
	size = strings.TrimSpace(size)

	fmt.Println("Enter New Cloth Color:")
	color, _ := input.ReadString('\n')
	color = strings.TrimSpace(color)

	fmt.Println("Enter New Cloth Stock:")
	stockStr, _ := input.ReadString('\n')
	stockStr = strings.TrimSpace(stockStr)
	stock, err := strconv.Atoi(stockStr)
	if err != nil {
		fmt.Println("Invalid stock:", err)
		return
	}

	fmt.Println("Enter New Cloth Description:")
	description, _ := input.ReadString('\n')
	description = strings.TrimSpace(description)

	query := "UPDATE cloths SET name = ?, price = ?, size = ?, color = ?, stock = ?, description = ? WHERE id = ?"
	_, err = config.InitDB().Exec(query, name, price, size, color, stock, description, id)
	if err != nil {
		fmt.Println("Error updating cloth:", err)
		return
	}
	fmt.Println("Cloth updated successfully!")
}

func FavoriteCloth() {
	var id int
	fmt.Println("Enter Cloth ID to favorite:")
	fmt.Scan(&id)

	var userId int
	fmt.Println("Enter User ID:")
	fmt.Scan(&userId)

	db := config.InitDB()

	// Check if userId exists in users table
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE id=?)", userId).Scan(&exists)
	if err != nil {
		fmt.Println("Error checking user ID:", err)
		return
	}
	if !exists {
		fmt.Println("User ID does not exist.")
		return
	}

	query := "INSERT INTO user_cloth_favorites (userId, clothId) VALUES (?, ?)"
	_, err = db.Exec(query, userId, id)
	if err != nil {
		fmt.Println("Error favoriting cloth:", err)
		return
	}

	fmt.Println("Cloth favorited successfully!")
}

func ReportCloth() {
	query := "SELECT id, name , stock FROM cloths"
	rows, err := config.InitDB().Query(query)
	if err != nil {
		fmt.Println("Error retrieving cloths:", err)
		return
	}
	defer rows.Close()

	var cloths []entity.Cloth
	fmt.Println("Stock Product Report :")
	for rows.Next() {
		var id int
		var name, stock string
		if err := rows.Scan(&id, &name, &stock); err != nil {
			fmt.Println("Error scanning row:", err)
			return
		}
		cloths = append(cloths, entity.Cloth{
			ID:    id,
			Name:  name,
			Stock: stock,
		})
		fmt.Printf("ID: %d, Name: %s, Stock: %s \n", id, name, stock)
	}
	ExportExcel(cloths)
	ExportPDF(cloths)
}

func ExportExcel(cloths []entity.Cloth) {
	f := excelize.NewFile()
	sheet := "Sheet1"

	f.SetCellValue(sheet, "A1", "ID")
	f.SetCellValue(sheet, "B1", "Name")
	f.SetCellValue(sheet, "C1", "Stock")

	for i, c := range cloths {
		row := i + 2
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), c.ID)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), c.Name)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), c.Stock)
	}

	err := f.SaveAs("stock_report.xlsx")
	if err != nil {
		fmt.Println("Error saving Excel file:", err)
		return
	}
	fmt.Println("Excel report saved as stock_report.xlsx")
}

func ExportPDF(cloths []entity.Cloth) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(40, 10, "Stock Product Report")
	pdf.Ln(12)

	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(20, 10, "ID")
	pdf.Cell(80, 10, "Name")
	pdf.Cell(30, 10, "Stock")
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 12)
	for _, c := range cloths {
		pdf.Cell(20, 10, fmt.Sprintf("%d", c.ID))
		pdf.Cell(80, 10, c.Name)
		pdf.Cell(30, 10, fmt.Sprintf("%s", c.Stock))
		pdf.Ln(8)
	}

	err := pdf.OutputFileAndClose("stock_report.pdf")
	if err != nil {
		fmt.Println("Error saving PDF file:", err)
		return
	}
	fmt.Println("PDF report saved as stock_report.pdf")
}

func CreateClothTest(db *sql.DB, name string, price int, size string, color string, stock int, description string) error {
	query := "INSERT INTO cloths (name, price, size, color, stock, description) VALUES (?, ?, ?, ?, ?, ?)"
	_, err := db.Exec(query, name, price, size, color, stock, description)
	return err
}

func FavoriteClothTest(db *sql.DB, userId int, clothId int) error {
	var exists bool

	// cek userId
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE id=?)", userId).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("user ID tidak ditemukan")
	}

	// cek clothId
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM cloths WHERE id=?)", clothId).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("cloth ID tidak ditemukan")
	}

	// insert favorite
	_, err = db.Exec("INSERT INTO user_cloth_favorites (userId, clothId) VALUES (?, ?)", userId, clothId)
	if err != nil {
		return err
	}

	return nil
}

// ListCloth mengambil semua cloth dari DB dan mengembalikan slice hasilnya
func ListClothTest(db *sql.DB) ([]entity.Cloth, error) {
	query := "SELECT id, name, price, size, color, stock, description FROM cloths"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cloths []entity.Cloth
	for rows.Next() {
		var c entity.Cloth
		err := rows.Scan(&c.ID, &c.Name, &c.Price, &c.Size, &c.Color, &c.Stock, &c.Description)
		if err != nil {
			return nil, err
		}
		cloths = append(cloths, c)
	}
	return cloths, nil
}

// DeleteCloth menerima db dan id, menghapus data cloth berdasar id
func DeleteClothTest(db *sql.DB, id int) error {
	query := "DELETE FROM cloths WHERE id = ?"
	_, err := db.Exec(query, id)
	return err
}

func EditClothTest(db *sql.DB, id int, name string, price int, size, color string, stock int, description string) error {
	query := "UPDATE cloths SET name = ?, price = ?, size = ?, color = ?, stock = ?, description = ? WHERE id = ?"
	_, err := db.Exec(query, name, price, size, color, stock, description, id)
	return err
}
