package main

import (
	"main/clipboard"
	"main/types"
)

func main() {
	c := clipboard.GetCB()
	c.Write(types.Item{})
}
