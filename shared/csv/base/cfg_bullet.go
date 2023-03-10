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

type CfgBullet struct {
	Id               int32
	Type             int32
	SpawnCenter      []float64
	SpawnPosition    []float64
	Delay            float64
	Duration         float64
	ParabolaHigh     float64
	ParabolaAngle    int32
	ParabolaOffset   []float64
	DisorderAngle    []float64
	DisorderSpdAngle float64
	Speed            int32
	Acceleration     int32
	Length           float64
	PassEffects      []int32
	HitEffects       []int32
	EndEffects       []int32
	Pass             int32
	PassTimes        int32
}

type CfgBulletConfig struct {
	data map[int32]*CfgBullet
}

func NewCfgBulletConfig() *CfgBulletConfig {
	return &CfgBulletConfig{
		data: make(map[int32]*CfgBullet),
	}
}

func (c *CfgBulletConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(CfgBullet)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse CfgBullet.Id field error,value:", vId)
			return false
		}

		/* parse Type field */
		vType, _ := parse.GetFieldByName(uint32(i), "type")
		var TypeRet bool
		data.Type, TypeRet = String2Int32(vType)
		if !TypeRet {
			glog.Error("Parse CfgBullet.Type field error,value:", vType)
			return false
		}

		/* parse SpawnCenter field */
		vecSpawnCenter, _ := parse.GetFieldByName(uint32(i), "spawnCenter")
		arraySpawnCenter := strings.Split(vecSpawnCenter, ",")
		for j := 0; j < len(arraySpawnCenter); j++ {
			v, ret := String2Float(arraySpawnCenter[j])
			if !ret {
				glog.Error("Parse CfgBullet.SpawnCenter field error,value:", arraySpawnCenter[j])
				return false
			}
			data.SpawnCenter = append(data.SpawnCenter, v)
		}

		/* parse SpawnPosition field */
		vecSpawnPosition, _ := parse.GetFieldByName(uint32(i), "spawnPosition")
		arraySpawnPosition := strings.Split(vecSpawnPosition, ",")
		for j := 0; j < len(arraySpawnPosition); j++ {
			v, ret := String2Float(arraySpawnPosition[j])
			if !ret {
				glog.Error("Parse CfgBullet.SpawnPosition field error,value:", arraySpawnPosition[j])
				return false
			}
			data.SpawnPosition = append(data.SpawnPosition, v)
		}

		/* parse Delay field */
		vDelay, _ := parse.GetFieldByName(uint32(i), "delay")
		var DelayRet bool
		data.Delay, DelayRet = String2Float(vDelay)
		if !DelayRet {
			glog.Error("Parse CfgBullet.Delay field error,value:", vDelay)
		}

		/* parse Duration field */
		vDuration, _ := parse.GetFieldByName(uint32(i), "duration")
		var DurationRet bool
		data.Duration, DurationRet = String2Float(vDuration)
		if !DurationRet {
			glog.Error("Parse CfgBullet.Duration field error,value:", vDuration)
		}

		/* parse ParabolaHigh field */
		vParabolaHigh, _ := parse.GetFieldByName(uint32(i), "parabolaHigh")
		var ParabolaHighRet bool
		data.ParabolaHigh, ParabolaHighRet = String2Float(vParabolaHigh)
		if !ParabolaHighRet {
			glog.Error("Parse CfgBullet.ParabolaHigh field error,value:", vParabolaHigh)
		}

		/* parse ParabolaAngle field */
		vParabolaAngle, _ := parse.GetFieldByName(uint32(i), "parabolaAngle")
		var ParabolaAngleRet bool
		data.ParabolaAngle, ParabolaAngleRet = String2Int32(vParabolaAngle)
		if !ParabolaAngleRet {
			glog.Error("Parse CfgBullet.ParabolaAngle field error,value:", vParabolaAngle)
			return false
		}

		/* parse ParabolaOffset field */
		vecParabolaOffset, _ := parse.GetFieldByName(uint32(i), "parabolaOffset")
		arrayParabolaOffset := strings.Split(vecParabolaOffset, ",")
		for j := 0; j < len(arrayParabolaOffset); j++ {
			v, ret := String2Float(arrayParabolaOffset[j])
			if !ret {
				glog.Error("Parse CfgBullet.ParabolaOffset field error,value:", arrayParabolaOffset[j])
				return false
			}
			data.ParabolaOffset = append(data.ParabolaOffset, v)
		}

		/* parse DisorderAngle field */
		vecDisorderAngle, _ := parse.GetFieldByName(uint32(i), "disorderAngle")
		arrayDisorderAngle := strings.Split(vecDisorderAngle, ",")
		for j := 0; j < len(arrayDisorderAngle); j++ {
			v, ret := String2Float(arrayDisorderAngle[j])
			if !ret {
				glog.Error("Parse CfgBullet.DisorderAngle field error,value:", arrayDisorderAngle[j])
				return false
			}
			data.DisorderAngle = append(data.DisorderAngle, v)
		}

		/* parse DisorderSpdAngle field */
		vDisorderSpdAngle, _ := parse.GetFieldByName(uint32(i), "disorderSpdAngle")
		var DisorderSpdAngleRet bool
		data.DisorderSpdAngle, DisorderSpdAngleRet = String2Float(vDisorderSpdAngle)
		if !DisorderSpdAngleRet {
			glog.Error("Parse CfgBullet.DisorderSpdAngle field error,value:", vDisorderSpdAngle)
		}

		/* parse Speed field */
		vSpeed, _ := parse.GetFieldByName(uint32(i), "speed")
		var SpeedRet bool
		data.Speed, SpeedRet = String2Int32(vSpeed)
		if !SpeedRet {
			glog.Error("Parse CfgBullet.Speed field error,value:", vSpeed)
			return false
		}

		/* parse Acceleration field */
		vAcceleration, _ := parse.GetFieldByName(uint32(i), "acceleration")
		var AccelerationRet bool
		data.Acceleration, AccelerationRet = String2Int32(vAcceleration)
		if !AccelerationRet {
			glog.Error("Parse CfgBullet.Acceleration field error,value:", vAcceleration)
			return false
		}

		/* parse Length field */
		vLength, _ := parse.GetFieldByName(uint32(i), "length")
		var LengthRet bool
		data.Length, LengthRet = String2Float(vLength)
		if !LengthRet {
			glog.Error("Parse CfgBullet.Length field error,value:", vLength)
		}

		/* parse PassEffects field */
		vecPassEffects, _ := parse.GetFieldByName(uint32(i), "passEffects")
		if vecPassEffects != "" {
			arrayPassEffects := strings.Split(vecPassEffects, ",")
			for j := 0; j < len(arrayPassEffects); j++ {
				v, ret := String2Int32(arrayPassEffects[j])
				if !ret {
					glog.Error("Parse CfgBullet.PassEffects field error, value:", arrayPassEffects[j])
					return false
				}
				data.PassEffects = append(data.PassEffects, v)
			}
		}

		/* parse HitEffects field */
		vecHitEffects, _ := parse.GetFieldByName(uint32(i), "hitEffects")
		if vecHitEffects != "" {
			arrayHitEffects := strings.Split(vecHitEffects, ",")
			for j := 0; j < len(arrayHitEffects); j++ {
				v, ret := String2Int32(arrayHitEffects[j])
				if !ret {
					glog.Error("Parse CfgBullet.HitEffects field error, value:", arrayHitEffects[j])
					return false
				}
				data.HitEffects = append(data.HitEffects, v)
			}
		}

		/* parse EndEffects field */
		vecEndEffects, _ := parse.GetFieldByName(uint32(i), "endEffects")
		if vecEndEffects != "" {
			arrayEndEffects := strings.Split(vecEndEffects, ",")
			for j := 0; j < len(arrayEndEffects); j++ {
				v, ret := String2Int32(arrayEndEffects[j])
				if !ret {
					glog.Error("Parse CfgBullet.EndEffects field error, value:", arrayEndEffects[j])
					return false
				}
				data.EndEffects = append(data.EndEffects, v)
			}
		}

		/* parse Pass field */
		vPass, _ := parse.GetFieldByName(uint32(i), "pass")
		var PassRet bool
		data.Pass, PassRet = String2Int32(vPass)
		if !PassRet {
			glog.Error("Parse CfgBullet.Pass field error,value:", vPass)
			return false
		}

		/* parse PassTimes field */
		vPassTimes, _ := parse.GetFieldByName(uint32(i), "passTimes")
		var PassTimesRet bool
		data.PassTimes, PassTimesRet = String2Int32(vPassTimes)
		if !PassTimesRet {
			glog.Error("Parse CfgBullet.PassTimes field error,value:", vPassTimes)
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

func (c *CfgBulletConfig) Clear() {
}

func (c *CfgBulletConfig) Find(id int32) (*CfgBullet, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *CfgBulletConfig) GetAllData() map[int32]*CfgBullet {
	return c.data
}

func (c *CfgBulletConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id, ",", v.Type, ",", v.SpawnCenter, ",", v.SpawnPosition, ",", v.Delay, ",", v.Duration, ",", v.ParabolaHigh, ",", v.ParabolaAngle, ",", v.ParabolaOffset, ",", v.DisorderAngle, ",", v.DisorderSpdAngle, ",", v.Speed, ",", v.Acceleration, ",", v.Length, ",", v.PassEffects, ",", v.HitEffects, ",", v.EndEffects, ",", v.Pass, ",", v.PassTimes)
	}
}
