package wechat

import (
	"cooller/server/global"
	"time"
)

type TeamRecord struct {
	global.GVA_MODEL
	UserId        int    `json:"userId" form:"userId" gorm:"not null;"`
	CaptainId     int    `json:"captainId" form:"captainId" gorm:"not null;comment:队长ID;"`
	InviteCode    string `json:"inviteCode" form:"inviteCode" gorm:"not null;comment:邀请码;size:6"`
	TeamId        int    `json:"teamId" form:"teamId" gorm:"not null;comment:团队编号;"`
	TeamSequence  int    `json:"teamSequence" form:"teamSequence" gorm:"not null;comment:在团队中的序号(1或2)"`
	IsActivated   int    `json:"isActivated" form:"isActivated" gorm:"not null;comment:是否激活，0：未激活，1：激活;"`
	IsSettled     int    `json:"isSettled" form:"isSettled" gorm:"not null;default:0;comment:是否已参与成团结算，0：否，1：是;"`
	WXUser        WXUser `json:"wxUser" form:"wxUser" gorm:"foreignKey:UserId;references：ID;comment:管理WXUser"`
	CaptainWXUser WXUser `json:"captainWXUser" form:"captainWXUser" gorm:"foreignKey:CaptainId;references：ID;comment:关联队长"`
}

func (TeamRecord) TableName() string { return "ums_team_record" }

type TeamSequenceNum struct {
	CaptainId   int `json:"captainId" form:"captainId" gorm:"not null;comment:队长;"`
	MaxSequence int `json:"maxSequence" form:"maxSequence" gorm:"not null;comment:最大序列编号;"`
}

func (TeamSequenceNum) TableName() string { return "ums_team_sequence_num" }

// TeamConsumeRecord 团队消费记录
type TeamConsumeRecord struct {
	global.GVA_MODEL
	UserId                int        `json:"userId" form:"userId" gorm:"not null;"`
	CaptainId             int        `json:"captainId" form:"captainId" gorm:"not null;comment:队长;"`
	OrderId               int        `json:"orderId" form:"orderId" gorm:"not null;"`
	Amount                float32    `json:"amount" form:"amount" gorm:"not null;comment:金额;"`
	IsFirst               int        `json:"isFirst" form:"isFirst" gorm:"not null;comment:是否结算，0：否，1：是;"`
	IsFirstRewardSettled  int        `json:"isFirstRewardSettled" form:"isFirstRewardSettled" gorm:"not null;default:0;comment:首推佣金是否结算，0：否，1：是;"`
	IsGroupRewardSettled  int        `json:"isGroupRewardSettled" form:"isGroupRewardSettled" gorm:"not null;default:0;comment:成团奖励是否结算，0：否，1：是;"`
	FirstRewardSettleTime *time.Time `json:"firstRewardSettleTime" form:"firstRewardSettleTime" gorm:"comment:首推佣金结算时间"`
	GroupRewardSettleTime *time.Time `json:"groupRewardSettleTime" form:"groupRewardSettleTime" gorm:"comment:成团奖励结算时间"`
	//TeamId                int        `json:"teamId" form:"teamId" gorm:"index;comment:所属团队ID"`
	WXUser WXUser `json:"wxUser" form:"wxUser" gorm:"foreignKey:UserId;references：ID;comment:管理WXUser"`
}

func (TeamConsumeRecord) TableName() string { return "ums_team_consume_record" }

// TeamSettlement 团队结算记录，用于汇总每次结算
type TeamSettlement struct {
	global.GVA_MODEL
	UserId         int       `json:"userId" form:"userId" gorm:"not null;index"`
	SettlementNo   string    `json:"settlementNo" form:"settlementNo" gorm:"not null;unique;comment:结算单号"`
	TotalAmount    float32   `json:"totalAmount" form:"totalAmount" gorm:"not null;comment:结算总金额"`
	SettlementTime time.Time `json:"settlementTime" form:"settlementTime" gorm:"not null;comment:结算时间"`
	Status         int       `json:"status" form:"status" gorm:"not null;default:0;comment:结算状态，0:待处理,1:已完成"`
}

func (TeamSettlement) TableName() string { return "ums_team_settlement" }
