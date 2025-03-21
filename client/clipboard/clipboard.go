package clipboard

import "client/types"

type Clipboard interface {
	Init()
	GetChan() chan *types.Item
	Write(types.Item)
}
