// ===================================== //
// author:       gavingqf                //
// == Please don'g change me by hand ==  //
//====================================== //

/*you have defined the following interface:
type IConfig interface {
	// load interface
	Load(path string) bool

	// clear interface
	Clear()
}
*/

package base

import (
	"shared/utility/glog"
)

type CfgExploreChapterReward struct {
	Id         int32
	ChapterId  int32
	Difficulty int32
	Number     int32
	DropId     int32
}

type CfgExploreChapterRewardConfig struct {
	data map[int32]*CfgExploreChapterReward
}

func NewCfgExploreChapterRewardConfig() *CfgExploreChapterRewardConfig {
	return &CfgExploreChapterRewardConfig{
		data: make(map[int32]*CfgExploreChapterReward),
	}
}

func (c *CfgExploreChapterRewardConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(CfgExploreChapterReward)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse CfgExploreChapterReward.Id field error,value:", vId)
			return false
		}

		/* parse ChapterId field */
		vChapterId, _ := parse.GetFieldByName(uint32(i), "chapterId")
		var ChapterIdRet bool
		data.ChapterId, ChapterIdRet = String2Int32(vChapterId)
		if !ChapterIdRet {
			glog.Error("Parse CfgExploreChapterReward.ChapterId field error,value:", vChapterId)
			return false
		}

		/* parse Difficulty field */
		vDifficulty, _ := parse.GetFieldByName(uint32(i), "difficulty")
		var DifficultyRet bool
		data.Difficulty, DifficultyRet = String2Int32(vDifficulty)
		if !DifficultyRet {
			glog.Error("Parse CfgExploreChapterReward.Difficulty field error,value:", vDifficulty)
			return false
		}

		/* parse Number field */
		vNumber, _ := parse.GetFieldByName(uint32(i), "number")
		var NumberRet bool
		data.Number, NumberRet = String2Int32(vNumber)
		if !NumberRet {
			glog.Error("Parse CfgExploreChapterReward.Number field error,value:", vNumber)
			return false
		}

		/* parse DropId field */
		vDropId, _ := parse.GetFieldByName(uint32(i), "dropId")
		var DropIdRet bool
		data.DropId, DropIdRet = String2Int32(vDropId)
		if !DropIdRet {
			glog.Error("Parse CfgExploreChapterReward.DropId field error,value:", vDropId)
			return false
		}

		if _, ok := c.data[data.Id]; ok {
			glog.Errorf("Find %d repeated", data.Id)
			return false
		}
		c.data[data.Id] = data
	}

	return true
}

func (c *CfgExploreChapterRewardConfig) Clear() {
}

func (c *CfgExploreChapterRewardConfig) Find(id int32) (*CfgExploreChapterReward, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *CfgExploreChapterRewardConfig) GetAllData() map[int32]*CfgExploreChapterReward {
	return c.data
}

func (c *CfgExploreChapterRewardConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id, ",", v.ChapterId, ",", v.Difficulty, ",", v.Number, ",", v.DropId)
	}
}
