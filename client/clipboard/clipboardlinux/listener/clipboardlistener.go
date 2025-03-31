//go:build linux

package listener

/*
#cgo LDFLAGS: -lX11 -lXfixes
#include "clipboardlistener.h"
*/
import "C"
import (
	"client/types"
	"errors"
	"slices"
	"unsafe"
)

var c chan *types.Item

func addItemToChan(i types.Item) {
	if c == nil {
		panic(errors.New("adding item to nil channel"))
	}
	i.Type = types.X11
	c <- &i
}

func GetChan() chan *types.Item {
	return c
}

func SetWriter(w uint64) {
	C.set_writer(C.Window(w))
}

func Init() (uint64, error) {
	w := C.Init()
	if w == 0 {
		panic(errors.New("error initing listener"))
	}
	go C.StartListening() // TODO: add return
	c = make(chan *types.Item)
	return uint64(w), nil
}

//export handleChange
func handleChange(i C.Item) {
	item := ConvertItemC2G(i)

	go addItemToChan(item)
}

func ConvertItemC2G(i C.Item) types.Item {
	cValues := unsafe.Slice((*C.Value)(i.values), i.len)
	goValues := make([]types.Value, i.len)

	// Convert each C Value to a Go Value
	for i, cVal := range cValues {
		name := C.GoStringN(cVal.format, cVal.format_len) // Convert format to Name

		// Convert the C uint8_t* array (data) to a Go []byte
		goData := C.GoBytes(unsafe.Pointer(cVal.data), cVal.data_len)

		goValues[i] = types.Value{
			Format: name,
			Data:   goData,
		}
	}

	name := FindName(goValues)

	return types.Item{
		Text:   name,
		Values: goValues,
	}
}

func FindName(values []types.Value) string {
	STRINGi := slices.IndexFunc(values, func(v types.Value) bool {
		return v.Format == "STRING"
	})
	if STRINGi != -1 {
		return string(values[STRINGi].Data)
	}

	return ""
}
