package session

import (
	"context"
	"shared/common"
	"shared/protobuf/pb"
	"shared/utility/errors"
)

func (s *Session) GetSignInInfo(ctx context.Context, req *pb.C2SGetSignInInfo) (*pb.S2CGetSignInInfo, error) {
	ids, err := s.User.CheckHasSignInWrap()
	if err != nil {
		return nil, err
	}

	infos := []*pb.VOSignInInfo{}

	for _, id := range ids {
		signinGroup, ok := s.User.Info.SignIn.SignInGroups[int32(id)]
		if !ok {
			return nil, errors.Swrapf(common.ErrSignInWrongIDForSignInGroups, id)
		}
		info := pb.VOSignInInfo{
			SigninId:      int32(id),
			RecordAndType: signinGroup.VOSignInRecordAndType(),
		}
		infos = append(infos, &info)
	}

	return &pb.S2CGetSignInInfo{SigninInfos: infos}, nil
}

func (s *Session) SignIn(ctx context.Context, req *pb.C2SSignIn) (*pb.S2CSignIn, error) {
	id := req.SigninId
	err := s.User.SignIn(id)
	if err != nil {
		return nil, err
	}

	signinGroup, ok := s.User.Info.SignIn.SignInGroups[id]
	if !ok {
		return nil, errors.Swrapf(common.ErrSignInWrongIDForSignInGroups, id)
	}

	info := pb.VOSignInInfo{
		SigninId:      id,
		RecordAndType: signinGroup.VOSignInRecordAndType(),
	}

	return &pb.S2CSignIn{
		SigninInfo:     &info,
		ResourceResult: s.User.VOResourceResult(),
	}, nil
}
