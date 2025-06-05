package entity

type Transaction struct {
	ID      int
	UserID  int
	User    string
	Total   int
	Details []TransactionDetail
}
