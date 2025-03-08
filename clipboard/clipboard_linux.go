//go:build linux

package clipboard

import "fmt"

func Write() {
	fmt.Println("Writing to linux clipboard")
}
