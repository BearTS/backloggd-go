package sdk

// Gets the user's current playiing
func (sdk *BackloggdSDK) UserGamesPlayed() (*[]Game, error) {
	query := GamesQueryReq{
		Username: sdk.username,
		ListType: UserGamesListTypePlayed,
		PageSort: UserGamesQueryPageSortWhenAdded,
	}

	games, err := sdk.GetGamesListFromUserPage(query)
	if err != nil {
		return nil, err
	}

	return games, nil
}
