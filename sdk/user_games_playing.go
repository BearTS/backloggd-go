package sdk

// Gets the user's current playiing
func (sdk *BackloggdSDK) UserGamesPlaying() (*[]Game, error) {
	query := GamesQueryReq{
		Username: sdk.Username,
		Filter: GamesQueryFilter{
			ListType: []GamesListType{
				UserGamesListTypePlaying,
			},
		}, PageSort: UserGamesQueryPageSortWhenAdded,
	}

	games, err := sdk.GetGamesListFromUserPage(query)
	if err != nil {
		return nil, err
	}

	return games, nil
}
