//go:build windows

package clipboard

import "fmt"

func Write() {
	fmt.Println("Writing to windows clipboard")
}
