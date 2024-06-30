package main

import (
	"fmt"
	"log"

	"github.com/BearTS/backloggd-go/sdk"
)

func main() {
	// Create a new instance of the BackloggdSDK
	client, err := sdk.NewBackloggdSDK("test", "test")
	if err != nil {
		log.Fatal(err)
	}

	// Example usage: Login with credentials
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// var user sdk.User
	// user.Bio = ptr.String("Sample test from api")

	// err = client.UpdateUser(user)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// data, err := client.Autocomplete("Spiderman")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	data, err := client.Wishlist()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(len(*data))
}
