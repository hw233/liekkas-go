package controller

// var testSession *session.Session
// var testHandler *Handler
// var testReq *pb.TestReq
// var testFunc interface{}
//
// func TestMain(m *testing.M) {
// 	gsRouter = router.NewRouter()
//
// 	config := map[interface{}]string{
// 		1: "Test",
// 	}
//
// 	err := gsRouter.RegisterHandler(&session.Session{}, router.WithConfig(config))
// 	if err != nil {
// 		return
// 	}
//
// 	testSession = &session.Session{
// 		User: &model.User{
// 			UID:  2100,
// 			Target: "雪辙",
// 		},
// 	}
//
// 	testReq = &pb.TestReq{
// 		Uid:  0,
// 		Target: "",
// 	}
//
// 	val, err := proto.Marshal(testReq)
// 	if err != nil {
// 		return
// 	}
//
// 	testHandler = NewPortalHandler(val, testSession)
//
// 	f, ok := gsRouter.Route(1)
// 	if !ok {
// 		return
// 	}
//
// 	testFunc = f
//
// 	m.Run()
// }
//
// func TestHandler_Call(t *testing.T) {
// 	err := testHandler.Call(context.Background(), testFunc)
// 	if err != nil {
// 		t.Errorf("proto.Marshal(testReq) error: %v", err)
// 	}
//
// 	testResp := &pb.TestResp{}
// 	err = proto.Unmarshal(testHandler.Out(), testResp)
// 	if err != nil {
// 		t.Errorf("proto.Unmarshal(testResp) error: %v", err)
// 	}
//
// 	if testResp.Uid != 2100 || testResp.Target != "雪辙" {
// 		t.Errorf("response data wrong, testResp.Uid<%d> !=2100, testResp.Target<%s> != 雪辙", testResp.Uid, testResp.Target)
// 	}
// }
//
// // goos: darwin
// // goarch: amd64
// // pkg: gamesvr/controller
// // cpu: Intel(R) Core(TM) i7-1068NG7 CPU @ 2.30GHz
//
// // BenchmarkHandler_DirectCall-8    	28591557	        41.88 ns/op
// func BenchmarkHandler_DirectCall(b *testing.B) {
// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		_, _ = testFunc.(func(*session.Session, *pb.TestReq) (*pb.TestResp, error))(testSession, testReq)
// 	}
// }
//
// // BenchmarkHandler_ReflectCall-8   	 1889206	       628.8 ns/op
// func BenchmarkHandler_ReflectCall(b *testing.B) {
// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		_ = testHandler.Call(context.Background(), testFunc)
// 	}
// }
