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

type Protocol struct {
	Id       int32
	Function string
	Explain  string
	Type     int32
}

type ProtocolConfig struct {
	data map[int32]*Protocol
}

func NewProtocolConfig() *ProtocolConfig {
	return &ProtocolConfig{
		data: make(map[int32]*Protocol),
	}
}

func (c *ProtocolConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(Protocol)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse Protocol.Id field error,value:", vId)
			return false
		}

		/* parse Function field */
		data.Function, _ = parse.GetFieldByName(uint32(i), "function")

		/* parse Explain field */
		data.Explain, _ = parse.GetFieldByName(uint32(i), "explain")

		/* parse Type field */
		vType, _ := parse.GetFieldByName(uint32(i), "type")
		var TypeRet bool
		data.Type, TypeRet = String2Int32(vType)
		if !TypeRet {
			glog.Error("Parse Protocol.Type field error,value:", vType)
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

func (c *ProtocolConfig) Clear() {
}

func (c *ProtocolConfig) Find(id int32) (*Protocol, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *ProtocolConfig) GetAllData() map[int32]*Protocol {
	return c.data
}

func (c *ProtocolConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id, ",", v.Function, ",", v.Explain, ",", v.Type)
	}
}
