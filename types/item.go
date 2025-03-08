package types

type Item struct {
	Text   string
	Values []Value
}

type Value struct {
	Format string
	Data   []byte
}
