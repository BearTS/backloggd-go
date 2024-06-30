package main

import (
	"log"

	"github.com/BearTS/backloggd-go/sdk"
)

func main() {
	// Create a new instance of the BackloggdSDK
	client, err := sdk.NewBackloggdSDK("bearts", "test")
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

	req := sdk.LogReq{
		Slug:    "metroid-dread",
		LogType: sdk.Played,
		Enable:  true,
	}

	err = client.LogGame(req)
	if err != nil {
		log.Fatal(err)
	}

	req2 := sdk.LogStatusReq{
		Slug:   "metroid-dread",
		Status: sdk.Mastered,
	}

	err = client.LogStatus(req2)
	if err != nil {
		log.Fatal(err)
	}

}
