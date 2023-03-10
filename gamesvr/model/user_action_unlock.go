package model

import (
	"gamesvr/manager"
	"shared/utility/errors"
)

// 入参id为cons_action_id_type.go中各个模块对应的ActionIdType
func (u *User) CheckActionUnlock(id int32) error {
	unlockConditions, err := manager.CSV.ActionUnlock.GetUnlockConditions(id)
	if err != nil {
		return errors.WrapTrace(err)
	}

	err = u.CheckUserConditions(unlockConditions.UnlockCondition)
	if err != nil {
		//return false
		return errors.WrapTrace(err)
	}

	return nil
}
