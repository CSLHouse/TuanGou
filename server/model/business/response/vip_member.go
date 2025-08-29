package response

import (
	"cooller/server/model/business"
)

//	type VIPMemberResponseModel struct {
//		CardID             string                   `json:"cardID" form:"cardID" gorm:"comment:会员卡号"`                           // 客户名
//		Telephone          string                   `json:"telephone" form:"telephone" gorm:"comment:会员手机号"`                    // 客户手机号
//		MemberName         string                   `json:"memberName" form:"memberName" gorm:"comment:会员名"`                    // 客户名
//		MemberType         uint                     `json:"memberType" form:"memberType" gorm:"comment:会员类型"`                   // 管理ID
//		RemainTimes        uint                     `json:"remainTimes" form:"remainTimes" gorm:"comment:剩余次数/金额"`              // 管理ID
//		Deadline           uint                     `json:"deadline" form:"deadline" gorm:"comment:截止日期"`                       // 管理ID
//		State              uint                     `json:"state" form:"state" gorm:"comment:状态"`                               // 管理角色ID
//		SysUserAuthorityID uint                     `json:"sysUserAuthorityID" form:"sysUserAuthorityID" gorm:"comment:管理角色ID"` // 管理角色ID
//		VIPMembers         []VIPMemberResponseModel `json:"vipmembers" gorm:"foreignKey:sys_user_authority_id;"`
//	}

type WXLoginResponse struct {
	Customer  business.Customer `json:"customer"`
	Token     string            `json:"token"`
	ExpiresAt int64             `json:"expiresAt"`
}

type VIPMemberResponse struct {
	Member business.Customer `json:"member"`
}
