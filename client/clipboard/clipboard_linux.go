//go:build linux

package clipboard

import "client/clipboard/clipboardlinux"

func GetCB() Clipboard {
	return &clipboardlinux.ClipboardLinux{}
}
