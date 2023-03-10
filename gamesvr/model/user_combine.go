package model

import (
	"encoding/json"

	"shared/utility/number"
)

type UserCombine struct {
	GuildID   int64  `json:"guild_id"`
	GuildName string `json:"guild_name"`

	GuideID int32 `json:"guide_id"`

	Gold        int32 `json:"gold"`
	DiamondGift int32 `json:"diamond_gift"`
	DiamondCash int32 `json:"diamond_cash"`
	Energy      int32 `json:"energy"`

	Drop         *Drop             `json:"drop_info"`
	SignIn       *SignIn           `json:"sign_in"`
	SurveyRecord *number.BitNumber `json:"survey_record"`
	GachaRecords *GachaRecords     `json:"-"`
}

func NewUserCombine() *UserCombine {
	return &UserCombine{
		GuildID:      0,
		GuildName:    "",
		GuideID:      0,
		Gold:         0,
		DiamondGift:  0,
		DiamondCash:  0,
		Energy:       0,
		Drop:         NewDrop(),
		SignIn:       NewSignIn(),
		SurveyRecord: number.NewBitNumber(),
		GachaRecords: NewGachaRecords(),
	}
}

func (ui *UserCombine) Marshal() ([]byte, error) {
	return json.Marshal(ui)
}

func (ui *UserCombine) Unmarshal(data []byte) error {
	if len(data) == 0 {
		return nil
	}

	err := json.Unmarshal(data, ui)
	if err != nil {
		return err
	}

	return nil
}
