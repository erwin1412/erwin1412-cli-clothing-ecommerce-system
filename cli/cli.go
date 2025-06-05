// cli/file.go
package cli

import (
	"database/sql"
	"fmt"
	"pairingproject/handler"
)

func PrintMenu() {
	fmt.Println("\n=== Login / Register ===")
	fmt.Println("1. Login")
	fmt.Println("2. Register")
	fmt.Println("3. Exit")
}

func HandleChoice(choice int, db *sql.DB) bool {
	switch choice {
	case 1:
		return handler.Login()
	case 2:
		handler.Register()
		return false
	default:
		fmt.Println("invalid choice")
		return false
	}
}

func PrintDashboard() {
	fmt.Println("\n=== Dashboard ===")
	fmt.Println("1. List Product")
	fmt.Println("2. Add Product")
	fmt.Println("3. Delete Product")
	fmt.Println("4. Edit Product")
	fmt.Println("5. Favoite Product")
	fmt.Println("6. Transaction")
	fmt.Println("7. Stock Product Report ")
	fmt.Println("8. Transaction Report ")
	fmt.Println("9. User Report ")
	fmt.Println("10. Exit ")
}

func HandleChoiceDashboard(dashboardChoice int, db *sql.DB) {

	switch dashboardChoice {
	case 1:
		handler.ListCloth()
	case 2:
		handler.CreateCloth()
	case 3:
		handler.DeleteCloth()
	case 4:
		handler.EditCloth()
	case 5:
		handler.FavoriteCloth()
	case 6:
		handler.CreateTransaction()
	case 7:
		handler.ReportCloth()
	case 8:
		handler.ReportTransactionDetail()
	case 9:
		handler.UserReport()
	default:
		break
	}
}
