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
	c       clipboard.Clipboard
	history []common.Item
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	common.InitConfig()

	a.ctx = ctx

	a.c = clipboard.GetCB()
	a.c.Init()

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

func (a *App) WriteToCB(id int32) {
	i, err := storage.GetItemByID(id)
	if err != nil {
		fmt.Printf("Error getting item with id %d: %s\n", id, err.Error())
		return
	}

	a.c.Write(i)
}

func (a *App) startListeningForClipboard() {
	for {
		select {
		case i := <-a.c.GetChan():
			err := storage.SaveItem(i)
			if err != nil {
				fmt.Printf("Error saving item to storage: %s\n", err.Error())
			}
			runtime.EventsEmit(a.ctx, CB_UPDATE_EVENT)
		}
	}
}

func (a *App) UpdateConfig(conf common.Config) {
	common.SetConf(conf)
}

func (a *App) GetConfig() common.Config {
	return common.GetConf()
}
