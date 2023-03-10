package entry

import "testing"

func TestUserSettingEntry_GetDailyRewardByIndex(t *testing.T) {
	reward, err := CSV.UserSetting.GetDailyRewardByIndex(1)
	t.Log(reward)
	t.Log(err)

}
