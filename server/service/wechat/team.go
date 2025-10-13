package wechat

import (
	"cooller/server/global"
	"cooller/server/model/wechat"
	"gorm.io/gorm"
	"time"
)

type TeamService struct{}

func (exa *TeamService) CreateInviteCodeRecode(inviteCodeData *wechat.TeamRecord) (err error) {
	err = global.GVA_DB.Debug().Model(&wechat.TeamRecord{}).Create(inviteCodeData).Error
	return err
}

func (exa *TeamService) GetUserJoinedTeam(userId int, captainId int) (teamRecord *wechat.TeamRecord, err error) {
	err = global.GVA_DB.Debug().Where("user_id = ? AND captain_id = ?", userId, captainId).Preload("WXUser").
		Preload("CaptainWXUser").First(&teamRecord).Error
	if err != nil {
		return nil, err
	}
	return teamRecord, nil
}

func (exa *TeamService) GetUserJoinedTeamByDBSync(tx *gorm.DB, userId int) (*wechat.TeamRecord, error) {
	var teamRecord wechat.TeamRecord
	err := tx.Where("user_id = ? AND is_activated = 1", userId).First(&teamRecord).Error
	if err != nil {
		return nil, err
	}
	return &teamRecord, nil
}

func (exa *TeamService) GetMyTeamsRecordList(userId int) (teamRecordList []wechat.TeamRecord, err error) {
	db := global.GVA_DB.Model(&wechat.TeamRecord{})
	err = db.Where("captain_id = ?", userId).Preload("WXUser").Preload("CaptainWXUser").Find(&teamRecordList).Error
	return teamRecordList, err
}

func (exa *TeamService) GetCaptainTeamConsumeRecordList(captainId int) (list []wechat.TeamConsumeRecord, err error) {
	db := global.GVA_DB.Model(&wechat.TeamConsumeRecord{})
	err = db.Where("captain_id = ? and is_settlement = 0", captainId).Find(&list).Error
	return list, err
}

func (exa *TeamService) GetUnsettledConsumesByCaptainSync(tx *gorm.DB, captainId int) (list []wechat.TeamConsumeRecord, err error) {
	err = tx.Where("captain_id = ? AND is_first_reward_settled = 0", captainId).Find(&list).Error
	return list, err
}

func (exa *TeamService) GetTeamConsumeRecordListByUsers(captain_id int) (list []wechat.TeamConsumeRecord, err error) {
	db := global.GVA_DB.Model(&wechat.TeamConsumeRecord{})
	err = db.Where("captain_id = ?", captain_id).Find(&list).Error
	return list, err
}

func (exa *TeamService) GetTeamFirstConsumeRecord(userId int) (item wechat.TeamConsumeRecord, err error) {
	db := global.GVA_DB.Model(&wechat.TeamConsumeRecord{})
	err = db.Where("user_id = ? and is_first = 1 and is_settlement = 0", userId).First(&item).Error
	return item, err
}

// GetTeamFirstConsumeRecordSync 获取该用户在这个团队中的未结算成团奖励的首购记录
func (exa *TeamService) GetTeamFirstConsumeRecordSync(tx *gorm.DB, userId int) (groupReward *wechat.TeamConsumeRecord, err error) {
	err = tx.Where("user_id = ? AND is_first = 1 AND  is_group_reward_settled = 0", userId).First(&groupReward).Error
	if err != nil {
		return nil, err
	}
	return groupReward, err
}

func (exa *TeamService) GetLoseFirstConsumeRecord(userId int) (item wechat.TeamConsumeRecord, err error) {
	db := global.GVA_DB.Model(&wechat.TeamConsumeRecord{})
	err = db.Where("user_id = ? and is_first = 1", userId).First(&item).Error
	return item, err
}

func (exa *TeamService) GetTeamRecordList(teamId int, captainId int) (list []wechat.TeamRecord, err error) {
	db := global.GVA_DB.Model(&wechat.TeamRecord{})
	err = db.Where("team_id = ? and captain_id = ?", teamId, captainId).Preload("WXUser").Find(&list).Error
	return list, err
}

func (exa *TeamService) GetTeamMembersByDBSync(tx *gorm.DB, teamId int) ([]wechat.TeamRecord, error) {
	var members []wechat.TeamRecord
	return members, tx.Where("team_id = ? AND is_activated = 1", teamId).Find(&members).Error
}

func (exa *TeamService) GetTeamMembers(captainId int) (members []wechat.TeamRecord, err error) {
	db := global.GVA_DB.Model(&wechat.TeamRecord{})
	err = db.Where("captain_id = ?", captainId).Find(&members).Error
	return members, err
}

// GetUserUnsettledConsumes 从数据库获取指定用户的未结算消费记录
func (exa *TeamService) GetUserUnsettledConsumes(userId int) (consumes []wechat.TeamConsumeRecord, err error) {
	db := global.GVA_DB.Model(&wechat.TeamConsumeRecord{})
	err = db.Where("user_id = ? AND is_settlement = 0", userId).Find(&consumes).Error
	return consumes, err
}

// GetUserUnsettledFirstConsumes 从数据库获取指定用户的未结算首购消费记录
func (exa *TeamService) GetUserUnsettledFirstConsumes(userId int) (consumes *wechat.TeamConsumeRecord, err error) {
	db := global.GVA_DB.Model(&wechat.TeamConsumeRecord{})
	err = db.Where("user_id = ? AND is_first = 1 AND is_group_reward_settled = 0", userId).
		First(&consumes).Error
	if err != nil {
		return nil, err
	}
	return consumes, nil
}

func (exa *TeamService) UpdateFirstAndRepurchaseStatusSync(tx *gorm.DB, consumesId []int, now time.Time) (err error) {
	err = tx.Model(&wechat.TeamConsumeRecord{}).Where("id IN (?)", consumesId).
		Updates(map[string]interface{}{
			"is_first_reward_settled":  1,
			"first_reward_settle_time": &now,
		}).Error
	return err
}

// UpdateTeamSettleStatusSync 更新成团奖励结算状态
func (exa *TeamService) UpdateTeamSettleStatusSync(tx *gorm.DB, consumeId int, now time.Time) (err error) {
	err = tx.Model(&wechat.TeamConsumeRecord{}).Where("id = ?", consumeId).
		Updates(map[string]interface{}{
			"is_group_reward_settled":  1,
			"group_reward_settle_time": &now,
		}).Error
	return err
}

func (exa *TeamService) GetTeamsConsumeList(captainId int, userIds []int) (list []wechat.TeamConsumeRecord, err error) {
	db := global.GVA_DB.Model(&wechat.TeamConsumeRecord{})
	err = db.Where("user_id IN ? AND captain_id = ?", userIds, captainId).Preload("WXUser").Find(&list).Error
	return list, err
}
