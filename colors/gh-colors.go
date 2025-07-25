package colors

var (
	DefaultBg1 Color = *FromHex("#033a16")
	DefaultBg2 Color = *FromHex("#196c2e")
	DefaultBg3 Color = *FromHex("#2ea043")
	DefaultBg4 Color = *FromHex("#56d364")
)

func FromIntensity(i uint8) *Color {
	switch i {
	case 0:
		return nil
	case 1:
		return &DefaultBg1
	case 2:
		return &DefaultBg2
	case 3:
		return &DefaultBg3
	case 4:
		return &DefaultBg4
	default:
		panic("UNREACHABLE")
	}
	return nil
}
