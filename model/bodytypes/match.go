package bodytypes

type ConfidenceLevel int

const (
	Negative ConfidenceLevel = iota
	None
	Assigned
	DistantlyPossible
	Possible
	Likely
	Perfect
)

func (c ConfidenceLevel) String() string {
	switch c {
	case Negative:
		return "negative"
	case None:
		return "none"
	case Assigned:
		return "assigned"
	case DistantlyPossible:
		return "distantly"
	case Possible:
		return "possible"
	case Likely:
		return "likely"
	case Perfect:
		return "perfect"
	default:
		panic("unknown ConfidenceLevel")
	}
}

type Match struct {
	Confidence ConfidenceLevel
	Rationale  string
}
