package clipboard

import "client/clipboard/clipboardwindows"

func GetCB() Clipboard {
	return &clipboardwindows.ClipboardWindows{}
}
