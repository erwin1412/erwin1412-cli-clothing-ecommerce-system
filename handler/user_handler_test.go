package handler

import (
	"fmt"
	"os"
	"pairingproject/entity"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"golang.org/x/crypto/bcrypt"
)

func TestGenerateUserReportPDF(t *testing.T) {
	reports := []entity.UserReport{
		{UserID: 1, Name: "Alice", Email: "alice@example.com", TotalTransaksi: 2, TotalFavorite: 1},
	}

	filename := "test_user_report.pdf"
	err := GenerateUserReportPDF(reports, filename)
	if err != nil {
		t.Errorf("Failed to generate PDF: %v", err)
	}

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		t.Errorf("PDF file was not created: %s", filename)
	}

	// Cleanup
	_ = os.Remove(filename)
	fmt.Println("TestGenerateUserReportPDF passed")
}

func TestGenerateUserReportExcel(t *testing.T) {
	reports := []entity.UserReport{
		{UserID: 1, Name: "Alice", Email: "alice@example.com", TotalTransaksi: 2, TotalFavorite: 1},
		{UserID: 2, Name: "Alicea", Email: "alice@example.com", TotalTransaksi: 2, TotalFavorite: 1},
	}

	filename := "test_user_report.xlsx"
	err := GenerateUserReportExcel(reports, filename)
	if err != nil {
		t.Errorf("Failed to generate Excel: %v", err)
	}

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		t.Errorf("Excel file was not created: %s", filename)
	}

	// Cleanup
	_ = os.Remove(filename)

	fmt.Println("TestGenerateUserReportExcel passed")
}

func TestRegisterAndLogin(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %v", err)
	}
	defer db.Close()

	name := "Test User"
	email := "testing@gmail.com"
	password := "test123"

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	// Mock Register
	mock.ExpectExec("INSERT INTO users").
		WithArgs(email, sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectQuery("SELECT id FROM users WHERE email = ?").
		WithArgs(email).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectExec("INSERT INTO user_profiles").
		WithArgs(1, name).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = RegisterTest(db, name, email, password)
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	// Mock Login
	rows := sqlmock.NewRows([]string{"password"}).AddRow(hashedPassword)
	mock.ExpectQuery("SELECT password FROM users WHERE email = ?").
		WithArgs(email).
		WillReturnRows(rows)

	success := LoginTest(db, email, password)
	if !success {
		t.Fatalf("Login failed for registered user")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
