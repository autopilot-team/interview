package types

type Mode string

const (
	DebugMode   Mode = "debug"
	ReleaseMode Mode = "release"
)

func (m Mode) IsValid() bool {
	return m == DebugMode || m == ReleaseMode
}

func (m Mode) String() string {
	return string(m)
}
