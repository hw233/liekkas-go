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

type CfgBattlePosition struct {
	Id         int32
	IsPosition bool
	Pos1       []float64
	Enter1     []float64
	Facing1    bool
	Buff1      int32
	Pos2       []float64
	Enter2     []float64
	Facing2    bool
	Buff2      int32
	Pos3       []float64
	Enter3     []float64
	Facing3    bool
	Buff3      int32
	Pos4       []float64
	Enter4     []float64
	Facing4    bool
	Buff4      int32
	Pos5       []float64
	Enter5     []float64
	Facing5    bool
	Buff5      int32
	Pos6       []float64
	Enter6     []float64
	Facing6    bool
	Buff6      int32
	Delay      float64
}

type CfgBattlePositionConfig struct {
	data map[int32]*CfgBattlePosition
}

func NewCfgBattlePositionConfig() *CfgBattlePositionConfig {
	return &CfgBattlePositionConfig{
		data: make(map[int32]*CfgBattlePosition),
	}
}

func (c *CfgBattlePositionConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(CfgBattlePosition)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse CfgBattlePosition.Id field error,value:", vId)
			return false
		}

		/* parse IsPosition field */
		vIsPosition, _ := parse.GetFieldByName(uint32(i), "isPosition")
		var IsPositionRet bool
		data.IsPosition, IsPositionRet = String2Bool(vIsPosition)
		if !IsPositionRet {
			glog.Error("Parse CfgBattlePosition.IsPosition field error,value:", vIsPosition)
		}

		/* parse Pos1 field */
		vecPos1, _ := parse.GetFieldByName(uint32(i), "pos1")
		arrayPos1 := strings.Split(vecPos1, ",")
		for j := 0; j < len(arrayPos1); j++ {
			v, ret := String2Float(arrayPos1[j])
			if !ret {
				glog.Error("Parse CfgBattlePosition.Pos1 field error,value:", arrayPos1[j])
				return false
			}
			data.Pos1 = append(data.Pos1, v)
		}

		/* parse Enter1 field */
		vecEnter1, _ := parse.GetFieldByName(uint32(i), "enter1")
		arrayEnter1 := strings.Split(vecEnter1, ",")
		for j := 0; j < len(arrayEnter1); j++ {
			v, ret := String2Float(arrayEnter1[j])
			if !ret {
				glog.Error("Parse CfgBattlePosition.Enter1 field error,value:", arrayEnter1[j])
				return false
			}
			data.Enter1 = append(data.Enter1, v)
		}

		/* parse Facing1 field */
		vFacing1, _ := parse.GetFieldByName(uint32(i), "facing1")
		var Facing1Ret bool
		data.Facing1, Facing1Ret = String2Bool(vFacing1)
		if !Facing1Ret {
			glog.Error("Parse CfgBattlePosition.Facing1 field error,value:", vFacing1)
		}

		/* parse Buff1 field */
		vBuff1, _ := parse.GetFieldByName(uint32(i), "buff1")
		var Buff1Ret bool
		data.Buff1, Buff1Ret = String2Int32(vBuff1)
		if !Buff1Ret {
			glog.Error("Parse CfgBattlePosition.Buff1 field error,value:", vBuff1)
			return false
		}

		/* parse Pos2 field */
		vecPos2, _ := parse.GetFieldByName(uint32(i), "pos2")
		arrayPos2 := strings.Split(vecPos2, ",")
		for j := 0; j < len(arrayPos2); j++ {
			v, ret := String2Float(arrayPos2[j])
			if !ret {
				glog.Error("Parse CfgBattlePosition.Pos2 field error,value:", arrayPos2[j])
				return false
			}
			data.Pos2 = append(data.Pos2, v)
		}

		/* parse Enter2 field */
		vecEnter2, _ := parse.GetFieldByName(uint32(i), "enter2")
		arrayEnter2 := strings.Split(vecEnter2, ",")
		for j := 0; j < len(arrayEnter2); j++ {
			v, ret := String2Float(arrayEnter2[j])
			if !ret {
				glog.Error("Parse CfgBattlePosition.Enter2 field error,value:", arrayEnter2[j])
				return false
			}
			data.Enter2 = append(data.Enter2, v)
		}

		/* parse Facing2 field */
		vFacing2, _ := parse.GetFieldByName(uint32(i), "facing2")
		var Facing2Ret bool
		data.Facing2, Facing2Ret = String2Bool(vFacing2)
		if !Facing2Ret {
			glog.Error("Parse CfgBattlePosition.Facing2 field error,value:", vFacing2)
		}

		/* parse Buff2 field */
		vBuff2, _ := parse.GetFieldByName(uint32(i), "buff2")
		var Buff2Ret bool
		data.Buff2, Buff2Ret = String2Int32(vBuff2)
		if !Buff2Ret {
			glog.Error("Parse CfgBattlePosition.Buff2 field error,value:", vBuff2)
			return false
		}

		/* parse Pos3 field */
		vecPos3, _ := parse.GetFieldByName(uint32(i), "pos3")
		arrayPos3 := strings.Split(vecPos3, ",")
		for j := 0; j < len(arrayPos3); j++ {
			v, ret := String2Float(arrayPos3[j])
			if !ret {
				glog.Error("Parse CfgBattlePosition.Pos3 field error,value:", arrayPos3[j])
				return false
			}
			data.Pos3 = append(data.Pos3, v)
		}

		/* parse Enter3 field */
		vecEnter3, _ := parse.GetFieldByName(uint32(i), "enter3")
		arrayEnter3 := strings.Split(vecEnter3, ",")
		for j := 0; j < len(arrayEnter3); j++ {
			v, ret := String2Float(arrayEnter3[j])
			if !ret {
				glog.Error("Parse CfgBattlePosition.Enter3 field error,value:", arrayEnter3[j])
				return false
			}
			data.Enter3 = append(data.Enter3, v)
		}

		/* parse Facing3 field */
		vFacing3, _ := parse.GetFieldByName(uint32(i), "facing3")
		var Facing3Ret bool
		data.Facing3, Facing3Ret = String2Bool(vFacing3)
		if !Facing3Ret {
			glog.Error("Parse CfgBattlePosition.Facing3 field error,value:", vFacing3)
		}

		/* parse Buff3 field */
		vBuff3, _ := parse.GetFieldByName(uint32(i), "buff3")
		var Buff3Ret bool
		data.Buff3, Buff3Ret = String2Int32(vBuff3)
		if !Buff3Ret {
			glog.Error("Parse CfgBattlePosition.Buff3 field error,value:", vBuff3)
			return false
		}

		/* parse Pos4 field */
		vecPos4, _ := parse.GetFieldByName(uint32(i), "pos4")
		arrayPos4 := strings.Split(vecPos4, ",")
		for j := 0; j < len(arrayPos4); j++ {
			v, ret := String2Float(arrayPos4[j])
			if !ret {
				glog.Error("Parse CfgBattlePosition.Pos4 field error,value:", arrayPos4[j])
				return false
			}
			data.Pos4 = append(data.Pos4, v)
		}

		/* parse Enter4 field */
		vecEnter4, _ := parse.GetFieldByName(uint32(i), "enter4")
		arrayEnter4 := strings.Split(vecEnter4, ",")
		for j := 0; j < len(arrayEnter4); j++ {
			v, ret := String2Float(arrayEnter4[j])
			if !ret {
				glog.Error("Parse CfgBattlePosition.Enter4 field error,value:", arrayEnter4[j])
				return false
			}
			data.Enter4 = append(data.Enter4, v)
		}

		/* parse Facing4 field */
		vFacing4, _ := parse.GetFieldByName(uint32(i), "facing4")
		var Facing4Ret bool
		data.Facing4, Facing4Ret = String2Bool(vFacing4)
		if !Facing4Ret {
			glog.Error("Parse CfgBattlePosition.Facing4 field error,value:", vFacing4)
		}

		/* parse Buff4 field */
		vBuff4, _ := parse.GetFieldByName(uint32(i), "buff4")
		var Buff4Ret bool
		data.Buff4, Buff4Ret = String2Int32(vBuff4)
		if !Buff4Ret {
			glog.Error("Parse CfgBattlePosition.Buff4 field error,value:", vBuff4)
			return false
		}

		/* parse Pos5 field */
		vecPos5, _ := parse.GetFieldByName(uint32(i), "pos5")
		arrayPos5 := strings.Split(vecPos5, ",")
		for j := 0; j < len(arrayPos5); j++ {
			v, ret := String2Float(arrayPos5[j])
			if !ret {
				glog.Error("Parse CfgBattlePosition.Pos5 field error,value:", arrayPos5[j])
				return false
			}
			data.Pos5 = append(data.Pos5, v)
		}

		/* parse Enter5 field */
		vecEnter5, _ := parse.GetFieldByName(uint32(i), "enter5")
		arrayEnter5 := strings.Split(vecEnter5, ",")
		for j := 0; j < len(arrayEnter5); j++ {
			v, ret := String2Float(arrayEnter5[j])
			if !ret {
				glog.Error("Parse CfgBattlePosition.Enter5 field error,value:", arrayEnter5[j])
				return false
			}
			data.Enter5 = append(data.Enter5, v)
		}

		/* parse Facing5 field */
		vFacing5, _ := parse.GetFieldByName(uint32(i), "facing5")
		var Facing5Ret bool
		data.Facing5, Facing5Ret = String2Bool(vFacing5)
		if !Facing5Ret {
			glog.Error("Parse CfgBattlePosition.Facing5 field error,value:", vFacing5)
		}

		/* parse Buff5 field */
		vBuff5, _ := parse.GetFieldByName(uint32(i), "buff5")
		var Buff5Ret bool
		data.Buff5, Buff5Ret = String2Int32(vBuff5)
		if !Buff5Ret {
			glog.Error("Parse CfgBattlePosition.Buff5 field error,value:", vBuff5)
			return false
		}

		/* parse Pos6 field */
		vecPos6, _ := parse.GetFieldByName(uint32(i), "pos6")
		arrayPos6 := strings.Split(vecPos6, ",")
		for j := 0; j < len(arrayPos6); j++ {
			v, ret := String2Float(arrayPos6[j])
			if !ret {
				glog.Error("Parse CfgBattlePosition.Pos6 field error,value:", arrayPos6[j])
				return false
			}
			data.Pos6 = append(data.Pos6, v)
		}

		/* parse Enter6 field */
		vecEnter6, _ := parse.GetFieldByName(uint32(i), "enter6")
		arrayEnter6 := strings.Split(vecEnter6, ",")
		for j := 0; j < len(arrayEnter6); j++ {
			v, ret := String2Float(arrayEnter6[j])
			if !ret {
				glog.Error("Parse CfgBattlePosition.Enter6 field error,value:", arrayEnter6[j])
				return false
			}
			data.Enter6 = append(data.Enter6, v)
		}

		/* parse Facing6 field */
		vFacing6, _ := parse.GetFieldByName(uint32(i), "facing6")
		var Facing6Ret bool
		data.Facing6, Facing6Ret = String2Bool(vFacing6)
		if !Facing6Ret {
			glog.Error("Parse CfgBattlePosition.Facing6 field error,value:", vFacing6)
		}

		/* parse Buff6 field */
		vBuff6, _ := parse.GetFieldByName(uint32(i), "buff6")
		var Buff6Ret bool
		data.Buff6, Buff6Ret = String2Int32(vBuff6)
		if !Buff6Ret {
			glog.Error("Parse CfgBattlePosition.Buff6 field error,value:", vBuff6)
			return false
		}

		/* parse Delay field */
		vDelay, _ := parse.GetFieldByName(uint32(i), "delay")
		var DelayRet bool
		data.Delay, DelayRet = String2Float(vDelay)
		if !DelayRet {
			glog.Error("Parse CfgBattlePosition.Delay field error,value:", vDelay)
		}

		if _, ok := c.data[data.Id]; ok {
			glog.Errorf("Find %d repeated", data.Id)
			return false
		}
		c.data[data.Id] = data
	}

	return true
}

func (c *CfgBattlePositionConfig) Clear() {
}

func (c *CfgBattlePositionConfig) Find(id int32) (*CfgBattlePosition, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *CfgBattlePositionConfig) GetAllData() map[int32]*CfgBattlePosition {
	return c.data
}

func (c *CfgBattlePositionConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id, ",", v.IsPosition, ",", v.Pos1, ",", v.Enter1, ",", v.Facing1, ",", v.Buff1, ",", v.Pos2, ",", v.Enter2, ",", v.Facing2, ",", v.Buff2, ",", v.Pos3, ",", v.Enter3, ",", v.Facing3, ",", v.Buff3, ",", v.Pos4, ",", v.Enter4, ",", v.Facing4, ",", v.Buff4, ",", v.Pos5, ",", v.Enter5, ",", v.Facing5, ",", v.Buff5, ",", v.Pos6, ",", v.Enter6, ",", v.Facing6, ",", v.Buff6, ",", v.Delay)
	}
}
