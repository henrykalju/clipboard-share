package clipboard

import "client/common"

type Clipboard interface {
	Init()
	GetChan() chan *common.Item
	Write(common.Item)
}
