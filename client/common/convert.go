package common

import (
	"errors"
	"fmt"
	"slices"
)

func ConvertItem(item Item, to Type) (Item, error) {
	if item.Type == to {
		return item, nil
	}

	if item.Type == X11 && to == WINDOWS {
		return convertX11ToWindows(item)
	}

	if item.Type == WINDOWS && to == X11 {
		return convertWindowsToX11(item)
	}

	return item, fmt.Errorf("cannot convert from type %s to %s", item.Type.Text, to.Text)
}

func convertX11ToWindows(item Item) (Item, error) {
	i := Item{
		Type: WINDOWS,
		Text: item.Text,
	}

	i.Values = make([]Value, 0)

	textI := slices.IndexFunc(item.Values, func(v Value) bool {
		return v.Format == "STRING"
	})
	if textI != -1 {
		v := Value{
			Format: "CF_TEXT",
			Data:   append(item.Values[textI].Data, 0),
		}
		i.Values = append(i.Values, v)
	}

	if len(i.Values) == 0 {
		return i, errors.New("no STRING found converting X11 to Windows")
	}

	return i, nil
}

func convertWindowsToX11(item Item) (Item, error) {
	i := Item{
		Type: X11,
		Text: item.Text,
	}

	i.Values = make([]Value, 0)

	textI := slices.IndexFunc(item.Values, func(v Value) bool {
		return v.Format == "CF_TEXT"
	})
	if textI != -1 {
		v := Value{
			Format: "STRING",
			Data:   item.Values[textI].Data[:len(item.Values[textI].Data)-1],
		}
		i.Values = append(i.Values, v)
	}

	if len(i.Values) == 0 {
		return i, errors.New("no CF_TEXT found converting Windows to X11")
	}

	return i, nil
}
