package game

type aiLevel byte

const (
	humanLikeLevel aiLevel = iota
	beginnerLevel
	perfectLevel
)

func (l aiLevel) String() string {
	switch l {
	case beginnerLevel:
		return "BEGINNER"
	case perfectLevel:
		return "PERFECT"
	default:
		return "HUMAN-LIKE"
	}
}
