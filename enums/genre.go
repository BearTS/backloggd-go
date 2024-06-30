package enums

type GameGenre int

const (
	GameGenreDefaultNone GameGenre = iota
	GameGenreAdventure
	GameGenreArcade
	GameGenreBrawler
	GameGenreCardAndBoardGame
	GameGenreFighting
	GameGenreIndie
	GameGenreMOBA
	GameGenreMusic
	GameGenrePinball
	GameGenrePlatform
	GameGenrePointAndClick
	GameGenrePuzzle
	GameGenreQuizTrivia
	GameGenreRacing
	GameGenreRealTimeStrategy
	GameGenreRolePlaying
	GameGenreShooter
	GameGenreSimulator
	GameGenreSport
	GameGenreStrategy
	GameGenreTactical
	GameGenreTurnBasedStrategy
	GameGenreVisualNovel
)

func (g GameGenre) String() string {
	switch g {
	case GameGenreAdventure:
		return "adventure"
	case GameGenreArcade:
		return "arcade"
	case GameGenreBrawler:
		return "hack-and-slash-beat-em-up"
	case GameGenreCardAndBoardGame:
		return "card-and-board-game"
	case GameGenreFighting:
		return "fighting"
	case GameGenreIndie:
		return "indie"
	case GameGenreMOBA:
		return "moba"
	case GameGenreMusic:
		return "music"
	case GameGenrePinball:
		return "pinball"
	case GameGenrePlatform:
		return "platform"
	case GameGenrePointAndClick:
		return "point-and-click"
	case GameGenrePuzzle:
		return "puzzle"
	case GameGenreQuizTrivia:
		return "quiz-trivia"
	case GameGenreRacing:
		return "racing"
	case GameGenreRealTimeStrategy:
		return "real-time-strategy-rts"
	case GameGenreRolePlaying:
		return "role-playing-rpg"
	case GameGenreShooter:
		return "shooter"
	case GameGenreSimulator:
		return "simulator"
	case GameGenreSport:
		return "sport"
	case GameGenreStrategy:
		return "strategy"
	case GameGenreTactical:
		return "tactical"
	case GameGenreTurnBasedStrategy:
		return "turn-based-strategy-tbs"
	case GameGenreVisualNovel:
		return "visual-novel"
	default:
		return GameGenreDefaultNone.String()
	}
}
