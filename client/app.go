package main

import (
	"client/clipboard"
	"client/common"
	"client/storage"
	"context"
	"fmt"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

const CB_UPDATE_EVENT = "CB_UPDATE_EVENT"

var CB = "asd"

// App struct
type App struct {
	ctx     context.Context
	cbChan  chan *common.Item
	history []common.Item
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
	//c.Write(common.Item{Text: "test", Values: []common.Value{{Format: "STRING", Data: []byte("test\n")}}})
	a.cbChan = c.GetChan()
	go a.startListeningForClipboard()

	a.history = make([]common.Item, 0)
	a.history = append(a.history, common.Item{Text: "test", Values: []common.Value{{Format: "STRING", Data: []byte("test\n")}}})
}

func (a *App) GetHistory() []common.ItemWithID {
	r, err := storage.GetItems()
	if err != nil {
		fmt.Printf("Error getting items from storage: %s\n", err.Error())
		return nil
	}
	return r
}

func (a *App) startListeningForClipboard() {
	for {
		select {
		case i := <-a.cbChan:
			err := storage.SaveItem(i)
			if err != nil {
				fmt.Printf("Error saving item to storage: %s\n", err.Error())
			}
			runtime.EventsEmit(a.ctx, CB_UPDATE_EVENT)
		}
	}
}
