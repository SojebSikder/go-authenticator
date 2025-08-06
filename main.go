package main

import (
	"fmt"
	"log"
	"os"
	"sojebsikder/qrcodeutil"
	"sojebsikder/storage"
	"sojebsikder/totp"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	pass := os.Getenv("AUTH_MASTER_PASSPHRASE")

	if pass == "" {
		fmt.Println("Set AUTH_MASTER_PASSPHRASE environment variable first.")
		return
	}
	passphrase := []byte(pass)

	st, err := storage.Load(passphrase)
	if err != nil {
		fmt.Println("Failed to load store:", err)
		return
	}

	for {
		fmt.Println("\n1) Add account (manual secret)")
		fmt.Println("2) Add account (QR image)")
		fmt.Println("3) List accounts & codes")
		fmt.Println("4) Validate code")
		fmt.Println("0) Exit")
		fmt.Print("Select: ")

		var choice int
		fmt.Scan(&choice)

		switch choice {
		case 1:
			fmt.Print("Account name: ")
			var name string
			fmt.Scan(&name)
			fmt.Print("Secret: ")
			var secret string
			fmt.Scan(&secret)

			st.Accounts = append(st.Accounts, storage.Account{Name: name, Secret: secret})
			_ = storage.Save(st, passphrase)
			fmt.Println("Added.")

		case 2:
			fmt.Print("QR image path: ")
			var path string
			fmt.Scan(&path)
			content, err := qrcodeutil.DecodeFromFile(path)
			if err != nil {
				fmt.Println("Decode error:", err)
				continue
			}
			name, secret, err := qrcodeutil.ParseOtpauth(content)
			if err != nil {
				fmt.Println("Parse error:", err)
				continue
			}
			st.Accounts = append(st.Accounts, storage.Account{Name: name, Secret: secret})
			_ = storage.Save(st, passphrase)
			fmt.Println("Added from QR.")

		case 3:
			for _, acc := range st.Accounts {
				code, err := totp.CurrentCode(acc.Secret)
				if err != nil {
					fmt.Printf("%s\t(error)\n", acc.Name)
				} else {
					fmt.Printf("%s\t%s\n", acc.Name, code)
				}
			}

		case 4:
			fmt.Print("Account name: ")
			var name string
			fmt.Scan(&name)
			fmt.Print("Code: ")
			var code string
			fmt.Scan(&code)

			for _, acc := range st.Accounts {
				if acc.Name == name {
					if totp.ValidateCode(acc.Secret, code) {
						fmt.Println("Valid")
					} else {
						fmt.Println("Invalid")
					}
				}
			}

		case 0:
			return

		default:
			fmt.Println("Invalid choice")
		}
	}
}
