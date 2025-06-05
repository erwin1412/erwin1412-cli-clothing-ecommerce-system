package handler

import (
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestFavoriteCloth(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to open mock db: %v", err)
	}
	defer db.Close()

	userId := 1
	clothId := 10

	// Mock cek userId ada
	mock.ExpectQuery("SELECT EXISTS\\(SELECT 1 FROM users WHERE id=\\?\\)").
		WithArgs(userId).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

	// Mock cek clothId ada
	mock.ExpectQuery("SELECT EXISTS\\(SELECT 1 FROM cloths WHERE id=\\?\\)").
		WithArgs(clothId).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

	// Mock insert favorite sukses
	mock.ExpectExec("INSERT INTO user_cloth_favorites").
		WithArgs(userId, clothId).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = FavoriteClothTest(db, userId, clothId)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Pastikan semua mock terpanggil
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}

	fmt.Println("FavoriteClothTest passed")
}

func TestCreateCloth(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %s", err)
	}
	defer db.Close()

	// Setup expected query dan args sesuai data dummy
	mock.ExpectExec("INSERT INTO cloths").
		WithArgs("Kain Batik", 150000, "M", "Biru", 20, "Kain batik bagus").
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Panggil fungsi dengan data dummy
	err = CreateClothTest(db, "Kain Batik", 150000, "M", "Biru", 20, "Kain batik bagus")
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	// Cek semua expectation sudah terpenuhi
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	fmt.Println("CreateClothTest passed")
}

func TestListCloth(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "price", "size", "color", "stock", "description"}).
		AddRow(1, "Batik", 150000, "M", "Merah", 10, "Kain batik merah").
		AddRow(2, "Jeans", 200000, "L", "Biru", 5, "Celana jeans biru")

	mock.ExpectQuery("SELECT id, name, price, size, color, stock, description FROM cloths").
		WillReturnRows(rows)

	clothes, err := ListClothTest(db)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(clothes) != 2 {
		t.Errorf("expected 2 rows, got %d", len(clothes))
	}

	mock.ExpectationsWereMet()

	fmt.Println("ListClothTest passed")
}

func TestDeleteCloth(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	mock.ExpectExec("DELETE FROM cloths WHERE id = ?").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := DeleteClothTest(db, 1)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	mock.ExpectationsWereMet()

	fmt.Println("DeleteClothTest passed")
}

func TestEditCloth(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %v", err)
	}
	defer db.Close()

	mock.ExpectExec(`UPDATE cloths SET name = \?, price = \?, size = \?, color = \?, stock = \?, description = \? WHERE id = \?`).
		WithArgs("Batik Baru", 180000, "L", "Hijau", 15, "Batik hijau", 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = EditClothTest(db, 1, "Batik Baru", 180000, "L", "Hijau", 15, "Batik hijau")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	fmt.Println("EditClothTest passed")
}
