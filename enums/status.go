package enums

type GameStatus int

const (
	GameStatusDefaultNone GameStatus = iota
	GameStatusCompleted
	GameStatusMastered
	GameStatusAbandoned
	GameStatusRetired
	GameStatusShelved
	GameStatusPlayedNothingSpecific
)

func (s GameStatus) String() string {
	switch s {
	case GameStatusCompleted:
		return "completed"
	case GameStatusMastered:
		return "mastered"
	case GameStatusAbandoned:
		return "abandoned"
	case GameStatusRetired:
		return "retired"
	case GameStatusShelved:
		return "shelved"
	case GameStatusPlayedNothingSpecific:
		return "played-nothing-specific"
	case GameStatusDefaultNone:
		return ""
	default:
		return ""
	}
}

// For Logstatus API
func (s GameStatus) Int() int {
	switch s {
	case GameStatusCompleted:
		return 0
	case GameStatusMastered:
		return 1
	case GameStatusAbandoned:
		return 2
	case GameStatusRetired:
		return 3
	case GameStatusShelved:
		return 4
	case GameStatusPlayedNothingSpecific:
		return 5
	default:
		return -1
	}
}
