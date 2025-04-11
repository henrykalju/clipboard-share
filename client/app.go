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
	err := common.InitConfig()
	if err != nil {
		fmt.Printf("Error initing config: %s\n", err.Error())
		panic("Error initing config")
	}

	a.ctx = ctx

	a.c = clipboard.GetCB()
	err = a.c.Init()
	if err != nil {
		fmt.Printf("Error initing clipboard: %s\n", err.Error())
		panic("Error initing clipboard")
	}

	go a.startListeningForClipboard()

	a.history = make([]common.Item, 0)
	a.history = append(a.history, common.Item{Text: "test", Values: []common.Value{{Format: "STRING", Data: []byte("test\n")}}})
}

func (a *App) GetHistory() ([]common.ItemWithID, error) {
	r, err := storage.GetItems()
	if err != nil {
		return nil, fmt.Errorf("error getting items from storage: %w", err)
	}
	return r, nil
}

func (a *App) WriteToCB(id int32) error {
	i, err := storage.GetItemByID(id)
	if err != nil {
		return fmt.Errorf("error getting item with id %d: %w", id, err)
	}

	err = a.c.Write(i)
	if err != nil {
		return fmt.Errorf("error writing to clipboard: %w", err)
	}
	return nil
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

func (a *App) UpdateConfig(conf common.Config) error {
	err := common.SetConf(conf)
	if err != nil {
		fmt.Printf("Error updating conf: %s\n", err.Error())
	}
	return err
}

func (a *App) GetConfig() common.Config {
	return common.GetConf()
}
