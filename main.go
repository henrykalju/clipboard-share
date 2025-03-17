package main

import (
	"main/clipboard"
	"main/types"
	"time"
)

func main() {
	c := clipboard.GetCB()
	c.Init()
	c.Write(types.Item{Text: "test", Values: []types.Value{{Format: "STRING", Data: []byte("test\n")}}})
	time.Sleep(time.Second * 10)
}
