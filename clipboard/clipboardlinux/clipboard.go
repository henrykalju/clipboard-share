//go:build linux

package clipboardlinux

import (
	"main/clipboard/clipboardlinux/listener"
	"main/clipboard/clipboardlinux/writer"
	"main/types"
)

type ClipboardLinux struct {
}

func (c *ClipboardLinux) Init() {
	l, err := listener.Init()
	if err != nil {
		panic(err)
	}

	w, err := writer.Init()
	if err != nil {
		panic(err)
	}

	listener.SetWriter(w)
	writer.SetListener(l)
}

func (c *ClipboardLinux) GetChan() chan *types.Item {
	return listener.GetChan()
}

func (c *ClipboardLinux) Write(i types.Item) {
	writer.Write(i)
}
