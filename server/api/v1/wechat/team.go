package wechat

import (
	"context"
	"cooller/server/global"
	"cooller/server/model/common/response"
	"cooller/server/model/wechat"
	wechatReq "cooller/server/model/wechat/request"
	wechatRes "cooller/server/model/wechat/response"
	"cooller/server/utils"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"go.uber.org/zap"
	"math/rand"
	"time"
)

type TeamApi struct{}

// 缓存键定义
const (
	// 团队奖励缓存，过期时间1分钟
	teamRewardCacheKey = "team:reward:%d"
	// 团队成员缓存，过期时间5分钟
	teamMembersCacheKey = "team:members:%d"
	// 消费记录缓存，过期时间5分钟
	consumeRecordsCacheKey = "team:consumes:%d"
	userJoinedTeamCacheKey = "user:joined:team:%d"
)

// 佣金比例常量定义
const (
	firstRewardRate  = 0.12 // 首推佣金比例
	repeatRewardRate = 0.20 // 复购佣金比例
	groupRewardRate  = 0.05 // 成团奖励池比例
	captainShare     = 0.5  // 队长成团奖励分成
	memberShare      = 0.25 // 成员成团奖励分成
)

func (t *TeamApi) GetTeamRecordList(c *gin.Context) {
	userId := utils.GetUserID(c)
	if userId < 1 {
		response.FailWithMessage("获取UserID失败", c)
		return
	}
	captainId := utils.GetCaptainId(c)

	joinedTeam := make([]wechatRes.TeamShowItem, 0)
	if captainId > 0 {
		teamInfo, err := teamService.GetUserJoinedTeam(userId, captainId)
		if err == nil && teamInfo != nil && teamInfo.TeamId > 0 {
			var captionInfo wechat.WXUser
			if teamInfo.CaptainWXUser.ID > 0 {
				captionInfo = teamInfo.CaptainWXUser
			} else {
				captionInfo, _ = accountService.GetWXAccountByID(captainId)
			}
			// 队长信息
			var caption wechatRes.TeamShowItem
			caption.ID = captionInfo.ID
			caption.Name = captionInfo.UserName
			caption.AvatarUrl = captionInfo.AvatarUrl
			caption.IsCaptain = 1
			caption.IsActivated = captionInfo.IsFirstPurchase
			joinedTeam = append(joinedTeam, caption)

			// 另一个队员
			teammateList, err := teamService.GetTeamRecordList(teamInfo.TeamId, captainId)
			var teammate wechatRes.TeamShowItem
			if err == nil && len(teammateList) > 1 {
				for _, item := range teammateList {
					if item.UserId != userId {
						teammate.ID = item.ID
						teammate.Name = item.WXUser.UserName
						teammate.AvatarUrl = item.WXUser.AvatarUrl
						teammate.IsActivated = item.IsActivated
						joinedTeam = append(joinedTeam, teammate)
					}
				}
			}
		}
	}

	var teamRecordList [][]wechatRes.TeamShowItem
	myTeamRecordList, err := teamService.GetMyTeamsRecordList(userId)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}

	// 初始化为空切片，而非容量为2的切片
	teamActivatedMap := map[int][]wechatRes.TeamShowItem{}
	teamNotActivatedList := make([]wechatRes.TeamShowItem, 0)

	for _, teamRecord := range myTeamRecordList {
		teamItem := wechatRes.TeamShowItem{
			ID:          teamRecord.WXUser.ID,
			AvatarUrl:   teamRecord.WXUser.AvatarUrl,
			Name:        teamRecord.WXUser.UserName,
			IsActivated: teamRecord.IsActivated,
		}
		if teamRecord.TeamId > 0 {
			teamActivatedMap[teamRecord.TeamId] = append(teamActivatedMap[teamRecord.TeamId], teamItem)
		} else {
			teamNotActivatedList = append(teamNotActivatedList, teamItem)
		}
	}
	// 先添加成团的团员
	for _, teamItems := range teamActivatedMap {
		teamRecordList = append(teamRecordList, teamItems)
	}
	// 再添加未成团的团员
	groupCount := (len(teamNotActivatedList) + 1) / 2 // 向上取整

	for i := 0; i < groupCount; i++ {
		start := i * 2
		end := start + 2

		// 防止越界
		if end > len(teamNotActivatedList) {
			end = len(teamNotActivatedList)
		}

		// 截取当前组并添加到TeamRecordList
		group := teamNotActivatedList[start:end]
		teamRecordList = append(teamRecordList, group)
	}

	response.OkWithData(wechatRes.TeamShowResponse{
		JoinedTeam: joinedTeam,
		MyTeams:    teamRecordList,
	}, c)
}

// GetTeamReward 获取团队奖励
func (t *TeamApi) GetTeamReward(c *gin.Context) {
	userId := utils.GetUserID(c)
	if userId < 1 {
		global.GVA_LOG.Error("获取userId错误!")
		response.FailWithMessage("获取userId错误!", c)
		return
	}
	captainId := utils.GetCaptainId(c)
	// 从缓存或计算获取奖励（性能优化）
	totalAmount, err := t.FetchTeamReward(userId, captainId)
	if err != nil {
		global.GVA_LOG.Error("计算团队奖励失败!", zap.Error(err))
		response.FailWithMessage("计算团队奖励失败: "+err.Error(), c)
		return
	}

	response.OkWithData(totalAmount, c)
}

// FetchTeamReward 获取团队奖励（带缓存）
func (t *TeamApi) FetchTeamReward(userId int, captainId int) (float64, error) {
	// 1. 尝试从缓存获取
	cacheKey := fmt.Sprintf(teamRewardCacheKey, userId)
	var cachedReward float64
	if err := global.GVA_REDIS.Get(context.Background(), cacheKey).Scan(&cachedReward); err == nil {
		return cachedReward, nil
	}

	// 2. 缓存未命中，计算奖励
	reward, err := t.calculateTeamReward(userId, captainId)
	if err != nil {
		return 0, err
	}
	reward = utils.FloatTo2(reward)

	// 3. 写入缓存（1分钟过期，平衡实时性和性能）
	global.GVA_REDIS.SetEx(context.Background(), cacheKey, reward, time.Minute)

	return reward, nil
}

// calculateTeamReward 实际计算奖励的内部函数
func (t *TeamApi) calculateTeamReward(userId int, captainId int) (float64, error) {
	// 批量获取数据，减少数据库查询次数
	// 1. 获取用户邀请的所有团队成员（作为队长）
	allMembers, err := t.getCachedTeamMembers(userId)
	if err != nil {
		return 0, err
	}

	// 2. 按团队ID分组
	teamMap := make(map[int][]wechat.TeamRecord)
	for _, member := range allMembers {
		if member.TeamId > 0 && member.IsActivated == 1 {
			teamMap[member.TeamId] = append(teamMap[member.TeamId], member)
		}
	}

	// 3. 通过CaptainId获取旗下所有队员的消费记录（作为队长）
	consumeMap, err := t.getCachedConsumeRecordsByCaptain(userId)
	if err != nil {
		return 0, err
	}

	// 4. 计算作为队长的奖励
	totalReward := float32(0)
	captainReward := t.calculateCaptainReward(teamMap, consumeMap)
	totalReward += captainReward

	// 5. 计算作为成员的成团奖励
	memberGroupReward, err := t.calculateMemberGroupRewardInMemory(userId, captainId)
	if err != nil {
		return 0, err
	}
	totalReward += memberGroupReward

	return float64(totalReward), nil
}

// calculateCaptainReward 计算作为队长的奖励
func (t *TeamApi) calculateCaptainReward(teamMap map[int][]wechat.TeamRecord, consumeMap map[int][]wechat.TeamConsumeRecord) float32 {
	totalReward := float32(0)
	// 处理每个团队的奖励
	for _, members := range teamMap {
		teamFirstTotal := float32(0)
		memberFirstConsumed := make(map[int]bool)

		// 计算成员带来的首推和复购奖励
		for _, member := range members {
			memberRewards, firstConsume := t.calculateMemberRewardInMemory(consumeMap[member.UserId])
			// 只有成团之后才能获得团员的复购奖励
			if len(members) > 1 {
				totalReward += memberRewards
			}

			if firstConsume > 0 {
				teamFirstTotal += firstConsume
				memberFirstConsumed[member.UserId] = true
			}
		}
		// 计算成团奖励（2人成团）
		if len(members) >= 2 && t.allMembersFirstConsumed(memberFirstConsumed) {
			groupRewardPool := teamFirstTotal * groupRewardRate
			totalReward += groupRewardPool * captainShare // 队长获得50%
		}
	}
	return totalReward
}

// 从缓存或数据库获取团队成员
func (t *TeamApi) getCachedTeamMembers(captainId int) ([]wechat.TeamRecord, error) {
	cacheKey := fmt.Sprintf(teamMembersCacheKey, captainId)
	var members []wechat.TeamRecord

	// 尝试从缓存获取
	if err := global.GVA_REDIS.Get(context.Background(), cacheKey).Scan(&members); err == nil {
		return members, nil
	}

	// 缓存未命中，从数据库获取（一次查询）
	members, err := teamService.GetTeamMembers(captainId)
	if err != nil {
		return nil, err
	}

	// 写入缓存（5分钟过期）
	global.GVA_REDIS.SetEx(context.Background(), cacheKey, members, 5*time.Minute)
	return members, nil
}

// 从缓存或数据库通过CaptainId获取旗下所有队员的消费记录
func (t *TeamApi) getCachedConsumeRecordsByCaptain(captainId int) (map[int][]wechat.TeamConsumeRecord, error) {
	consumeMap := make(map[int][]wechat.TeamConsumeRecord)
	cacheKey := fmt.Sprintf(consumeRecordsCacheKey, captainId)

	// 1. 尝试从缓存获取
	if err := global.GVA_REDIS.Get(context.Background(), cacheKey).Scan(&consumeMap); err == nil {
		return consumeMap, nil
	}

	// 2. 缓存未命中，从数据库通过CaptainId批量查询 不筛选结算状态，所有记录都可能用到
	allConsumes, err := teamService.GetTeamConsumeRecordListByUsers(captainId)
	if err != nil {
		return nil, err
	}

	// 3. 按用户ID分组
	for _, consume := range allConsumes {
		consumeMap[consume.UserId] = append(consumeMap[consume.UserId], consume)
	}

	// 4. 写入缓存（5分钟过期）
	global.GVA_REDIS.SetEx(context.Background(), cacheKey, consumeMap, 5*time.Minute)

	return consumeMap, nil
}

// 纯内存计算单个成员的奖励（无数据库操作）
func (t *TeamApi) calculateMemberRewardInMemory(consumes []wechat.TeamConsumeRecord) (float32, float32) {
	total := float32(0)
	firstConsume := float32(0)
	firstRewardCalculated := false // 标记首购奖励是否已计算

	for _, consume := range consumes {
		// 只处理未结算的普通消费记录（非成团奖励）
		// 首推佣金：未结算的首购记录
		if consume.IsFirst == 1 && !firstRewardCalculated && consume.IsFirstRewardSettled == 0 {
			reward := consume.Amount * firstRewardRate
			total += reward
			firstConsume = consume.Amount
			firstRewardCalculated = true // 确保只计算一次
		} else if consume.IsFirst == 0 && consume.IsFirstRewardSettled == 0 { // 复购佣金：未结算的非首购记录
			reward := consume.Amount * repeatRewardRate
			total += reward
		}
	}
	return total, firstConsume
}

// 检查团队所有成员是否都已首购
func (t *TeamApi) allMembersFirstConsumed(memberMap map[int]bool) bool {
	for _, consumed := range memberMap {
		if !consumed {
			return false
		}
	}
	return true
}

// 纯内存计算作为成员的成团奖励
func (t *TeamApi) calculateMemberGroupRewardInMemory(userId int, captainId int) (float32, error) {
	// 1. 直接从数据库获取用户所在的团队（用户作为成员）
	myTeam, err := teamService.GetUserJoinedTeam(userId, captainId)
	if err != nil {
		return 0, err
	}

	// 未加入任何团队，没有成团奖励
	if myTeam == nil {
		return 0, nil
	}

	// 2. 获取该团队的所有成员
	teamMembers, err := teamService.GetTeamMembers(myTeam.TeamId)
	if err != nil {
		return 0, err
	}

	// 团队成员不足2人，未成团
	if len(teamMembers) < 2 {
		return 0, nil
	}
	// 3. 检查自己是否已结算过该团队的奖励
	if myTeam.IsSettled == 1 {
		return 0, nil
	}
	// 4. 分别获取团队中每个成员的首购记录（从数据库直接查询）
	teamFirstTotal := float32(0)
	allFirstConsumed := true

	for _, member := range teamMembers {
		firstConsume, err := teamService.GetUserUnsettledFirstConsumes(member.UserId)
		if err != nil {
			return 0, err
		}
		if firstConsume == nil {
			allFirstConsumed = false
			break
		}

		teamFirstTotal += firstConsume.Amount
	}

	// 团队所有成员都完成首购才能获得成团奖励
	if !allFirstConsumed {
		return 0, nil
	}

	// 成员获得成团奖池的25%
	groupRewardPool := teamFirstTotal * captainShare
	return groupRewardPool * memberShare, nil
}

func (t *TeamApi) GetTeamReward1(c *gin.Context) {
	userId := utils.GetUserID(c)
	if userId < 1 {
		global.GVA_LOG.Error("获取userId错误!")
		response.FailWithMessage("获取userId错误!", c)
		return
	}
	captainId := utils.GetCaptainId(c)
	// 获取我邀请的所有人
	myTeamRecordList, err := teamService.GetMyTeamsRecordList(userId)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}

	totalAmount := float32(0)
	teamActivatedMap := map[int][]wechat.TeamRecord{}
	// 记录我邀请的所有激活队伍的人
	for _, teamRecord := range myTeamRecordList {
		if teamRecord.TeamId > 0 {
			teamActivatedMap[teamRecord.TeamId] = append(teamActivatedMap[teamRecord.TeamId], teamRecord)
		}
	}

	// 获取所有成员的消费记录并存入字典
	consumeRecordMap := map[int][]wechat.TeamConsumeRecord{}
	consumeRecord, err := teamService.GetTeamConsumeRecordListByUsers(userId)
	if err != nil {
		global.GVA_LOG.Error("获取消费记录失败!", zap.Error(err))
		response.FailWithMessage("获取消费记录失败!", c)
		return
	}
	for _, consume := range consumeRecord {
		consumeRecordMap[consume.UserId] = append(consumeRecordMap[consume.UserId], consume)
	}

	for _, teamItems := range teamActivatedMap {
		if len(teamItems) > 1 { // team有两个人，已经成团
			// 成团，totalAmount = 首推奖励（首购金额x12%）+ 成团奖励：成员首购总额×5%x×50% + 复购奖励（复购金额x20%）
			firstConsumeTotal := float32(0)
			loseUserMap := map[int]bool{}
			for _, teamItem := range teamItems { // 2个成员分别计算
				loseUserMap[teamItem.UserId] = true
				for _, consumeItem := range consumeRecordMap[teamItem.UserId] {
					if consumeItem.IsFirst > 0 {
						totalAmount += consumeItem.Amount * 0.12 // 首推奖励
						firstConsumeTotal += consumeItem.Amount
						loseUserMap[teamItem.UserId] = false
					} else {
						totalAmount += consumeItem.Amount * 0.2 // 复购奖励
					}
				}
			}
			for id, state := range loseUserMap {
				if state {
					firstConsume, err := teamService.GetLoseFirstConsumeRecord(id)
					if err != nil {
						global.GVA_LOG.Error("获取丢失的首购失败!", zap.Error(err))
						response.FailWithMessage("获取丢失的首购失败!", c)
						return
					}
					firstConsumeTotal += firstConsume.Amount
				}
			}
			totalAmount += firstConsumeTotal * 0.05 * 0.5 // 队长获得成团奖池的50%
		} else {
			//	未成团，totalAmount = 首推奖励（首购金额x12%）+ 复购奖励（复购金额x20%）
			for _, consumeItem := range consumeRecordMap[teamItems[0].UserId] {
				if consumeItem.IsFirst > 0 {
					totalAmount += consumeItem.Amount * 0.12 // 首推奖励
				} else {
					totalAmount += consumeItem.Amount * 0.2 // 复购奖励
				}
			}
		}
	}

	// 计算我所在team是否成团，是否分得成团奖励
	joinedTeam, err := teamService.GetUserJoinedTeam(userId, captainId)
	if err != nil {
		global.GVA_LOG.Error("获取所在队伍失败!", zap.Error(err))
		response.FailWithMessage("获取所在队伍失败!", c)
		return
	}

	teamList, err := teamService.GetTeamRecordList(joinedTeam.TeamId, captainId)
	if err != nil {
		global.GVA_LOG.Error("获取队伍失败!", zap.Error(err))
		response.FailWithMessage("获取队伍失败!", c)
		return
	}
	if len(teamList) > 1 {
		firstConsumeTotal := float32(0)
		loseUserMap := map[int]bool{}
		for _, team := range teamList {
			loseUserMap[team.UserId] = true
			firstConsume, err := teamService.GetTeamFirstConsumeRecord(team.UserId)
			if err != nil {
				global.GVA_LOG.Error("获取首购失败!", zap.Error(err))
				response.FailWithMessage("获取首购失败!", c)
				return
			}
			if firstConsume.Amount > 0 {
				loseUserMap[team.UserId] = false
				firstConsumeTotal += firstConsume.Amount
			}
		}
		for id, state := range loseUserMap {
			if state {
				firstConsume, err := teamService.GetLoseFirstConsumeRecord(id)
				if err != nil {
					global.GVA_LOG.Error("获取丢失的首购失败!", zap.Error(err))
					response.FailWithMessage("获取丢失的首购失败!", c)
					return
				}
				firstConsumeTotal += firstConsume.Amount
			}
		}
		totalAmount += firstConsumeTotal * 0.05 * 0.25 // 成员获得奖池的25%
	}
	response.OkWithData(totalAmount, c)
}

// SettlementTeamReward 执行团队奖励结算（含事务和缓存同步）
func (t *TeamApi) SettlementTeamReward(userId int, captainId int) (float32, error) {
	// 1. 先查询当前可结算金额（确保使用最新数据，跳过缓存）
	rewardAmount, err := t.calculateTeamRewardWithoutCache(userId, captainId)
	if err != nil {
		return 0, fmt.Errorf("计算可结算金额失败: %w", err)
	}

	if rewardAmount <= 0 {
		return 0, errors.New("没有可结算的奖励金额")
	}

	// 2. 开启数据库事务执行结算
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		// 2.1 创建结算记录
		settlement := wechat.TeamSettlement{
			UserId:         userId,
			SettlementNo:   t.generateSettlementNo(userId),
			TotalAmount:    rewardAmount,
			SettlementTime: time.Now(),
			Status:         1, // 直接标记为已完成
		}
		if err := tx.Create(&settlement).Error; err != nil {
			return fmt.Errorf("创建结算记录失败: %w", err)
		}
		now := settlement.SettlementTime

		// 2.2 获取用户作为队长的所有未结算消费记录
		captainConsumes, err := teamService.GetUnsettledConsumesByCaptainSync(tx, userId)
		if err != nil {
			return fmt.Errorf("查询队长消费记录失败: %w", err)
		}

		// 2.3 获取用户作为成员的未结算成团奖励记录
		memberGroupRewards, err := t.getUnsettledMemberGroupRewards(tx, userId)
		if err != nil {
			return fmt.Errorf("查询成员成团奖励失败: %w", err)
		}

		// 2.4 更新首推/复购奖励结算状态
		if len(captainConsumes) > 0 {
			err = teamService.UpdateFirstAndRepurchaseStatusSync(tx, t.extractIds(captainConsumes), now)
			if err != nil {
				return fmt.Errorf("更新首推奖励结算状态失败: %w", err)
			}
		}

		// 2.5 更新成团奖励结算状态
		if memberGroupRewards != nil {
			err = teamService.UpdateTeamSettleStatusSync(tx, memberGroupRewards.ID, now)
			if err != nil {
				return fmt.Errorf("更新成团奖励结算状态失败: %w", err)
			}

			// 更新团队结算状态
			if err := t.updateUserTeamSettledStatus(tx, userId); err != nil {
				return fmt.Errorf("更新团队结算状态失败: %w", err)
			}
		}
		return nil
	})

	if err != nil {
		return 0, err
	}

	// 3. 结算成功后刷新所有相关缓存（关键步骤）
	t.refreshSettlementRelatedCache(userId, captainId)

	return rewardAmount, nil
}

func (t *TeamApi) extractIds(consumes []wechat.TeamConsumeRecord) []int {
	ids := make([]int, len(consumes))
	for i, c := range consumes {
		ids[i] = c.ID
	}
	return ids
}

// 不使用缓存计算可结算金额（确保结算金额准确）
func (t *TeamApi) calculateTeamRewardWithoutCache(userId int, captainId int) (float32, error) {
	// 计算最新值
	members, err := teamService.GetTeamMembers(userId)
	if err != nil {
		return 0, err
	}

	teamMap := make(map[int][]wechat.TeamRecord)
	for _, member := range members {
		if member.TeamId > 0 && member.IsActivated == 1 {
			teamMap[member.TeamId] = append(teamMap[member.TeamId], member)
		}
	}

	allConsumes, err := teamService.GetTeamConsumeRecordListByUsers(userId)
	if err != nil {
		return 0, err
	}

	// 3. 按用户ID分组
	consumeMap := make(map[int][]wechat.TeamConsumeRecord)
	for _, consume := range allConsumes {
		consumeMap[consume.UserId] = append(consumeMap[consume.UserId], consume)
	}

	// 计算队长奖励
	totalReward := t.calculateCaptainReward(teamMap, consumeMap)

	// 计算成员奖励
	memberReward, err := t.calculateMemberGroupRewardInMemory(userId, captainId)
	if err != nil {
		return 0, err
	}

	return totalReward + memberReward, nil
}

// 刷新结算相关的所有缓存
func (t *TeamApi) refreshSettlementRelatedCache(userId int, captainId int) {
	// 1. 刷新用户自身的奖励缓存
	t.RefreshTeamRewardCache(userId)

	// 2. 刷新用户作为队长的消费记录缓存
	t.RefreshConsumeRecordCache(userId)

	// 3. 刷新用户作为队长的团队成员缓存
	t.RefreshTeamMemberCache(userId)

	// 4. 刷新用户所在团队的缓存（如果用户是成员）
	//t.RefreshUserTeamCache(userId)

	// 5. 如果用户是成员，还需要刷新其所在团队的队长的相关缓存
	if team, err := teamService.GetUserJoinedTeam(userId, captainId); err == nil && team != nil {
		t.RefreshTeamRewardCache(team.CaptainId)
		t.RefreshConsumeRecordCache(team.CaptainId)
	}
}

// 辅助函数：获取用户作为成员的已满足条件的成团奖励记录
func (t *TeamApi) getUnsettledMemberGroupRewards(tx *gorm.DB, userId int) (*wechat.TeamConsumeRecord, error) {
	// 1. 获取用户所在的团队
	myTeam, err := teamService.GetUserJoinedTeamByDBSync(tx, userId)
	if err != nil || myTeam == nil {
		return nil, err
	}

	// 2. 检查团队是否满足成团条件
	teamMembers, err := teamService.GetTeamMembersByDBSync(tx, myTeam.TeamId)
	if err != nil {
		return nil, err
	}

	// 团队成员不足2人，不满足成团条件
	if len(teamMembers) < 2 {
		return nil, nil
	}

	// 3. 获取该用户在这个团队中的未结算成团奖励的首购记录
	groupReward, err := teamService.GetTeamFirstConsumeRecordSync(tx, userId)
	if err != nil {
		return nil, err
	}
	return groupReward, err

}

// 辅助函数：更新用户参与的团队结算状态
func (t *TeamApi) updateUserTeamSettledStatus(tx *gorm.DB, userId int) error {
	// 更新用户作为成员的团队结算状态
	return tx.Model(&wechat.TeamRecord{}).
		Where("user_id = ?", userId).
		Update("is_settled", 1).Error
}

// 辅助函数：生成唯一结算单号
func (t *TeamApi) generateSettlementNo(userId int) string {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("SETTLE%d%s%04d",
		userId,
		time.Now().Format("20060102150405"),
		10000+random.Intn(90000))
}

// RefreshTeamRewardCache 当数据更新时主动刷新缓存（在添加消费记录、成团等操作后调用）
func (t *TeamApi) RefreshTeamRewardCache(userId int) {
	global.GVA_REDIS.Del(context.Background(), fmt.Sprintf(teamRewardCacheKey, userId))
}

func (t *TeamApi) RefreshConsumeRecordCache(captainId int) {
	global.GVA_REDIS.Del(context.Background(), fmt.Sprintf(consumeRecordsCacheKey, captainId))
}

// RefreshTeamMemberCache 批量刷新成员相关缓存
func (t *TeamApi) RefreshTeamMemberCache(captainId int) {
	global.GVA_REDIS.Del(context.Background(), fmt.Sprintf(teamMembersCacheKey, captainId))
}

// GetTeamConsumeDetails 获取团队奖励详情
func (t *TeamApi) GetTeamConsumeDetails(c *gin.Context) {
	var userInfo wechatReq.DetailsInfo
	err := c.ShouldBindJSON(&userInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	userId := utils.GetUserID(c)
	if userId < 1 {
		global.GVA_LOG.Error("获取userId错误!")
		response.FailWithMessage("获取userId错误!", c)
		return
	}

	fmt.Println(userInfo.UserIds)
	consumeList, err := teamService.GetTeamsConsumeList(userId, userInfo.UserIds)
	if err != nil {
		global.GVA_LOG.Error("获取消费列表错误!")
		response.FailWithMessage("获取消费列表错误!", c)
		return
	}
	detailsMap := make(map[string][]wechatRes.ConsumeData)
	for _, consume := range consumeList {
		reward := float64(0)
		if consume.IsFirst > 0 {
			reward = float64(consume.Amount) * firstRewardRate
		} else {
			reward = float64(consume.Amount) * repeatRewardRate
		}
		reward = utils.FloatTo2(reward)
		// 格式化日期（统一格式，确保后续排序正确）
		dateStr := consume.CreatedAt.Format("01-02 15:04")
		consumeData := wechatRes.ConsumeData{
			Date:  dateStr,
			Value: reward,
		}
		detailsMap[consume.WXUser.UserName] = append(detailsMap[consume.WXUser.UserName], consumeData)
	}

	response.OkWithData(wechatRes.ConsumeDetailsResponse{
		Details: detailsMap,
	}, c)
}
