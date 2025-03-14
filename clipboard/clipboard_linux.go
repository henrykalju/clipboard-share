//go:build linux

package clipboard

import "main/clipboard/clipboardlinux"

func GetCB() Clipboard {
	return &clipboardlinux.ClipboardLinux{}
}
