package entity

type User struct {
	id       int
	email    string
	password string
}

type UserReport struct {
	UserID         int
	Name           string
	Email          string
	TotalTransaksi int
	TotalFavorite  int
}
