package request

import (
	"cooller/server/model/common/request"
	"cooller/server/model/system"
)

type SysOperationRecordSearch struct {
	system.SysOperationRecord
	request.PageInfo
}
