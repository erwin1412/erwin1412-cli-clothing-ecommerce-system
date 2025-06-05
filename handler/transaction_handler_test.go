package handler

import (
	"testing"

	"pairingproject/entity"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestCreateTransaction(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error initializing mock db: %s", err)
	}
	defer db.Close()

	// Mulai transaksi mock
	mock.ExpectBegin()

	// Expect insert ke transactions dan return last insert id 10
	mock.ExpectExec("INSERT INTO transactions").
		WithArgs(1, 300000).
		WillReturnResult(sqlmock.NewResult(10, 1)) // lastInsertId=10, rowsAffected=1

	// Expect insert ke transaction_details untuk setiap detail
	mock.ExpectExec("INSERT INTO transaction_details").
		WithArgs(10, 1, 100000, 2, 200000).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec("INSERT INTO transaction_details").
		WithArgs(10, 2, 100000, 1, 100000).
		WillReturnResult(sqlmock.NewResult(2, 1))

	// Commit transaksi
	mock.ExpectCommit()

	transaction := entity.Transaction{
		UserID: 1,
		Total:  300000,
		Details: []entity.TransactionDetail{
			{ClothId: 1, Price: 100000, Qty: 2, Total: 200000},
			{ClothId: 2, Price: 100000, Qty: 1, Total: 100000},
		},
	}

	err = CreateTransactionTest(db, transaction)
	if err != nil {
		t.Errorf("error: %v", err)
	}

	// Pastikan semua ekspektasi terpenuhi
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}

}
