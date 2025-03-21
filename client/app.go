package main

import (
	"client/clipboard"
	"client/types"
	"context"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

const CB_UPDATE_EVENT = "CB_UPDATE_EVENT"

var CB = "asd"

type Item struct {
	ID   int
	Text string
}

// App struct
type App struct {
	ctx     context.Context
	cbChan  chan *types.Item
	history []types.Item
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	c := clipboard.GetCB()
	c.Init()
	c.Write(types.Item{Text: "test", Values: []types.Value{{Format: "STRING", Data: []byte("test\n")}}})
	a.cbChan = c.GetChan()
	go a.startListeningForClipboard()

	a.history = make([]types.Item, 0)
	a.history = append(a.history, types.Item{Text: "test", Values: []types.Value{{Format: "STRING", Data: []byte("test\n")}}})
}

// Greet returns a greeting for the given name
func (a *App) GetHistory() []Item {
	r := make([]Item, 0, len(a.history))
	for i, v := range a.history {
		r = append(r, Item{ID: i, Text: v.Text})
	}
	return r
}

func (a *App) startListeningForClipboard() {
	for {
		select {
		case <-a.cbChan:
			//fmt.Printf("%v\n", i)
			runtime.EventsEmit(a.ctx, CB_UPDATE_EVENT)
		}
	}
}
