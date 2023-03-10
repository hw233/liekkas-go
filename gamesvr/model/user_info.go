package model

import (
	"encoding/json"
	"math"
	"shared/common"
	"shared/csv/static"
	"shared/utility/servertime"
	"time"

	"gamesvr/manager"
	"shared/utility/number"
)

type UserInfo struct {
	*DailyRefreshChecker

	Exp   *common.BIEventNumber `json:"exp"`
	Level *common.BIEventNumber `json:"level"`

	NameIndex     int32  `json:"name_index"`      // 昵称后缀
	IsFirstRename bool   `json:"is_first_rename"` // 是否是第一次改名，第一次改名免费
	Avatar        int32  `json:"avatar"`          // 头像
	Frame         int32  `json:"frame"`           // 选择的头像框
	Signature     string `json:"signature"`       // 签名
	Kanban        int32  `json:"kanban"`          // 看板娘
	Power         int32  `json:"power"`           // 总战力
	RegisterAt    int64  `json:"register_at"`     // 注册日期
	LastLoginTime int64  `json:"last_login_time"` // 上次登录时间

	Gold        *common.BIEventNumber `json:"gold"`
	DiamondGift *common.BIEventNumber `json:"diamond_gift"`
	DiamondCash *common.BIEventNumber `json:"diamond_cash"`
	Energy      *number.TimerNumber   `json:"energy"`

	Drop           *Drop                    `json:"drop_info"`
	SignIn         *SignIn                  `json:"sign_in"`
	SurveyRecord   *number.BitNumber        `json:"survey_record"`
	StaminaRecord  *StaminaRecord           `json:"stamina_record"`
	RewardedVnIds  *number.NonRepeatableArr `json:"rewarded_vn_ids"`
	PassedGuideIds *number.NonRepeatableArr `json:"passed_guide_ids"`
	Ap             *number.TimerNumber      `json:"ap"`

	CardShow               *common.VisitingCardShow `json:"card_show"` // 名片展示
	LatestExploreLevel     int32                    `json:"latest_explore_level"`
	ExploredLevelVns       *number.NonRepeatableArr `json:"explored_level_vns"`
	DailyGiftInfoRewarded  []bool                   `json:"daily_gift_info_rewarded"` //日常奖励领取情况
	AchievementFinishCount int32                    `json:"achievement_finish_count"`
	ManualCount            int32                    `json:"manual_count"` // 总获得的图鉴数（角色和世界级道具）
	LastDailyRefreshTime   int64                    `json:"last_daily_refresh_time"`
	TotalLoginDay          int32                    `json:"total_login_day"`
	TotalLoginTime         int64                    `json:"total_login_time"`
}

func NewUserInfo() *UserInfo {
	itemEntry := manager.CSV.Item
	GoldLimit := itemEntry.GetLimit(static.CommonResourceTypeMoney)
	DiamondGiftLimit := itemEntry.GetLimit(static.CommonResourceTypeDiamondGift)
	DiamondCashLimit := itemEntry.GetLimit(static.CommonResourceTypeDiamondCash)
	EnergyLimit := itemEntry.GetLimit(static.CommonResourceTypeEnergy)
	ApLimit := itemEntry.GetLimit(static.CommonResourceTypeAp)

	return &UserInfo{
		DailyRefreshChecker:    NewDailyRefreshChecker(),
		NameIndex:              0,
		IsFirstRename:          true,
		Exp:                    common.NewBIEventNumber(0, 0, math.MaxInt32),
		Level:                  common.NewBIEventNumber(1, 0, math.MaxInt32),
		Avatar:                 1811012, // 初始头像
		Frame:                  1831001, //初始头像框
		Signature:              "",
		Kanban:                 1012, // 初始看板娘
		Power:                  0,
		RegisterAt:             servertime.Now().Unix(),
		LastLoginTime:          servertime.Now().Unix(),
		Gold:                   common.NewBIEventNumber(0, 0, GoldLimit),
		DiamondGift:            common.NewBIEventNumber(0, 0, DiamondGiftLimit),
		DiamondCash:            common.NewBIEventNumber(0, 0, DiamondCashLimit),
		Energy:                 number.NewTimerNumber(0, 0, EnergyLimit),
		Drop:                   NewDrop(),
		SignIn:                 NewSignIn(),
		SurveyRecord:           number.NewBitNumber(),
		StaminaRecord:          NewStaminaRecord(),
		RewardedVnIds:          number.NewNonRepeatableArr(),
		PassedGuideIds:         number.NewNonRepeatableArr(),
		Ap:                     number.NewTimerNumber(0, 0, ApLimit),
		CardShow:               common.NewVisitingCardShow(),
		LatestExploreLevel:     0,
		ExploredLevelVns:       number.NewNonRepeatableArr(),
		AchievementFinishCount: 0,
		ManualCount:            0,
		DailyGiftInfoRewarded:  []bool{false, false, false},
	}
}

func (ui *UserInfo) Marshal() ([]byte, error) {
	return json.Marshal(ui)
}

func (ui *UserInfo) Unmarshal(data []byte) error {
	if len(data) == 0 {
		return nil
	}

	err := json.Unmarshal(data, ui)
	if err != nil {
		return err
	}

	return nil
}

func (ui *UserInfo) GetLevel() int32 {
	return ui.Level.Value()
}

func (ui *UserInfo) GetExp() int32 {
	return ui.Exp.Value()
}

func (ui *UserInfo) SetLatestExploreLevel(levelId int32) {
	ui.LatestExploreLevel = levelId

}
func (ui *UserInfo) SetExploredLevelVns(levelId int32) {
	level, err := manager.CSV.LevelsEntry.GetLevel(levelId)
	if err != nil {
		return
	}
	if level.ChapterStory > 0 {
		ui.ExploredLevelVns.Append(level.ChapterStory)
	}
	if level.WinStory > 0 {
		ui.ExploredLevelVns.Append(level.WinStory)
	}

}

func (ui *UserInfo) GetLatestExploreLevel() int32 {
	return ui.LatestExploreLevel
}

func (ui *UserInfo) AddAchievementFinishCount(count int32) {
	ui.AchievementFinishCount = ui.AchievementFinishCount + count
}

func (ui *UserInfo) SetAchievementFinishCount(count int32) {
	ui.AchievementFinishCount = count
}

func (ui *UserInfo) GetAchievementFinishCount() int32 {
	return ui.AchievementFinishCount
}

func (ui *UserInfo) RecordLoginDay() {
	ui.TotalLoginDay = ui.TotalLoginDay + 1
}

func (ui *UserInfo) OnOffline() {
	duration := time.Now().Unix() - ui.LastLoginTime
	ui.TotalLoginTime = ui.TotalLoginTime + duration
}
