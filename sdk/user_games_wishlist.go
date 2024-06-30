package sdk

type Game struct {
	Name string `json:"name"`
	Id   int    `json:"id"`
	Link string `json:"link"`
}

// Gets the user's games wishlist
func (sdk *BackloggdSDK) UserGamesWishlist() (*[]Game, error) {
	query := GamesQueryReq{
		Username: sdk.Username,
		Filter: GamesQueryFilter{
			ListType: []GamesListType{
				UserGamesListTypeWishlist,
			},
		}, PageSort: UserGamesQueryPageSortWhenAdded,
	}

	games, err := sdk.GetGamesListFromUserPage(query)
	if err != nil {
		return nil, err
	}

	return games, nil
}
