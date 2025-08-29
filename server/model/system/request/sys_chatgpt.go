package request

import (
	"cooller/server/model/common/request"
	"cooller/server/model/system"
)

type ChatGptRequest struct {
	system.ChatGpt
	request.PageInfo
}
