package model // number event

import (
	"shared/common"
)

func (u *User) expChangeEvent(opt *common.ChangedEventOpt) {
	u.BIExpLevel("exp", opt.NowValue, u.Info.GetExp(), opt.Reason)
}

func (u *User) levelChangeEvent(opt *common.ChangedEventOpt) {
	u.BIExpLevel("level", opt.NowValue, u.Info.GetExp(), opt.Reason)
}

func (u *User) genResourceEvent(resourceId int32) func(*common.ChangedEventOpt) {
	return func(opt *common.ChangedEventOpt) {
		u.BIResourceChange(resourceId, int64(opt.Diff), int64(opt.NowValue), opt.Reason)
	}
}
