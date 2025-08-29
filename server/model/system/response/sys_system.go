package response

import "cooller/server/config"

type SysConfigResponse struct {
	Config config.Server `json:"config"`
}
