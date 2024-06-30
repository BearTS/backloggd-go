package main

import (
	"log"

	"github.com/BearTS/backloggd-go/sdk"
)

func main() {
	// Create a new instance of the BackloggdSDK
	sdk, err := sdk.NewBackloggdSDK()
	if err != nil {
		log.Fatal(err)
	}

	// Example usage: Login with credentials
	err = sdk.Login("username", "password")
	if err != nil {
		log.Fatal(err)
	}

	// Example usage: Export cookies to a JSON file
	err = sdk.ExportCookies()
	if err != nil {
		log.Fatal(err)
	}
}
