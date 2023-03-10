package grpc

// import (
// 	"context"
//
// 	"github.com/golang/protobuf/proto"
// 	"google.golang.org/grpc"
//
// 	"shared/protobuf/pb"
// 	"shared/utility/errors"
// )
//
// // TODO：暂时不用，写着玩的
// type innerServer struct {
// 	*grpc.Server
// 	methods map[string]grpc.MethodDesc
// 	desc    *grpc.ServiceDesc
// 	svr     interface{}
// 	*pb.UnimplementedRPCShellServer
// }
//
// func NewInnerServer() *innerServer {
// 	s := &innerServer{
// 		Server:  grpc.NewServer(),
// 		methods: map[string]grpc.MethodDesc{},
// 	}
//
// 	pb.RegisterRPCShellServer(s.Server, s)
//
// 	return s
// }
//
// func (s *innerServer) RegisterService(desc *grpc.ServiceDesc, impl interface{}) {
// 	s.desc = desc
// 	s.svr = impl
// 	for i, method := range desc.Methods {
// 		s.methods[method.MethodName] = desc.Methods[i]
// 	}
// 	s.Server.RegisterService(desc, impl)
// }
//
// func (s *innerServer) Call(ctx context.Context, req *pb.RPCCallReq) (*pb.RPCCallResp, error) {
// 	resp, err := s.call(ctx, req.Method, req.In)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	out, err := proto.Marshal(resp.(proto.Message))
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	return &pb.RPCCallResp{
// 		Out: out,
// 	}, nil
// }
//
// func (s *innerServer) call(ctx context.Context, methodName string, in []byte) (interface{}, error) {
// 	method, ok := s.methods[methodName]
// 	if !ok {
// 		return nil, errors.New("not found method")
// 	}
//
// 	dec := func(v interface{}) error {
// 		err := proto.Unmarshal(in, v.(proto.Message))
// 		if err != nil {
// 			return err
// 		}
//
// 		return nil
// 	}
//
// 	return method.Handler(s.svr, ctx, dec, nil)
// }
