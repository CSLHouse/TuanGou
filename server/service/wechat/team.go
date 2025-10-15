package wechat

import (
	"cooller/server/global"
	"cooller/server/model/common/request"
	"cooller/server/model/wechat"
	"errors"
	"fmt"
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

func (exa *TeamService) GetUserJoinedTeamByDBSync(userId int) (*wechat.TeamRecord, error) {
	var teamRecord wechat.TeamRecord
	err := global.GVA_DB.Debug().Where("user_id = ?", userId).First(&teamRecord).Error
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

func (exa *TeamService) GetUnsettledConsumesByCaptainSync(captainId int) (list []wechat.TeamConsumeRecord, err error) {
	err = global.GVA_DB.Debug().Where("captain_id = ? AND is_first_reward_settled = 0", captainId).Find(&list).Error
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
func (exa *TeamService) GetTeamFirstConsumeRecordSync(userId int) (groupReward *wechat.TeamConsumeRecord, err error) {
	result := global.GVA_DB.Where("user_id = ? AND is_first = 1 AND  is_group_reward_settled = 0", userId).First(&groupReward)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
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

func (exa *TeamService) GetTeamMembers(captainId int, teamId int) ([]wechat.TeamRecord, error) {
	var members []wechat.TeamRecord
	err := global.GVA_DB.Where("captain_id = ? AND team_id = ?", captainId, teamId).Find(&members).Error
	return members, err
}

func (exa *TeamService) GetAllTeamMembers(captainId int) (members []wechat.TeamRecord, err error) {
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

// CreateTeamSettlementSync 创建结算单
// consumesIds 用户作为队长的所有未结算消费记录的id
// 用户作为成员的未结算成团奖励记录
func (exa *TeamService) CreateTeamSettlementSync(consumesIds []int, consumesId int, data *wechat.TeamSettlement) (err error) {
	tx := global.GVA_DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if tx.Error != nil {
		return err
	}

	// 创建结算单
	err = tx.Model(&wechat.TeamSettlement{}).Create(data).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	// 2.4 更新首推/复购奖励结算状态
	if len(consumesIds) > 0 {
		err = exa.UpdateFirstAndRepurchaseStatusSync(tx, consumesIds, data.SettlementTime)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("更新首推奖励结算状态失败: %w", err)
		}
	}

	// 2.5 更新成团奖励结算状态
	if consumesId != 0 {
		err = exa.UpdateTeamSettleStatusSync(tx, consumesId, data.SettlementTime)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("更新成团奖励结算状态失败: %w", err)
		}

		// 更新团队结算状态
		if err := exa.UpdateUserTeamSettledStatus(tx, data.UserId); err != nil {
			tx.Rollback()
			return fmt.Errorf("更新团队结算状态失败: %w", err)
		}
	}
	return tx.Commit().Error
}

func (exa *TeamService) UpdateUserTeamSettledStatus(tx *gorm.DB, userId int) (err error) {
	err = tx.Model(&wechat.TeamRecord{}).Where("user_id = ?", userId).Update("is_settled", 1).Error
	return err
}

func (exa *TeamService) GetTeamSettlementList(searchInfo request.PageInfo, userId int) (list []wechat.TeamSettlement, total int64, err error) {
	limit := searchInfo.PageSize
	offset := searchInfo.PageSize * (searchInfo.Page - 1)

	db := global.GVA_DB.Model(&wechat.TeamSettlement{})
	err = db.Count(&total).Error
	if err != nil {
		return list, total, err
	} else {
		err = db.Limit(limit).Offset(offset).Where("user_id = ?", userId).Preload("WXUser").Find(&list).Error
	}

	return list, total, err
}

// 辅助函数：获取用户作为成员的已满足条件的成团奖励记录
func (exa *TeamService) GetUnsettledMemberGroupRewards(userId int) (*wechat.TeamConsumeRecord, error) {
	// 1. 获取用户所在的团队
	myTeam, err := exa.GetUserJoinedTeamByDBSync(userId)
	if err != nil || myTeam == nil {
		return nil, err
	}

	// 2. 检查团队是否满足成团条件
	teamMembers, err := exa.GetTeamMembers(myTeam.CaptainId, myTeam.TeamId)
	if err != nil {
		return nil, err
	}

	// 团队成员不足2人，不满足成团条件
	if len(teamMembers) < 2 {
		return nil, errors.New("must have at least two team members")
	}

	// 3. 获取该用户在这个团队中的未结算成团奖励的首购记录
	groupReward, err := exa.GetTeamFirstConsumeRecordSync(userId)
	if err != nil {
		return nil, err
	}
	return groupReward, err
}
