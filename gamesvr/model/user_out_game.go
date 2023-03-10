package model

type UserOutGameInfo struct {
	udid      string //硬件设备号
	sdkUdid   string //sdk设备号
	sdkUid    string //sdk账号
	isTourist bool   //是否游客
	device    string //机型
	osVersion string //系统
	network   string //网络类型
	mac       string //mac地址
	ip        string
}

func NewUserOutGameInfo() *UserOutGameInfo {
	return &UserOutGameInfo{}
}

func (uogi *UserOutGameInfo) SetUdid(udid string) {
	uogi.udid = udid
}

func (uogi *UserOutGameInfo) GetUdid() string {
	return uogi.udid
}

func (uogi *UserOutGameInfo) SetSdkUdid(sdkUdid string) {
	uogi.sdkUdid = sdkUdid
}

func (uogi *UserOutGameInfo) GetSdkUdid() string {
	return uogi.sdkUdid
}

func (uogi *UserOutGameInfo) SetSdkUid(sdkUid string) {
	uogi.sdkUid = sdkUid
}

func (uogi *UserOutGameInfo) GetSdkUid() string {
	return uogi.sdkUid
}

func (uogi *UserOutGameInfo) SetIsTourist(isTourist bool) {
	uogi.isTourist = isTourist
}

func (uogi *UserOutGameInfo) GetIsTourist() bool {
	return uogi.isTourist
}

func (uogi *UserOutGameInfo) SetDevice(device string) {
	uogi.device = device
}

func (uogi *UserOutGameInfo) GetDevice() string {
	return uogi.device
}

func (uogi *UserOutGameInfo) SetOsVersion(osVersion string) {
	uogi.osVersion = osVersion
}

func (uogi *UserOutGameInfo) GetOsVersion() string {
	return uogi.osVersion
}

func (uogi *UserOutGameInfo) SetNetwork(network string) {
	uogi.network = network
}

func (uogi *UserOutGameInfo) GetNetwork() string {
	return uogi.network
}

func (uogi *UserOutGameInfo) SetMac(mac string) {
	uogi.mac = mac
}

func (uogi *UserOutGameInfo) GetMac() string {
	return uogi.mac
}

func (uogi *UserOutGameInfo) SetIp(ip string) {
	uogi.ip = ip
}

func (uogi *UserOutGameInfo) GetIp() string {
	return uogi.ip
}
