package main

import (
	"bufio"
	"fmt"
	"os"
	"pairingproject/cli"
	"pairingproject/config"
	"strconv"
	"strings"
)

func main() {

	db := config.InitDB()
	reader := bufio.NewReader(os.Stdin)
	defer db.Close()
	for {
		cli.PrintMenu()
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		choice, err := strconv.Atoi(input)

		if err != nil {
			fmt.Println("invalid input")
		}

		if choice == 3 {
			fmt.Println("Selamat Tinggal")
			break
		}
		loginSuccess := cli.HandleChoice(choice, db)

		if loginSuccess {
			for {
				cli.PrintDashboard()
				inputDashboard, _ := reader.ReadString('\n')
				inputDashboard = strings.TrimSpace(inputDashboard)
				dashboardChoice, err := strconv.Atoi(inputDashboard)

				if err != nil {
					fmt.Println("Invalid input")
				}

				cli.HandleChoiceDashboard(dashboardChoice, db)

				if dashboardChoice == 10 {
					break
				}

			}
			break
		}
	}

}
