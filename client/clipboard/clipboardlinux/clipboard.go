//go:build linux

package clipboardlinux

import (
	"client/clipboard/clipboardlinux/listener"
	"client/clipboard/clipboardlinux/writer"
	"client/common"
	"fmt"
)

type ClipboardLinux struct {
}

func (c *ClipboardLinux) Init() error {
	l, err := listener.Init()
	if err != nil {
		return fmt.Errorf("error initing linux clipboard listener: %w", err)
	}

	w, err := writer.Init()
	if err != nil {
		return fmt.Errorf("error initing linux clipboard writer: %w", err)
	}

	listener.SetWriter(w)
	writer.SetListener(l)
	return nil
}

func (c *ClipboardLinux) GetChan() chan *common.Item {
	return listener.GetChan()
}

func (c *ClipboardLinux) Write(i common.Item) error {
	return writer.Write(i)
}
