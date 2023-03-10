package model

import (
	"shared/protobuf/pb"
	"shared/utility/servertime"
	"time"
)

type SignIn struct {
	SignInGroups map[int32]SignInGroup
}

type SignInGroup struct {
	HasSignIn    bool  `json:"has_signin"`
	SignInType   int32 `json:"type"`
	Record       int32 `json:"record"`
	LastTimeSign int64 `json:"last_time_sign"`
}

func NewSignIn() *SignIn {
	return &SignIn{
		SignInGroups: map[int32]SignInGroup{},
	}
}

func NewSignInGroup() *SignInGroup {
	return &SignInGroup{
		HasSignIn:    false,
		SignInType:   0,
		Record:       0,
		LastTimeSign: 0,
	}
}

//----------------------------------------
//SignIn
//----------------------------------------
// 根据ID来检查这个ID下，是否今天已经签过到了，hasSigned为true，代表已经签到
func (s *SignIn) CheckForSignIn(id int32, signinType int32, tm int64) {

	group, ok := s.SignInGroups[id]
	if !ok {
		group = *NewSignInGroup()
	}
	group.SignInType = signinType

	todayRefreshTime := DailyRefreshTime(time.Unix(tm, 0)).Unix()
	lastTimeSign := group.LastTimeSign
	if lastTimeSign >= todayRefreshTime {
		group.HasSignIn = true
	} else {
		group.HasSignIn = false
	}

	s.SignInGroups[id] = group
	//return group.HasSignIn
}

//----------------------------------------
//SignInGroup
//----------------------------------------
func (sg *SignInGroup) DoSignIn() {
	sg.LastTimeSign = servertime.Now().Unix()
	sg.HasSignIn = true
	sg.Record++
}

func (sg *SignInGroup) VOSignInRecordAndType() *pb.VOSignInRecordAndType {
	//fmt.Printf("signinmodel--------------->record: %v, type: %v, flag: %v", sg.Record, sg.SignInType, sg.HasSignIn)
	return &pb.VOSignInRecordAndType{
		SigninType: sg.SignInType,
		Record:     sg.Record,
		HasSignIn:  sg.HasSignIn,
	}
}
