package handler

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"pairingproject/config"
	"pairingproject/entity"
	"strings"

	"github.com/jung-kurt/gofpdf"
	"github.com/xuri/excelize/v2"
	"golang.org/x/crypto/bcrypt"
)

func Login() bool {

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("---- Login Page ----")
	fmt.Println("Enter Your Email")
	email, _ := reader.ReadString('\n')
	email = strings.TrimSpace(email)
	fmt.Println("Enter Your Password:")
	password, _ := reader.ReadString('\n')
	password = strings.TrimSpace(password)

	var storedPassword []byte
	query := "SELECT password FROM users WHERE email = ?"
	err := config.InitDB().QueryRow(query, email).Scan(&storedPassword)
	if err != nil {
		fmt.Println("Invalid email or password")
		return false
	}

	err = bcrypt.CompareHashAndPassword(storedPassword, []byte(password))
	if err != nil {
		fmt.Println("Invalid email or password")
		return false
	}

	fmt.Println("Login successful!")
	return true
}

func Register() {

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("---- Register Page ----")
	fmt.Println("Enter Your Name : ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	fmt.Println("Enter Your Email : ")
	email, _ := reader.ReadString('\n')
	email = strings.TrimSpace(email)

	fmt.Println("Enter Your Password : ")
	password, _ := reader.ReadString('\n')
	password = strings.TrimSpace(password)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("Error hashing password:", err)
		return
	}

	query := "INSERT INTO users(email , password) VALUES (? , ?)"

	_, errs := config.InitDB().Exec(query, email, hashedPassword)
	if errs != nil {
		fmt.Println("Fail : ", errs)
		return
	}

	var id int
	querySelect := "SELECT id FROM users WHERE email = ?"
	errSelect := config.InitDB().QueryRow(querySelect, email).Scan(&id)

	if errSelect != nil {
		fmt.Println("Fail : ", errSelect)
	}

	queryInsert := "INSERT INTO user_profiles(userId , name) VALUES (?,?)"
	_, errInsert := config.InitDB().Exec(queryInsert, id, name)

	if errInsert != nil {
		fmt.Println("Fail")
	}

	fmt.Println("Register Success")
}

func UserReport() {
	db := config.InitDB()
	defer db.Close()

	query := `
		SELECT 
			u.id AS user_id,
			up.name,
			u.email,
			COALESCE(COUNT(t.id), 0) AS total_transaksi,
			(
				SELECT COUNT(*) 
				FROM user_cloth_favorites ucf 
				WHERE ucf.userId = u.id
			) AS total_favorite
		FROM users u
		JOIN user_profiles up ON up.userId = u.id
		LEFT JOIN transactions t ON t.userId = u.id
		GROUP BY u.id, up.name, u.email
	`

	rows, err := db.Query(query)
	if err != nil {
		fmt.Println("Error executing query:", err)
		return
	}
	defer rows.Close()

	var reports []entity.UserReport

	for rows.Next() {
		var r entity.UserReport
		if err := rows.Scan(&r.UserID, &r.Name, &r.Email, &r.TotalTransaksi, &r.TotalFavorite); err != nil {
			fmt.Println("Error scanning row:", err)
			return
		}
		reports = append(reports, r)
	}

	fmt.Println("=== User Report ===")
	for _, r := range reports {
		fmt.Printf("ID: %d | Name: %s | Email: %s | Total Transaksi: %d | Total Favorite Cloths: %d\n",
			r.UserID, r.Name, r.Email, r.TotalTransaksi, r.TotalFavorite)
	}

	// Generate PDF
	err = GenerateUserReportPDF(reports, "user_report.pdf")
	if err != nil {
		fmt.Println("Failed to generate PDF:", err)
	} else {
		fmt.Println("PDF generated: user_report.pdf")
	}

	// Generate Excel
	err = GenerateUserReportExcel(reports, "user_report.xlsx")
	if err != nil {
		fmt.Println("Failed to generate Excel:", err)
	} else {
		fmt.Println("Excel generated: user_report.xlsx")
	}
}

func ExportUsersToPDF(reports []entity.UserReport) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(40, 10, "User Report")
	pdf.Ln(12)

	pdf.SetFont("Arial", "B", 10)
	headers := []string{"ID", "Name", "Email", "Total Transaksi", "Total Favorite"}
	for _, h := range headers {
		pdf.CellFormat(40, 10, h, "1", 0, "C", false, 0, "")
	}
	pdf.Ln(-1)

	pdf.SetFont("Arial", "", 10)
	for _, r := range reports {
		pdf.CellFormat(40, 10, fmt.Sprintf("%d", r.UserID), "1", 0, "", false, 0, "")
		pdf.CellFormat(40, 10, r.Name, "1", 0, "", false, 0, "")
		pdf.CellFormat(40, 10, r.Email, "1", 0, "", false, 0, "")
		pdf.CellFormat(40, 10, fmt.Sprintf("%d", r.TotalTransaksi), "1", 0, "", false, 0, "")
		pdf.CellFormat(40, 10, fmt.Sprintf("%d", r.TotalFavorite), "1", 0, "", false, 0, "")
		pdf.Ln(-1)
	}

	err := pdf.OutputFileAndClose("user_report.pdf")
	if err != nil {
		fmt.Println("Error exporting PDF:", err)
	}
}

func GenerateUserReportPDF(reports []entity.UserReport, filename string) error {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(40, 10, "User Report")
	pdf.Ln(12)

	pdf.SetFont("Arial", "B", 10)
	// Header tabel
	headers := []string{"ID", "Name", "Email", "Total Transaksi", "Total Favorite"}
	for _, header := range headers {
		pdf.CellFormat(38, 7, header, "1", 0, "C", false, 0, "")
	}
	pdf.Ln(-1)

	pdf.SetFont("Arial", "", 10)
	for _, r := range reports {
		pdf.CellFormat(38, 7, fmt.Sprintf("%d", r.UserID), "1", 0, "C", false, 0, "")
		pdf.CellFormat(38, 7, r.Name, "1", 0, "L", false, 0, "")
		pdf.CellFormat(38, 7, r.Email, "1", 0, "L", false, 0, "")
		pdf.CellFormat(38, 7, fmt.Sprintf("%d", r.TotalTransaksi), "1", 0, "C", false, 0, "")
		pdf.CellFormat(38, 7, fmt.Sprintf("%d", r.TotalFavorite), "1", 0, "C", false, 0, "")
		pdf.Ln(-1)
	}

	return pdf.OutputFileAndClose(filename)
}

func GenerateUserReportExcel(reports []entity.UserReport, filename string) error {
	f := excelize.NewFile()
	sheet := "Sheet1"

	// Header
	headers := []string{"ID", "Name", "Email", "Total Transaksi", "Total Favorite"}
	for i, h := range headers {
		col, _ := excelize.ColumnNumberToName(i + 1)
		f.SetCellValue(sheet, col+"1", h)
	}

	// Data
	for i, r := range reports {
		row := i + 2
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), r.UserID)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), r.Name)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), r.Email)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", row), r.TotalTransaksi)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", row), r.TotalFavorite)
	}

	// Save file
	if err := f.SaveAs(filename); err != nil {
		return err
	}
	return nil
}

func LoginTest(db *sql.DB, email, password string) bool {
	var storedPassword []byte
	query := "SELECT password FROM users WHERE email = ?"
	err := db.QueryRow(query, email).Scan(&storedPassword)
	if err != nil {
		fmt.Println("Invalid email or password")
		return false
	}

	err = bcrypt.CompareHashAndPassword(storedPassword, []byte(password))
	if err != nil {
		fmt.Println("Invalid email or password")
		return false
	}

	fmt.Println("Login successful!")
	return true
}

func RegisterTest(db *sql.DB, name, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	query := "INSERT INTO users(email , password) VALUES (? , ?)"
	_, errs := db.Exec(query, email, hashedPassword)
	if errs != nil {
		return errs
	}

	var id int
	querySelect := "SELECT id FROM users WHERE email = ?"
	errSelect := db.QueryRow(querySelect, email).Scan(&id)
	if errSelect != nil {
		return errSelect
	}

	queryInsert := "INSERT INTO user_profiles(userId , name) VALUES (?,?)"
	_, errInsert := db.Exec(queryInsert, id, name)
	if errInsert != nil {
		return errInsert
	}

	fmt.Println("Register Success")
	return nil
}
