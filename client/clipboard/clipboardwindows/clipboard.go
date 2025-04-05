//go:build windows

package clipboardwindows

import (
	"client/common"
	"errors"
	"fmt"
	"slices"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

type ClipboardWindows struct {
}

func (c *ClipboardWindows) Init() {
	user32 = windows.MustLoadDLL("user32.dll")
	procRegisterClass = user32.MustFindProc("RegisterClassW")
	procCreateWindowEx = user32.MustFindProc("CreateWindowExW")
	procAddClipboardListener = user32.MustFindProc("AddClipboardFormatListener")
	procGetMessage = user32.MustFindProc("GetMessageW")
	procTranslateMessage = user32.MustFindProc("TranslateMessage")
	procDispatchMessage = user32.MustFindProc("DispatchMessageW")
	procDefWindowProc = user32.MustFindProc("DefWindowProcW")
	ch = make(chan *common.Item)

	kernel32 = windows.MustLoadDLL("kernel32.dll")
	globalAlloc = kernel32.MustFindProc("GlobalAlloc")
	globalSize = kernel32.MustFindProc("GlobalSize")
	globalLock = kernel32.MustFindProc("GlobalLock")
	globalUnlock = kernel32.MustFindProc("GlobalUnlock")
	openClipboard = user32.MustFindProc("OpenClipboard")
	closeClipboard = user32.MustFindProc("CloseClipboard")

	emptyClipboard = user32.MustFindProc("EmptyClipboard")
	setClipboardData = user32.MustFindProc("SetClipboardData")
	getClipboardData = user32.MustFindProc("GetClipboardData")
	enumClipboardFormats = user32.MustFindProc("EnumClipboardFormats")
	getClipboardFormatName = user32.MustFindProc("GetClipboardFormatNameW")

	// Register the window class
	className, _ := syscall.UTF16PtrFromString("ClipboardListenerClass")

	wc := WNDCLASS{
		lpfnWndProc:   syscall.NewCallback(wndProc),
		lpszClassName: uintptr(unsafe.Pointer(className)),
		hInstance:     windows.Handle(0), // Use the current module handle
	}

	ret, _, err := procRegisterClass.Call(uintptr(unsafe.Pointer(&wc)))
	if ret == 0 {
		fmt.Println("Error registering class:", err)
		return
	}

	windowName, _ := syscall.UTF16PtrFromString("Clipboard Listener")
	// Create the window
	hwnd, _, err := procCreateWindowEx.Call(
		0,
		uintptr(unsafe.Pointer(className)),
		uintptr(unsafe.Pointer(windowName)),
		0,
		CW_USEDEFAULT, CW_USEDEFAULT, CW_USEDEFAULT, CW_USEDEFAULT,
		0, 0, 0, 0,
	)
	if hwnd == 0 {
		fmt.Println("Error creating window:", err)
		return
	}

	// Register for clipboard updates
	ret, _, err = procAddClipboardListener.Call(hwnd)
	if ret == 0 {
		fmt.Println("Error registering clipboard listener:", err)
		return
	}

	text := "testing"
	i := common.Item{
		Text:   text,
		Values: []common.Value{{Format: "CF_TEXT", Data: append([]byte(text), 0)}},
	}

	//if write first, and then listen, it works
	c.Write(i)

	go listen()
}

func (c *ClipboardWindows) GetChan() chan *common.Item {
	return ch
}

func (c *ClipboardWindows) Write(i common.Item) {
	fmt.Printf("writing %v\n", i)
	err := open()
	if err != nil {
		panic(err)
	}
	defer close()

	ok, _, err := emptyClipboard.Call()
	if !errors.Is(err, windows.ERROR_SUCCESS) {
		panic(fmt.Errorf("error emptying cb: %w", err))
	}
	if ok == 0 {
		panic("error emptying clipboard")
	}

	fmt.Println("adding values")
	for _, v := range i.Values {
		setValueToCB(v)
	}
	fmt.Println("values added")
}

// TODO add if getClipboardOwner == hwnd, dont look at it

var (
	ch                       chan *common.Item
	user32                   *windows.DLL
	procRegisterClass        *windows.Proc
	procCreateWindowEx       *windows.Proc
	procAddClipboardListener *windows.Proc
	procGetMessage           *windows.Proc
	procTranslateMessage     *windows.Proc
	procDispatchMessage      *windows.Proc
	procDefWindowProc        *windows.Proc

	kernel32       *windows.DLL
	globalAlloc    *windows.Proc
	globalLock     *windows.Proc
	globalUnlock   *windows.Proc
	globalSize     *windows.Proc
	openClipboard  *windows.Proc
	closeClipboard *windows.Proc

	emptyClipboard         *windows.Proc
	setClipboardData       *windows.Proc
	getClipboardData       *windows.Proc
	enumClipboardFormats   *windows.Proc
	getClipboardFormatName *windows.Proc

	hwnd uintptr
)

const (
	WM_CLIPBOARDUPDATE = 0x031D
	CW_USEDEFAULT      = 0x80000000
	GMEM_MOVEABLE      = 0x0002

	CF_BITMAP          = 2
	CF_DIB             = 8
	CF_DIBV5           = 17
	CF_DIF             = 5
	CF_DSPBITMAP       = 0x0082
	CF_DSPENHMETAFILE  = 0x008E
	CF_DSPMETAFILEPICT = 0x0083
	CF_DSPTEXT         = 0x0081
	CF_ENHMETAFILE     = 14
	CF_GDIOBJFIRST     = 0x0300
	CF_GDIOBJLAST      = 0x03FF
	CF_HDROP           = 15
	CF_LOCALE          = 16
	CF_METAFILEPICT    = 3
	CF_OEMTEXT         = 7
	CF_OWNERDISPLAY    = 0x0080
	CF_PALETTE         = 9
	CF_PENDATA         = 10
	CF_PRIVATEFIRST    = 0x0200
	CF_PRIVATELAST     = 0x02FF
	CF_RIFF            = 11
	CF_SYLK            = 4
	CF_TEXT            = 1
	CF_TIFF            = 6
	CF_UNICODETEXT     = 13
	CF_WAVE            = 12
)

type MSG struct {
	Hwnd    windows.Handle
	Message uint32
	WParam  uintptr
	LParam  uintptr
	Time    uint32
	Pt      struct {
		X, Y int32
	}
}

type WNDCLASS struct {
	style         uint32
	lpfnWndProc   uintptr
	cbClsExtra    int32
	cbWndExtra    int32
	hInstance     windows.Handle
	hIcon         windows.Handle
	hCursor       windows.Handle
	hbrBackground windows.Handle
	lpszMenuName  uintptr
	lpszClassName uintptr
}

// Window procedure to handle messages
func wndProc(hwnd windows.Handle, msg uint32, wParam, lParam uintptr) uintptr {
	switch msg {
	case WM_CLIPBOARDUPDATE:
		fmt.Println("Clipboard content updated!")
		contentToChan()
	default:
		ret, _, _ := procDefWindowProc.Call(uintptr(hwnd), uintptr(msg), wParam, lParam)
		return ret
	}
	return 0
}

func listen() {
	// Message loop
	var msg MSG
	for {
		fmt.Println("waiting for message")
		ret, _, err := procGetMessage.Call(uintptr(unsafe.Pointer(&msg)), 0, 0, 0)
		if ret == 0 {
			break // WM_QUIT message, exit the loop
		}
		if ret == ^uintptr(0) {
			fmt.Println("Error getting message:", err)
			return
		}
		procTranslateMessage.Call(uintptr(unsafe.Pointer(&msg)))
		procDispatchMessage.Call(uintptr(unsafe.Pointer(&msg)))
	}
}

func contentToChan() {
	i := common.Item{}
	open()
	defer close()

	var format uintptr = 0
	for {
		format, _, _ = enumClipboardFormats.Call(format)
		if format == 0 {
			break
		}
		addValue(&i, format)
	}
	//fmt.Printf("%v\n", i.Values[slices.IndexFunc(i.Values, func(e Value) bool { return e.Format == "CF_TEXT" })])
	fmt.Println("adding content")

	i.Text = findText(i.Values)
	i.Type = common.WINDOWS

	ch <- &i
}

func findText(values []common.Value) string {
	STRINGi := slices.IndexFunc(values, func(v common.Value) bool {
		return v.Format == "CF_TEXT"
	})
	if STRINGi != -1 {
		return string(values[STRINGi].Data[:len(values[STRINGi].Data)-1])
	}

	return ""
}

func forbiddenFormat(f string) bool {
	return f == "EnterpriseDataProtectionId"
}

func addValue(i *common.Item, f uintptr) {
	//fmt.Printf("%v\n", f)

	var formatName string
	switch f {
	case CF_BITMAP:
		formatName = "CF_BITMAP"
	case CF_DIB:
		formatName = "CF_DIB"
	case CF_DIBV5:
		formatName = "CF_DIBV5"
	case CF_DIF:
		formatName = "CF_DIF"
	case CF_DSPBITMAP:
		formatName = "CF_DSPBITMAP"
	case CF_DSPENHMETAFILE:
		formatName = "CF_DSPENHMETAFILE"
	case CF_DSPMETAFILEPICT:
		formatName = "CF_DSPMETAFILEPICT"
	case CF_DSPTEXT:
		formatName = "CF_DSPTEXT"
	case CF_ENHMETAFILE:
		formatName = "CF_ENHMETAFILE"
	case CF_GDIOBJFIRST:
		formatName = "CF_GDIOBJFIRST"
	case CF_GDIOBJLAST:
		formatName = "CF_GDIOBJLAST"
	case CF_HDROP:
		formatName = "CF_HDROP"
	case CF_LOCALE:
		formatName = "CF_LOCALE"
	case CF_METAFILEPICT:
		formatName = "CF_METAFILEPICT"
	case CF_OEMTEXT:
		formatName = "CF_OEMTEXT"
	case CF_OWNERDISPLAY:
		formatName = "CF_OWNERDISPLAY"
	case CF_PALETTE:
		formatName = "CF_PALETTE"
	case CF_PENDATA:
		formatName = "CF_PENDATA"
	case CF_PRIVATEFIRST:
		formatName = "CF_PRIVATEFIRST"
	case CF_PRIVATELAST:
		formatName = "CF_PRIVATELAST"
	case CF_RIFF:
		formatName = "CF_RIFF"
	case CF_SYLK:
		formatName = "CF_SYLK"
	case CF_TEXT:
		formatName = "CF_TEXT"
	case CF_TIFF:
		formatName = "CF_TIFF"
	case CF_UNICODETEXT:
		formatName = "CF_UNICODETEXT"
	case CF_WAVE:
		formatName = "CF_WAVE"
	default:
		var formatNameW [256]uint16
		ret, _, _ := getClipboardFormatName.Call(f, uintptr(unsafe.Pointer(&formatNameW[0])), uintptr(len(formatNameW)))
		if ret == 0 {
			//panic("failed to get format name")
			fmt.Println("failed to get format name")
			return
		}

		formatName = syscall.UTF16ToString(formatNameW[:])
	}

	if forbiddenFormat(formatName) {
		return
	}
	if ptr == 0 {
		panic("ptr = 0")
	}

	size, _, _ := globalSize.Call(ptr)
	if size == 0 {
		panic("failed to get size")
	}

	mem, _, _ := globalLock.Call(ptr)
	if mem == 0 {
		panic("failed to lock g mem")
	}
	defer globalUnlock.Call(ptr)
	data := unsafe.Slice((*byte)(unsafe.Pointer(mem)), size)

	i.Values = append(i.Values, common.Value{Format: formatName, Data: data})
}

func open() error {
	for tries := range 5 {
	_, _, err := openClipboard.Call(hwnd)
		if errors.Is(err, windows.ERROR_SUCCESS) {
			fmt.Printf("Clipboard opened in %d tries\n", tries+1)
			break
		}
		if tries == 4 {
			fmt.Printf("Couldn't open clipboard in 5 tries: %s\n", err.Error())
		panic(err)
	}
		//time.Sleep
	}

	return nil
}

func close() {
	_, _, err := closeClipboard.Call()
	if !errors.Is(err, windows.ERROR_SUCCESS) {
		fmt.Printf("Error opening clipboard: %s\n", err.Error())
	}
}
func getWinFormat(f string) uint32 {
	switch f {
	case "CF_TEXT":
		return CF_TEXT
	case "CF_UNICODETEXT":
		return CF_UNICODETEXT
	}

	// TODO add format registering

	return 1
}

func setValueToCB(v common.Value) {
	hMem, _, _ := globalAlloc.Call(GMEM_MOVEABLE, uintptr(len(v.Data)))
	if hMem == 0 {
		panic("failed to allocate global memory")
	}

	ptr, _, _ := globalLock.Call(hMem)
	if ptr == 0 {
		panic("failed to lock global memory")
	}

	copy((*[1 << 30]byte)(unsafe.Pointer(ptr))[:len(v.Data)], v.Data)
	globalUnlock.Call(hMem)

	if r, _, _ := setClipboardData.Call(uintptr(getWinFormat(v.Format)), hMem); r == 0 {
		panic("failed to set clipboard data")
	}
}
