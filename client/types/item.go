package types

type Type struct {
	Type string
}

var (
	X11     = Type{"X11"}
	WINDOWS = Type{"WINDOWS"}
)

type Item struct {
	Type   Type
	Text   string
	Values []Value
}

type Value struct {
	Format string
	Data   []byte
}
