package business

type VIPCertificate struct {
	ID        int    `json:"Id" gorm:"primarykey"`
	Telephone string `json:"telephone" form:"telephone" gorm:"comment:会员手机号"` // 客户手机号
	StoreName string `json:"storeName" form:"storeName" gorm:"comment:店名"`
	IsFirst   bool   `json:"isFirst" form:"isFirst" gorm:"comment:第一个凭证"`
	SysUserId int    `json:"sysUserId" form:"sysUserId" gorm:"comment:管理角色ID"` // 管理角色ID
	Count     int    `json:"count" form:"count" gorm:"comment:消费次数"`
}

func (VIPCertificate) TableName() string {
	return "bus_certificate"
}
