package request

import (
	"cooller/server/model/common/request"
	"cooller/server/model/system"
)

type SysDictionaryDetailSearch struct {
	system.SysDictionaryDetail
	request.PageInfo
}
