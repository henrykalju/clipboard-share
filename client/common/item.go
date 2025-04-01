package common

type Type struct {
	Text string
}

var (
	X11     = Type{"X11"}
	WINDOWS = Type{"WINDOWS"}
	NONE    = Type{}
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

type ItemWithID struct {
	Item
	ID int32
}

func GetType(text string) Type {
	if text == X11.Text {
		return X11
	}
	if text == WINDOWS.Text {
		return WINDOWS
	}
	return NONE
}
