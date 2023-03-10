package model

import (
	"gamesvr/manager"
	"shared/common"
	"shared/csv/static"
	"shared/statistic/logreason"
	"shared/utility/errors"
	"shared/utility/servertime"
	"time"
)

// 首先获得当前时间点有可能需要签到的id，然后逐个检查每个id今天是否需要签到
func (u *User) CheckHasSignInWrap() ([]int32, error) {
	// 没有解锁签到功能，直接返回nil
	err := u.CheckActionUnlock(static.ActionIdTypeSigninunlock)
	if err != nil {
		return nil, nil
	}

	tm := servertime.Now().Unix()
	idBeforeCheck := manager.CSV.SignIn.GetID(tm, int64(manager.CSV.GlobalEntry.DailyRefreshTimeOffset()/time.Second))

	// if len(idBeforeCheck) < 1 {
	// 	return nil, errors.Swrapf(common.ErrSignInNoDataForToday, time.Unix(tm, 0))
	// }

	ids := make([]int32, 0, len(idBeforeCheck))

	for id, signInType := range idBeforeCheck {

		if signInType != 1 {
			continue // 如果是其他类型，暂时不作处理
		}

		// 更新标志位hasSigned
		u.Info.SignIn.CheckForSignIn(id, signInType, tm)
		ids = append(ids, id)

	}

	return ids, nil
}

func (u *User) ReceiveDailyReward(id int32) error {

	index := u.Info.SignIn.SignInGroups[id].Record
	dropID, err := manager.CSV.SignIn.GetDropID(id, index)
	if err != nil {
		return errors.WrapTrace(err)
	}

	reason := logreason.NewReason(logreason.SignIn)
	_, err = u.AddRewardsByDropId(dropID, reason)
	if err != nil {
		return errors.WrapTrace(err)
	}

	return nil

}

func (u *User) SignIn(id int32) error {
	// 没有解锁签到功能，如果还试图签到，会报错
	// if !u.checkActionUnlock(static.ActionIdTypeSigninunlock) {
	// 	return errors.Swrapf(common.ErrActionUnlockNotsatisfied, static.ActionIdTypeSigninunlock)
	// }
	err := u.CheckActionUnlock(static.ActionIdTypeSigninunlock)
	if err != nil {
		return err
	}

	signinGroup, ok := u.Info.SignIn.SignInGroups[id]
	if !ok {
		return errors.Swrapf(common.ErrSignInWrongIDForSignInGroups, id)
	}

	if signinGroup.HasSignIn {
		return errors.Swrapf(common.ErrSignInUserHasSigned, id)
	}

	err = u.ReceiveDailyReward(id)
	if err != nil {
		return err
	}

	signinGroup.DoSignIn()

	u.Info.SignIn.SignInGroups[id] = signinGroup

	return nil
}
