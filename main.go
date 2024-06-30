package main

import (
	"fmt"
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
	query := sdk.GamesQueryReq{
		Username: client.Username,
		Filter: sdk.GamesQueryFilter{
			ListType: []sdk.GamesListType{
				sdk.UserGamesListTypePlayed,
				sdk.UserGamesListTypePlaying,
			},
		}, PageSort: sdk.UserGamesQueryPageSortWhenAdded,
	}

	games, err := client.GetGamesListFromUserPage(query)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(len(*games))

}
