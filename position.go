package bubbleup

type Position string

func (p Position) IsValid() bool {
	return p.String() != "unknown"
}

func (p Position) String() string {
	switch p {
	case TopLeftPosition:
		return "top-left"
	case TopCenterPosition:
		return "top-center"
	case TopRightPosition:
		return "top-right"
	case BottomLeftPosition:
		return "bottom-left"
	case BottomCenterPosition:
		return "bottom-center"
	case BottomRightPosition:
		return "bottom-right"
	default:
		return "unknown"
	}
}
func (p Position) Label() string {
	switch p {
	case TopLeftPosition:
		return "Top Left"
	case TopCenterPosition:
		return "Top Center"
	case TopRightPosition:
		return "Top Right"
	case BottomLeftPosition:
		return "Bottom Left"
	case BottomCenterPosition:
		return "Bottom Center"
	case BottomRightPosition:
		return "Bottom Right"
	default:
		return "Unknown"
	}
}

const (
	TopLeftPosition      Position = "TL"
	TopCenterPosition    Position = "TC"
	TopRightPosition     Position = "TR"
	BottomLeftPosition   Position = "BL"
	BottomCenterPosition Position = "BC"
	BottomRightPosition  Position = "BR"
	UnspecifiedPosition  Position = ""
)
