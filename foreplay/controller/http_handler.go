package controller

import (
	"encoding/json"
	"net/http"

	"foreplay/csv/entry"
	"foreplay/manager"
	"shared/common"
	"shared/utility/errors"
	"shared/utility/glog"
)

func RegisterHttpHandler() {
	http.Handle("/patch", NewPatchHandler())
	http.Handle("/announcement", NewAnnouncementHandler())
	http.Handle("/caution", NewCautionHandler())
	http.Handle("/maintain", NewMaintainHandler())
}

func writerErrorResponse(writer http.ResponseWriter, err error) {
	errCode := errors.Code(err)
	if errCode == 0 {
		errCode = -1
	}

	resp := struct {
		ErrCode int    `json:"err_code"`
		ErrMsg  string `json:"err_msg"`
	}{
		ErrCode: errCode,
		ErrMsg:  err.Error(),
	}

	bytes, err := json.Marshal(resp)
	if err != nil {
		glog.Errorf("marshal response error: %+v", err)
		return
	}

	writer.Write(bytes)
}

type PatchHandler struct{}

func NewPatchHandler() *PatchHandler {
	return &PatchHandler{}
}

func (ph *PatchHandler) ServeHTTP(respWriter http.ResponseWriter, req *http.Request) {
	req.ParseForm()

	channel := req.Form.Get("channel")
	appVersion := req.Form.Get("app_version")

	glog.Debugf("new patch request, channel: %s, app_version: %s", channel, appVersion)

	appCfg, err := manager.CSV.AppList.GetAppCfg(channel)
	if err != nil {
		writerErrorResponse(respWriter, err)
		return
	}

	resp := struct {
		ErrCode int32           `json:"err_code"`
		ErrMsg  string          `json:"err_msg"`
		App     *entry.AppCfg   `json:"app"`
		Patch   *entry.PatchCfg `json:"patch"`
	}{
		ErrCode: 0,
		ErrMsg:  "success",
		App:     appCfg,
	}

	patchCfg, err := manager.CSV.Patches.GetPatchCfg(channel, appVersion)
	if err == nil {
		resp.Patch = patchCfg
	}

	bytes, err := json.Marshal(resp)
	if err != nil {
		glog.Errorf("marshal response error: %+v", err)
		return
	}

	respWriter.Write(bytes)
}

type AnnouncementHandler struct {
}

func NewAnnouncementHandler() *AnnouncementHandler {
	return &AnnouncementHandler{}
}

func (ah *AnnouncementHandler) ServeHTTP(respWriter http.ResponseWriter, req *http.Request) {
	anncs, banners := manager.Announcements.GetAnnouncements()

	resp := struct {
		ErrCode       int32                  `json:"err_code"`
		ErrMsg        string                 `json:"err_msg"`
		Announcements []*common.Announcement `json:"announcements"`
		Banners       []*common.Banner       `json:"banners"`
	}{
		ErrCode:       0,
		ErrMsg:        "success",
		Announcements: anncs,
		Banners:       banners,
	}

	bytes, err := json.Marshal(resp)
	if err != nil {
		glog.Errorf("marshal response error: %+v", err)
		return
	}

	respWriter.Write(bytes)
}

type CautionHandler struct {
}

func NewCautionHandler() *CautionHandler {
	return &CautionHandler{}
}

func (ah *CautionHandler) ServeHTTP(respWriter http.ResponseWriter, req *http.Request) {
	cautions := manager.Announcements.GetCautions()

	resp := struct {
		ErrCode  int32             `json:"err_code"`
		ErrMsg   string            `json:"err_msg"`
		Cautions []*common.Caution `json:"cautions"`
	}{
		ErrCode:  0,
		ErrMsg:   "success",
		Cautions: cautions,
	}

	bytes, err := json.Marshal(resp)
	if err != nil {
		glog.Errorf("marshal response error: %+v", err)
		return
	}

	respWriter.Write(bytes)
}

type MaintainHandler struct {
}

func NewMaintainHandler() *MaintainHandler {
	return &MaintainHandler{}
}

func (ah *MaintainHandler) ServeHTTP(respWriter http.ResponseWriter, req *http.Request) {
	isMaintain := manager.Server.IsMaintain()

	resp := struct {
		ErrCode    int32  `json:"err_code"`
		ErrMsg     string `json:"err_msg"`
		IsMaintain bool   `json:"is_maintain"`
	}{
		ErrCode:    0,
		ErrMsg:     "success",
		IsMaintain: isMaintain,
	}

	bytes, err := json.Marshal(resp)
	if err != nil {
		glog.Errorf("marshal response error: %+v", err)
		return
	}

	respWriter.Write(bytes)
}
