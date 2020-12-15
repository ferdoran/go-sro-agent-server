package party

type PurposeType byte

const (
	Hunting PurposeType = iota
	Quest
	Trader
	Thief
)

const (
	JoinRequestTimeout int = 10
)

const (
	JoinRequestResponseDenied   uint16 = 0
	JoinRequestResponseAccepted uint16 = 1
	JoinRequestResponseTimeout  uint16 = 2
)