package clipboard

import "client/common"

type Clipboard interface {
	Init() error
	GetChan() chan *common.Item
	Write(common.Item) error
}
