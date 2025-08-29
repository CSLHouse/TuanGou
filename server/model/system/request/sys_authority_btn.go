package request

type SysAuthorityBtnReq struct {
	MenuID      int   `json:"menuID"`
	AuthorityId int   `json:"authorityId"`
	Selected    []int `json:"selected"`
}
