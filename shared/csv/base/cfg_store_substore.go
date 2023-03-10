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

type CfgStoreSubstore struct {
	Id                int32
	UnlockCondition   []string
	Cell0             []int32
	UnlockCondition0  []string
	Cell1             []int32
	UnlockCondition1  []string
	Cell2             []int32
	UnlockCondition2  []string
	Cell3             []int32
	UnlockCondition3  []string
	Cell4             []int32
	UnlockCondition4  []string
	Cell5             []int32
	UnlockCondition5  []string
	Cell6             []int32
	UnlockCondition6  []string
	Cell7             []int32
	UnlockCondition7  []string
	Cell8             []int32
	UnlockCondition8  []string
	Cell9             []int32
	UnlockCondition9  []string
	Cell10            []int32
	UnlockCondition10 []string
	Cell11            []int32
	UnlockCondition11 []string
	Cell12            []int32
	UnlockCondition12 []string
	Cell13            []int32
	UnlockCondition13 []string
	Cell14            []int32
	UnlockCondition14 []string
	Cell15            []int32
	UnlockCondition15 []string
	Cell16            []int32
	UnlockCondition16 []string
	Cell17            []int32
	UnlockCondition17 []string
	Cell18            []int32
	UnlockCondition18 []string
	Cell19            []int32
	UnlockCondition19 []string
}

type CfgStoreSubstoreConfig struct {
	data map[int32]*CfgStoreSubstore
}

func NewCfgStoreSubstoreConfig() *CfgStoreSubstoreConfig {
	return &CfgStoreSubstoreConfig{
		data: make(map[int32]*CfgStoreSubstore),
	}
}

func (c *CfgStoreSubstoreConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(CfgStoreSubstore)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse CfgStoreSubstore.Id field error,value:", vId)
			return false
		}

		/* parse UnlockCondition field */
		vecUnlockCondition, _ := parse.GetFieldByName(uint32(i), "unlockCondition")
		arrayUnlockCondition := strings.Split(vecUnlockCondition, ",")
		for j := 0; j < len(arrayUnlockCondition); j++ {
			v := arrayUnlockCondition[j]
			data.UnlockCondition = append(data.UnlockCondition, v)
		}

		/* parse Cell0 field */
		vecCell0, _ := parse.GetFieldByName(uint32(i), "cell0")
		if vecCell0 != "" {
			arrayCell0 := strings.Split(vecCell0, ",")
			for j := 0; j < len(arrayCell0); j++ {
				v, ret := String2Int32(arrayCell0[j])
				if !ret {
					glog.Error("Parse CfgStoreSubstore.Cell0 field error, value:", arrayCell0[j])
					return false
				}
				data.Cell0 = append(data.Cell0, v)
			}
		}

		/* parse UnlockCondition0 field */
		vecUnlockCondition0, _ := parse.GetFieldByName(uint32(i), "unlockCondition0")
		arrayUnlockCondition0 := strings.Split(vecUnlockCondition0, ",")
		for j := 0; j < len(arrayUnlockCondition0); j++ {
			v := arrayUnlockCondition0[j]
			data.UnlockCondition0 = append(data.UnlockCondition0, v)
		}

		/* parse Cell1 field */
		vecCell1, _ := parse.GetFieldByName(uint32(i), "cell1")
		if vecCell1 != "" {
			arrayCell1 := strings.Split(vecCell1, ",")
			for j := 0; j < len(arrayCell1); j++ {
				v, ret := String2Int32(arrayCell1[j])
				if !ret {
					glog.Error("Parse CfgStoreSubstore.Cell1 field error, value:", arrayCell1[j])
					return false
				}
				data.Cell1 = append(data.Cell1, v)
			}
		}

		/* parse UnlockCondition1 field */
		vecUnlockCondition1, _ := parse.GetFieldByName(uint32(i), "unlockCondition1")
		arrayUnlockCondition1 := strings.Split(vecUnlockCondition1, ",")
		for j := 0; j < len(arrayUnlockCondition1); j++ {
			v := arrayUnlockCondition1[j]
			data.UnlockCondition1 = append(data.UnlockCondition1, v)
		}

		/* parse Cell2 field */
		vecCell2, _ := parse.GetFieldByName(uint32(i), "cell2")
		if vecCell2 != "" {
			arrayCell2 := strings.Split(vecCell2, ",")
			for j := 0; j < len(arrayCell2); j++ {
				v, ret := String2Int32(arrayCell2[j])
				if !ret {
					glog.Error("Parse CfgStoreSubstore.Cell2 field error, value:", arrayCell2[j])
					return false
				}
				data.Cell2 = append(data.Cell2, v)
			}
		}

		/* parse UnlockCondition2 field */
		vecUnlockCondition2, _ := parse.GetFieldByName(uint32(i), "unlockCondition2")
		arrayUnlockCondition2 := strings.Split(vecUnlockCondition2, ",")
		for j := 0; j < len(arrayUnlockCondition2); j++ {
			v := arrayUnlockCondition2[j]
			data.UnlockCondition2 = append(data.UnlockCondition2, v)
		}

		/* parse Cell3 field */
		vecCell3, _ := parse.GetFieldByName(uint32(i), "cell3")
		if vecCell3 != "" {
			arrayCell3 := strings.Split(vecCell3, ",")
			for j := 0; j < len(arrayCell3); j++ {
				v, ret := String2Int32(arrayCell3[j])
				if !ret {
					glog.Error("Parse CfgStoreSubstore.Cell3 field error, value:", arrayCell3[j])
					return false
				}
				data.Cell3 = append(data.Cell3, v)
			}
		}

		/* parse UnlockCondition3 field */
		vecUnlockCondition3, _ := parse.GetFieldByName(uint32(i), "unlockCondition3")
		arrayUnlockCondition3 := strings.Split(vecUnlockCondition3, ",")
		for j := 0; j < len(arrayUnlockCondition3); j++ {
			v := arrayUnlockCondition3[j]
			data.UnlockCondition3 = append(data.UnlockCondition3, v)
		}

		/* parse Cell4 field */
		vecCell4, _ := parse.GetFieldByName(uint32(i), "cell4")
		if vecCell4 != "" {
			arrayCell4 := strings.Split(vecCell4, ",")
			for j := 0; j < len(arrayCell4); j++ {
				v, ret := String2Int32(arrayCell4[j])
				if !ret {
					glog.Error("Parse CfgStoreSubstore.Cell4 field error, value:", arrayCell4[j])
					return false
				}
				data.Cell4 = append(data.Cell4, v)
			}
		}

		/* parse UnlockCondition4 field */
		vecUnlockCondition4, _ := parse.GetFieldByName(uint32(i), "unlockCondition4")
		arrayUnlockCondition4 := strings.Split(vecUnlockCondition4, ",")
		for j := 0; j < len(arrayUnlockCondition4); j++ {
			v := arrayUnlockCondition4[j]
			data.UnlockCondition4 = append(data.UnlockCondition4, v)
		}

		/* parse Cell5 field */
		vecCell5, _ := parse.GetFieldByName(uint32(i), "cell5")
		if vecCell5 != "" {
			arrayCell5 := strings.Split(vecCell5, ",")
			for j := 0; j < len(arrayCell5); j++ {
				v, ret := String2Int32(arrayCell5[j])
				if !ret {
					glog.Error("Parse CfgStoreSubstore.Cell5 field error, value:", arrayCell5[j])
					return false
				}
				data.Cell5 = append(data.Cell5, v)
			}
		}

		/* parse UnlockCondition5 field */
		vecUnlockCondition5, _ := parse.GetFieldByName(uint32(i), "unlockCondition5")
		arrayUnlockCondition5 := strings.Split(vecUnlockCondition5, ",")
		for j := 0; j < len(arrayUnlockCondition5); j++ {
			v := arrayUnlockCondition5[j]
			data.UnlockCondition5 = append(data.UnlockCondition5, v)
		}

		/* parse Cell6 field */
		vecCell6, _ := parse.GetFieldByName(uint32(i), "cell6")
		if vecCell6 != "" {
			arrayCell6 := strings.Split(vecCell6, ",")
			for j := 0; j < len(arrayCell6); j++ {
				v, ret := String2Int32(arrayCell6[j])
				if !ret {
					glog.Error("Parse CfgStoreSubstore.Cell6 field error, value:", arrayCell6[j])
					return false
				}
				data.Cell6 = append(data.Cell6, v)
			}
		}

		/* parse UnlockCondition6 field */
		vecUnlockCondition6, _ := parse.GetFieldByName(uint32(i), "unlockCondition6")
		arrayUnlockCondition6 := strings.Split(vecUnlockCondition6, ",")
		for j := 0; j < len(arrayUnlockCondition6); j++ {
			v := arrayUnlockCondition6[j]
			data.UnlockCondition6 = append(data.UnlockCondition6, v)
		}

		/* parse Cell7 field */
		vecCell7, _ := parse.GetFieldByName(uint32(i), "cell7")
		if vecCell7 != "" {
			arrayCell7 := strings.Split(vecCell7, ",")
			for j := 0; j < len(arrayCell7); j++ {
				v, ret := String2Int32(arrayCell7[j])
				if !ret {
					glog.Error("Parse CfgStoreSubstore.Cell7 field error, value:", arrayCell7[j])
					return false
				}
				data.Cell7 = append(data.Cell7, v)
			}
		}

		/* parse UnlockCondition7 field */
		vecUnlockCondition7, _ := parse.GetFieldByName(uint32(i), "unlockCondition7")
		arrayUnlockCondition7 := strings.Split(vecUnlockCondition7, ",")
		for j := 0; j < len(arrayUnlockCondition7); j++ {
			v := arrayUnlockCondition7[j]
			data.UnlockCondition7 = append(data.UnlockCondition7, v)
		}

		/* parse Cell8 field */
		vecCell8, _ := parse.GetFieldByName(uint32(i), "cell8")
		if vecCell8 != "" {
			arrayCell8 := strings.Split(vecCell8, ",")
			for j := 0; j < len(arrayCell8); j++ {
				v, ret := String2Int32(arrayCell8[j])
				if !ret {
					glog.Error("Parse CfgStoreSubstore.Cell8 field error, value:", arrayCell8[j])
					return false
				}
				data.Cell8 = append(data.Cell8, v)
			}
		}

		/* parse UnlockCondition8 field */
		vecUnlockCondition8, _ := parse.GetFieldByName(uint32(i), "unlockCondition8")
		arrayUnlockCondition8 := strings.Split(vecUnlockCondition8, ",")
		for j := 0; j < len(arrayUnlockCondition8); j++ {
			v := arrayUnlockCondition8[j]
			data.UnlockCondition8 = append(data.UnlockCondition8, v)
		}

		/* parse Cell9 field */
		vecCell9, _ := parse.GetFieldByName(uint32(i), "cell9")
		if vecCell9 != "" {
			arrayCell9 := strings.Split(vecCell9, ",")
			for j := 0; j < len(arrayCell9); j++ {
				v, ret := String2Int32(arrayCell9[j])
				if !ret {
					glog.Error("Parse CfgStoreSubstore.Cell9 field error, value:", arrayCell9[j])
					return false
				}
				data.Cell9 = append(data.Cell9, v)
			}
		}

		/* parse UnlockCondition9 field */
		vecUnlockCondition9, _ := parse.GetFieldByName(uint32(i), "unlockCondition9")
		arrayUnlockCondition9 := strings.Split(vecUnlockCondition9, ",")
		for j := 0; j < len(arrayUnlockCondition9); j++ {
			v := arrayUnlockCondition9[j]
			data.UnlockCondition9 = append(data.UnlockCondition9, v)
		}

		/* parse Cell10 field */
		vecCell10, _ := parse.GetFieldByName(uint32(i), "cell10")
		if vecCell10 != "" {
			arrayCell10 := strings.Split(vecCell10, ",")
			for j := 0; j < len(arrayCell10); j++ {
				v, ret := String2Int32(arrayCell10[j])
				if !ret {
					glog.Error("Parse CfgStoreSubstore.Cell10 field error, value:", arrayCell10[j])
					return false
				}
				data.Cell10 = append(data.Cell10, v)
			}
		}

		/* parse UnlockCondition10 field */
		vecUnlockCondition10, _ := parse.GetFieldByName(uint32(i), "unlockCondition10")
		arrayUnlockCondition10 := strings.Split(vecUnlockCondition10, ",")
		for j := 0; j < len(arrayUnlockCondition10); j++ {
			v := arrayUnlockCondition10[j]
			data.UnlockCondition10 = append(data.UnlockCondition10, v)
		}

		/* parse Cell11 field */
		vecCell11, _ := parse.GetFieldByName(uint32(i), "cell11")
		if vecCell11 != "" {
			arrayCell11 := strings.Split(vecCell11, ",")
			for j := 0; j < len(arrayCell11); j++ {
				v, ret := String2Int32(arrayCell11[j])
				if !ret {
					glog.Error("Parse CfgStoreSubstore.Cell11 field error, value:", arrayCell11[j])
					return false
				}
				data.Cell11 = append(data.Cell11, v)
			}
		}

		/* parse UnlockCondition11 field */
		vecUnlockCondition11, _ := parse.GetFieldByName(uint32(i), "unlockCondition11")
		arrayUnlockCondition11 := strings.Split(vecUnlockCondition11, ",")
		for j := 0; j < len(arrayUnlockCondition11); j++ {
			v := arrayUnlockCondition11[j]
			data.UnlockCondition11 = append(data.UnlockCondition11, v)
		}

		/* parse Cell12 field */
		vecCell12, _ := parse.GetFieldByName(uint32(i), "cell12")
		if vecCell12 != "" {
			arrayCell12 := strings.Split(vecCell12, ",")
			for j := 0; j < len(arrayCell12); j++ {
				v, ret := String2Int32(arrayCell12[j])
				if !ret {
					glog.Error("Parse CfgStoreSubstore.Cell12 field error, value:", arrayCell12[j])
					return false
				}
				data.Cell12 = append(data.Cell12, v)
			}
		}

		/* parse UnlockCondition12 field */
		vecUnlockCondition12, _ := parse.GetFieldByName(uint32(i), "unlockCondition12")
		arrayUnlockCondition12 := strings.Split(vecUnlockCondition12, ",")
		for j := 0; j < len(arrayUnlockCondition12); j++ {
			v := arrayUnlockCondition12[j]
			data.UnlockCondition12 = append(data.UnlockCondition12, v)
		}

		/* parse Cell13 field */
		vecCell13, _ := parse.GetFieldByName(uint32(i), "cell13")
		if vecCell13 != "" {
			arrayCell13 := strings.Split(vecCell13, ",")
			for j := 0; j < len(arrayCell13); j++ {
				v, ret := String2Int32(arrayCell13[j])
				if !ret {
					glog.Error("Parse CfgStoreSubstore.Cell13 field error, value:", arrayCell13[j])
					return false
				}
				data.Cell13 = append(data.Cell13, v)
			}
		}

		/* parse UnlockCondition13 field */
		vecUnlockCondition13, _ := parse.GetFieldByName(uint32(i), "unlockCondition13")
		arrayUnlockCondition13 := strings.Split(vecUnlockCondition13, ",")
		for j := 0; j < len(arrayUnlockCondition13); j++ {
			v := arrayUnlockCondition13[j]
			data.UnlockCondition13 = append(data.UnlockCondition13, v)
		}

		/* parse Cell14 field */
		vecCell14, _ := parse.GetFieldByName(uint32(i), "cell14")
		if vecCell14 != "" {
			arrayCell14 := strings.Split(vecCell14, ",")
			for j := 0; j < len(arrayCell14); j++ {
				v, ret := String2Int32(arrayCell14[j])
				if !ret {
					glog.Error("Parse CfgStoreSubstore.Cell14 field error, value:", arrayCell14[j])
					return false
				}
				data.Cell14 = append(data.Cell14, v)
			}
		}

		/* parse UnlockCondition14 field */
		vecUnlockCondition14, _ := parse.GetFieldByName(uint32(i), "unlockCondition14")
		arrayUnlockCondition14 := strings.Split(vecUnlockCondition14, ",")
		for j := 0; j < len(arrayUnlockCondition14); j++ {
			v := arrayUnlockCondition14[j]
			data.UnlockCondition14 = append(data.UnlockCondition14, v)
		}

		/* parse Cell15 field */
		vecCell15, _ := parse.GetFieldByName(uint32(i), "cell15")
		if vecCell15 != "" {
			arrayCell15 := strings.Split(vecCell15, ",")
			for j := 0; j < len(arrayCell15); j++ {
				v, ret := String2Int32(arrayCell15[j])
				if !ret {
					glog.Error("Parse CfgStoreSubstore.Cell15 field error, value:", arrayCell15[j])
					return false
				}
				data.Cell15 = append(data.Cell15, v)
			}
		}

		/* parse UnlockCondition15 field */
		vecUnlockCondition15, _ := parse.GetFieldByName(uint32(i), "unlockCondition15")
		arrayUnlockCondition15 := strings.Split(vecUnlockCondition15, ",")
		for j := 0; j < len(arrayUnlockCondition15); j++ {
			v := arrayUnlockCondition15[j]
			data.UnlockCondition15 = append(data.UnlockCondition15, v)
		}

		/* parse Cell16 field */
		vecCell16, _ := parse.GetFieldByName(uint32(i), "cell16")
		if vecCell16 != "" {
			arrayCell16 := strings.Split(vecCell16, ",")
			for j := 0; j < len(arrayCell16); j++ {
				v, ret := String2Int32(arrayCell16[j])
				if !ret {
					glog.Error("Parse CfgStoreSubstore.Cell16 field error, value:", arrayCell16[j])
					return false
				}
				data.Cell16 = append(data.Cell16, v)
			}
		}

		/* parse UnlockCondition16 field */
		vecUnlockCondition16, _ := parse.GetFieldByName(uint32(i), "unlockCondition16")
		arrayUnlockCondition16 := strings.Split(vecUnlockCondition16, ",")
		for j := 0; j < len(arrayUnlockCondition16); j++ {
			v := arrayUnlockCondition16[j]
			data.UnlockCondition16 = append(data.UnlockCondition16, v)
		}

		/* parse Cell17 field */
		vecCell17, _ := parse.GetFieldByName(uint32(i), "cell17")
		if vecCell17 != "" {
			arrayCell17 := strings.Split(vecCell17, ",")
			for j := 0; j < len(arrayCell17); j++ {
				v, ret := String2Int32(arrayCell17[j])
				if !ret {
					glog.Error("Parse CfgStoreSubstore.Cell17 field error, value:", arrayCell17[j])
					return false
				}
				data.Cell17 = append(data.Cell17, v)
			}
		}

		/* parse UnlockCondition17 field */
		vecUnlockCondition17, _ := parse.GetFieldByName(uint32(i), "unlockCondition17")
		arrayUnlockCondition17 := strings.Split(vecUnlockCondition17, ",")
		for j := 0; j < len(arrayUnlockCondition17); j++ {
			v := arrayUnlockCondition17[j]
			data.UnlockCondition17 = append(data.UnlockCondition17, v)
		}

		/* parse Cell18 field */
		vecCell18, _ := parse.GetFieldByName(uint32(i), "cell18")
		if vecCell18 != "" {
			arrayCell18 := strings.Split(vecCell18, ",")
			for j := 0; j < len(arrayCell18); j++ {
				v, ret := String2Int32(arrayCell18[j])
				if !ret {
					glog.Error("Parse CfgStoreSubstore.Cell18 field error, value:", arrayCell18[j])
					return false
				}
				data.Cell18 = append(data.Cell18, v)
			}
		}

		/* parse UnlockCondition18 field */
		vecUnlockCondition18, _ := parse.GetFieldByName(uint32(i), "unlockCondition18")
		arrayUnlockCondition18 := strings.Split(vecUnlockCondition18, ",")
		for j := 0; j < len(arrayUnlockCondition18); j++ {
			v := arrayUnlockCondition18[j]
			data.UnlockCondition18 = append(data.UnlockCondition18, v)
		}

		/* parse Cell19 field */
		vecCell19, _ := parse.GetFieldByName(uint32(i), "cell19")
		if vecCell19 != "" {
			arrayCell19 := strings.Split(vecCell19, ",")
			for j := 0; j < len(arrayCell19); j++ {
				v, ret := String2Int32(arrayCell19[j])
				if !ret {
					glog.Error("Parse CfgStoreSubstore.Cell19 field error, value:", arrayCell19[j])
					return false
				}
				data.Cell19 = append(data.Cell19, v)
			}
		}

		/* parse UnlockCondition19 field */
		vecUnlockCondition19, _ := parse.GetFieldByName(uint32(i), "unlockCondition19")
		arrayUnlockCondition19 := strings.Split(vecUnlockCondition19, ",")
		for j := 0; j < len(arrayUnlockCondition19); j++ {
			v := arrayUnlockCondition19[j]
			data.UnlockCondition19 = append(data.UnlockCondition19, v)
		}

		if _, ok := c.data[data.Id]; ok {
			glog.Errorf("Find %d repeated", data.Id)
			return false
		}
		c.data[data.Id] = data
	}

	return true
}

func (c *CfgStoreSubstoreConfig) Clear() {
}

func (c *CfgStoreSubstoreConfig) Find(id int32) (*CfgStoreSubstore, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *CfgStoreSubstoreConfig) GetAllData() map[int32]*CfgStoreSubstore {
	return c.data
}

func (c *CfgStoreSubstoreConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id, ",", v.UnlockCondition, ",", v.Cell0, ",", v.UnlockCondition0, ",", v.Cell1, ",", v.UnlockCondition1, ",", v.Cell2, ",", v.UnlockCondition2, ",", v.Cell3, ",", v.UnlockCondition3, ",", v.Cell4, ",", v.UnlockCondition4, ",", v.Cell5, ",", v.UnlockCondition5, ",", v.Cell6, ",", v.UnlockCondition6, ",", v.Cell7, ",", v.UnlockCondition7, ",", v.Cell8, ",", v.UnlockCondition8, ",", v.Cell9, ",", v.UnlockCondition9, ",", v.Cell10, ",", v.UnlockCondition10, ",", v.Cell11, ",", v.UnlockCondition11, ",", v.Cell12, ",", v.UnlockCondition12, ",", v.Cell13, ",", v.UnlockCondition13, ",", v.Cell14, ",", v.UnlockCondition14, ",", v.Cell15, ",", v.UnlockCondition15, ",", v.Cell16, ",", v.UnlockCondition16, ",", v.Cell17, ",", v.UnlockCondition17, ",", v.Cell18, ",", v.UnlockCondition18, ",", v.Cell19, ",", v.UnlockCondition19)
	}
}
