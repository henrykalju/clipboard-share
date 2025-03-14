package clipboard

import "main/types"

type Clipboard interface {
	Init()
	GetChan() chan *types.Item
	Write(types.Item)
}
