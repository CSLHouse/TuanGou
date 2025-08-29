package initialize

import (
	_ "cooller/server/source/example"
	_ "cooller/server/source/product"
	_ "cooller/server/source/system"
	_ "cooller/server/source/wechat"
)

func init() {
	// do nothing,only import source package so that inits can be registered
}
