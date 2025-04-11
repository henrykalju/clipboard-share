//go:build linux

package writer

/*
#cgo LDFLAGS: -lX11 -lXfixes
#include "clipboardwriter.h"
*/
import "C"
import (
	"client/common"
	"errors"
	"fmt"
	"unsafe"
)

func SetListener(l uint64) {
	C.set_listener(C.Window(l))
}

func Init() (uint64, error) {
	w := C.init_clipboard()
	if w == 0 {
		return 0, errors.New("error initing writer")
	}
	go C.start_clipboard()
	return uint64(w), nil
}

func Write(i common.Item) error {
	i, err := common.ConvertItem(i, common.X11)
	if err != nil {
		return fmt.Errorf("error converting: %w", err)
	}
	C.set_clipboard_item(ItemGoToC(i))
	return nil
}

func newCValue(format string, data []byte) C.Value {
	cFormat := C.CString(format)
	cData := (*C.uint8_t)(C.CBytes(data))

	return C.Value{
		format:     cFormat,
		format_len: C.int(len(format)),
		data:       cData,
		data_len:   C.int(len(data)),
	}
}

func newCItem(values []C.Value) C.Item {
	cValues := (*C.Value)(C.malloc(C.size_t(len(values)) * C.size_t(unsafe.Sizeof(C.Value{}))))
	for i, v := range values {
		*(*C.Value)(unsafe.Pointer(uintptr(unsafe.Pointer(cValues)) + uintptr(i)*unsafe.Sizeof(C.Value{}))) = v
	}

	return C.Item{
		values: cValues,
		len:    C.int(len(values)),
	}
}

func ItemGoToC(i common.Item) C.Item {
	values := []C.Value{}
	for _, v := range i.Values {
		values = append(values, newCValue(v.Format, v.Data))
	}
	return newCItem(values)
}
