package enums

type GameCategory int

const (
	GameCategoryDefaultNone GameCategory = iota
	GameCategoryMainGame
	GameCategoryDLCAddon
	GameCategoryExpansion
	GameCategoryBundle
	GameCategoryStandaloneExpansion
	GameCategoryMod
	GameCategoryEpisode
	GameCategorySeason
	GameCategoryRemake
	GameCategoryRemaster
	GameCategoryExpandedGame
	GameCategoryPort
	GameCategoryFork
	GameCategoryPackAddon
	GameCategoryGameUpdate
)

func (c GameCategory) String() string {
	switch c {
	case GameCategoryMainGame:
		return "main_game"
	case GameCategoryDLCAddon:
		return "dlc_addon"
	case GameCategoryExpansion:
		return "expansion"
	case GameCategoryBundle:
		return "bundle"
	case GameCategoryStandaloneExpansion:
		return "standalone_expansion"
	case GameCategoryMod:
		return "mod"
	case GameCategoryEpisode:
		return "episode"
	case GameCategorySeason:
		return "season"
	case GameCategoryRemake:
		return "remake"
	case GameCategoryRemaster:
		return "remaster"
	case GameCategoryExpandedGame:
		return "expanded_game"
	case GameCategoryPort:
		return "port"
	case GameCategoryFork:
		return "fork"
	case GameCategoryPackAddon:
		return "pack"
	case GameCategoryGameUpdate:
		return "game_update"
	case GameCategoryDefaultNone:
		return ""
	default:
		return GameCategoryDefaultNone.String()
	}
}
