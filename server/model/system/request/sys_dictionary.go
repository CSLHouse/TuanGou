package request

import (
	"cooller/server/model/common/request"
	"cooller/server/model/system"
)

type SysDictionarySearch struct {
	system.SysDictionary
	request.PageInfo
}
