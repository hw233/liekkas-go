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
	"strings"
)

type CfgExploreChapter struct {
	Id              int32
	ChapterName     string
	ChapterType     int32
	UnlockCondition []string
	InitPos         []float64
	MapSize         []int32
}

type CfgExploreChapterConfig struct {
	data map[int32]*CfgExploreChapter
}

func NewCfgExploreChapterConfig() *CfgExploreChapterConfig {
	return &CfgExploreChapterConfig{
		data: make(map[int32]*CfgExploreChapter),
	}
}

func (c *CfgExploreChapterConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(CfgExploreChapter)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse CfgExploreChapter.Id field error,value:", vId)
			return false
		}

		/* parse ChapterName field */
		data.ChapterName, _ = parse.GetFieldByName(uint32(i), "chapterName")

		/* parse ChapterType field */
		vChapterType, _ := parse.GetFieldByName(uint32(i), "chapterType")
		var ChapterTypeRet bool
		data.ChapterType, ChapterTypeRet = String2Int32(vChapterType)
		if !ChapterTypeRet {
			glog.Error("Parse CfgExploreChapter.ChapterType field error,value:", vChapterType)
			return false
		}

		/* parse UnlockCondition field */
		vecUnlockCondition, _ := parse.GetFieldByName(uint32(i), "unlockCondition")
		arrayUnlockCondition := strings.Split(vecUnlockCondition, ",")
		for j := 0; j < len(arrayUnlockCondition); j++ {
			v := arrayUnlockCondition[j]
			data.UnlockCondition = append(data.UnlockCondition, v)
		}

		/* parse InitPos field */
		vecInitPos, _ := parse.GetFieldByName(uint32(i), "initPos")
		arrayInitPos := strings.Split(vecInitPos, ",")
		for j := 0; j < len(arrayInitPos); j++ {
			v, ret := String2Float(arrayInitPos[j])
			if !ret {
				glog.Error("Parse CfgExploreChapter.InitPos field error,value:", arrayInitPos[j])
				return false
			}
			data.InitPos = append(data.InitPos, v)
		}

		/* parse MapSize field */
		vecMapSize, _ := parse.GetFieldByName(uint32(i), "mapSize")
		if vecMapSize != "" {
			arrayMapSize := strings.Split(vecMapSize, ",")
			for j := 0; j < len(arrayMapSize); j++ {
				v, ret := String2Int32(arrayMapSize[j])
				if !ret {
					glog.Error("Parse CfgExploreChapter.MapSize field error, value:", arrayMapSize[j])
					return false
				}
				data.MapSize = append(data.MapSize, v)
			}
		}

		if _, ok := c.data[data.Id]; ok {
			glog.Errorf("Find %d repeated", data.Id)
			return false
		}
		c.data[data.Id] = data
	}

	return true
}

func (c *CfgExploreChapterConfig) Clear() {
}

func (c *CfgExploreChapterConfig) Find(id int32) (*CfgExploreChapter, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *CfgExploreChapterConfig) GetAllData() map[int32]*CfgExploreChapter {
	return c.data
}

func (c *CfgExploreChapterConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id, ",", v.ChapterName, ",", v.ChapterType, ",", v.UnlockCondition, ",", v.InitPos, ",", v.MapSize)
	}
}
