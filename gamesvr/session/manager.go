package session

// var (
// 	SessManager session.Manager = session.NewManager(&Builder{}, &session.ManagerConfig{
// 		Expire:   15 * time.Minute,
// 		Capacity: 100000,
// 	})
// )

// func GetSession(uid int64) (*Session, bool) {
// 	lock.RLock()
// 	defer lock.RUnlock()
//

// 	s, ok := active[uid]
// 	return s, ok
// }
//
// func GetActiveSession() map[int64]*Session {
// 	return active
// }
//
// func NewSessionIfNotExist(uid int64) (*Session, error) {
// 	lock.Lock()
// 	defer lock.Unlock()
//
// 	_, ok := active[uid]
// 	if !ok {
// 		err := manager.RPCGameServer.SetRecord(uid)
// 		if err != nil {
// 			log.Printf("GetSession: record error: %v", err)
// 			return nil, common.ErrConnectUserInOtherGameSvr
// 		}
//
// 		newS, err := NewSession(uid)
// 		if err != nil {
// 			log.Printf("ERROR: new session error: %v", err)
// 			return nil, err
// 		}
//
// 		active[uid] = newS
//
// 		// err = manager.RPCGameServer.SetBalancerCount(len(active))
// 		// if err != nil {
// 		// 	log.Printf("set active count error: %v", err)
// 		// }
// 	}
//
// 	return active[uid], nil
// }
