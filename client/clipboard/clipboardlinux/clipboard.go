//go:build linux

package clipboardlinux

import (
	"client/clipboard/clipboardlinux/listener"
	"client/clipboard/clipboardlinux/writer"
	"client/common"
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

func (c *ClipboardLinux) GetChan() chan *common.Item {
	return listener.GetChan()
}

func (c *ClipboardLinux) Write(i common.Item) {
	writer.Write(i)
}
